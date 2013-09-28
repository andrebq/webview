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
	"testing"
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
	if len(files) == 0 {
		t.Fatalf("files cannot be empty")
	}

	if files[0].Name() != "index.html" {
		t.Errorf("name should be index.html but got %v", files[0].Name())
	}

	dirs, err := vfs.ReadDirs()
	if err != nil {
		t.Errorf("unable to read vfs dirs %v", err)
	}
	if len(dirs) == 0 {
		t.Fatalf("dirs cannot be empty")
	}

	if dirs[0].Name() != "layout" {
		t.Errorf("name should be layout but got %v", dirs[0].Name())
	}

	files, err = dirs[0].ReadFiles()
	if err != nil {
		t.Fatalf("unable to read child files %v", err)
	}
	if len(files) == 0 {
		t.Fatalf("subdir files cannot be empty")
	}

	if files[0].Name() != "index.html" {
		t.Errorf("name should be layout but got %v", files[0].Name())
	}
}
