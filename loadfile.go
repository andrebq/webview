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
	"html/template"
	"strings"
	tt "text/template"
	"text/template/parse"
	"unicode/utf8"
)

// Represent any entry in the filesystem that can have a parent
type Parenter interface {
	// Return the parent of this directory, might return nil,
	// in this case, this Dir is the root
	Parent() Dir
}

// A item that have a name
type Namer interface {
	// The name of the item, only text that is allowed inside a URL
	// can be returned.
	//
	// Unicode MUST BE encoded as UTF-8 instead of using %xx
	//
	// The name here should only return its name without any information
	// from its parents
	//
	// Directories MUST NOT include the "/"
	Name() string
}

// A virtual directory
type Dir interface {
	Namer
	Parenter
	// Return the files that hold any content that can be parsed from
	// the template engine
	ReadFiles() ([]File, error)
	// Return the list of sub-directories
	ReadDirs() ([]Dir, error)
}

// A virtual file
type File interface {
	Namer
	// Return the contents of this file
	Contents() ([]byte, error)
	Parenter
}

// Load all files under root and return a root template
//
// Each template can be accessed by its full path from root,
// that means "layout/body.html" represents a file under
// "layout" with a name of "body.html"
//
// This function is equivalent to
//
//	t, _ := template.New("root").Parse("")
//	return t, LoadDirInto(t, root)
//
// This means you cannot have a template named root, and all
// calls to t.Execute will result in a empty result
//
// You should call t.ExecuteTemplate(reader, "name/of/your/template", data)
func LoadDir(root Dir, funcs tt.FuncMap) (*template.Template, error) {
	t, _ := template.New("root").Parse("")
	return t, LoadDirInto(t, root, funcs)
}

// Load all files from the given Dir into the given template
func LoadDirInto(t *template.Template, dir Dir, funcs tt.FuncMap) error {
	files, err := dir.ReadFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		err = LoadFileInto(t, file, funcs)
		if err != nil {
			return err
		}
	}
	dirs, err := dir.ReadDirs()
	if err != nil {
		return err
	}
	for _, cdir := range dirs {
		err = LoadDirInto(t, cdir, funcs)
		if err != nil {
			return err
		}
	}
	return nil
}

// Read a file and register a new template under the filename
//
// The name is given by TemplateName(f)
func LoadFileInto(t *template.Template, f File, funcs tt.FuncMap) error {
	name := TemplateName(f)
	contents, err := f.Contents()
	if err != nil {
		return err
	}
	if !utf8.Valid(contents) {
		return fmt.Errorf("file %v must be encoded using utf8", name)
	}
	leftDelim, rightDelim := DiscoverDelim(name)
	treeSet, err := parse.Parse(name,
		string(contents),
		leftDelim,
		rightDelim,
		funcs)
	if err != nil {
		return err
	}
	_, err = t.AddParseTree(name, treeSet[name])
	return err
}

// Return the unique name of the object.
//
// If the object is also a Parenter, the name will contain the name of it's parents
//
func TemplateName(f Namer) string {
	myName := f.Name()
	if p, ok := f.(Parenter); ok {
		parent := p.Parent()
		if parent != nil {
			myName = TemplateName(parent) + myName
		}
	}
	if _, ok := f.(Dir); ok {
		myName = myName + "/"
	}
	return myName
}

// Discover the best delimiter for the given filename
//
// HTML => {{ / }}
// JS => <% / %>
// CSS => <% / %>
func DiscoverDelim(name string) (string, string) {
	if strings.HasSuffix(name, ".html") {
		return "{{", "}}"
	} else if strings.HasSuffix(name, ".js") || strings.HasSuffix(name, ".json") {
		return "<%", "%>"
	} else if strings.HasSuffix(name, ".css") {
		return "<%", "%>"
	}
	return "{{", "}}"
}
