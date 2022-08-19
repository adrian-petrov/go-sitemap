package htmlparser

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseLinks(t *testing.T) {
	t.Run("it returns one link tag", func(t *testing.T) {
		htmlFile := ReadHTMLFile("../html/ex1.html")
		htmlReader := strings.NewReader(string(htmlFile))

		got, err := ParseLinks(htmlReader)
		if err != nil {
			t.Fatal("Should parse html file:", err)
		}
		want := []Link{{Href: "/other-page", Text: "A link to another page"}}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %#v, but want %#v", got, want)
		}
	})

	t.Run("it returns two link tags without nested html", func(t *testing.T) {
		htmlFile := ReadHTMLFile("../html/ex2.html")
		htmlReader := strings.NewReader(string(htmlFile))

		got, err := ParseLinks(htmlReader)
		if err != nil {
			t.Fatal("Should parse html file:", err)
		}
		want := []Link{
			{Href: "https://www.twitter.com/joncalhoun", Text: "Check me out on twitter"},
			{Href: "https://github.com/gophercises", Text: "Gophercises is on Github!"},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, but want %v", got, want)
		}
	})

	t.Run("it returns one link tag from html with comments", func(t *testing.T) {
		htmlFile := ReadHTMLFile("../html/ex3.html")
		htmlReader := strings.NewReader(string(htmlFile))

		got, err := ParseLinks(htmlReader)
		if err != nil {
			t.Fatal("Should parse html file:", err)
		}
		want := []Link{
			{Href: "#", Text: "Login"},
			{Href: "/lost", Text: "Lost? Need help?"},
			{Href: "https://twitter.com/marcusolsson", Text: "@marcusolsson"},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %#v\n, but want %#v", got, want)
		}
	})
}
