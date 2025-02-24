package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"

	"golang.org/x/net/html"
)

func analyze(url, baseURL string, visited *map[string]string) {
	// parsing the page
	// and checking for the error as often the last value is error
	page, err := parse(url)
	if err != nil {
		fmt.Printf("error fetching page %s: %v\n", url, err) // Fixed newline character
		return
	}
	// getting the page title
	title := pageTitle(page)

	// store url and title in the map
	(*visited)[url] = title

	// get all the links in the page
	// using nil becsause we are using recursion here
	links := GetLinks(nil, page)

	//visiting each link from the same page
	for _, link := range links {
		if (*visited)[link] == "" && strings.HasPrefix(link, baseURL) {
			analyze(link, baseURL, visited)
		}
	}
}

// parse fetches and parses an html page
func parse(url string) (*html.Node, error) {
	// making a http get request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch page: %v", err) // Fixed error message
	}
	// closing the body
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: received status code %d", resp.StatusCode)
	}

	node, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot parse pages")
	}
	return node, nil
}

// page title extracts from the html page
func pageTitle(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "title" {
		if n.FirstChild != nil { // Added nil check
			return n.FirstChild.Data
		}
	}

	// recursivley for the children nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling { // Fixed missing assignment
		if title := pageTitle(c); title != "" {
			return title
		}
	}
	return "No Title"
}

// page links extracts all the linsks from the page
func GetLinks(links []string, n *html.Node) []string {
	// if current node is the anchor tag
	if n.Type == html.ElementNode && n.Data == "a" {
		// loo for href attribute
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				// Add link if it's not already in the slice -> to avoid dupliacte
				if !sliceContains(links, attr.Val) {
					links = append(links, attr.Val)
				}
			}
		}
	}

	//recursove for children nowlinks
	for c := n.FirstChild; c != nil; c = c.NextSibling { // Fixed missing assignment
		links = GetLinks(links, c) // Fixed incorrect function name
	}
	return links
}

// _ means it discards the index
func sliceContains(sortedSlice []string, value string) bool {
	sort.Strings(sortedSlice)

	index := sort.SearchStrings(sortedSlice, value)
	return index < len(sortedSlice) && sortedSlice[index] == value
}

// check the duplicates finds and prints pages with the duplicate titles
func checkDuped(visited *map[string]string) {
	found := false
	uniqueTitles := make(map[string]string)

	fmt.Println("\nChecking for Duplicate Titles:") // Fixed incorrect newline

	for link, title := range *visited {
		if firstAppearance, exists := uniqueTitles[title]; exists {
			found = true
			fmt.Printf("Duplicate: \"%s\" found in %s (first in %s)\n",
				title, link, firstAppearance)
		} else {
			uniqueTitles[title] = link
		}
	}

	if !found {
		fmt.Println("No duplicate titles found!")
	}
}

// this is the main function
func main() {
	// command line flags
	// store url to crwal
	var url string
	//will store weather to check for updates or not
	var duplicate bool

	flag.StringVar(&url, "url", "", "URL to start crawling")                    // string flag to a variable
	flag.BoolVar(&duplicate, "duplicate", false, "Check for duplicate content") // boll flag to a variable

	flag.Parse()

	if url == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	// creating a map to see all are visited

	visited := make(map[string]string)

	// crawling process
	analyze(url, url, &visited)

	fmt.Println("\nDiscovered pages:")
	for link, title := range visited {
		if title == "No Title" {
			fmt.Printf("%s -> No Title\n", link)
		} else {
			fmt.Printf("%s -> %s\n", link, title)
		}
	}

	// check if duplicates are requested
	if duplicate {
		checkDuped(&visited)
	}
}
