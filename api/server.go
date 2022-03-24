package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/teguhatma/simple-bank/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	v1 := router.Group("/api/v1")
	{
		v1.POST("/accounts", server.createAccount)
		v1.GET("/accounts/:id", server.getAccount)
		v1.GET("/accounts", server.listAccount)

		v1.POST("/transfers", server.createTransfer)

		v1.POST("/users", server.createUser)
	}

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
