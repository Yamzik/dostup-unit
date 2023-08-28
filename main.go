package main

import (
	"core"
	"fmt"
	"interop"
	"os"

	docs "dostup/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	requiredOptions, errors := new(RequiredOptionsProvider).Load()
	if len(errors) != 0 {
		for _, e := range errors {
			log.Fatal().
				Err(e).
				Msg("could not load env")
		}
		os.Exit(1)
	}
	options := new(OptionsProvider).Load()
	mementoProvider := new(MementoProvider).New(
		options.Port,
		options.MementoPath,
		options.ConfigPath)

	mem, err := mementoProvider.Load()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("could not load memento")
		os.Exit(1)
	}

	peerRepo := new(core.PeerRepository).New().SetMemento(mem)
	userProvider := new(core.UserProvider).New(peerRepo)
	peerService := new(core.PeerService).
		New(peerRepo, userProvider, requiredOptions.Host, options.Port).
		SetMemento(mem)

	config, err := peerService.MakeConfig()
	if err != nil {
		os.Exit(1)
	}
	mementoProvider.SaveConfig(config)

	interop.Down()
	interop.Up()

	cancel := make(chan struct{})
	go SyncTransfer(peerService, mementoProvider, cancel)

	r := gin.Default()
	r.Use(AuthMiddleware(requiredOptions.PwdHash))
	r.Use(LoggerMiddleware)
	(&UserController{}).New(peerService, mementoProvider).Use(r)
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Host = "example.com"
	docs.SwaggerInfo.Description = "Unit"
	docs.SwaggerInfo.Title = "Unit"
	docs.SwaggerInfo.Version = "0.0.1"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(fmt.Sprintf(":%v", options.UnitPort)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
