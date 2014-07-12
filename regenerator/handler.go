// This package wrap the npm module regenerator from Facebook
// to allow more people to use ES6 without having to compile
// lot's of files every new deploy.
//
// If the user-agent is capable of handling es6 files, then no
// compilation is made.
package regenerator

import (
	"github.com/andrebq/gas"
	"net/http"
	"strings"
)

// Regenerator wrap the compiler under the net/http.Handler interface
type Regenerator struct {
	Compiler
}

func (r *Regenerator) wantRuntime(req *http.Request) bool {
	return strings.HasSuffix(req.URL.Path, "regenerator-runtime.js")
}

func (r *Regenerator) serveRuntime(w http.ResponseWriter, req *http.Request) {
	file, err := gas.FromDirs(r.Compiler.Dir).Abs("regenerator-min.js", false)
	if err != nil {
		file, err = gas.Abs("github.com/andrebq/webview/regenerator/regenerator-min.js")
		if err != nil {
			if gas.IsNotFound(err) {
				http.Error(w, "not found", http.StatusNotFound)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	http.ServeFile(w, req, file)
}

func (r *Regenerator) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if strings.HasSuffix(req.URL.Path, ".js") {
		if r.wantRuntime(req) {
			r.serveRuntime(w, req)
			return
		}
		err := r.Compile(w, req.URL.Path)
		if err != nil {
			if gas.IsNotFound(err) {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func New(prefix, regenerator string, dirs []string, args ...string) *Regenerator {
	return &Regenerator{
		Compiler{
			Prefix:            prefix,
			Dir:               dirs,
			Regenerator:       regenerator,
			RegeneratorParams: args,
		},
	}
}
