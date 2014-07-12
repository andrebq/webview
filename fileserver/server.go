package fileserver

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
	"errors"
	"github.com/robertkrimen/otto"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"
)

// Handler implements the net/http.Handler interface and uses
// a javascript program to select which file should be loaded
// for a given request
type Handler struct {
	sync.RWMutex
	// Rules holds the code used to match the requests
	Rules *otto.Otto
	// Base is the base dir used to search for files
	Base string
	// Fallback is called if the Rules don't return
	// a valid path
	Fallback http.Handler
}

// ServeHTTP implements the Handler interface
func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.RLock()
	defer h.RUnlock()

	if h.Rules == nil {
		h.useFallback(w, req)
	} else {
		dest, err := h.processRequest(req)
		if err != nil {
			h.handleError(w, req, err)
			return
		}
		h.serveFile(w, req, dest)
	}
}

func (h *Handler) useFallback(w http.ResponseWriter, req *http.Request) {
	if h.Fallback == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	h.Fallback.ServeHTTP(w, req)
}

func (h *Handler) handleError(w http.ResponseWriter, req *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (h *Handler) serveFile(w http.ResponseWriter, req *http.Request, name string) {
	_, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		} else {
			h.handleError(w, req, err)
			return
		}
	}
	http.ServeFile(w, req, name)
}

func (h *Handler) processRequest(req *http.Request) (string, error) {
	vm := h.Rules.Copy()
	finalPath := filepath.Join(h.Base, filepath.FromSlash(path.Clean(req.URL.Path)))
	dest, err := vm.Get("processRequest")
	if err != nil {
		return "", err
	}
	if !dest.IsFunction() {
		return "", errors.New("processRequest should be a function")
	}
	finalPathValue, _ := otto.ToValue(finalPath)
	reqPathValue, _ := otto.ToValue(req.URL.Path)
	userAgentValue, _ := otto.ToValue(req.Header.Get("User-Agent"))
	dest, err = dest.Call(finalPathValue, reqPathValue, userAgentValue)
	return dest.String(), err
}
