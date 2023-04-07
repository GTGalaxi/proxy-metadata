# Proxy Browser

A cloud metadata enumeration and scraping tool used in conjunction with an exploited SSRF through a web proxy. It currently supports AWS, with GCP and Azure to be implemented.

Written as a personal tooling project to get more familiar with Go

![proxyaws](https://user-images.githubusercontent.com/10473238/230639960-e085d414-fc5d-467e-9ba4-788d6c0fd06e.gif)

## Installation

### Installing from Source

```plaintext
go get -u github.com/GTGalaxi/proxy-metadata
```

You can then build the package with `go build`, or just run it with `go run main.go [platform] -[args]` (See [below](#usage)).

### Pre-built

You can download a Windows Binary from the [releases](https://github.com/GTGalaxi/proxy-metadata/releases) page.

## Usage

### Command Line Tool

```yaml
Example Usage:

    go run main.go <platform> -proxy <proxy-url> -out <output-path> -all -v

    proxy-browser.exe <platform> -proxy <proxy-url> -out <output-path> -all -v

Sub Commands:

    aws, gcp, azure

Flags:

    -all                Enumerate all versions of the instance metadata. 
                        If not specified, only 'latest/' will be enumerated

    -out string         Path of the output .json file (default "./out/{platform}.json")

    -proxy string       URL of the proxy that will be used for SSRF and enumeration.
                        {0} will be replaced with the proxied URL (default "http://example.com/proxy?{0}")

    -v                  Verbose output
```

### Library

`proxy-browser` can be used as a library and provides easy to use functions and custom encoding functionality

```go
package main

import (
    "proxy-browser/metadata"
)

func main() {
    // Enumerate all versions in AWS through 'https://example.com/proxy?{0}' and save to file './out/aws.json
    metadata.AWS.Enumerate("https://example.com/proxy?{0}", true, metadata.PassthroughEncoder, "./out/aws.json")
}
```

Simply pass through your chosen proxy, a bool to enumerate all versions or just 'latest/', use the built in 'metadata.PassthroughEncoder' or create your own and provide an output path.

`Enumerate()` is a temporary stand in that extends on `Platform` types like AWS, GCP and Azure and calls their internal Enumeration functions like `EnumerateAWS().ToFile()`.

The `Encoder` type can be applied to your own functions to pass them through to the `Enumerate()` function and be used to encode the proxied URL in whatever way you need.

#### Example `Encoder` function

```go
func simpleBase64Encoder(endpoint string) string {
    // Base64 Encoded String
    encodedString := base64.StdEncoding.EncodeToString([]byte(endpoint))
    return encodedString
}

func main() {
    // Implement Custom Encoder
    metadata.AWS.Enumerate("https://example.com/proxy?{0}", true, simpleBase64Encoder, "./out/aws.json")
}
```

## License

This tool is licensed under the MIT License. Feel free to use, modify, and distribute the code as you see fit. If you make any improvements or bugfixes, please consider contributing back to the project by opening a pull request.
