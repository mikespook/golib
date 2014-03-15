package session

import (
	"net/http"
	"strings"
	"testing"
)

const (
	sessionKey = "TEST_SESSION"
	cookieKey  = "1234567890123456"
)

func init() {
	SetKey([]byte("0000000000000000"))
}

type testResponse struct {
	header http.Header
}

func newTestResponse() *testResponse                            { return &testResponse{make(http.Header)} }
func (resp *testResponse) Header() http.Header                  { return resp.header }
func (resp *testResponse) Write(data []byte) (n int, err error) { return }
func (resp *testResponse) WriteHeader(status int)               {}
func splitCookie(cookie string) *http.Cookie {
	str := strings.Split(cookie, ";")
	str = strings.SplitN(str[0], "=", 2)
	return &http.Cookie{Name: str[0], Value: str[1]}
}

func TestCookieStorage(t *testing.T) {
	r, err := http.NewRequest("GET", "http://127.0.0.1", nil)
	if err != nil {
		t.Error(err)
		return
	}
	storage := CookieStorage(sessionKey, nil)
	s := &Session{}
	if err := storage.LoadTo(r, s); err == nil {
		t.Errorf("No-named cookie error should be presented.")
		return
	}
	s.id = cookieKey

	resp := newTestResponse()
	s.w = resp
	s.Set("foo", 123)
	s.SetOption("max-age", 1234)
	storage.Flush(s)
	for _, setCookie := range resp.header["Set-Cookie"] {
		if strings.Index(setCookie, "Max-Age=1234") == -1 {
			t.Errorf("Options not effective: %s", setCookie)
			return
		}
		r.AddCookie(splitCookie(setCookie))
	}

	if err := storage.LoadTo(r, s); err != nil {
		t.Error(err)
		return
	}

	if s.Id() != cookieKey {
		t.Errorf("text[%s] != origin[%s]", s.Id(), cookieKey)
		return
	}

	foo := s.Get("foo")
	if v, ok := foo.(int); !ok || v != 123 {
		t.Errorf("Session load issue")
		return
	}
}
