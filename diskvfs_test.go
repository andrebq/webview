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
	"testing"
	"github.com/andrebq/gas"
)

func TestDiskVFS(t *testing.T) {
	abs, err := gas.Abs("github.com/andrebq/webview/testdata")
	if err != nil {
		t.Fatalf("unable to open gas resource %v", err)
	}

	vfs, err := DiskVFS(abs)
	if err != nil {
		t.Fatalf("unable to open vfs %v", err)
	}

	files, err := vfs.ReadFiles()
	if err != nil {
		t.Errorf("unable to read vfs files %v", err)
	}
	if files == nil {
		t.Fatalf("files cannot be null")
	}

	if files[0].Name() != "index.html" {
		t.Errorf("name should be index.html but got %v", files[0].Name())
	}
}
