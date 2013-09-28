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
package webview

import (
	"io"
)

// A item that have a name
type Namer interface {
	// The name of the item, only text that is allowed inside a URL
	// can be returned.
	//
	// Unicode MUST BE encoded as UTF-8 instead of using %xx
	Name() string
}

// A virtual directory
type Dir interface {
	Namer
	// Return the files that hold any content that can be parsed from
	// the template engine
	ReadFiles() ([]Files, error)
	// Return the list of sub-directories
	ReadDirs() ([]Dir, error)
	// Return the parent of this directory, might return nil,
	// in this case, this Dir is the root
	Parent() (Dir, error)
}

// A virtual file
type File interface {
	Namer
	io.Reader
	// Return the parent of this directory, might return nil,
	// in this case, this Dir is the root
	Parent() (Dir, error)
}
