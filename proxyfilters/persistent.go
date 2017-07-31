package proxyfilters

import (
	"net/http"
	"strconv"

	"github.com/getlantern/proxy/filters"
)

const (
	// XLanternPersistent is the X-Lantern-Persistent header that indicates
	// persistent connections are to be used.
	xLanternPersistent = "X-Lantern-Persistent"
)

// DiscardInitialPersistentRequest discards the initial request for persistent
// HTTP connections from the Lantern client.
var DiscardInitialPersistentRequest = filters.FilterFunc(func(ctx filters.Context, req *http.Request, next filters.Next) (*http.Response, filters.Context, error) {
	isInitialPersistent, _ := strconv.ParseBool(req.Header.Get(xLanternPersistent))
	if isInitialPersistent {
		return filters.Discard(ctx, req)
	}
	return next(ctx, req)
})
