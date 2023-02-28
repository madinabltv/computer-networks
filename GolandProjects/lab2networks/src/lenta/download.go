package main

import (
	"github.com/mgutz/logxi/v1"
	"golang.org/x/net/html"
	"net/http"
)

//const URL = "https://kruzhok.org/news/"

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func getChildren(node *html.Node) []*html.Node {
	var children []*html.Node
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		children = append(children, c)
	}
	return children
}

func isElem(node *html.Node, tag string) bool {
	return node != nil && node.Type == html.ElementNode && node.Data == tag
}

func isText(node *html.Node) bool {
	return node != nil && node.Type == html.TextNode
}

func isDiv(node *html.Node, class string) bool {
	return isElem(node, "div") && getAttr(node, "class") == class
}

type Item struct {
	Ref, Title string
}

func readItem(item *html.Node) *Item {
	children := getChildren(item)
	for _, child := range children {
		if isElem(child, "a") && getAttr(child, "class") == "n-news_title" {
			return &Item{
				Ref:   getAttr(child, "href"),
				Title: getChildren(child)[0].Data,
			}
		}
	}
	return nil
}

func search(node *html.Node) []*Item {
	if isDiv(node, "n-news_list") {
		var items []*Item
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if isDiv(c, "n-news_item") {
				if item := readItem(c); item != nil {
					items = append(items, item)
				}
			}
		}
		return items
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if items := search(c); items != nil {
			return items
		}
	}
	return nil
}

func downloadNews() []*Item {
	log.Info("sending request to kruzhok.org")
	if response, err := http.Get("https://kruzhok.org/news"); err != nil {
		log.Error("request to kruzhok.org failed", "error", err)
	} else {
		defer response.Body.Close()
		status := response.StatusCode
		log.Info("got response from kruzhok.org", "status", status)
		if status == http.StatusOK {
			if doc, err := html.Parse(response.Body); err != nil {
				log.Error("invalid HTML from kruzhok.org", "error", err)
			} else {
				log.Info("HTML from kruzhok.org parsed successfully")
				return search(doc)
			}
		}
	}
	return nil
}
