package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

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
		return nil, fmt.Errorf("invalid REDIS_SIZE")
	}
	maxAge, err := strconv.Atoi(os.Getenv("SESSION_MAX_AGE"))
	if err != nil {
		return nil, fmt.Errorf("invalid SESSION_MAX_AGE")
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
		return fmt.Errorf("(session configuration) %v", err)
	}
	store, err := redis.NewStore(
		cfg.RedisSize, "tcp", cfg.RedisAddress, cfg.RDPassword, []byte(cfg.SecretKey))
	if err != nil {
		return fmt.Errorf("(session initialization) %v", err)
	}

	store.Options(sessions.Options{
		MaxAge: cfg.MaxAge,
	})

	SessionStore = store
	return nil
}

// Create a new session for the user
func CreateSession(ctx *gin.Context, id uint) error {
	session := sessions.Default(ctx)

	// Encode value before setting
	value, err := json.Marshal(id)
	if err != nil {
		return fmt.Errorf("failed to marshal session value: %v", err)
	}

	session.Set("user", value)
	if err := session.Save(); err != nil {
		return fmt.Errorf("failed to save session: %v", err)
	}
	return nil
}

// Delete the user's session
func DeleteSession(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Delete("user")
	// Invalidate the cookie
	session.Options(sessions.Options{MaxAge: -1})
	if err := session.Save(); err != nil {
		log.Printf("Failed to delete session: %v", err)
	}
}
