package http

import (
	"encoding/json"
	"errors"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"mesaggio-test/internal/messages"
	"mesaggio-test/internal/models"
	pkgHttp "mesaggio-test/pkg/http"
	"net/http"
	"strings"
)

type messageHandler struct {
	log *logrus.Logger
	ms  messages.Service
}

// NewMessageHandler Messages handlers constructor
func NewMessageHandler(log *logrus.Logger, ms messages.Service) messages.Handler {
	return &messageHandler{log: log, ms: ms}
}

// ReceiveMessage
// @Summary Receive new message
// @Tags Comments
// @Accept  json
// @Produce  json
// @Param input body models.Request true "Message content"
// @Success 201 {object} models.CreatedResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /msg [post]
func (m *messageHandler) ReceiveMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req := &models.Request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			m.log.Errorf("Error: %s", err)
			pkgHttp.RespondError(w, http.StatusBadRequest, err)
			return
		}

		if strings.TrimSpace(req.Content) == "" {
			m.log.Error("Content is empty")
			pkgHttp.RespondError(w, http.StatusBadRequest, errors.New("content cannot be empty"))
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

// GetStatistics
// @Summary Get messages statistics
// @Tags Comments
// @Produce  json
// @Success 200 {object} models.Statistics
// @Failure 400 {object} models.ErrorResponse
// @Router /stats [get]
func (m *messageHandler) GetStatistics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		stats, err := m.ms.GetMessagesStatistics(ctx)
		if err != nil {
			m.log.Errorf("Error: %s", err)
			pkgHttp.RespondError(w, http.StatusBadRequest, err)
			return
		}

		pkgHttp.Respond(w, http.StatusOK, stats)
	}
}
