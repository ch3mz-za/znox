package main

import (
	"log"
	"os"

	"github.com/ch3mz-za/znox/app"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	var srcPath, dstPath string
	if len(os.Args) > 1 {
		srcPath = os.Args[1]
	}
	if len(os.Args) > 2 {
		srcPath = os.Args[2]
	}

	p := tea.NewProgram(app.New(srcPath, dstPath), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}
