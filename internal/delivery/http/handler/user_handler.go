package handler

import (
	"evently/internal/domain/usecase/user"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	usecase user.UserUsecase
}

func NewUserHandler(u user.UserUsecase) *userHandler {
	return &userHandler{usecase: u}
}

func (h *userHandler) RegisterHandler(c *gin.Context) {

}

func (h *userHandler) LoginHandler(c *gin.Context) {

}

func (h *userHandler) GetProfileHandler(c *gin.Context) {

}

func (h *userHandler) GetBookingHistoryHandler(c *gin.Context) {

}
