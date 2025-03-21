package api

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/newtoallofthis123/easy_parse/db"
	"github.com/newtoallofthis123/easy_parse/parser"
	"google.golang.org/genai"
)

type ApiServer struct {
	listenAddr string
	logger     *slog.Logger
	store      *db.Store
	gemini     *parser.GeminiAPI
}

func NewApiServer(port string, logger *slog.Logger, store *db.Store, gemini *parser.GeminiAPI) *ApiServer {
	logger.Debug("Initialized API Server")
	return &ApiServer{
		listenAddr: fmt.Sprintf(":%s", port),
		logger:     logger,
		store:      store,
		gemini:     gemini,
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
	authToken := c.Request.Header.Get("Authorization")
	if authToken == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// strip Bearer
	authToken = strings.TrimPrefix(authToken, "Bearer ")
	token, err := api.store.GetToken(authToken)
	if err != nil {
		api.logger.Error(fmt.Sprintln("Error Getting Token", err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// get form data
	file, err := c.FormFile("file")
	if err != nil {
		api.logger.Error(fmt.Sprintln("Error Getting Form File", err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	data, err := file.Open()
	if err != nil {
		api.logger.Error(fmt.Sprintln("Error Opening File", err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer data.Close()
	contentType := c.PostForm("Content-Type")
	if contentType == "" {
		contentType = file.Header.Get("Content-Type")
	}
	fileData, err := io.ReadAll(data)
	if err != nil {
		log.Fatal(err)
	}

	parts := []*genai.Part{
		{Text: "Parse this PDF"},
		{InlineData: &genai.Blob{Data: fileData, MIMEType: contentType}},
	}
	res, err := api.gemini.Send(parts)
	if err != nil {
		api.logger.Error(fmt.Sprintln("Error Sending Request", err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	_, err = api.store.CreateRequest(db.CreateRequestRequest{UserId: token.UserId}, "success")
	if err != nil {
		api.logger.Error(fmt.Sprintln("Error Creating Request", err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": res})
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
