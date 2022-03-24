package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func errorResponse(ctx *gin.Context, httpStatus int, err error) {
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code.Name() {
		case "unique_violation", "foreign_key_violation":
			ctx.JSON(http.StatusForbidden, map[string]string{
				"error": err.Error(),
			})
			return
		}
	}

	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	} else {
		ctx.JSON(httpStatus, map[string]string{
			"error": err.Error(),
		})
	}
}
