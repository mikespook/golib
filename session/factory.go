package session

import (
	"net/http"
)

type FactoryFunc func(http.ResponseWriter, *http.Request) *Session

func NewFactory(storage Storage) FactoryFunc {
	return func(w http.ResponseWriter, r *http.Request) (s *Session) {
		s = &Session{
			storage: storage,
			w:       w,
		}
		storage.LoadTo(r, s)
		return
	}
}
