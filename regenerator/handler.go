// This package wrap the npm module regenerator from Facebook
// to allow more people to use ES6 without having to compile
// lot's of files every new deploy.
//
// If the user-agent is capable of handling es6 files, then no
// compilation is made.
package regenerator

import (
	"errors"
	"github.com/andrebq/gas"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CompilerCache links a source file to a regenerated file
type CompilerCache map[string]string

// Store will link source file to the regenerated one
func (cc *CompilerCache) Store(source, regen string) {
	(*cc)[source] = regen
}

// Open will try to open a regenerated file for the given source file
// if the mtime from source is after the regenerated, then this will
// result in a error.
//
// If source isn't found, an error is also returned.
func (cc *CompilerCache) Open(source string) (io.ReadCloser, error) {
	if regen, has := (*cc)[source]; has {
		return cc.openIfHot(source, regen)
	}
	return nil, errors.New("source not found")
}

func (cc *CompilerCache) openIfHot(src, regen string) (io.ReadCloser, error) {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return nil, err
	}
	regenInfo, err := os.Stat(regen)
	if err != nil {
		return nil, err
	}

	if srcInfo.ModTime().After(regenInfo.ModTime()) {
		return nil, errors.New("source is newer than regenerated")
	}
	return os.OpenFile(regen, os.O_RDONLY, 0644)
}

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

	// Cache for regenerated files
	Cache CompilerCache

	tmpDir string
}

// Compile will read the input file and output the result
// of callin regenerator on it.
//
// input should be separated by "/" instead of "\\" on Windows,
func (c *Compiler) Compile(out io.Writer, input string) error {
	input, err := c.FindFile(input)
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

func (c *Compiler) FindFile(in string) (string, error) {
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
	reader, err := c.fromCache(input)
	if err == nil {
		defer reader.Close()
		_, err := io.Copy(out, reader)
		return err
	}
	npmCmd := exec.Command(c.Regenerator, append(c.RegeneratorParams, input)...)
	pipedOut, err := npmCmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err = npmCmd.Start(); err != nil {
		return err
	}

	tmpOut, filename, err := c.makeTmpOut(input)
	if err != nil {
		// unable to create cached file
		// write to output
		io.Copy(out, pipedOut)
		return npmCmd.Wait()
	}
	defer tmpOut.Close()

	// we created the temporary output
	io.Copy(tmpOut, pipedOut)
	cmdOut := npmCmd.Wait()

	c.ensureCache()
	c.Cache.Store(input, filename)

	// ensure that everything is on disk
	tmpOut.Sync()
	tmpOut.Seek(0, 0)
	// copy
	io.Copy(out, tmpOut)
	return cmdOut
}

func (c *Compiler) makeTmpOut(input string) (*os.File, string, error) {
	c.ensureCache()
	if c.tmpDir == "" {
		var err error
		c.tmpDir, err = ioutil.TempDir("", "rengenerator")
		if err != nil {
			return nil, "", err
		}
	}
	tmpFile, err := ioutil.TempFile(c.tmpDir, filepath.Base(input))
	if err != nil {
		return tmpFile, "", err
	}
	return tmpFile, filepath.Join(c.tmpDir, tmpFile.Name()), err
}

func (c *Compiler) fromCache(input string) (io.ReadCloser, error) {
	c.ensureCache()
	return c.Cache.Open(input)
}

func (c *Compiler) ensureCache() {
	if c.Cache == nil {
		c.Cache = make(CompilerCache)
	}
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
