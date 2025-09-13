package handler

import (
	"net/http"
	"strconv"

	"evently/internal/domain/model"
	"evently/internal/domain/usecase"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	eventUsecase usecase.EventUsecase
}

func NewEventHandler(eventUsecase usecase.EventUsecase) *EventHandler {
	return &EventHandler{
		eventUsecase: eventUsecase,
	}
}

func (h *EventHandler) ListUpcomingEvents(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	events, err := h.eventUsecase.ListUpcomingEvents(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

func (h *EventHandler) GetEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event ID is required"})
		return
	}

	event, err := h.eventUsecase.GetEvent(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": event})
}

func (h *EventHandler) CreateEvent(c *gin.Context) {
	var event model.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	event.CreatedBy = userID.(string)

	if err := h.eventUsecase.CreateEvent(c.Request.Context(), &event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "event created successfully", "event": event})
}

func (h *EventHandler) UpdateEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event ID is required"})
		return
	}

	var event model.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event.ID = eventID
	if err := h.eventUsecase.UpdateEvent(c.Request.Context(), &event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event updated successfully", "event": event})
}

func (h *EventHandler) DeleteEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event ID is required"})
		return
	}

	if err := h.eventUsecase.DeleteEvent(c.Request.Context(), eventID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event deleted successfully"})
}
