package main

import (
	"crypto/sha256"
	"net/http"
	"time"
	"util"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func LoggerMiddleware(c *gin.Context) {
	t := time.Now()
	c.Next()
	// after request
	latency := time.Since(t)

	if l := len(c.Errors); l != 0 {
		err := c.Errors[l-1].Err
		code := 0
		level := zerolog.WarnLevel
		switch derr := err.(type) {
		case *util.DostupError:
			switch derr.Kind() {
			case util.InvalidParameter:
				level = zerolog.ErrorLevel
				code = http.StatusBadRequest
			case util.Conflict:
				level = zerolog.ErrorLevel
				code = http.StatusConflict
			case util.NotFound:
				level = zerolog.ErrorLevel
				code = http.StatusNotFound
			case util.Unavailable:
				level = zerolog.FatalLevel
				code = http.StatusServiceUnavailable
			default:
				level = zerolog.FatalLevel
				code = http.StatusInternalServerError
			}

			if derr.Message() != "" {
				c.String(code, derr.Message())
			} else {
				c.Status(code)
			}
		}

		log.WithLevel(level).
			Timestamp().
			Dur("dur", latency).
			Str("method", c.Request.Method).
			Str("uri", c.Request.RequestURI).
			Err(err).
			Int("status", c.Writer.Status()).
			Send()
	}
}

func AuthMiddleware(ph [32]byte) func(c *gin.Context) {
	return func(c *gin.Context) {
		pwd := c.Request.Header.Get("Authorization")
		if pwd != "" {
			pwdHash := sha256.Sum256([]byte(pwd))
			i := 0
			for i < 32 {
				if ph[i] != pwdHash[i] {
					break
				}
				i++
			}
			if i == 32 {
				c.Next()
			} else {
				c.AbortWithStatus(http.StatusUnauthorized)
			}
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
