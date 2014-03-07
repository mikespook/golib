package session

import (
	"net/http"
)

type FactoryFunc func(http.ResponseWriter, *http.Request) (*Session, error)

func NewFactory(storage Storage) FactoryFunc {
	return func(w http.ResponseWriter, r *http.Request) (s *Session, err error) {
		s = &Session{
			storage: storage,
			w:       w,
		}
		err = storage.LoadTo(r, s)
		return
	}
}
