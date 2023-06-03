package ws

import (
	"context"
	"net/http"
	"time"

	"github.com/alash3al/go-pubsub"
	"github.com/f0rmul/sensor-control/internal/models"
	"github.com/f0rmul/sensor-control/pkg/logger"
	"nhooyr.io/websocket"
)

const (
	subscriptionTopic = "snapshots"
)

type Service interface {
	Broker() *pubsub.Broker
}

type PushNotifier struct {
	service Service
	logger  logger.Logger
}

func NewPushNotifier(service Service, logger logger.Logger) *PushNotifier {
	return &PushNotifier{service: service, logger: logger}
}

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Write(ctx, websocket.MessageText, msg)
}

func (h *PushNotifier) HandlePush(w http.ResponseWriter, r *http.Request) {
	wsConn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		h.logger.Errorf("websocket.Accept(): %v", err)
		return
	}
	defer wsConn.Close(websocket.StatusInternalError, "handshake failed")

	h.logger.Info("new websocket connection\n")
	newSub, err := pubsub.NewSubscriber()

	if err != nil {
		h.logger.Errorf("pubsub.NewSubscriber(): %v", err)
		return
	}

	h.service.Broker().Subscribe(newSub, subscriptionTopic)

	messages := newSub.GetMessages()
	for msg := range messages {

		item, ok := msg.GetPayload().(*models.Snapshot)

		if !ok {
			h.logger.Errorf("msg.GetPayload().(*models.Snapshot): %s", "invalid type assertion")
			return
		}
		h.logger.Infof("Sending snapshot: %s", item.Stringify())

		err := writeTimeout(r.Context(), 5*time.Second, wsConn, []byte(item.Stringify()))

		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			h.logger.Infof("writeTimeout(): %v", err)
			return
		}

		if err != nil {
			h.logger.Errorf("writeTimeout(): %v", err)
			return
		}
	}
}
