package memory

import (
	"time"

	kcache "github.com/koding/cache"
)

type Session struct {
	SessCache *kcache.MemoryTTL
}

func New() *Session {
	var s Session

	s.SessCache = kcache.NewMemoryWithTTL(24 * time.Hour)
	s.SessCache.StartGC(1 * time.Second)

	return &s
}

func (s *Session) Set(name string, value string) bool {
	if name == "" {
		return false
	}

	setErr := s.SessCache.Set(name, value)

	if setErr != nil {
		return false
	}

	return true
}

func (s *Session) Delete(name string) bool {
	if name == "" {
		return false
	}

	delErr := s.SessCache.Delete(name)

	if delErr != nil {
		return false
	}

	return true
}

func (s *Session) Get(name string) string {
	if name == "" {
		return ""
	}

	getToken, getErr := s.SessCache.Get(name)

	if getErr != nil {
		return ""
	}

	return getToken.(string)
}
