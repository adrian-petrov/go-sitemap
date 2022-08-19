package htmlparser

import (
	"io"
	"io/ioutil"
	"log"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func ParseLinks(reader io.Reader) ([]Link, error) {
	node, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}

	var result []Link
	nodes := getLinkNodes(node)
	for _, node := range nodes {
		result = append(result, buildLink(node))
	}

	return result, nil
}

func getLinkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var result []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result = append(result, getLinkNodes(c)...)
	}
	return result
}

func buildLink(n *html.Node) Link {
	var result Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			result.Href = attr.Val
			break
		}
	}
	result.Text = getTextFromNode(n)
	return result
}

func getTextFromNode(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var result string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result += getTextFromNode(c)
	}

	return strings.Join(strings.Fields(result), " ")
}

func ReadHTMLFile(file string) []byte {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return content
}
