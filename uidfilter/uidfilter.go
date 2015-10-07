package uidfilter

import (
	"net/http"
	"net/http/httputil"

	"github.com/gorilla/context"

	"../utils"
)

type UIDFilter struct {
	log  utils.Logger
	next http.Handler
}

type optSetter func(f *UIDFilter) error

func Logger(l utils.Logger) optSetter {
	return func(f *UIDFilter) error {
		f.log = l
		return nil
	}
}

func New(next http.Handler, setters ...optSetter) (*UIDFilter, error) {
	f := &UIDFilter{
		log:  utils.NullLogger,
		next: next,
	}
	for _, s := range setters {
		if err := s(f); err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (f *UIDFilter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	reqStr, _ := httputil.DumpRequest(req, true)
	f.log.Debugf("UIDFilter Middleware received request:\n%s", reqStr)

	lanternUID := req.Header.Get(utils.UIDHeader)

	// An UID must be provided always by the client.  Respond 404 otherwise.
	if lanternUID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get the client and attach it as request context
	key := []byte(lanternUID)
	client := utils.GetClient(key)
	context.Set(req, utils.ClientKey, client)

	req.Header.Del(utils.UIDHeader)

	f.next.ServeHTTP(w, req)
}