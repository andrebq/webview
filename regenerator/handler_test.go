package regenerator

import (
	"bytes"
	"github.com/andrebq/gas"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRuntime(t *testing.T) {
	handler := &Regenerator{
		Compiler{
			Prefix:            "/src/",
			Dir:               []string{gas.MustAbs("github.com/andrebq/webview/regenerator")},
			Regenerator:       "regenerator",
			RegeneratorParams: []string{},
		},
	}

	server := httptest.NewServer(handler)
	defer server.Close()

	res, err := http.Get(server.URL + "/src/regenerator-runtime.js")
	if err != nil {
		t.Fatalf("unexptected error: %v", err)
	}

	fromRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	res.Body.Close()

	runtimeFile, err := gas.ReadFile("github.com/andrebq/webview/regenerator/regenerator-min.js")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if bytes.Compare(fromRes, runtimeFile) != 0 {
		t.Errorf("unexpected difference between files")
	}
}

func TestCompilation(t *testing.T) {
	handler := &Regenerator{
		Compiler{
			Prefix: "/src/",
			Dir: []string{
				gas.MustAbs("github.com/andrebq/webview/regenerator"),
				gas.MustAbs("github.com/andrebq/webview/regenerator/testdata")},
			Regenerator:       "regenerator",
			RegeneratorParams: []string{},
		},
	}

	server := httptest.NewServer(handler)
	defer server.Close()

	res, err := http.Get(server.URL + "/src/sample.js")
	if err != nil {
		t.Fatalf("unexptected error: %v", err)
	}

	fromRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	res.Body.Close()

	expected, err := gas.ReadFile("github.com/andrebq/webview/regenerator/testdata/sample-regenerated.js")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if bytes.Compare(fromRes, expected) != 0 {
		t.Errorf("unexpected difference between files")
		ioutil.WriteFile("fromres.out", fromRes, 0644)
		ioutil.WriteFile("expected.out", expected, 0644)
		t.Logf("formRes: \n%v", string(fromRes))
		t.Logf("expected: \n%v", string(expected))
	}
}
