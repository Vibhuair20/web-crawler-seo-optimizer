package main

import (
	"flag"
	"fmt"
	"html"
	"net/http"
	"os"
	"strings"
)

func analyze(url, baseURL string, visited* map[string]string){
	// parsing the page 
	// and checking for the error as often the last value is error 
	page, err := parse(url)
	if err != nil{
		fmt.Printf("error fetching page%s: %v/n", url, err)
		return
	}
	// getting the page tilte
	title := pageTitle(page)

	// store url and title in the map
	(*visited)[url] = title

	// get all the links in the page
	// using nil becsause we are using recursion here
	links := GetLinks(nil, page)

	//visiting each link from the same page
	for _, link := range link{
		if(*visited)[link] == "" && strings.HasPrefix(link, baseURL){
			analyze(link, baseUrl, visited)
		}
	}
}

// parse fetches and parses an html page
func parse(url string) (*html.Node, error){
	// making a http get request
	resp, err := http.Get(url)
	if err != nil{
		return nil, fmt.Errorf("cannot parse pages")
	}
	// closing the body
	defer resp.Body.Close()

	// parsing the html
	node, err := html.Parse(resp.Body)
	if err != nil{
		return nil, fmt.Errorf("cannot parse pages")
	}
	return node, nil
}

// page title extracts from the html page
func pageTitle(n string) (*html.Node){
	// base case
	if n.Type == html.ElementNode && n.Data == "title"{
		return n.FirstChild.Data
	}

	// recursivley for the children nodes
	for c := n.FirstChild; c != nil; c.NextSibling{
		if title := pageTitle(c); title != ""{
			return title
		}
	}
	return "No Title"
}

// page links extracts all the linsks from the page
func GetLinks(links []string, n *html.Node) []string{
	// if current node is the anchor tag
	if n.Type == html.ElementNode && n.Data == "a"{
		// loo for href attribute
		for _, attr := range n.attr{
			if attr.Key == "href" {
                // Add link if it's not already in the slice
                if !sliceContains(links, attr.Val) {
                    links = append(links, attr.Val)
                }
            }
		}
	}

	//recursove for children nowlinks
	for c := n.FirstChild; c != nil; c.NextSibling{
		links = pageLinks(links, c)
	}
	return links
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

	if url == ""(
		flag.PrintDefaults()
		os.Exit(1)
	)
	// creating a map to see all are visited

	visited := make(map[string]string)

	// crawling process
	analyze(url, url, &visited)

	//print the resuts
	fmt.Println("/n discovered pages")
	for link, title := range visited{
		fmt.Println("%s -> %s\n", link, title)
	}

	// check if duplicates are requested
	if dup{
		checkDuped(&visited)
	}

}
