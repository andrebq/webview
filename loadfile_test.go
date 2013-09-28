package webview

// The MIT License (MIT)
//
// Copyright (c) 2013 Andre Luiz Alves Moraes
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:

import (
	"github.com/andrebq/gas"
	"io/ioutil"
	"testing"
)

func TestLoadDir(t *testing.T) {
	abs, err := gas.Abs("github.com/andrebq/webview/testdata")
	if err != nil {
		t.Fatalf("unable to load abs testdata dir %v", err)
	}

	vfs, err := DiskVFS(abs)
	if err != nil {
		t.Fatalf("unable to load vfs dir %v", err)
	}

	tmpl, err := LoadDir(vfs, nil, FilterFunc(AllowHtmlJsAndCss))
	if err != nil {
		t.Fatalf("unable to load template %v", err)
	}

	err = tmpl.ExecuteTemplate(ioutil.Discard, "index.html", nil)
	if err != nil {
		t.Fatalf("Unexpected error while rendering template %v", err)
	}
}
