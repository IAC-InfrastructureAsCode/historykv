package session

import (
	"log"

	historySessionMemory "github.com/tokopedia/historykv/src/session/memory"
	historySessionRedis "github.com/tokopedia/historykv/src/session/redis"
)

type Session interface {
	Set(name string, value string) bool
	Delete(name string) bool
	Get(name string) string
}

func New(saveType string, redisAddress string) Session {
	var s Session

	if saveType == "redis" {
		log.Println("> Session  : Using Redis. Addr:", redisAddress)
		s = historySessionRedis.New(redisAddress)
	} else {
		log.Println("> Session  : Using Memory.")
		s = historySessionMemory.New()
	}

	return s
}
