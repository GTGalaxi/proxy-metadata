package main

import (
	"fmt"
	"os"
	"proxy-browser/metadata"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("\033[31m[!]\033[0m Expected platform (aws, gcp or azure).")
		os.Exit(1)
	}
	metadata.SelectPlatform(os.Args[1])
	metadata.SelectedPlatform.Enumerate(metadata.BaseURL, metadata.AllVersions, metadata.PassthroughEncoder, metadata.OutPath)
}
