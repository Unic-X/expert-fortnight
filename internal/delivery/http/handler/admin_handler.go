package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"evently/internal/domain/model"
	"evently/internal/domain/usecase/admin"
)

type adminHandler struct {
	usecase admin.AdminUsecase
}

func NewAdminHandler(u admin.AdminUsecase) *adminHandler {
	return &adminHandler{usecase: u}
}

func (h *adminHandler) AddEventHandler(c *gin.Context) {
	name := c.Query("name")
	venue := c.Query("venue")
	timeStr := c.Query("time")
	capacity := c.Query("capacity")

	if name == "" || venue == "" || timeStr == "" || capacity == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	eventTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format. Please use RFC3339 format (e.g., 2006-01-02T15:04:05Z07:00)"})
		return
	}

	capacityInt, err := strconv.Atoi(capacity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Capacity must be a number"})
		return
	}

	newEvent := &model.Event{
		Name:     name,
		Venue:    venue,
		Time:     eventTime,
		Capacity: capacityInt,
	}

	json.NewDecoder(c.Request.Body).Decode(newEvent)

	err = h.usecase.CreateEvent(newEvent)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event created successfully"})
}

func (h *adminHandler) UpdateEventHandler(c *gin.Context) {

}

func (h *adminHandler) DeleteEventHandler(c *gin.Context) {

}

func (h *adminHandler) GetEventAnalytics(c *gin.Context) {

}

func (h *adminHandler) GetOverallAnalytics(c *gin.Context) {

}
