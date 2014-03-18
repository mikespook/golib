package session

import (
	"net/http"
)

type FactoryFunc func(http.ResponseWriter, *http.Request, M) (*Session, error)

func NewFactory(storage Storage) FactoryFunc {
	return func(w http.ResponseWriter, r *http.Request, options M) (s *Session, err error) {
		s = &Session{
			storage: storage,
			w:       w,
			options: options,
		}
		err = storage.LoadTo(r, s)
		return
	}
}
