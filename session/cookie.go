package session

import (
	"net/http"
	"time"
)

var optionsDefault = M{
	"domain":    "",
	"expires":   time.Now(),
	"http-only": false,
	"max-age":   3600,
	"path":      "/",
	"secure":    false,
}

var optionsClean = M{
	"expires": time.Now(),
	"max-age": 0,
}

type cookieStorage struct {
	keyName string
	options M
}

func fillCookie(options M, cookie *http.Cookie) {
	if domain, ok := options["domain"]; ok {
		cookie.Domain = domain.(string)
	}
	if expires, ok := options["expires"]; ok {
		cookie.Expires = expires.(time.Time)
	}
	if httpOnly, ok := options["http-only"]; ok {
		cookie.HttpOnly = httpOnly.(bool)
	}
	if maxAge, ok := options["max-age"]; ok {
		cookie.MaxAge = maxAge.(int)
	}
	if path, ok := options["path"]; ok {
		cookie.Path = path.(string)
	}
	if secure, ok := options["secure"]; ok {
		cookie.Secure = secure.(bool)
	}
}

func (storage *cookieStorage) Clean(s *Session) error {
	key := &http.Cookie{Name: storage.keyName}
	fillCookie(optionsClean, key)
	http.SetCookie(s.w, key)
	value := &http.Cookie{Name: s.id}
	fillCookie(optionsClean, value)
	http.SetCookie(s.w, value)
	s.Init()
	return nil
}

func (storage *cookieStorage) Flush(s *Session) error {
	key := &http.Cookie{
		Name:  storage.keyName,
		Value: s.id,
	}
	fillCookie(storage.options, key)
	http.SetCookie(s.w, key)
	v, err := encoding([]byte(s.id), s.data)
	value := &http.Cookie{
		Name:  s.id,
		Value: v,
	}
	fillCookie(storage.options, value)
	http.SetCookie(s.w, value)
	return err
}

func (storage *cookieStorage) LoadTo(r *http.Request, s *Session) error {
	cookie, err := r.Cookie(storage.keyName)
	if err != nil {
		s.Init()
		return err
	}
	s.id = cookie.Value
	cookie, err = r.Cookie(cookie.Value)
	if err != nil {
		s.Init()
		return err
	}
	if err := decoding([]byte(s.id), cookie.Value, &s.data); err != nil {
		s.data = make(M)
		return err
	}
	return nil
}

func CookieStorage(keyName string, options M) Storage {
	if options == nil {
		options = optionsDefault
	}
	return &cookieStorage{keyName, options}
}
