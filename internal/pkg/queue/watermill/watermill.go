package watermillqueue

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/email"
	"github.com/ayinke-llc/malak/internal/pkg/queue"
	wotelfloss "github.com/dentech-floss/watermill-opentelemetry-go-extra/pkg/opentelemetry"
	"github.com/garsue/watermillzap"
	"github.com/google/uuid"
	redis "github.com/redis/go-redis/v9"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/client"
	wotel "github.com/voi-oss/watermill-opentelemetry/pkg/opentelemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

var tracer = otel.Tracer("watermill")

type WatermillClient struct {
	publisher  *message.Router
	susbcriber *redisstream.Subscriber
	messager   message.Publisher
	logger     *zap.Logger

	userRepo      malak.UserRepository
	workspaceRepo malak.WorkspaceRepository
	updateRepo    malak.UpdateRepository
	contactRepo   malak.ContactRepository
	cfg           config.Config
	emailClient   email.Client
	stripeClient  *client.API
}

func New(redisClient *redis.Client,
	cfg config.Config,
	logger *zap.Logger,
	emailClient email.Client,
	userRepo malak.UserRepository,
	workspaceRepo malak.WorkspaceRepository,
	updateRepo malak.UpdateRepository,
	contactRepo malak.ContactRepository,
	stripeClient *client.API) (queue.QueueHandler, error) {

	p, err := redisstream.NewPublisher(
		redisstream.PublisherConfig{
			Client:     redisClient,
			Marshaller: redisstream.DefaultMarshallerUnmarshaller{},
		},
		watermillzap.NewLogger(logger))
	if err != nil {
		return nil, err
	}

	publisher := wotel.NewNamedPublisherDecorator("queue.Publish",
		wotelfloss.NewTracePropagatingPublisherDecorator(p))

	subscriber, err := redisstream.NewSubscriber(
		redisstream.SubscriberConfig{
			Client:       redisClient,
			Unmarshaller: redisstream.DefaultMarshallerUnmarshaller{},
		},
		watermillzap.NewLogger(logger))
	if err != nil {
		return nil, err
	}

	router, err := message.NewRouter(message.RouterConfig{},
		watermillzap.NewLogger(logger))
	if err != nil {
		return nil, err
	}

	router.AddPlugin(plugin.SignalsHandler)

	poisionQueue, err := middleware.PoisonQueue(publisher, "poision.queue")
	if err != nil {
		return nil, err
	}

	router.AddMiddleware(
		middleware.CorrelationID,

		middleware.Retry{
			MaxRetries:      3,
			InitialInterval: time.Millisecond * 100,
			Logger:          watermill.NewStdLogger(false, false),
		}.Middleware,
		poisionQueue,
		// Recoverer handles panics from handlers.
		// In this case, it passes them as errors to the Retry middleware.
		middleware.Recoverer,
		// OTEL
		wotelfloss.ExtractRemoteParentSpanContext(),
		wotel.Trace(),
	)

	t := &WatermillClient{
		cfg:           cfg,
		publisher:     router,
		logger:        logger,
		messager:      publisher,
		susbcriber:    subscriber,
		userRepo:      userRepo,
		workspaceRepo: workspaceRepo,
		updateRepo:    updateRepo,
		contactRepo:   contactRepo,
		emailClient:   emailClient,
		stripeClient:  stripeClient,
	}

	t.setUpRoutes(router, subscriber)

	return t, nil
}

func (t *WatermillClient) setUpRoutes(router *message.Router,
	subscriber *redisstream.Subscriber) {

	router.AddNoPublisherHandler(
		queue.QueueTopicBillingCreateCustomer.String(),
		queue.QueueTopicBillingCreateCustomer.String(),
		subscriber,
		t.createStripeCustomer,
	)
}

func (t *WatermillClient) Add(ctx context.Context,
	topic queue.QueueTopic, data any) error {
	return t.messager.Publish(
		topic.String(), message.NewMessage(uuid.NewString(),
			queue.ToPayload(data)))
}

func (t *WatermillClient) Start(context.Context) {
	_ = t.publisher.Run(context.Background())
}

func (t *WatermillClient) Close() error { return t.publisher.Close() }

func (t *WatermillClient) createStripeCustomer(msg *message.Message) error {

	if !t.cfg.Billing.IsEnabled {
		return nil
	}

	var opts queue.BillingCreateCustomerOptions

	ctx, span := tracer.Start(context.Background(),
		"createStripeCustomer")

	defer span.End()

	if err := json.NewDecoder(bytes.NewBuffer(msg.Payload)).
		Decode(&opts); err != nil {
		return err
	}

	logger := t.logger.With(zap.String("method", "createStripeCustomer"),
		zap.String("workspace_id", opts.Workspace.ID.String()))

	logger.Debug("creating stripe customer")

	params := &stripe.CustomerParams{
		Name:  hermes.Ref(opts.Workspace.WorkspaceName),
		Email: hermes.Ref(opts.Email.String()),
	}

	customer, err := t.stripeClient.Customers.New(params)
	if err != nil {
		logger.Error("could not create new customer for this workspace",
			zap.Error(err))

		span.RecordError(err)
		span.SetStatus(codes.Error, "could not create stripe customer")
		return err
	}

	opts.Workspace.StripeCustomerID = customer.ID
	if err := t.workspaceRepo.Update(ctx, opts.Workspace); err != nil {
		logger.Error("could not update workspace with stripe customer id",
			zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, "could not update workspace")
		return err
	}

	return nil
}
