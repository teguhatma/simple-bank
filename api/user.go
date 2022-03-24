package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/teguhatma/simple-bank/db/sqlc"
	"github.com/teguhatma/simple-bank/request"
	"github.com/teguhatma/simple-bank/response"
	"github.com/teguhatma/simple-bank/util"
)

func (server *Server) createUser(ctx *gin.Context) {
	var req request.UserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		FullName:       req.FullName,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	resp := userResponse(user)

	ctx.JSON(http.StatusCreated, resp)
}

func userResponse(user db.User) *response.UserResponse {
	return &response.UserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}
