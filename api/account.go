package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/teguhatma/simple-bank/db/sqlc"
	"github.com/teguhatma/simple-bank/request"
)

func (server *Server) createAccount(ctx *gin.Context) {
	var req request.AccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	id, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, map[string]int64{
		"data": id,
	})
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req request.GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (server *Server) listAccount(ctx *gin.Context) {
	var req request.ListAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	account, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, account)
}
