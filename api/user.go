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

func (server *Server) loginUser(ctx *gin.Context) {
	var req request.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		errorResponse(ctx, http.StatusUnauthorized, err)
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	resp := response.LoginUserResponse{
		AccessToken: accessToken,
		User:        *userResponse(user),
	}

	ctx.JSON(http.StatusOK, resp)
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
