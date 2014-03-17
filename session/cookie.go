package session

import (
	"net/http"
	"time"
)

const (
	CookieDomain   = "domain"
	CookieExpires  = "expires"
	CookieHttpOnly = "http-only"
	CookieMaxAge   = "max-age"
	CookiePath     = "path"
	CookieSecure   = "secure"
)

var DefaultCookieOptions = M{
	CookieDomain:   "",
	CookieHttpOnly: false,
	CookieMaxAge:   3600,
	CookiePath:     "/",
	CookieSecure:   false,
}

var CleanCookieOptions = M{
	CookieExpires: time.Now(),
	CookieMaxAge:  0,
}

type cookieStorage struct {
	keyName string
	options M
}

func fillCookie(options M, cookie *http.Cookie) {
	if domain, ok := options[CookieDomain]; ok {
		cookie.Domain = domain.(string)
	}
	if maxAge, ok := options[CookieMaxAge]; ok {
		maxAge := maxAge.(int)
		cookie.MaxAge = maxAge
		cookie.Expires = time.Now().Add(time.Duration(maxAge) * time.Second)
	}
	if expires, ok := options[CookieExpires]; ok {
		cookie.Expires = expires.(time.Time)
	}
	if httpOnly, ok := options[CookieHttpOnly]; ok {
		cookie.HttpOnly = httpOnly.(bool)
	}
	if path, ok := options[CookiePath]; ok {
		cookie.Path = path.(string)
	}
	if secure, ok := options[CookieSecure]; ok {
		cookie.Secure = secure.(bool)
	}
}

func (storage *cookieStorage) Clean(s *Session) error {
	key := &http.Cookie{Name: storage.keyName}
	fillCookie(CleanCookieOptions, key)
	http.SetCookie(s.w, key)
	value := &http.Cookie{Name: s.id}
	fillCookie(CleanCookieOptions, value)
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
	s.storage = storage
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

func (storage *cookieStorage) SetOption(key string, value interface{}) {
	storage.options[key] = value
}

func CookieStorage(keyName string, options M) Storage {
	if options == nil {
		options = DefaultCookieOptions
	}
	return &cookieStorage{keyName, options}
}
