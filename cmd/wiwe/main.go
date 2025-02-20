package main

import (
	"fmt"
	"os"

	gemini "github.com/secretval/wiwe/cmd/wiwe/protocols/gemini"
)

func main() {
	prog := os.Args[0]
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Printf("Usage: <%s> <url>\n", prog)
		fmt.Printf("ERROR: Did not specify url\n")
		os.Exit(1)
	}
	url := args[0]

	req, err := gemini.ParseGeminiRequest(url, gemini.PORT)
	if err != true {
		fmt.Print("ERROR: Invalid URL")
		os.Exit(1)
	}

	res := gemini.MakeGeminiQuery(req)
	fmt.Println(res.Body)
}
