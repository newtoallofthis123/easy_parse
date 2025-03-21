package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/newtoallofthis123/easy_parse/api"
	"github.com/newtoallofthis123/easy_parse/db"
	"github.com/newtoallofthis123/easy_parse/parser"
	"github.com/newtoallofthis123/easy_parse/utils"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	logger.Info("Initialized Logger to STDERR")
	env := utils.ReadEnv()
	logger.Info("Initialized Environment Variables")

	store, err := db.NewStore(env.ConnString)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to initialize database: %s", err.Error()))
		return
	}
	logger.Info("Initialized Database")

	gemini, err := parser.NewGeminiAPI(env.GeminiAPIKey, parser.SystemPrompt())
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to initialize Gemini API: %s", err.Error()))
		return
	}
	logger.Info("Initialized Gemini API")

	api := api.NewApiServer(env.Port, logger, store, gemini)
	err = api.Start()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to start API server: %s", err.Error()))
		return
	}
}
