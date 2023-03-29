package main

import (
	"fmt"
	"proxy-browser/metadata"
)

func main() {
	baseURL := "http://example.com/proxy?{0}"
	fmt.Println(metadata.EnumerateAWS(baseURL, metadata.PassthroughEncoder, false))
}
