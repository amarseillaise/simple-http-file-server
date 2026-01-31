package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/amarseillaise/simple-http-file-server/internal/service"
	"github.com/amarseillaise/simple-http-file-server/internal/storage"
	"github.com/gorilla/mux"
)

type VideoHandler struct {
	service *service.VideoService
}

func NewVideoHandler(service *service.VideoService) *VideoHandler {
	return &VideoHandler{
		service: service,
	}
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

func writeError(w http.ResponseWriter, statusCode int, errType string, message string) {
	writeJSON(w, statusCode, ErrorResponse{
		Error:   errType,
		Message: message,
	})
}

func handleServiceError(w http.ResponseWriter, err error, operation string) bool {
	if err == nil {
		return false
	}

	switch {
	case errors.Is(err, storage.ErrInvalidShortcode):
		writeError(w, http.StatusBadRequest, "bad_request", "invalid shortcode format")
	case errors.Is(err, storage.ErrDirectoryExists):
		writeError(w, http.StatusConflict, "conflict", "video with this shortcode already exists")
	case errors.Is(err, storage.ErrDirectoryNotFound):
		writeError(w, http.StatusNotFound, "not_found", "video not found")
	case errors.Is(err, service.ErrReelDoesNotExistOrYtdlpBroken):
		writeError(w, http.StatusNotFound, "not_found", "video not found or yt-dlp error")
	default:
		log.Printf("Error %s: %v", operation, err)
		writeError(w, http.StatusInternalServerError, "internal_error", "internal server error")
	}
	return true
}

func (h *VideoHandler) CreateVideo(w http.ResponseWriter, r *http.Request) {
	shortcode := mux.Vars(r)["shortcode"]
	if shortcode == "" {
		writeError(w, http.StatusBadRequest, "bad_request", "shortcode is required")
		return
	}

	log.Printf("CreateVideo request for shortcode: %s", shortcode)

	if handleServiceError(w, h.service.CreateVideo(shortcode), "creating video") {
		return
	}

	writeJSON(w, http.StatusCreated, SuccessResponse{
		Success: true,
		Message: "video created successfully",
	})
}

func (h *VideoHandler) GetVideo(w http.ResponseWriter, r *http.Request) {
	shortcode := mux.Vars(r)["shortcode"]
	if shortcode == "" {
		writeError(w, http.StatusBadRequest, "bad_request", "shortcode is required")
		return
	}

	log.Printf("GetVideo request for shortcode: %s", shortcode)

	videoPath, err := h.service.GetVideoPath(shortcode)
	if handleServiceError(w, err, "getting video path") {
		return
	}

	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Content-Disposition", "inline; filename=\"video.mp4\"")
	http.ServeFile(w, r, videoPath)
}

func (h *VideoHandler) DeleteVideo(w http.ResponseWriter, r *http.Request) {
	shortcode := mux.Vars(r)["shortcode"]
	if shortcode == "" {
		writeError(w, http.StatusBadRequest, "bad_request", "shortcode is required")
		return
	}

	log.Printf("DeleteVideo request for shortcode: %s", shortcode)

	if handleServiceError(w, h.service.DeleteVideo(shortcode), "deleting video") {
		return
	}

	writeJSON(w, http.StatusOK, SuccessResponse{
		Success: true,
		Message: "video deleted successfully",
	})
}

func (h *VideoHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/video/{shortcode}", h.CreateVideo).Methods(http.MethodPost)
	router.HandleFunc("/api/video/{shortcode}", h.DeleteVideo).Methods(http.MethodDelete)
	router.HandleFunc("/api/video/{shortcode}", h.GetVideo).Methods(http.MethodGet)
}
