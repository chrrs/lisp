# clisp
A simple interpreted lisp language written in Go.

I decided I wanted to learn Go, and what better way to do that than to write my own lisp-like programming language.

## Building and running

### Building

```bash
$ go build lisp
```

An executable file for your operating system should appear in the project directory. Running this file will start the REPL.

### Running without building

```bash
$ go run lisp
```

This runs the REPL on your computer directly without compiling.

> Note: You probably want to use the standard library, so you should import that using `import "lib/std"`

## Syntax

Every call in clisp follows the `[func] [args...]` pattern. For example:

```
> + 5 2
7
```

For nested evaluation, you can create an S-Expression using `()`:

```
> + 5 (* 4 2)
13
```

While S-Expressions will always get evaluated, Q-Expressions (`{}`) will be left as-is. This creates lists!

```
> {+ 5 2}
{+ 5 2}
```

You can still evaluate these expressions using the `eval` function though:

```
> eval {+ 5 2}
7
```

Combine these with some of the built-in functions...

```
> list 1 2 3
{1 2 3}
> head {1 2 3}
{1}
> tail {1 2 3}
{2 3}
> join {1 2} {3 4}
{1 2 3 4}
```

and variables of course...

```
> def {x} 5
()
> + x 2
7
```

not forgetting lambdas...

```
> fn {a b} {+ a (a b)}
(fn [a b] {+ a (a b)})
> (fn {a b} {+ a (a b)}) 1 2
4
> (fn {a b & c} {list a b c}) 1 2 3 4
{1 2 {3 4}}
```

we can create some pretty cool stuff!

(this is an incomplete list of functionality)

Don't forget to check out the standard library at `lib/std.clsp` for some common useful functions.

_I won't be creating a comprehensive documentation for this language and its standard library, since it is nothing more than a hobby / learning project and not intended to be used by anyone._

## Examples

### Fibonacci

```
(fun {fibonacci n} {
    select
        {(= n 0) 0}
        {(= n 1) 1}
        {else (+ (fibonacci (- n 1)) (fibonacci (- n 2)))}
})
```

```
> map fibonacci (range 0 20)
{0 1 1 2 3 5 8 13 21 34 55 89 144 233 377 610 987 1597 2584 4181 6765}
```

### Factorial
```
(fun {factorial n} {
    if (= n 0)
        {1}
        {* n (factorial (- n 1))}
})
```

```
> map factorial (range 0 10)
{1 1 2 6 24 120 720 5040 40320 362880 3.6288e+06}
```

## Credits

- Syntax inspired by [lispy](http://www.buildyourownlisp.com/).
- A few small implementation details inspired by [go-lisp](https://github.com/janne/go-lisp/).
