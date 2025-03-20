package api

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/newtoallofthis123/easy_parse/db"
)

type ApiServer struct {
	listenAddr string
	logger     *slog.Logger
	store      *db.Store
}

func NewApiServer(port string, logger *slog.Logger, store *db.Store) *ApiServer {
	logger.Debug("Initialized API Server")
	return &ApiServer{
		listenAddr: fmt.Sprintf(":%s", port),
		logger:     logger,
		store:      store,
	}
}

func (api *ApiServer) handleUserCreate(c *gin.Context) {
	var user db.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		api.logger.Error(fmt.Sprintln("Error Creating User", err))
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err = api.store.CreateUser(user)
	if err != nil {
		api.logger.Error(fmt.Sprintln("Error Creating User", err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, user)
}

func (api *ApiServer) handleUserGet(c *gin.Context) {
	id := c.Param("id")
	user, err := api.store.GetUser(id)
	if err != nil {
		api.logger.Error(fmt.Sprintln("Error Getting User", err))
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, user)
}

func (api *ApiServer) handleUserDelete(c *gin.Context) {
	id := c.Param("id")
	err := api.store.DeleteUser(id)
	if err != nil {
		api.logger.Error(fmt.Sprintln("Error Deleting User", err))
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User deleted"})
}

func (api *ApiServer) handleTokenCreate(c *gin.Context) {
	var tokenReq db.CreateTokenRequest
	err := c.ShouldBindJSON(&tokenReq)
	if err != nil {
		api.logger.Error(fmt.Sprintln("Error Creating Token", err))
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := api.store.CreateToken(tokenReq)
	if err != nil {
		api.logger.Error(fmt.Sprintln("Error Creating Token", err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, token)
}

func (api *ApiServer) handleTokenGet(c *gin.Context) {
	id := c.Param("id")
	token, err := api.store.GetToken(id)
	if err != nil {
		api.logger.Error(fmt.Sprintln("Error Getting Token", err))
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, token)
}

func (api *ApiServer) handleTokenDelete(c *gin.Context) {
	id := c.Param("id")
	err := api.store.DeleteToken(id)
	if err != nil {
		api.logger.Error(fmt.Sprintln("Error Deleting Token", err))
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Token deleted"})
}

func (api *ApiServer) handleParse(c *gin.Context) {
	var reqReq db.CreateRequestRequest
	err := c.ShouldBindJSON(&reqReq)
	if err != nil {
		api.logger.Error(fmt.Sprintln("Error Creating Request", err))
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//TODO: Implement the actual request create
	_, err = api.store.CreateRequest(reqReq, "success")
	if err != nil {
		api.logger.Error(fmt.Sprintln("Error Creating Request", err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"todo": "I have to do this later!"})
}

func (api *ApiServer) Start() error {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	api.logger.Info("Starting API Server", "addr", api.listenAddr)

	// User routes
	user := r.Group("/users")
	user.POST("/create", api.handleUserCreate)
	user.GET("/:id", api.handleUserGet)
	user.DELETE("/:id", api.handleUserDelete)

	// Token routes
	token := r.Group("/tokens")
	token.POST("/create", api.handleTokenCreate)
	token.GET("/:id", api.handleTokenGet)
	token.DELETE("/:id", api.handleTokenDelete)

	// The main parse route
	// FIXME: This doesn't work
	r.POST("/parse", api.handleParse)

	err := r.Run(api.listenAddr)
	return err
}
