igo
===

A simple interactive Go interpreter built on [go-eval](https://github.com/sbinet/go-eval) with some readline-like refinements

## Setting up

```sh

$ git clone git@github.com:sbinet/igo.git
$ cd igo
$ go install
$ go build main.go
$ ./main

********************************
** Interactive Go interpreter **
********************************


igo>

```

You can add compiled main file to your bash file for easy access from anywhere:

``` sh

$ echo 'alias igo="$PWD/main"' >> ~/.bash_profile
$ source  ~/.bash_profile
$ igo

********************************
** Interactive Go interpreter **
********************************

igo> println("Hello, World!")
Hello, World!
igo>

```



## Example

```sh
$ igo
igo> func f() { println("hello world") }
igo> f()
hello world
igo> type Foo struct {
...   A int
...  }
...
igo> foo := Foo{A:32}
igo> foo
{32}
igo> foo.A
32
```

## Documentation

http://godoc.org/github.com/sbinet/igo


## TODO

- implement code completion

  - with gocode ?

- code colorization ?

- see TODOs of [go-eval](https://github.com/sbinet/go-eval)
