package session

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"sync"
)

const IdLength = 32

type M map[string]interface{}

type Session struct {
	sync.RWMutex
	data    M
	id      string
	storage Storage
	w       http.ResponseWriter
	options M
}

func (s *Session) Id() string {
	return s.id
}

func (s *Session) Set(key string, value interface{}) {
	defer s.Unlock()
	s.Lock()
	s.data[key] = value
}

func (s *Session) Get(key string) (value interface{}) {
	defer s.RUnlock()
	s.RLock()
	return s.data[key]
}

func (s *Session) Del(key string) (value interface{}) {
	defer s.Unlock()
	s.Lock()
	value = s.data[key]
	delete(s.data, key)
	return
}

func (s *Session) Init() error {
	defer s.Unlock()
	s.Lock()
	s.data = make(M)
	s.id = fmt.Sprintf("%x", genKey(IdLength))
	return nil
}

func (s *Session) Clean() error {
	defer s.Unlock()
	s.Lock()
	return s.storage.Clean(s)
}

func (s *Session) Flush() error {
	defer s.Unlock()
	s.Lock()
	return s.storage.Flush(s)
}

func genKey(size int) []byte {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return nil
	}
	return b
}
