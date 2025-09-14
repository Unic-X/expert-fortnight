package handler

import (
	"net/http"
	"strconv"
	"strings"

	"evently/internal/domain/booking"
	"evently/internal/domain/waitlist"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingUsecase  booking.BookingUsecase
	waitlistUsecase waitlist.WaitlistUsecase
}

func NewBookingHandler(bookingUsecase booking.BookingUsecase, waitlistUsecase waitlist.WaitlistUsecase) *BookingHandler {
	return &BookingHandler{
		bookingUsecase:  bookingUsecase,
		waitlistUsecase: waitlistUsecase,
	}
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var booking booking.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	booking.UserID = userID.(string)

	err := h.bookingUsecase.CreateBooking(c.Request.Context(), &booking)
	if err != nil {
		if strings.Contains(err.Error(), "insufficient seats available") {

			waitlistErr := h.waitlistUsecase.JoinWaitlist(c.Request.Context(), userID.(string), booking.EventID, booking.Quantity)
			if waitlistErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "event is full and failed to join waitlist: " + waitlistErr.Error()})
				return
			}

			position, posErr := h.waitlistUsecase.GetWaitlistPosition(c.Request.Context(), userID.(string), booking.EventID)
			if posErr != nil {
				position = 0
			}

			c.JSON(http.StatusAccepted, gin.H{
				"message":           "Event is full. You have been added to the waitlist.",
				"waitlist_position": position,
				"status":            "waitlisted",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "booking created successfully", "booking": booking, "status": "confirmed"})
}

func (h *BookingHandler) CancelBooking(c *gin.Context) {
	bookingID := c.Param("id")
	if bookingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "booking ID is required"})
		return
	}

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

// TODO: Shorten GetBooking and move everything inside usecase
func (h *BookingHandler) GetBooking(c *gin.Context) {
	bookingID := c.Param("id")
	if bookingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "booking ID is required"})
		return
	}

	booking, err := h.bookingUsecase.GetBooking(c.Request.Context(), bookingID)
	if err != nil {
		if h.waitlistUsecase != nil {
			wl, wlErr := h.waitlistUsecase.GetWaitlistByID(c.Request.Context(), bookingID)
			if wlErr == nil && wl != nil {
				// Auth check
				userID, exists := c.Get("user_id")
				if !exists {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
					return
				}
				userRole, _ := c.Get("user_role")
				if wl.UserID != userID.(string) && userRole != "admin" {
					c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
					return
				}

				position := 0
				if wl.Status == waitlist.WaitlistStatusActive {
					pos, perr := h.waitlistUsecase.GetWaitlistPosition(c.Request.Context(), wl.UserID, wl.EventID)
					if perr == nil {
						position = pos
					}
				}

				c.JSON(http.StatusOK, gin.H{
					"status":   "waitlisted",
					"waitlist": wl,
					"position": position,
				})
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}

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

	wl, _ := h.waitlistUsecase.GetUserWaitlist(c.Request.Context(), userID.(string), limit, offset)

	c.JSON(http.StatusOK, gin.H{"bookings": bookings, "waitlist": wl})
}
