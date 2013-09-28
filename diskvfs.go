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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// use the filepath to implement the vfs
type diskvfs struct {
	root, cd string
	isfile   bool
}

// create a new VFS at the given root
func DiskVFS(root string) (Dir, error) {
	stat, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("%v isn't a valid directory", root)
	}
	return &diskvfs{root: root, cd: "", isfile: false}, nil
}

func (vfs *diskvfs) Name() string {
	if strings.EqualFold(vfs.root, vfs.cd) {
		return ""
	}
	return filepath.Base(vfs.cd)
}

func (vfs *diskvfs) Parent() Dir {
	if strings.EqualFold(vfs.root, vfs.cd) {
		return nil
	}
	return &diskvfs{root: vfs.root, cd: filepath.Dir(vfs.cd), isfile: false}
}

func (vfs *diskvfs) ReadDirs() ([]Dir, error) {
	f, err := os.Open(vfs.cd)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	childs, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}
	ret := make([]Dir, len(childs), 0)
	for _, v := range childs {
		if !v.IsDir() {
			continue
		}
		ret = append(ret, &diskvfs{
			root:   vfs.root,
			cd:     filepath.Join(vfs.cd, f.Name()),
			isfile: false,
		})
	}
	return ret, nil
}

func (vfs *diskvfs) ReadFiles() ([]File, error) {
	f, err := os.Open(vfs.cd)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	childs, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}
	ret := make([]File, len(childs), 0)
	for _, v := range childs {
		if v.IsDir() {
			continue
		}
		ret = append(ret, &diskvfs{
			root:   vfs.root,
			cd:     filepath.Join(vfs.cd, f.Name()),
			isfile: true,
		})
	}
	return ret, nil
}

func (vfs *diskvfs) Contents() ([]byte, error) {
	f, err := os.Open(vfs.cd)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
