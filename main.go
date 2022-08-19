package main

import (
	s "github.com/adrian-petrov/go-sitemap/sitemap"
)

func main() {
	sitemap := new(s.Sitemap)
	sitemap.Build()
}
