// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"os"
	"bitbucket.org/binet/go-readline"
	"exp/eval"
	"path"
)

var fset = token.NewFileSet()
var filename = flag.String("f", "", "file to run")

func init() {
	readline.ParseAndBind("tab: complete")
	readline.ParseAndBind("set show-all-if-ambiguous On")
	fmt.Println(`
********************************
** Interactive Go interpreter **
********************************

`)
	readline.ReadHistoryFile(path.Join(os.Getenv("HOME"),".go.history"))
}

func atexit() {
	readline.WriteHistoryFile(path.Join(os.Getenv("HOME"),".go.history"))
}

func main() {
	defer atexit()

	flag.Parse()
	w := eval.NewWorld()
	if *filename != "" {
		data, err := ioutil.ReadFile(*filename)
		if err != nil {
			fmt.Println(err.String())
			os.Exit(1)
		}
		file, err := parser.ParseFile(fset, *filename, data, 0)
		if err != nil {
			fmt.Println(err.String())
			os.Exit(1)
		}
		code, err := w.CompileDeclList(fset, file.Decls)
		if err != nil {
			if list, ok := err.(scanner.ErrorList); ok {
				for _, e := range list {
					fmt.Println(e.String())
				}
			} else {
				fmt.Println(err.String())
			}
			os.Exit(1)
		}
		_, err = code.Run()
		if err != nil {
			fmt.Println(err.String())
			os.Exit(1)
		}
		code, err = w.Compile(fset, "init()")
		if code != nil {
			_, err := code.Run()
			if err != nil {
				fmt.Println(err.String())
				os.Exit(1)
			}
		}
		code, err = w.Compile(fset, "main()")
		if err != nil {
			fmt.Println(err.String())
			os.Exit(1)
		}
		_, err = code.Run()
		if err != nil {
			fmt.Println(err.String())
			os.Exit(1)
		}
		os.Exit(0)
	}

	prompt := "igo> "
	for {
		line := readline.ReadLine(&prompt)
		if line == nil {
			break; //os.Exit(0)
		}
		readline.AddHistory(*line)
		code, err := w.Compile(fset, *line+";")
		if err != nil {
			fmt.Println(err.String())
			continue
		}
		v,err := code.Run()
		if err != nil {
			fmt.Println(err.String())
			continue
		}
		if v!= nil {
			fmt.Println(v.String())
		}
	}
}
