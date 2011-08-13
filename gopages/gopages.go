// Copyright 2010 Abiola Ibrahim <abiola89@gmail.com>. All rights reserved.
// Use of this source code is governed by New BSD License
// http://www.opensource.org/licenses/bsd-license.php
// The content and logo is governed by Creative Commons Attribution 3.0
// The mascott is a property of Go governed by Creative Commons Attribution 3.0
// http://creativecommons.org/licenses/by/3.0/

package main

import (
	"util"
	"os"
	"strings"
	"flag"
	"path"
)

const (
	MAKE    = iota
	GOBUILD = iota
)

//where the build execution starts
func main() {
	make := flag.Bool("make", false, "build package with make (a Makefile must exist)")
	gobuild := flag.Bool("gobuild", true, "build package with gobuild (gobuild must be in $GOBIN)")
	flags := flag.String("flags", "", "flags for build tool selected e.g. '-lib' for gobuild or 'install' for make")
	cl := flag.Bool("clean", false, "don't build, just clean the generated pages")
	run := flag.Bool("run", false, "run the generated executable (only works with gobuild)")
	flag.Parse()
	if *cl {
		err := clean()
		if err != nil {
			println(err.String())
		}
		return
	}
	settings, err := util.LoadSettings() //inits the settings and generates the .go source files
	if err != nil {
		println(err.String())
		return
	}
	util.Config = settings.Data //stores settings to accessible variable
	println("generated", len(settings.Data["pages"]), "gopages")
	err = util.AddHandlers(settings.Data["handle"]) //add all handlers
	if err != nil {
		println(err.String())
		return
	}
	var do int
	if *run {
		err = build(GOBUILD, "-run") //build with corresponding build tool
		if err != nil {
			println(err.String())
		}
		return
	} else if *make {
		do = MAKE
	} else if *gobuild {
		do = GOBUILD
	}
	err = build(do, *flags) //build with corresponding build tool
	if err != nil {
		println(err.String())
	}
}

//create the pages directory to store generated source codes
func init() {
	err := os.MkdirAll(util.DIR, 0755)
	if err != nil {
		println(err.String())
		os.Exit(1)
	}
}

//to build the project with gobuild or make after generating .go source files
func build(b int, s string) (err os.Error) {
	p := strings.Fields(s)
	params := make([]string, len(p)+1)
	params[0] = ""
	for i, param := range p {
		params[i+1] = param
	}
	m := b == MAKE
	g := b == GOBUILD
	//fd := []*os.File{os.Stdin, os.Stdout, os.Stderr}
	if m {
		err = os.Exec("/usr/bin/make", params, os.Environ())
		if err != nil {
			panic("Cannot call make")
		} 
	} else if g {
		gobuild := os.Getenv("GOBIN")
		if len(gobuild) == 0 {
			gobuild = path.Join(os.Getenv("HOME"), "bin", "gobuild")
		} else {
			gobuild = path.Join(gobuild, "gobuild")
		}
		err := os.Exec(gobuild, params, os.Environ())
		if err != nil {
			panic("Cannot call gobuild")
		} 
	}
	return
}
//deletes the generated source codes
func clean() (err os.Error) {
	err = os.RemoveAll(util.DIR)
	if err != nil {
		println(err.String())
	}
	build(GOBUILD, "-clean")
	return
}
