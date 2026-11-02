package pkg

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/shivamkkas/rdpms25-go-rpc-service/pkg/config"
	"github.com/shivamkkas/rdpms25-go-rpc-service/pkg/util"
)

func Start(version, buildTime string) {
	util.DurationLogger().Start("app_setup")
	conf, logger, _, err := config.Initialise()
	logger.Info("starting application", "version", version, "build", buildTime)
	if err != nil {
		logger.Error("unable to start application", "err", err)
		os.Exit(1)
	}

	logger.Info("starting application", "version", version, "build", buildTime)

	// setup rest apis
	engine := gin.New()

	engine.Use(func(c *gin.Context) {
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), midi=(), sync-xhr=(), microphone=(), camera=(), magnetometer=(), gyroscope=(), fullscreen=(self), payment=()")
		c.Next()
	})

	util.DurationLogger().End("app_setup")

	slog.Info("starting server", "port", conf.App.Port)
	if err := engine.Run(fmt.Sprintf("%s:%d", "0.0.0.0", conf.App.Port)); err != nil {
		slog.Error("unable to start server", "err", err)
	}
}
