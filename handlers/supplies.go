package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/GDG-on-Campus-KHU/SC4_BE/config"
	"github.com/GDG-on-Campus-KHU/SC4_BE/models"
	"github.com/GDG-on-Campus-KHU/SC4_BE/services"
	"github.com/golang-jwt/jwt"
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
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		h.sendErrorResponse(w, http.StatusUnauthorized, "유효하지 않은 토큰입니다.")
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return h.config.JWTSecret, nil
	})

	if err != nil || !token.Valid {
		h.sendErrorResponse(w, http.StatusUnauthorized, "유효하지 않은 토큰입니다.")
		return
	}

	userID := int(claims["user_id"].(float64))
	username := claims["username"].(string)

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
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		h.sendErrorResponse(w, http.StatusUnauthorized, "유효하지 않은 토큰입니다.")
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return h.config.JWTSecret, nil
	})

	if err != nil || !token.Valid {
		h.sendErrorResponse(w, http.StatusUnauthorized, "유효하지 않은 토큰입니다.")
		return
	}

	userID := int(claims["user_id"].(float64))

	var req SuppliesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "잘못된 요청 형식입니다.")
		return
	}

	if err = h.suppliesService.SaveUserSupplies(userID, req.Supplies); err != nil {
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
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		h.sendErrorResponse(w, http.StatusUnauthorized, "유효하지 않은 토큰입니다.")
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return h.config.JWTSecret, nil
	})

	if err != nil || !token.Valid {
		h.sendErrorResponse(w, http.StatusUnauthorized, "유효하지 않은 토큰입니다.")
		return
	}

	userID := int(claims["user_id"].(float64))

	var req SuppliesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "잘못된 요청 형식입니다.")
		return
	}

	if err = h.suppliesService.UpdateUserSupplies(userID, req.Supplies); err != nil {
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
