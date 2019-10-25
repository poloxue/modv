package main

import (
	"fmt"
	"os"
	"runtime"
)

func PrintUsage() {
	fmt.Printf("\nUsages:\n\n")
	switch runtime.GOOS {
	case "darwin":
		fmt.Printf("\tgo mod graph | modv | dot -T svg | open -f -a /System/Applications/Preview.app")
	case "linux":
		fmt.Printf("\tgo mod graph | modv | dot -T svg -o /tmp/modv.svg | xdg-open /tmp/modv.svg")
	case "windows":
		fmt.Printf("\tgo mod graph | modv | dot -T png -o graph.png; start graph.png")
	}

	fmt.Printf("\n\n")
}

func main() {
	info, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println("os.Stdin.Stat:", err)
		PrintUsage()
		os.Exit(1)
	}

	if info.Mode()&os.ModeNamedPipe == 0 {
		fmt.Println("command err: command is intended to work with pipes.")
		PrintUsage()
		os.Exit(1)
	}

	mg := NewModuleGraph(os.Stdin)
	if err := mg.Parse(); err != nil {
		fmt.Println("mg.Parse: ", err)
		PrintUsage()
		os.Exit(1)
	}

	if err := mg.Render(os.Stdout); err != nil {
		fmt.Println("mg.Render: ", err)
		PrintUsage()
		os.Exit(1)
	}
}
