package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/teguhatma/simple-bank/db/sqlc"
	"github.com/teguhatma/simple-bank/request"
	"github.com/teguhatma/simple-bank/token"
)

func (server *Server) createTransfer(ctx *gin.Context) {
	var req request.TransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	if req.FromAccountID == req.ToAccountID {
		err := fmt.Errorf("you cannot transfer to your own account")
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Username == fromAccount.Owner {
		err := errors.New("from account does not belong to the authenticated user")
		errorResponse(ctx, http.StatusUnauthorized, err)
	}

	_, valid = server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			errorResponse(ctx, http.StatusNotFound, err)
			return account, false
		}
		errorResponse(ctx, http.StatusInternalServerError, err)
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		errorResponse(ctx, http.StatusBadRequest, err)
		return account, false
	}
	return account, true
}
