package watermillqueue

import (
	"context"
	"os"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/internal/pkg/queue"
	redis "github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
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
}

func New(redisClient *redis.Client,
	logger *zap.Logger,
	userRepo malak.UserRepository,
	workspaceRepo malak.WorkspaceRepository,
) (queue.QueueHandler, error) {
	publisher, err := redisstream.NewPublisher(
		redisstream.PublisherConfig{
			Client:     redisClient,
			Marshaller: redisstream.DefaultMarshallerUnmarshaller{},
		},
		watermill.NewStdLoggerWithOut(os.Stdout, false, false))
	if err != nil {
		return nil, err
	}

	subscriber, err := redisstream.NewSubscriber(
		redisstream.SubscriberConfig{
			Client:       redisClient,
			Unmarshaller: redisstream.DefaultMarshallerUnmarshaller{},
		},
		watermill.NewStdLoggerWithOut(os.Stdout, false, false),
	)
	if err != nil {
		return nil, err
	}

	router, err := message.NewRouter(message.RouterConfig{}, watermill.NewStdLogger(false, false))
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

func (t *WatermillClient) Add(ctx context.Context, topic string, msg *message.Message) error {
	return t.messager.Publish(topic, msg)
}

func (t *WatermillClient) Start(context.Context) {
	t.publisher.Run(context.Background())
}

func (t *WatermillClient) Close() error { return t.publisher.Close() }

func (t *WatermillClient) sendPreviewEmail(msg *message.Message) error {

	return nil
}
