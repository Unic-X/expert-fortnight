package handler

import (
	"net/http"
	"strconv"

	"evently/internal/domain/model"
	"evently/internal/domain/usecase"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingUsecase usecase.BookingUsecase
}

func NewBookingHandler(bookingUsecase usecase.BookingUsecase) *BookingHandler {
	return &BookingHandler{
		bookingUsecase: bookingUsecase,
	}
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var booking model.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from JWT token context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	booking.UserID = userID.(string)

	if err := h.bookingUsecase.CreateBooking(c.Request.Context(), &booking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "booking created successfully", "booking": booking})
}

func (h *BookingHandler) CancelBooking(c *gin.Context) {
	bookingID := c.Param("id")
	if bookingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "booking ID is required"})
		return
	}

	// Get user ID from JWT token context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.bookingUsecase.CancelBooking(c.Request.Context(), bookingID, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "booking cancelled successfully"})
}

func (h *BookingHandler) GetBooking(c *gin.Context) {
	bookingID := c.Param("id")
	if bookingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "booking ID is required"})
		return
	}

	booking, err := h.bookingUsecase.GetBooking(c.Request.Context(), bookingID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}

	// Check if user owns the booking or is admin
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userRole, _ := c.Get("user_role")
	if booking.UserID != userID.(string) && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"booking": booking})
}

func (h *BookingHandler) GetUserBookings(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Get user ID from JWT token context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	bookings, err := h.bookingUsecase.GetUserBookings(c.Request.Context(), userID.(string), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"bookings": bookings})
}
