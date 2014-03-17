package session

import (
	"crypto/rand"
	"fmt"
	"net/http"
)

const IdLength = 32

type M map[string]interface{}

type Session struct {
	data    M
	id      string
	storage Storage
	w       http.ResponseWriter
}

func (s *Session) Id() string {
	return s.id
}

func (s *Session) Set(key string, value interface{}) {
	s.data[key] = value
}

func (s *Session) Get(key string) (value interface{}) {
	return s.data[key]
}

func (s *Session) Del(key string) (value interface{}) {
	value = s.data[key]
	delete(s.data, key)
	return
}

func (s *Session) Init() {
	s.data = make(M)
	s.id = fmt.Sprintf("%x", genKey(IdLength))
}

func (s *Session) Clean() error {
	return s.storage.Clean(s)
}

func (s *Session) Flush(options M) error {
	return s.storage.Flush(s, options)
}

func genKey(size int) []byte {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return nil
	}
	return b
}
