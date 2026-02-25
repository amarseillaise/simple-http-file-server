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

type ApiResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func (h *VideoHandler) CreateReel(w http.ResponseWriter, r *http.Request) {
	shortcode, ok := extractShortcode(w, r, "CreateReel")
	if !ok {
		return
	}

	err := h.service.CreateReel(shortcode)
	var response ApiResponse
	var statusCode int
	switch {
	case err == nil:
		statusCode = http.StatusCreated
		response.Success = true
		response.Message = "video created successfully"
	case errors.Is(err, storage.ErrDirectoryExists):
		statusCode = http.StatusAccepted
		response.Success = true
		response.Message = "video with this shortcode already exists"
	default:
		statusCode = http.StatusBadRequest
		response.Success = false
		response.Message = "reel doesn't exist or yt-dlp error"
	}

	writeJSON(w, statusCode, response)
}

func (h *VideoHandler) GetReelVideo(w http.ResponseWriter, r *http.Request) {
	shortcode, ok := extractShortcode(w, r, "GetReelVideo")
	if !ok {
		return
	}
	videoPath, err := h.service.GetVideoPath(shortcode)
	var response ApiResponse
	var statusCode int
	switch {
	case err == nil:
		w.Header().Set("Content-Type", "video/mp4")
		w.Header().Set("Content-Disposition", "inline; filename=\"video.mp4\"")
		http.ServeFile(w, r, videoPath)
		return
	case errors.Is(err, storage.ErrInvalidShortcode):
		statusCode = http.StatusBadRequest
		response.Success = false
		response.Message = "invalid shortcode format"
	case errors.Is(err, storage.ErrDirectoryNotFound):
		statusCode = http.StatusNotFound
		response.Success = false
		response.Message = "directory not found"
	case errors.Is(err, storage.ErrVideoNotFound):
		statusCode = http.StatusNotFound
		response.Success = false
		response.Message = "video not found"
	default:
		statusCode = http.StatusInternalServerError
		response.Success = false
		response.Message = "Unexpected error occurred while retrieving video"
	}
	writeJSON(w, statusCode, response)
}

func (h *VideoHandler) GetReelDescription(w http.ResponseWriter, r *http.Request) {
	shortcode, ok := extractShortcode(w, r, "GetReelDescription")
	if !ok {
		return
	}
	var description string
	descriptionPath, err := h.service.GetDescriptionPath(shortcode)
	if err != nil {
		description = ""
	} else {
		description = h.service.GetReelDescription(descriptionPath)
	}

	writeJSON(w, http.StatusOK, ApiResponse{
		Success: true,
		Message: description,
	})
}

func (h *VideoHandler) DeleteReel(w http.ResponseWriter, r *http.Request) {
	shortcode, ok := extractShortcode(w, r, "DeleteReel")
	if !ok {
		return
	}
	err := h.service.DeleteReel(shortcode)
	var response ApiResponse
	var statusCode int
	switch {
	case err == nil:
		statusCode = http.StatusOK
		response.Success = true
		response.Message = "video deleted successfully"
	case errors.Is(err, storage.ErrInvalidShortcode):
		statusCode = http.StatusBadRequest
		response.Success = false
		response.Message = "invalid shortcode format"
	case errors.Is(err, storage.ErrDirectoryNotFound):
		statusCode = http.StatusNotFound
		response.Success = false
		response.Message = "directory not found"
	default:
		statusCode = http.StatusInternalServerError
		response.Success = false
		response.Message = "Unexpected error occurred while retrieving video"
	}
	writeJSON(w, statusCode, response)
}

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

func extractShortcode(w http.ResponseWriter, r *http.Request, operation string) (string, bool) {
	shortcode := mux.Vars(r)["shortcode"]
	if shortcode == "" {
		writeJSON(w, http.StatusBadRequest, ApiResponse{
			Success: false,
			Message: "shortcode is required",
		})
		return "", false
	}
	log.Printf("%s request for shortcode: %s", operation, shortcode)
	return shortcode, true
}

func (h *VideoHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/reel/{shortcode}", h.CreateReel).Methods(http.MethodPost)
	router.HandleFunc("/api/reel/{shortcode}", h.DeleteReel).Methods(http.MethodDelete)
	router.HandleFunc("/api/reel/{shortcode}/video.mp4", h.GetReelVideo).Methods(http.MethodGet)
	router.HandleFunc("/api/reel/{shortcode}/description", h.GetReelDescription).Methods(http.MethodGet)
}
