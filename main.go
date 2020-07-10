package main

import (
	"os"

	"github.com/remiflavien1/recon-archy/selenium"
)

func main() {

	comp := os.Args[1]
	selenium.Start(comp)
}
