package http

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"mesaggio-test/internal/messages"
	"mesaggio-test/internal/models"
	pkgHttp "mesaggio-test/pkg/http"
	"net/http"
)

type messageHandler struct {
	log *logrus.Logger
	ms  messages.Service
}

func NewMessageHandler(log *logrus.Logger, ms messages.Service) *messageHandler {
	return &messageHandler{log: log, ms: ms}
}

func (m *messageHandler) ReceiveMessage() http.HandlerFunc {
	type request struct {
		Content string `json:"content"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			m.log.Errorf("Error: %s", err)
			pkgHttp.RespondError(w, http.StatusBadRequest, err)
			return
		}

		message := &models.Message{
			MessageID: uuid.NewV4(),
			Content:   req.Content,
		}

		if err := m.ms.ReceiveMessage(ctx, message); err != nil {
			m.log.Errorf("Error: %s", err)
			pkgHttp.RespondError(w, http.StatusBadRequest, err)
			return
		}

		pkgHttp.Respond(w, http.StatusCreated, models.CreatedResponse{MessageID: message.MessageID})
	}
}

/*func (m *messageHandler) UpdateMessage() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {

	}
}*/
