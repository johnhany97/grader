# Automatic Programming Assignments Grader

```
Usage of main.go:
  -schema string
        Marking scheme to follow when grading the assignment (required)
```

## Installation
First off, you'll need to [Install Go](https://golang.org/doc/install) and make sure enviornment variables are set up properly. Following that, you'll need to [Install Docker](https://docs.docker.com/v17.09/engine/installation/).
Running the following two commands install two Go-based projects that run the grading.
```
go get github.com/docker-exec/dexec
go get github.com/johnhany97/grader
```

## Supported languages

Following is a list of supported languages with the required docker image.

- `C`: `dexec/lang-c`
- `Clojure`: `dexec/lang-clojure`
- `CoffeeScript`: `dexec/lang-coffee`
- `C++`: `dexec/lang-cpp`
- `C#`: `dexec/lang-csharp`
- `D`: `dexec/lang-d`
- `Erlang`: `dexec/lang-erlang`
- `F#`: `dexec/lang-fsharp`
- `Go`: `dexec/lang-go`
- `Groovy`: `dexec/lang-groovy`
- `Haskell`: `dexec/lang-haskell`
- `Java`:  `dexec/lang-java` AND `johnhany97/grader-junit`
- `Lisp`:  `dexec/lang-lisp`
- `Lua`: `dexec/lang-lua`
- `JavaScript`: `dexec/lang-node`
- `Nim`: `dexec/lang-nim`
- `Objective C`: `dexec/lang-objc`
- `OCaml`: `dexec/lang-ocaml`
- `Perl 6`: `dexec/lang-perl6`
- `Perl`: `dexec/lang-perl`
- `PHP`: `dexec/lang-php`
- `Python`: `dexec/lang-python`
- `R`: `dexec/lang-r`
- `Racket`: `dexec/lang-racket`
- `Ruby`: `dexec/lang-ruby`
- `Rust`: `dexec/lang-rust`
- `Scala`: `dexec/lang-scala`
- `Bash`: `dexec/lang-bash`

## Example usage
The following commands are to be ran from the project's root directory.
```
grader -schema="examples/inputOutput/javaSchema.json"
```
```
grader -schema="examples/python/pySchema.json"
```
```
grader -schema="examples/inputOutput/goSchema.json"
```
