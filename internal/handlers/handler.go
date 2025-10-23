package handlers

import (
	v1 "go.xoder.cn/shortener-server/internal/handlers/v1"
)

// Handler is the handler struct
type Handler struct {
	AccountHandler *v1.AccountHandler
	UserHandler    *v1.UserHandler
	ShortenHandler *v1.ShortenHandler
	HistoryHandler *v1.HistoryHandler
}

// Handle expose the handler to outside
var Handle *Handler

func init() {
	Handle = &Handler{
		AccountHandler: v1.NewAccountHandler(),
		UserHandler:    v1.NewUserHandler(),
		ShortenHandler: v1.NewShortenHandler(),
		HistoryHandler: v1.NewHistoryHandler(),
	}
}
