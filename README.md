igo
===

A simple interactive Go interpreter built on [go-eval](github.com/sbinet/go-eval) with some readline-like refinements


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
