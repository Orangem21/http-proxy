package tokenfilter

import (
	"net/http"
	"net/http/httputil"

	"../utils"
)

const (
	tokenHeader = "X-Lantern-Auth-Token"
)

type TokenFilter struct {
	log   utils.Logger
	next  http.Handler
	token string
}

type optSetter func(f *TokenFilter) error

func TokenSetter(token string) optSetter {
	return func(f *TokenFilter) error {
		f.token = token
		return nil
	}
}

func Logger(l utils.Logger) optSetter {
	return func(f *TokenFilter) error {
		f.log = l
		return nil
	}
}

func New(next http.Handler, setters ...optSetter) (*TokenFilter, error) {
	f := &TokenFilter{
		log:   utils.NullLogger,
		next:  next,
		token: "",
	}
	for _, s := range setters {
		if err := s(f); err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (f *TokenFilter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if f.log.IsLevel(utils.DEBUG) {
		reqStr, _ := httputil.DumpRequest(req, true)
		f.log.Debugf("Token Filter Middleware received request:\n%s", reqStr)
	}

	token := req.Header.Get(tokenHeader)
	if token == "" || token != f.token {
		w.WriteHeader(http.StatusNotFound)
	} else {
		req.Header.Del(tokenHeader)
		f.next.ServeHTTP(w, req)
	}
}
