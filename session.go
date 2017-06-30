package commons

import (
	"strconv"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/context"
	sessions "github.com/littletwolee/gin-sessions"
	//"github.com/gorilla/sessions"
)

var (
	consCSession                    *CSession
	maxAgeDefault                   = 7 * 24 * 3600
	address, password               string
	maxIdle, maxActive, idleTimeout int
	DefaultKey                      string
)

type CSession struct {
	Name   string
	maxAge int
	Store  sessions.RedisStore
}

func GetCSession() *CSession {
	var (
		err error
	)
	if consCSession == nil {
		consCSession = &CSession{}
	}
	DefaultKey = Config.GetString("session.defaultkey")
	consCSession.Name = Config.GetString("session.name")
	maxAgeStr := Config.GetString("session.maxage")
	consCSession.maxAge, err = strconv.Atoi(maxAgeStr)
	if err != nil {
		GetLogger().LogErr(err)
	}
	address = Config.GetString("redis.address")
	password = Config.GetString("redis.password")
	maxIdleStr := Config.GetString("redis.maxidle")
	maxIdle, err = strconv.Atoi(maxIdleStr)
	if err != nil {
		GetLogger().LogErr(err)
	}
	maxActiveStr := Config.GetString("redis.maxactive")
	maxActive, err = strconv.Atoi(maxActiveStr)
	if err != nil {
		GetLogger().LogErr(err)
	}
	idleTimeoutStr := Config.GetString("redis.idletimeout")
	idleTimeout, err = strconv.Atoi(idleTimeoutStr)
	if err != nil {
		GetLogger().LogErr(err)
	}
	return consCSession
}

func newPool() *redigo.Pool {
	return &redigo.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", address)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func (r *CSession) NewStore() (sessions.RedisStore, error) {
	if r.Store != nil {
		return r.Store, nil
	}
	kp := Config.GetString("session.keypaire")
	store, err := sessions.NewRediStoreWithPool(newPool(), []byte(kp))
	if err != nil {
		return nil, err
	}
	store.Options(sessions.Options{
		MaxAge: r.maxAge,
	})
	return store, nil
}

func (r *CSession) DefaultSession(c *gin.Context) error {
	store, err := r.NewStore()
	if err != nil {
		return err
	}
	s := sessions.GetIsession(r.Name, c.Request, store, nil, false, c.Writer)
	defer context.Clear(c.Request)
	session := s.Default(c, DefaultKey)
	var tt int64
	v := session.Get("deadline")
	if v == nil {
		tt = GetTimes().TimeToStamp(time.Now())
	}
	session.Set("deadline", tt)
	err = session.Save()
	if err != nil {
		return err
	}
	return nil
}