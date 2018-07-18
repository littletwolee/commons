package commons

// var (
// 	consCSession                    *CSession
// 	maxAgeDefault                   = 7 * 24 * 3600
// 	address, password               string
// 	maxIdle, maxActive, idleTimeout int
// 	DefaultKey                      string
// )

// type CSession struct {
// 	Name   string
// 	maxAge int
// 	Store  sessions.RedisStore
// }

// // @Title GetCSession
// // @Description init & get CSession point
// // @Parameters
// // @Returns csession:*CSession
// func GetCSession() *CSession {
// 	var (
// 		err error
// 	)
// 	if consCSession == nil {
// 		consCSession = &CSession{}
// 	}
// 	DefaultKey = GetConfig().GetString("session.keypaire")
// 	consCSession.Name = GetConfig().GetString("session.name")
// 	maxAgeStr := GetConfig().GetString("session.maxage")
// 	consCSession.maxAge, err = strconv.Atoi(maxAgeStr)
// 	if err != nil {
// 		GetLogger().LogErr(err)
// 	}
// 	address = GetConfig().GetString("redis.address")
// 	password = GetConfig().GetString("redis.password")
// 	maxIdleStr := GetConfig().GetString("redis.maxidle")
// 	maxIdle, err = strconv.Atoi(maxIdleStr)
// 	if err != nil {
// 		GetLogger().LogErr(err)
// 	}
// 	maxActiveStr := GetConfig().GetString("redis.maxactive")
// 	maxActive, err = strconv.Atoi(maxActiveStr)
// 	if err != nil {
// 		GetLogger().LogErr(err)
// 	}
// 	idleTimeoutStr := GetConfig().GetString("redis.idletimeout")
// 	idleTimeout, err = strconv.Atoi(idleTimeoutStr)
// 	if err != nil {
// 		GetLogger().LogErr(err)
// 	}
// 	return consCSession
// }

// // @Title newPool
// // @Description create a new connection pool
// // @Parameters
// // @Returns pool:*redigo.Pool
// func newPool() *redigo.Pool {
// 	return &redigo.Pool{
// 		MaxIdle:     maxIdle,
// 		MaxActive:   maxActive,
// 		IdleTimeout: time.Duration(idleTimeout) * time.Second,
// 		Dial: func() (redigo.Conn, error) {
// 			c, err := redigo.Dial("tcp", address)
// 			if err != nil {
// 				return nil, err
// 			}
// 			if password != "" {
// 				if _, err := c.Do("AUTH", password); err != nil {
// 					c.Close()
// 					return nil, err
// 				}
// 			}
// 			return c, err
// 		},
// 		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
// 			if time.Since(t) < time.Minute {
// 				return nil
// 			}
// 			_, err := c.Do("PING")
// 			return err
// 		},
// 	}
// }

// // @Title NewStore
// // @Description get a new store
// // @Parameters
// // @Returns store:sessions.RedisStore err:error
// func (r *CSession) NewStore() (sessions.RedisStore, error) {
// 	if r.Store != nil {
// 		return r.Store, nil
// 	}
// 	kp := GetConfig().GetString("session.keypaire")
// 	store, err := sessions.NewRediStoreWithPool(newPool(), []byte(kp))
// 	if err != nil {
// 		return nil, err
// 	}
// 	store.Options(sessions.Options{
// 		MaxAge: r.maxAge,
// 	})
// 	return store, nil
// }

// // @Title DefaultSession
// // @Description get default session
// // @Parameters
// //       c         *gin.Context          gin context
// // @Returns err:error
// func (r *CSession) DefaultSession(c *gin.Context) error {
// 	store, err := r.NewStore()
// 	if err != nil {
// 		return err
// 	}
// 	s := sessions.GetIsession(r.Name, c.Request, store, nil, false, c.Writer)
// 	defer context.Clear(c.Request)
// 	session := s.Default(c, DefaultKey)
// 	var tt int64
// 	v := session.Get("deadline")
// 	if v == nil {
// 		tt = GetTimes().TimeToStamp(time.Now())
// 	}
// 	session.Set("deadline", tt)
// 	err = session.Save()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
