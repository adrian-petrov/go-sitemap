package sitemap

import (
	"bytes"
	"container/list"
	"encoding/xml"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/adrian-petrov/go-sitemap/htmlparser"
)

type Sitemap struct {
	baseURL string
	depth   int
}

func (s *Sitemap) Build() {
	url, depth := s.readFlags()
	s.SetBaseURL(url)
	s.SetDepth(depth)

	sitemap := s.bfs(url)
	s.formatXML(&sitemap)
}

func (s *Sitemap) SetBaseURL(url string) {
	idx := 0
	slashes := 0
	for _, r := range url {
		if slashes == 3 {
			s.baseURL = url[:idx-1]
			return
		}
		if string(r) == "/" {
			slashes++
		}
		idx++
	}
	s.baseURL = url
}

func (s *Sitemap) SetDepth(depth int) {
	s.depth = depth
}

func (s *Sitemap) formatXML(sitemap *map[string]struct{}) {
	type XMLLink struct {
		XMLName xml.Name `xml:"url"`
		Loc     string   `xml:"loc"`
	}
	const (
		Header      = `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
		URLSetOpen  = `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n"
		URLSetClose = `</urlset>`
	)

	f, err := os.Create("data.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.WriteString(Header)
	f.WriteString(URLSetOpen)
	for key := range *sitemap {
		link := &XMLLink{Loc: key}
		out, _ := xml.MarshalIndent(link, "  ", "  ")
		f.Write(out)
		f.WriteString("\n")
	}
	f.WriteString(URLSetClose)
}

func (s *Sitemap) bfs(url string) map[string]struct{} {
	seen := make(map[string]struct{})
	var queue *list.List
	newQueue := list.New()
	// enqueue first level of links
	s.enqueueURLS(newQueue, url)

	for i := 0; i < s.depth; i++ {
		queue, newQueue = newQueue, list.New()
		for queue.Len() > 0 {
			first := queue.Front()
			url := first.Value.(string)
			if _, ok := seen[url]; !ok {
				seen[url] = struct{}{}
				s.enqueueURLS(newQueue, url)
			}
			queue.Remove(first)
		}
	}
	return seen
}

func (s *Sitemap) enqueueURLS(q *list.List, url string) {
	// get the body
	body, err := s.get(url)
	if err != nil {
		log.Fatal(err)
	}
	// create the links
	reader := bytes.NewReader(body)
	links, err := htmlparser.ParseLinks(reader)
	if err != nil {
		log.Fatal(err)
	}
	// normalise the links
	normalised := s.normaliseURLS(links)
	// enqueue the new normalised urls
	for _, v := range normalised {
		q.PushBack(v)
	}
}

func (s *Sitemap) readFlags() (string, int) {
	urlFlag := flag.String("url", "https://adrianpetrov.com", "url to use for sitemap")
	depthFlag := flag.Int("depth", 3, "depth level to search for links")

	flag.Parse()
	return *urlFlag, *depthFlag
}

func (s *Sitemap) get(req string) (body []byte, err error) {
	res, err := http.Get(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *Sitemap) normaliseURLS(links []htmlparser.Link) []string {
	length := len([]rune(s.baseURL))

	var result []string
	for _, v := range links {
		if strings.HasPrefix(v.Href, "/") {
			normalised := s.baseURL + v.Href
			result = append(result, normalised)
			continue
		}
		if len(v.Href) >= len(s.baseURL) &&
			v.Href[:length] == s.baseURL {
			result = append(result, v.Href)
		}

	}
	return result
}
