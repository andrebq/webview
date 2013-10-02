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

	// Should return if the current item represents a directory
	IsDir() bool
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

// A set of compiled templates
type TreeSet map[string]*parse.Tree

// A virtual file
type File interface {
	Namer
	// Return the contents of this file
	Contents() ([]byte, error)
	Parenter
}

// Use to select which files should be parsed as a template
type Filter interface {
	// return
	Filter(f Namer) bool
}

// Implements the Filter interface
type FilterFunc func(f Namer) bool

// Call the function
func (ff FilterFunc) Filter(f Namer) bool {
	return ff(f)
}

// Return a new template containing all templates
// from set.
//
// The alias map can be used to access one template
// with two or more names
//
// Consider this:
//
// 	layout/main.html
//	The contents come from {{ template "contents" }}
//
//	index/index.html
//	I have the contents
//
//	user/index.html
//	I also have the contents
//
// If you loaded all those files, you have a treeSet with three
// templates: "layout/main.html", "index/index.html" and "user/index.html"
//
// Now you can use the alias to map "index/index.html" to "contents"
//
//	alias := map[string]string {
//		"contents": "index/index.html"
//	}
//
// When you execute the template, instead of having a "template contents not found"
// the system will execute the "index/index.html" template and put it's result
// on "layout/main.html"
//
// This keeps all the safety from html/template but enable your to use more
// dynamic templates without having to parse them every single time.
//
// If you need a new template with a different alias, just call this function again
// passing a different alias map
func Template(set TreeSet, alias map[string]string) (*template.Template, error) {
	t, _ := template.New("_root").Parse("")
	for k, v := range set {
		if _, err := t.AddParseTree(k, v); err != nil {
			return t, err
		}
	}
	for k, v := range alias {
		if tree, has := set[v]; has {
			if _, err := t.AddParseTree(k, tree); err != nil {
				return t, err
			}
		}
	}
	return t, nil
}

// Load all files under root and return a set of all templates
// if two templates have the same name (let's say that file a.html
// and b.html both define the template "nice_button"). Only one
// of those definitions will be available (the last one returned by
// the vfs)
//
// Each template can be accessed by its full path from root,
// that means "layout/body.html" represents a file under
// "layout" with a name of "body.html"
func LoadDir(root Dir, funcs tt.FuncMap, filter Filter) (TreeSet, error) {
	set := make(TreeSet)
	return set, LoadDirInto(set, root, funcs, filter)
}

// Load all files from the given Dir into the given template
func LoadDirInto(t TreeSet, dir Dir, funcs tt.FuncMap, filter Filter) error {
	files, err := dir.ReadFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		if filter == nil || filter.Filter(file) {
			err = LoadFileInto(t, file, funcs)
			if err != nil {
				return err
			}
		}
	}
	dirs, err := dir.ReadDirs()
	if err != nil {
		return err
	}
	for _, cdir := range dirs {
		if filter == nil || filter.Filter(cdir) {
			err = LoadDirInto(t, cdir, funcs, filter)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Read a file and register a new template under the filename
//
// The name is given by TemplateName(f)
func LoadFileInto(t TreeSet, f File, funcs tt.FuncMap) error {
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
	for k, v := range treeSet {
		t[k] = v
	}
	return nil
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
	if f.IsDir() && len(myName) > 0 {
		myName += "/"
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

// Filter function that allow only html,(js/json) and css files
//
// All directories are allowed too
func AllowHtmlJsAndCss(f Namer) bool {
	if f.IsDir() {
		return true
	}
	n := f.Name()
	return strings.HasSuffix(n, ".html") ||
		strings.HasSuffix(n, ".js") ||
		strings.HasSuffix(n, ".json") ||
		strings.HasSuffix(n, ".css")
}
