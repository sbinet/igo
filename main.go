// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"exp/eval"
	"flag"
	"fmt"
	"go/parser"
	"go/scanner"
	"io/ioutil"
	"os"
	"readline"
	//"bitbucket.org/binet/go-readline" /* cgo packages not goinstall-able*/
	"path"
)

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
			println(err.String())
			os.Exit(1)
		}
		file, err := parser.ParseFile(*filename, data, 0)
		if err != nil {
			println(err.String())
			os.Exit(1)
		}
		code, err := w.CompileDeclList(file.Decls)
		if err != nil {
			if list, ok := err.(scanner.ErrorList); ok {
				for _, e := range list {
					println(e.String())
				}
			} else {
				println(err.String())
			}
			os.Exit(2)
		}
		_, err = code.Run()
		if err != nil {
			println(err.String())
			os.Exit(3)
		}
		code, err = w.Compile("init()")
		if code != nil {
			_, err := code.Run()
			if err != nil {
				println(err.String())
				os.Exit(4)
			}
		}
		code, err = w.Compile("main()")
		if err != nil {
			println(err.String())
			os.Exit(5)
		}
		_, err = code.Run()
		if err != nil {
			println(err.String())
			os.Exit(6)
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
		code, err := w.Compile(*line+";")
		if err != nil {
			println(err.String())
			continue
		}
		v,err := code.Run()
		if err != nil {
			println(err.String())
			continue
		}
		if v!= nil {
			println(v.String())
		}
	}
}
