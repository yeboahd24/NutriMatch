package handler

import (
	"github.com/rs/zerolog"
)

// BaseHandler contains common dependencies for all handlers
type BaseHandler struct {
	logger zerolog.Logger
}

// NewBaseHandler creates a new base handler
func NewBaseHandler(logger zerolog.Logger) BaseHandler {
	return BaseHandler{
		logger: logger,
	}
}
