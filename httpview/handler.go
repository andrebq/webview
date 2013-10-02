package httpview

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
	"github.com/andrebq/webview"
	"github.com/gorilla/context"
	"net/http"
)

type key byte

const (
	treeSetKey = key(0)
	aliasKey   = key(1)
)

var (
	emptyMap = map[string]string{}
)

// Render the given view using the alias and treeset registered for the current request
//
// Use the RegisterView in your http pipeline beforer calling RenderView
func RenderView(w http.ResponseWriter, req *http.Request, name string, data interface{}) error {
	if tree, ok := context.GetOk(req, treeSetKey); ok {
		alias := context.Get(req, aliasKey)
		return renderViewFromTreeSet(w, req, tree.(webview.TreeSet), alias.(map[string]string), name, data)
	} else {
		return fmt.Errorf("webview treeset not found. are your sure you called RegisterView")
	}
}

// Register the treeset and the alias name for the current request
func RegisterView(req *http.Request, set webview.TreeSet, alias map[string]string) {
	context.Set(req, treeSetKey, set)
	context.Set(req, aliasKey, alias)
}

// Render the given template from the treeset using the given alias map
func renderViewFromTreeSet(w http.ResponseWriter, req *http.Request, set webview.TreeSet, alias map[string]string, view string, data interface{}) error {
	if alias == nil {
		alias = emptyMap
	}
	tmpl, err := webview.Template(set, alias)
	if err != nil {
		return err
	}
	err = tmpl.ExecuteTemplate(w, view, data)
	return err
}
