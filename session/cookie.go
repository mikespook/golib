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
	options *http.Cookie
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

func fillOptions(cookie *http.Cookie) M {
	options := make(M)
	options[CookieDomain] = cookie.Domain
	options[CookieMaxAge] = cookie.MaxAge
	options[CookieExpires] = time.Now().Add(time.Duration(cookie.MaxAge) * time.Second)
	options[CookieHttpOnly] = cookie.HttpOnly
	options[CookiePath] = cookie.Path
	options[CookieSecure] = cookie.Secure
	return options
}

func mergeOptions(origin M, extend M) M {
	if origin == nil {
		return extend
	}
	for k, v := range extend {
		if _, ok := origin[k]; ok {
			continue
		}
		origin[k] = v
	}
	return origin
}

func cloneCookie(cookie *http.Cookie) *http.Cookie {
	newOne := &http.Cookie{
		Domain: cookie.Domain,
		MaxAge: cookie.MaxAge,
		Expires: cookie.Expires,
		HttpOnly: cookie.HttpOnly,
		Path: cookie.Path,
		Secure: cookie.Secure,
	}
	return newOne
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
	key := cloneCookie(storage.options)
	if s.options != nil {
		fillCookie(s.options, key)
	}
	key.Name = storage.keyName
	key.Value = s.id
	http.SetCookie(s.w, key)

	v, err := encoding([]byte(s.id), s.data)
	value := cloneCookie(storage.options)
	if s.options != nil {
		fillCookie(s.options, value)
	}
	value.Name = s.id
	value.Value = v
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
	s.options = mergeOptions(s.options, fillOptions(cookie))
	if err := decoding([]byte(s.id), cookie.Value, &s.data); err != nil {
		s.data = make(M)
		return err
	}
	return nil
}

func CookieStorage(keyName string, options M) Storage {
	if options == nil {
		options = DefaultCookieOptions
	}
	cookie := &http.Cookie{}
	fillCookie(options, cookie)
	storage := &cookieStorage{keyName, cookie}
	storage.options = &http.Cookie{}
	fillCookie(options, storage.options)
	return storage
}
