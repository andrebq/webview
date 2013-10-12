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
	"net/url"
)

var (
	emptyMap = map[string]string{}
)

// Set the name of the view that should be rendered
// at Render
func SetViewName(req *http.Request, name string) {
	context.Set(req, viewNameKey, name)
}

// Return the name of the View from the request
func GetViewName(req *http.Request) string {
	if v, ok := context.GetOk(req, viewNameKey); !ok {
		return "index/index.html"
	} else {
		return v.(string)
	}
}

// Set the layout that will be used.
//
// To render the view the layout should call {{ template "contents" }}
func SetLayoutName(req *http.Request, name string) {
	context.Set(req, layoutNameKey, name)
}

// Return the layout defined, if nothing was set, returns
// layout/main.html
func GetLayoutName(req *http.Request) string {
	if v, ok := context.GetOk(req, layoutNameKey); !ok {
		return "layout/main.html"
	} else {
		return v.(string)
	}
}

// Set the data that should be used to render the
// template
func SetViewData(req *http.Request, data interface{}) {
	context.Set(req, dataKey, data)
}

// Render the view configured to that request
//
// This method will redirect if any of the Redirect* methods were called
// or will try to render the view configured with SetView{Name/Data}
func Render(w http.ResponseWriter, req *http.Request) {
	if redirect, ok := context.GetOk(req, redirectInfoKey); ok {
		// should return a redirect
		http.Redirect(w, req, redirect.(*url.URL).String(), http.StatusFound)
	} else {
		// grab the name and the data
		// from the request
		name := GetViewName(req)

		data, _ := context.GetOk(req, dataKey)

		// do the actual rendering
		RenderView(w, req, name, data)
	}
}

// Render the given view using the alias and treeset registered for the current request
//
// Use the RegisterView in your http pipeline beforer calling RenderView
func RenderView(w http.ResponseWriter, req *http.Request, name string, data interface{}) error {
	if tree, ok := context.GetOk(req, treeSetKey); ok {
		alias := GetAliasMap(req)
		provideDefaults(alias, req)
		return renderViewFromTreeSet(w, req, tree.(webview.TreeSet), alias, "main", data)
	} else {
		return fmt.Errorf("webview treeset not found. are your sure you called RegisterView")
	}
}

func provideDefaults(alias map[string]string, req *http.Request) {
	if _, has := alias["main"]; !has {
		alias["main"] = GetLayoutName(req)
	}
	if _, has := alias["contents"]; !has {
		alias["contents"] = GetViewName(req)
	}
}

// Set the alias that will be used to render the template
func SetAliasMap(req *http.Request, alias map[string]string) {
	context.Set(req, aliasMapKey, alias)
}

// Get the alias that will be used to render the template
func GetAliasMap(req *http.Request) map[string]string {
	if v, ok := context.GetOk(req, aliasMapKey); !ok {
		alias := map[string]string{
			"contents": GetViewName(req),
			"main":     GetLayoutName(req),
		}
		SetAliasMap(req, alias)
		return alias
	} else {
		return v.(map[string]string)
	}
}

// Register the treeset and the alias name for the current request
func RegisterView(req *http.Request, set webview.TreeSet) {
	context.Set(req, treeSetKey, set)
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
