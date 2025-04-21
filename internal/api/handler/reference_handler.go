package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/yeboahd24/nutrimatch/internal/service"
	"github.com/yeboahd24/nutrimatch/pkg/response"
)

type ReferenceHandler struct {
	BaseHandler
	referenceService service.ReferenceService
}

func NewReferenceHandler(referenceService service.ReferenceService, logger zerolog.Logger) *ReferenceHandler {
	return &ReferenceHandler{
		BaseHandler:      NewBaseHandler(logger),
		referenceService: referenceService,
	}
}

func (h *ReferenceHandler) RegisterRoutes(r chi.Router) {
	r.Get("/allergens", h.GetAllergens)
	r.Get("/health-conditions", h.GetHealthConditions)
	r.Get("/dietary-patterns", h.GetDietaryPatterns)
}

// @Summary Get allergens
// @Description Get a list of all allergens
// @Tags reference
// @Accept json
// @Produce json
// @Success 200 {object} docs.Response{data=[]docs.ReferenceItem}
// @Failure 500 {object} docs.ErrorResponse
// @Router /api/v1/reference/allergens [get]
func (h *ReferenceHandler) GetAllergens(w http.ResponseWriter, r *http.Request) {
	allergens, err := h.referenceService.GetAllergens(r.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get allergens")
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, allergens)
}

// @Summary Get health conditions
// @Description Get a list of all health conditions
// @Tags reference
// @Accept json
// @Produce json
// @Success 200 {object} docs.Response{data=[]docs.ReferenceItem}
// @Failure 500 {object} docs.ErrorResponse
// @Router /api/v1/reference/health-conditions [get]
func (h *ReferenceHandler) GetHealthConditions(w http.ResponseWriter, r *http.Request) {
	conditions, err := h.referenceService.GetHealthConditions(r.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get health conditions")
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, conditions)
}

// @Summary Get dietary patterns
// @Description Get a list of all dietary patterns
// @Tags reference
// @Accept json
// @Produce json
// @Success 200 {object} docs.Response{data=[]docs.ReferenceItem}
// @Failure 500 {object} docs.ErrorResponse
// @Router /api/v1/reference/dietary-patterns [get]
func (h *ReferenceHandler) GetDietaryPatterns(w http.ResponseWriter, r *http.Request) {
	patterns, err := h.referenceService.GetDietaryPatterns(r.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get dietary patterns")
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, patterns)
}
