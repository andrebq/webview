package protoview

import (
	"encoding/json"
	"github.com/andrebq/webview/httpview"
	"log"
	"net/http"
	"os"
)

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

// Hold the configuration of the prototype
type Prototype struct {
	// Urls configuration
	Urls map[string]*Config
}

// The configuration of a given url
type Config struct {
	// The alias map
	Alias map[string]string
	// Data to be used inside the template
	Data interface{}
}

// Render the prototype to the http response
func (p *Prototype) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !p.CanHandle(req) {
		http.NotFound(w, req)
		return
	}
	config := p.Urls[req.URL.Path]
	if config.Alias != nil {
		httpview.SetAliasMap(req, config.Alias)
	}
	httpview.SetViewData(req, config.Data)
	httpview.Render(w, req)
}

// Check if the prototype can handle the given request
func (p *Prototype) CanHandle(req *http.Request) bool {
	if p.Urls == nil || len(p.Urls) == 0 {
		return false
	}
	_, has := p.Urls[req.URL.Path]
	return has
}

// Load the prototype form the given config file every time a request
// hits the returned handler.
//
// Usually the returned handler is registered under "/"
func PrototypeFromFile(file string, fallback http.Handler) http.Handler {
	handlefunc := func(w http.ResponseWriter, req *http.Request) {
		proto, err := loadProto(file)
		if err != nil {
			log.Printf("error loading prototype %v. cause: %v", file, err)
		}
		log.Printf("prototype %v loaded", file)
		if proto == nil || !proto.CanHandle(req) {
			fallback.ServeHTTP(w, req)
		} else {
			proto.ServeHTTP(w, req)
		}
	}
	return http.HandlerFunc(handlefunc)
}

func loadProto(path string) (*Prototype, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	proto := &Prototype{}
	err = json.NewDecoder(file).Decode(&proto)
	return proto, err
}
