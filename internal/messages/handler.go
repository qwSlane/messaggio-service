package messages

import "net/http"

type Handler interface {
	ReceiveMessage() http.HandlerFunc
	GetStatistics() http.HandlerFunc
}
