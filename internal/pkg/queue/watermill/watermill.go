package watermillqueue

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/internal/pkg/queue"
	wotelfloss "github.com/dentech-floss/watermill-opentelemetry-go-extra/pkg/opentelemetry"
	"github.com/garsue/watermillzap"
	redis "github.com/redis/go-redis/v9"
	wotel "github.com/voi-oss/watermill-opentelemetry/pkg/opentelemetry"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var _ = otel.Tracer("watermill")

type WatermillClient struct {
	publisher  *message.Router
	susbcriber *redisstream.Subscriber
	messager   message.Publisher
	logger     *zap.Logger

	userRepo      malak.UserRepository
	workspaceRepo malak.WorkspaceRepository
}

func New(redisClient *redis.Client,
	logger *zap.Logger,
	userRepo malak.UserRepository,
	workspaceRepo malak.WorkspaceRepository) (queue.QueueHandler, error) {

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
		publisher:     router,
		logger:        logger,
		messager:      publisher,
		susbcriber:    subscriber,
		userRepo:      userRepo,
		workspaceRepo: workspaceRepo,
	}

	router.AddNoPublisherHandler(
		"a",
		queue.QueueEventSubscriptionMessageUpdatePreview.String(),
		subscriber,
		t.sendPreviewEmail,
	)

	return t, nil
}

func (t *WatermillClient) Add(ctx context.Context,
	topic string, msg *queue.Message) error {

	newMsg := message.NewMessage(msg.ID, msg.Data)

	newMsg.Metadata = msg.Metadata
	return t.messager.Publish(topic, newMsg)
}

func (t *WatermillClient) Start(context.Context) {
	_ = t.publisher.Run(context.Background())
}

func (t *WatermillClient) Close() error { return t.publisher.Close() }

func (t *WatermillClient) sendPreviewEmail(msg *message.Message) error {
	return nil
}
