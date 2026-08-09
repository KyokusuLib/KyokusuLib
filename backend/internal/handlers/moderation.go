package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/lanxre/kyokusulib/internal/models/dto"
	service "github.com/lanxre/kyokusulib/internal/services"
	"github.com/lanxre/kyokusulib/internal/utils/response"
)

type ModerationHandler struct {
	service   *service.ModerationService
	Validator *validator.Validate
}

func NewModerationHandler(moderationService *service.ModerationService, validator *validator.Validate) *ModerationHandler {
	return &ModerationHandler{
		service:   moderationService,
		Validator: validator,
	}
}

func (h *ModerationHandler) GetPending(w http.ResponseWriter, r *http.Request) {
	content, err := h.service.GetPendingContent(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, content)
}

func (h *ModerationHandler) ApproveVolume(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.service.ApproveVolume(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessOkEmpty(w)
}

func (h *ModerationHandler) RejectVolume(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.service.RejectVolume(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessOkEmpty(w)
}

func (h *ModerationHandler) ApproveChapter(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.service.ApproveChapter(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessOkEmpty(w)
}

func (h *ModerationHandler) RejectChapter(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.service.RejectChapter(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessOkEmpty(w)
}

func (h *ModerationHandler) UpdateVolume(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var req dto.UpdateVolumeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := h.Validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}
	if err := h.service.UpdateVolume(r.Context(), id, req); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessOkEmpty(w)
}

func (h *ModerationHandler) UpdateChapter(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var req dto.UpdateChapterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := h.Validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}
	if err := h.service.UpdateChapter(r.Context(), id, req); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessOkEmpty(w)
}
