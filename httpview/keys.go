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

// the list of keys used accross this package

type key byte

const (
	treeSetKey      = key(0)
	aliasMapKey     = key(1)
	redirectInfoKey = key(2)
	viewNameKey     = key(3)
	dataKey         = key(4)
	layoutNameKey   = key(5)
)
