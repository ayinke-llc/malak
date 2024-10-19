package watermillqueue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"text/template"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/email"
	"github.com/ayinke-llc/malak/internal/pkg/queue"
	wotelfloss "github.com/dentech-floss/watermill-opentelemetry-go-extra/pkg/opentelemetry"
	"github.com/garsue/watermillzap"
	redis "github.com/redis/go-redis/v9"
	wotel "github.com/voi-oss/watermill-opentelemetry/pkg/opentelemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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
}

func New(redisClient *redis.Client,
	cfg config.Config,
	logger *zap.Logger,
	emailClient email.Client,
	userRepo malak.UserRepository,
	workspaceRepo malak.WorkspaceRepository,
	updateRepo malak.UpdateRepository,
	contactRepo malak.ContactRepository) (queue.QueueHandler, error) {

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

	if msg.Metadata == nil {
		msg.Metadata = map[string]string{}
	}

	newMsg := message.NewMessage(msg.ID, msg.Data)

	newMsg.Metadata = msg.Metadata
	return t.messager.Publish(topic, newMsg)
}

func (t *WatermillClient) Start(context.Context) {
	_ = t.publisher.Run(context.Background())
}

func (t *WatermillClient) Close() error { return t.publisher.Close() }

func (t *WatermillClient) sendPreviewEmail(msg *message.Message) error {

	var p queue.PreviewUpdateMessage

	if err := json.NewDecoder(bytes.NewBuffer(msg.Payload)).Decode(&p); err != nil {
		t.logger.Error("could not decode message queue request", zap.Error(err))
		return err
	}

	logger := t.logger.With(
		zap.String("queue.handler", "sendPreviewEmail"),
		zap.String("update_id", p.UpdateID.String()),
		zap.String("schedule_id", p.ScheduleID.String()),
	)

	ctx, span := tracer.Start(context.Background(), "queue.sendPreviewEmail")
	defer span.End()

	span.SetAttributes(
		attribute.Bool("preview", true),
		attribute.String("update_id", p.UpdateID.String()),
		attribute.String("schedule_id", p.ScheduleID.String()),
	)

	update, err := t.updateRepo.Get(ctx, malak.FetchUpdateOptions{
		ID: p.UpdateID,
	})
	if err != nil {
		span.RecordError(err)
		logger.Error("could not fetch update from database",
			zap.Error(err))
		return err
	}

	schedule, err := t.updateRepo.GetSchedule(ctx, p.ScheduleID)
	if err != nil {
		span.RecordError(err)
		logger.Error("could not fetch update schedule from database",
			zap.Error(err))
		return err
	}

	contact, err := t.contactRepo.Get(ctx, malak.FetchContactOptions{
		Email: p.Email,
	})
	if err != nil {
		span.RecordError(err)
		logger.Error("could not fetch contact from database",
			zap.Error(err))
		return err
	}

	span.SetAttributes(
		attribute.String("triggered_user_id", schedule.ScheduledBy.String()))

	templatedFile, err := template.New("template").
		Parse(email.UpdateHTMLEmailTemplate)
	if err != nil {
		span.RecordError(err)
		logger.Error("could not create html template",
			zap.Error(err))
		return err
	}

	var b = new(bytes.Buffer)
	err = templatedFile.Execute(b, map[string]string{
		"Content": update.Content.HTML(),
	})

	if err != nil {
		span.RecordError(err)
		logger.Error("could not parse html template",
			zap.Error(err))
		return err
	}

	sendOptions := email.SendOptions{
		HTML:      b.String(),
		Sender:    t.cfg.Email.Sender,
		Recipient: contact.Email,
		Subject:   fmt.Sprintf("[TEST] %s", update.Title),
		DKIM: struct {
			Sign       bool
			PrivateKey []byte
		}{
			Sign:       false,
			PrivateKey: []byte(""),
		},
	}

	if err := t.emailClient.Send(ctx, sendOptions); err != nil {
		span.RecordError(err)
		logger.Error("could not send preview email",
			zap.Error(err))
		return err
	}

	msg.Ack()
	return nil
}
