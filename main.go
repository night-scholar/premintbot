package main

import (
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/mux"
	"github.com/mager/premintbot/bot"
	"github.com/mager/premintbot/config"
	"github.com/mager/premintbot/database"
	"github.com/mager/premintbot/handler"
	"github.com/mager/premintbot/logger"
	"github.com/mager/premintbot/premint"
	"github.com/mager/premintbot/router"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Provide(
			config.Options,
			database.Options,
			logger.Options,
			router.Options,
			premint.Options,
		),
		fx.Invoke(Register),
	).Run()
}

func Register(
	lc fx.Lifecycle,
	cfg config.Config,
	database *firestore.Client,
	logger *zap.SugaredLogger,
	router *mux.Router,
	premintClient *premint.PremintClient,
) {
	// Setup Discord Bot
	token := fmt.Sprintf("Bot %s", cfg.DiscordAuthToken)
	dg, err := discordgo.New(token)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	bot.Start(dg, logger, database, premintClient)
	// Start the handler for health checks
	handler.New(logger, router)
}
