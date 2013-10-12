package httpview

import (
	"github.com/gorilla/context"
	"net/http"
	"net/url"
)

// Set the redirect information for the given request
func RedirectLocal(req *http.Request, path string) {
	context.Set(req, redirectInfoKey, makeRedirectFor(req, &url.URL{Path: path}))
}

// Return a URL from the given hos
func makeRedirectFor(req *http.Request, path *url.URL) *url.URL {
	return cleanUrl(req.URL.ResolveReference(path))
}

// Remove everything from the url except for
// Host/Port/Scheme/Path
func cleanUrl(dirty *url.URL) *url.URL {
	copy := *dirty
	copy.RawQuery = ""
	copy.Opaque = ""
	copy.Fragment = ""
	copy.User = nil
	return &copy
}
