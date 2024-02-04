package main

import (
	//"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Missing filename argument")
		return
	}

	for _, filename := range os.Args[1:] {

		basename := filepath.Base(filename)

		file, err := os.Open(filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		// Parse the HTML
		doc, err := html.Parse(file)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Print an anchor tag for the filename
		//fmt.Printf("<div id=\"%s\"></div>", basename)

		var traverse func(*html.Node)
		traverse = func(n *html.Node) {
			// Passthrough the chapter heading
			// if n.Type == html.ElementNode && n.Data == "h1" {
			// 	for _, attr := range n.Attr {
			// 		if attr.Key == "class" {
			// 			classes := strings.Split(attr.Val, " ")
			// 			for _, class := range classes {
			// 				if class == "heading-chapter-title" {
			// 					var buf strings.Builder
			// 					html.Render(&buf, n)
			// 					fmt.Println(buf.String())
			// 					break
			// 				}
			// 			}
			// 		}
			// 	}
			// }

			// Handle sections
			if n.Type == html.ElementNode && n.Data == "section" {
				// // Set the id attribute to basename
				// n.Attr = append(n.Attr, html.Attribute{
				// 	Key: "id",
				// 	Val: basename,
				// })

				var traverseSection func(*html.Node)
				traverseSection = func(n *html.Node) {
					sawFirstH1 := false

					// Rewrite heading and anchor IDs to be unique
					if n.Type == html.ElementNode && (n.Data == "a" || n.Data == "h1" || n.Data == "h2" || n.Data == "h3" || n.Data == "h4" || n.Data == "h5" || n.Data == "h6") {
						// Set the ID of the first <h1> in the section to the file basename
						if !sawFirstH1 && n.Data == "h1" {
							idExists := false
							for i, attr := range n.Attr {
								if attr.Key == "id" {
									// If the id attribute exists, modify its value
									n.Attr[i].Val = basename
									idExists = true
									break
								}
							}
							if !idExists {
								// If the id attribute doesn't exist, append a new one
								n.Attr = append(n.Attr, html.Attribute{
									Key: "id",
									Val: basename,
								})
							}
							sawFirstH1 = true
						} else {
							// Otherwise, prepend the basename to the ID
							for i, attr := range n.Attr {
								if attr.Key == "id" {
									n.Attr[i].Val = basename + "-" + attr.Val
									break
								}
							}
						}
					}

					// Rewrite links to be relative to the current file
					if n.Type == html.ElementNode && n.Data == "a" {
						for i, attr := range n.Attr {
							if attr.Key == "href" {
								link := attr.Val
								if !strings.HasPrefix(link, "http") {
									link = strings.TrimPrefix(link, "/")

									if strings.Contains(link, "#") {
										parts := strings.Split(link, "#")
										linkBasename := parts[0]
										if linkBasename == "" {
											linkBasename = basename
										}
										linkAnchor := parts[1]
										link = linkBasename + "-" + linkAnchor
									}
									link = "#" + link
								}
								n.Attr[i].Val = link
							}
						}
					}

					for c := n.FirstChild; c != nil; c = c.NextSibling {
						traverseSection(c)
					}
				}

				traverseSection(n)

				var buf strings.Builder
				html.Render(&buf, n)
				fmt.Println(buf.String())
			}

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				traverse(c)
			}
		}

		traverse(doc)

		file.Close()
	}
}
