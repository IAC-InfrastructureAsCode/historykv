package memory

import (
	"fmt"
	"log"
	"time"

	"github.com/xuyu/goredis"
)

type Session struct {
	SessCache *goredis.Redis
}

func New(address string) *Session {
	var s Session

	if address == "" {
		return nil
	}

	redisCache, redisErr := goredis.Dial(&goredis.DialConfig{
		Address: address,
		MaxIdle: 1,
		Timeout: 5 * time.Second,
	})

	if redisErr == nil && redisCache.Ping() == nil {
		s.SessCache = redisCache
	} else {
		log.Fatalln("Couldn't connect to redis server:", address)
	}

	return &s
}

func (s *Session) Set(name string, value string) bool {
	if name == "" {
		return false
	}

	name = fmt.Sprintf("hkv:%s", name)

	setErr := s.SessCache.Setex(name, 24*60*60, value)

	if setErr != nil {
		return false
	}

	return true
}

func (s *Session) Delete(name string) bool {
	if name == "" {
		return false
	}

	name = fmt.Sprintf("hkv:%s", name)

	delTotal, delErr := s.SessCache.Del(name)

	if delErr != nil || delTotal != 1 {
		return false
	}

	return true
}

func (s *Session) Get(name string) string {
	if name == "" {
		return ""
	}

	name = fmt.Sprintf("hkv:%s", name)

	getToken, getErr := s.SessCache.Get(name)

	if getErr != nil {
		return ""
	}

	return string(getToken)
}
