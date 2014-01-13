// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	eval "github.com/sbinet/go-eval"
	"github.com/sbinet/liner"
)

var fset = token.NewFileSet()
var filename = flag.String("f", "", "file to run")

var term *liner.State = nil

func init() {
	fmt.Println(`
*********************************************
** Interactive Go interpreter (with liner) **
*********************************************

`)
	term = liner.NewLiner()

	fname := path.Join(os.Getenv("HOME"), ".go.history")
	f, err := os.Open(fname)
	if err != nil {
		fmt.Printf("**warning: could not access history file [%s]\n", fname)
		return
	}
	defer f.Close()
	_, err = term.ReadHistory(f)
	if err != nil {
		fmt.Printf("**warning: could not read history file [%s]\n", fname)
		return
	}
}

func atexit() {
	fname := path.Join(os.Getenv("HOME"), ".go.history")
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("**warning: could not access history file [%s]\n", fname)
		return
	}
	defer f.Close()
	_, err = term.WriteHistory(f)
	if err != nil {
		fmt.Printf("**warning: could not write history file [%s]\n", fname)
		return
	}

	err = term.Close()
	if err != nil {
		fmt.Printf("**warning: problem closing term: %v\n", err)
		return
	}
}

func main() {
	defer atexit()

	flag.Parse()
	w := eval.NewWorld()
	if *filename != "" {
		data, err := ioutil.ReadFile(*filename)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		file, err := parser.ParseFile(fset, *filename, data, 0)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		files := []*ast.File{file}
		code, err := w.CompilePackage(fset, files, "main")
		if err != nil {
			if list, ok := err.(scanner.ErrorList); ok {
				for _, e := range list {
					fmt.Println(e.Error())
				}
			} else {
				fmt.Println(err.Error())
			}
			os.Exit(1)
		}
		code, err = w.Compile(fset, "main()")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		_, err = code.Run()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	var ierr error = nil // previous interpreter error
	ps1 := "igo> "
	ps2 := "...  "
	prompt := &ps1
	codelet := ""
	// initialize main package
	{
		codelet := "package main\n"
		f, err := parser.ParseFile(fset, "input", codelet, 0)
		code, err := w.CompilePackage(fset, []*ast.File{f}, "main")
		if err == nil {
			code.Run()
		}
	}

	for {
		line, err := term.Prompt(*prompt)
		if err != nil {
			if err != io.EOF {
				ierr = err
			} else {
				ierr = nil
			}
			break //os.Exit(0)
		}
		if line == "" || line == ";" {
			// no more input
			prompt = &ps1
		}

		codelet += line
		if codelet != "" {
			for _, ll := range strings.Split(codelet, "\n") {
				term.AppendHistory(ll)
			}
		}
		code, err := w.Compile(fset, codelet+";")
		if err != nil {
			if ierr != nil && prompt == &ps1 {
				fmt.Println(err.Error())
				fmt.Printf("(error %T)\n", err)
				// reset state
				codelet = ""
				ierr = nil
				prompt = &ps1
				continue
			}
			// maybe multi-line command ?
			prompt = &ps2
			ierr = err
			codelet += "\n"
			continue
		}
		v, err := code.Run()
		if err != nil {
			fmt.Println(err.Error())
			fmt.Printf("(error %T)\n", err)
			codelet = ""
			continue
		}
		if v != nil {
			fmt.Println(v.String())
		}
		//	resetstate:
		// reset state
		codelet = ""
		ierr = nil
	}
}
