package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/GDG-on-Campus-KHU/SC4_BE/config"
	"github.com/GDG-on-Campus-KHU/SC4_BE/models"
	"github.com/GDG-on-Campus-KHU/SC4_BE/services"
)

type SuppliesHandler struct {
	suppliesService *services.SuppliesService
	config          *config.Config
}

func NewSuppliesHandler(ss *services.SuppliesService, cfg *config.Config) *SuppliesHandler {
	return &SuppliesHandler{
		suppliesService: ss,
		config:          cfg,
	}
}

func (h *SuppliesHandler) GetSupplies(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)

	if !ok {
		h.sendErrorResponse(w, http.StatusUnauthorized, "유효하지 않은 사용자입니다.")
		return
	}
	log.Println("userID:", userID)
	username, ok := r.Context().Value("username").(string)
	if !ok {
		h.sendErrorResponse(w, http.StatusUnauthorized, "유효하지 않은 사용자입니다.")
		return
	}
	log.Println("username:", username)
	supplies, err := h.suppliesService.GetUserSupplies(userID)
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, "물품 조회에 실패하였습니다.")
		return
	}

	userData := &models.UserData{
		Username: username,
		Supplies: supplies,
	}

	response := models.Response{
		Status:  http.StatusOK,
		Message: "물품 조회에 성공하였습니다.",
		Data:    userData,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *SuppliesHandler) sendErrorResponse(w http.ResponseWriter, status int, message string) {
	response := models.Response{
		Status:  status,
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

type SuppliesRequest struct {
	Supplies map[string]bool `json:"supplies"`
}

func (h *SuppliesHandler) SaveSupplies(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)

	if !ok {
		h.sendErrorResponse(w, http.StatusUnauthorized, "유효하지 않은 사용자입니다.")
		return
	}
	var req SuppliesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "잘못된 요청 형식입니다.")
		return
	}

	if err := h.suppliesService.SaveUserSupplies(userID, req.Supplies); err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, "물품 저장에 실패했습니다.")
		return
	}

	response := models.Response{
		Status:  http.StatusOK,
		Message: "물품 저장에 성공했습니다.",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *SuppliesHandler) UpdateSupplies(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)

	if !ok {
		h.sendErrorResponse(w, http.StatusUnauthorized, "유효하지 않은 사용자입니다.")
		return
	}

	var req SuppliesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "잘못된 요청 형식입니다.")
		return
	}

	if err := h.suppliesService.UpdateUserSupplies(userID, req.Supplies); err != nil {
		if err == services.ErrNoExistingSupplies {
			h.sendErrorResponse(w, http.StatusNotFound, "수정할 물품이 존재하지 않습니다.")
			return
		}
		h.sendErrorResponse(w, http.StatusInternalServerError, "물품 수정에 실패했습니다.")
		return
	}

	response := models.Response{
		Status:  http.StatusOK,
		Message: "물품 수정에 성공했습니다.",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
