This is a module dependency visualizer for go mod.

## Linux

```bash
$ # Ubuntu/Debian
$ sudo apt-get install graphviz
$ # ArchLinux
$ sudo pacman -S --needed graphviz
$ go install github.com/poloxue/modv
$ go mod graph | modv | dot -T svg -o /tmp/modv.svg && xdg-open /tmp/modv.svg
```

## MacOS

```bash
$ brew install graphviz
$ go install github.com/poloxue/modv
$ go mod graph | modv | dot -T svg | open -f -a /Applications/Preview.app
```

## Windows

```bash
$ choco install graphviz.portable
$ # for MSYS2 https://www.msys2.org/
$ pacman -S mingw-w64-x86_64-graphviz
$ go install github.com/poloxue/modv
$ go mod graph | modv | dot -T svg -o graph.svg; start graph.svg
```
