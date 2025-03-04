package config

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

var SessionStore redis.Store

func InitSession(addr string, pwd string, secretKey string) error {
	store, err := redis.NewStore(10, "tcp", addr, pwd, []byte(secretKey))
	if err != nil {
		log.Fatal("Failed to Start Redis Server")
		return err
	}

	store.Options(sessions.Options{
		MaxAge: 600,
	})

	SessionStore = store
	return nil
}

func CreateSession(ctx *gin.Context, id uint) error {
	session := sessions.Default(ctx)

	// Encode value before setting
	value, err := json.Marshal(id)
	if err != nil {
		return err
	}

	session.Set("user", value)
	if err := session.Save(); err != nil {
		return err
	}
	return nil
}

func DeleteSession(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Delete("user")
	// Invalidate the cookie
	session.Options(sessions.Options{MaxAge: -1})
	if err := session.Save(); err != nil {
		fmt.Println(err)
	}
}
