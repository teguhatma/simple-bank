package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/teguhatma/simple-bank/db/sqlc"
	"github.com/teguhatma/simple-bank/token"
	"github.com/teguhatma/simple-bank/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.POST("/accounts", server.createAccount)
		v1.GET("/accounts/:id", server.getAccount)
		v1.GET("/accounts", server.listAccount)

		v1.POST("/transfers", server.createTransfer)

		v1.POST("/users", server.createUser)
		v1.POST("/users/login", server.loginUser)
	}
	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
