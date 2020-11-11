// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 241.

// Crawl2 crawls web links starting with the command-line arguments.
//
// This version uses a buffered channel as a counting semaphore
// to limit the number of concurrent calls to links.Extract.
//
// Crawl3 adds support for depth limiting.
//
package main

import (
	"bufio"
	"log"
	"os"
	"strconv"

	"gopl.io/ch5/links"
)

//!+sema
// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent requests.

var tokens = make(chan struct{}, 20)

type url struct {
	str   string
	depth int
}

func crawl(site url, fw *bufio.Writer) []url {
	if site.depth < 1 {
		urls := make([]url, 0)
		return urls
	}

	fw.WriteString(site.str + "\n")
	tokens <- struct{}{} // acquire a token
	list, err := links.Extract(site.str)
	urls := make([]url, 0)
	for _, n := range list {
		urls = append(urls, url{str: n, depth: site.depth - 1})
	}
	<-tokens // release the token

	if err != nil {
		log.Print(err)
	}
	return urls
}

//!-sema

//!+

func main() {
	if len(os.Args) < 3 {
		panic("Not enough arguments introduced")
	}

	worklist := make(chan []url)
	var n int // number of pending sends to worklist

	// Start with the command-line arguments.
	n++
	depth, err := strconv.Atoi(os.Args[1][7:])
	filename := os.Args[2][9:]
	f, err := os.Create(filename)
	defer f.Close()
	fw := bufio.NewWriter(f)
	if err != nil {
		panic("Invalid Arguments")
	}
	urls := make([]url, 0)
	urls = append(urls, url{str: os.Args[3], depth: depth})
	go func() { worklist <- urls }()
	// Crawl the web concurrently.
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link.str] {
				seen[link.str] = true
				n++
				go func(link url) {
					worklist <- crawl(link, fw)
				}(link)
			}
		}
	}
	fw.Flush()
}

//!-
