package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Adejare77/go/taskManager/internals/handlers"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

type SessionConfig struct {
	RedisSize    int
	RedisAddress string
	RDPassword   string
	SecretKey    string
	MaxAge       int
}

func loadSessionConfig() (*SessionConfig, error) {
	size, err := strconv.Atoi(os.Getenv("REDIS_SIZE"))
	if err != nil {
		handlers.Warning("invalid REDIS_SIZE. Defaults to size 10")
		size = 10
	}

	maxAge, err := strconv.Atoi(os.Getenv("SESSION_MAX_AGE"))
	if err != nil {
		handlers.Warning("invalid SESSION_MAX_AGE. Defaults to 600s")
		maxAge = 600
	}

	// Create and return RDConfig
	return &SessionConfig{
		MaxAge:       maxAge,
		RedisSize:    size,
		RedisAddress: os.Getenv("REDIS_ADDRESS"),
		RDPassword:   os.Getenv("REDIS_PASSWORD"),
		SecretKey:    os.Getenv("SECRET_KEY"),
	}, nil
}

var SessionStore redis.Store

func InitSession() error {
	cfg, err := loadSessionConfig()
	if err != nil {
		return fmt.Errorf("session configuration %v", err)
	}

	store, err := redis.NewStore(
		cfg.RedisSize, "tcp", cfg.RedisAddress, cfg.RDPassword, []byte(cfg.SecretKey))
	if err != nil {
		return fmt.Errorf("session initialization %v", err)
	}

	store.Options(sessions.Options{
		MaxAge:   cfg.MaxAge,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	})

	SessionStore = store
	return nil
}

// Create a new session for the user
func CreateSession(ctx *gin.Context, userID string) error {
	session := sessions.Default(ctx)

	session.Set("currentUser", userID)
	if err := session.Save(); err != nil {
		return fmt.Errorf("failed to save session: %v", err)
	}
	return nil
}

// Delete the user's session
func DeleteSession(ctx *gin.Context) {
	session := sessions.Default(ctx)

	session.Clear()
	// Invalidate the cookie
	session.Options(sessions.Options{
		MaxAge: -1,
		Path:   "/",
	})

	if err := session.Save(); err != nil {
		handlers.InternalServerError(ctx, "error deleting session", "Failed to Delete Session")
	}
}
