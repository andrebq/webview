// This package wrap the npm module regenerator from Facebook
// to allow more people to use ES6 without having to compile
// lot's of files every new deploy.
//
// If the user-agent is capable of handling es6 files, then no
// compilation is made.
package regenerator

import (
	"github.com/andrebq/gas"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

// Compiler is responsible for calling npm
type Compiler struct {
	// Prefix to remove from incoming http request
	// before searching for files
	Prefix string

	// List of directories to search for a given
	// file. The first match is used
	Dir []string

	// Command used to activate NPM
	//
	// When empty the env var "WEBVIEW_REGENERATOR"
	Regenerator string

	// Default parameters
	RegeneratorParams []string
}

// Compile will read the input file and output the result
// of callin regenerator on it.
//
// input should be separated by "/" instead of "\\" on Windows,
func (c *Compiler) Compile(out io.Writer, input string) error {
	input, err := c.findFile(input)
	if err != nil {
		return err
	}
	return c.callRegenerator(input, out)
}

func (c *Compiler) Runtime(out io.Writer) error {
	file, err := gas.Open("github.com/andrebq/webview/regenerator/regenerator-min.js")
	if err != nil {
		file, err = gas.FromDirs(c.Dir).Open("regenerator-min.js")
		if err != nil {
			return err
		}
	}
	defer file.Close()
	_, err = io.Copy(out, file)
	return err
}

func (c *Compiler) findFile(in string) (string, error) {
	fs := gas.FromDirs(c.Dir)
	return fs.Abs(c.stripPrefix(in), false)
}

func (c *Compiler) stripPrefix(in string) string {
	if strings.Index(in, c.Prefix) == 0 {
		return in[len(c.Prefix):]
	}
	return in
}

func (c *Compiler) callRegenerator(input string, out io.Writer) error {
	npmCmd := exec.Command(c.Regenerator, append(c.RegeneratorParams, input)...)
	pipedOut, err := npmCmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err = npmCmd.Start(); err != nil {
		return err
	}

	io.Copy(out, pipedOut)

	return npmCmd.Wait()
}

// Regenerator wrap the compiler under the net/http.Handler interface
type Regenerator struct {
	Compiler
}

func (r *Regenerator) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if strings.HasSuffix(req.URL.Path, "regenerator-runtime.js") {
		err := r.Runtime(w)
		if err != nil {
			panic(err)
		}
	} else if strings.HasSuffix(req.URL.Path, ".js") {
		err := r.Compile(w, req.URL.Path)
		if err != nil {
			panic(err)
		}
	}
}
