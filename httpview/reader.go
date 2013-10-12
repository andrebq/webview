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
	"net/url"
	"strconv"
)

// Read values from a url.Values object
type Reader struct {
	Values url.Values
}

// Similar to ReadInt but ignores the ok flag
func (r *Reader) Int(name string, def int64) int64 {
	i, _ := r.ReadInt(name, def)
	return i
}

// Similar to ReadFloat but ignores the ok flag
func (r *Reader) Float(name string, def float64) float64 {
	f, _ := r.ReadFloat(name, def)
	return f
}

// Similar to Read but ignores the ok flag
func (r *Reader) Str(name string, def string) string {
	s, _ := r.Read(name, def)
	return s
}

// Similar to ReadBytes but ignores the ok flag
func (r *Reader) Bytes(name string, def []byte) []byte {
	b, _ := r.ReadBytes(name, def)
	return b
}

// Read a integer value from the values object, if no value is found
// or a value cannot be converted to a integer, return the default
//
// The boolean flag is true only if the value was read from the map
// and correctly converted
func (r *Reader) ReadInt(name string, def int64) (int64, bool) {
	if value, has := r.Read(name, ""); has {
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			return i, true
		} else {
			return def, false
		}
	} else {
		return def, false
	}
}

// Similar to ReadInt but returns a uint value
func (r *Reader) ReadUint(name string, def uint64) (uint64, bool) {
	if value, has := r.Read(name, ""); has {
		if i, err := strconv.ParseUint(value, 10, 64); err == nil {
			return i, true
		} else {
			return def, false
		}
	} else {
		return def, false
	}
}

// Similar to ReadInt but returns a float value
func (r *Reader) ReadFloat(name string, def float64) (float64, bool) {
	if value, has := r.Read(name, ""); has {
		if i, err := strconv.ParseFloat(value, 64); err == nil {
			return i, true
		} else {
			return def, false
		}
	} else {
		return def, false
	}
}

// Read a string value form the values object, if no value is found
// return the default
//
// The boolean flag is true only if the value was read form the map.
func (r *Reader) Read(name, def string) (string, bool) {
	if values, has := r.Values[name]; has {
		return values[0], true
	} else {
		return def, false
	}
}

// Similar to Read but returns the string as an array of bytes
func (r *Reader) ReadBytes(name string, def []byte) ([]byte, bool) {
	if value, has := r.Read(name, ""); has {
		return []byte(value), true
	} else {
		return def, false
	}
}
