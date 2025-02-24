package main

import (
	"flag"
	"fmt"
	"os"
)

func analyze(url, baseURL string, visited* map[string]string){
	// parsing the page 
	// and checking for the error as often the last value is error 
	page, err := parse(url)
	if err != nil{
		fmt.Printf("error fetching page%s: %v/n", url, err)
		return
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
