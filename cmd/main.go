package main

import (
	"CompClubProject/app"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide the input file")
		return
	}

	filename := os.Args[1]

	app.Run(filename)
}
