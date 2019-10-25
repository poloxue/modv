This is a module dependency visualizer for go mod.

# Usage

Modv's usage is different in different systems.

## Linux

Install `graphviz`. For Ubuntu/Debian
```bash
$ sudo apt-get install graphviz
```

For ArchLinux

```
$ sudo pacman -S --needed graphviz
```

Install `modv` and use it.

```
$ go install github.com/poloxue/modv
$ go mod graph | modv | dot -T svg -o /tmp/modv.svg && xdg-open /tmp/modv.svg
```


## MacOS

```bash
$ brew install graphviz
$ go get github.com/poloxue/modv
```

Try the following.

```
$ go mod graph | modv | dot -T png | open -f -a /Applications/Preview.app
```

If error accured, for eaxmple，`FSPathMakeRef(/Applications/Preview.app) failed with error -43.`，try the command:

```
$ go mod graph | modv | dot -T png | open -f -a /System/Applications/Preview.app
```

## Windows

First, install `graphviz`:

```bash
$ choco install graphviz.portable
```
For [MSYS2](https://www.msys2.org/)

```bash
$ pacman -S mingw-w64-x86_64-graphviz
```

Try it.

```bash
$ go get github.com/poloxue/modv
$ go mod graph | modv | dot -T svg -o graph.svg; start graph.svg
```

# Demo

If MacOS, tye the following:

```bash
$ git clone https://github.com/poloxue/testmod
$ cd testmod
$ go mod graph | modv | dot -T png | open -f -a /System/Applications/Preview.app
```

Output:

![](http://blogimg.poloxue.com/0014-go-mod-graph-visible-04.png)
