package main

import (
	"proxy-browser/metadata"
)

func main() {
	baseURL := "http://example.com/proxy?{0}"
	metadata.Initialise()
	metadata.EnumerateAWS(baseURL, metadata.PassthroughEncoder).ToFile(metadata.DefaultOutputPath.AWS)
}
