package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func GetLoggerFromContext(c *gin.Context) (*zerolog.Logger, bool) {
	logger, exitst := c.Get("logger")
	if !exitst {
		return nil, false
	}
	loggerFromContext, ok := logger.(*zerolog.Logger)
	if !ok {
		return nil, false
	}
	return loggerFromContext, true
}

func LoggerMiddleware(log *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID, ok := c.Get("request_id")
		if !ok {
			log.Error().Msg("Request ID not found in context")
			c.AbortWithStatus(500)
			return
		}
		reqLog := log.With().
			Str("request_id", requestID.(string)).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Logger()

		c.Set("logger", &reqLog)

		defer func() {
			if rec := recover(); rec != nil {
				reqLog.Error().
					Str("request_id", requestID.(string)).
					Str("method", c.Request.Method).
					Str("path", c.Request.URL.Path).
					Interface("panic", rec).
					Msg("Panic recovered")
				c.AbortWithStatus(500)
			}
		}()

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		var logEvent *zerolog.Event
		if statusCode >= 500 {
			logEvent = reqLog.Error()
		} else if statusCode >= 400 {
			logEvent = reqLog.Warn()
		} else {
			logEvent = reqLog.Info()
		}

		logEvent.
			Int("status_code", statusCode).
			Dur("latency", latency).
			Msg("Request completed")
	}
}
