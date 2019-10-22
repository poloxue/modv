This is a module dependency visualizer for go mod.

## Linux 

```bash
$ sudo apt-get install graphviz
$ go install github.com/poloxue/modv
$ go mod graph | modv | dot -T png | display
```

## MacOS

```bash
$ brew install graphviz
$ go install github.com/poloxue/modv
$ go mod graph | modv | dot -T png | open -f -a /Applications/Preview.app
```

## Windows

```bash
$ choco install graphviz.portable
$ go install github.com/poloxue/modv
$ go mod graph | modv | dot -T png -o graph.png; start graph.png
```
