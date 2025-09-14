package handler

import (
	"net/http"
	"strconv"

	"evently/internal/domain/booking"
	"evently/internal/domain/events"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	eventUsecase   events.EventUsecase
	bookingUsecase booking.BookingUsecase
}

func NewAdminHandler(eventUsecase events.EventUsecase, bookingUsecase booking.BookingUsecase) *AdminHandler {
	return &AdminHandler{
		eventUsecase:   eventUsecase,
		bookingUsecase: bookingUsecase,
	}
}

func (h *AdminHandler) GetEventAnalytics(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	analytics, err := h.eventUsecase.GetMostPopularEvents(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"analytics": analytics})
}

func (h *AdminHandler) GetBookingAnalytics(c *gin.Context) {
	eventID := c.Param("eventId")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event ID is required"})
		return
	}

	analytics, err := h.bookingUsecase.GetBookingAnalytics(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"analytics": analytics})
}

func (h *AdminHandler) GetAllEvents(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	events, err := h.eventUsecase.ListAllEvents(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

func (h *AdminHandler) GetEventBookings(c *gin.Context) {
	eventID := c.Param("eventId")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event ID is required"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	bookings, err := h.bookingUsecase.GetEventBookings(c.Request.Context(), eventID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"bookings": bookings})
}
