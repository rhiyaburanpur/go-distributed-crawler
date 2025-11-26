package util

import (
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func ExtractLinks(htmlContent, baseURL string) []string {

	base, err := url.Parse(baseURL)
	if err != nil {
		return nil
	}

	reader := strings.NewReader(htmlContent)
	tokenizer := html.NewTokenizer(reader)
	var links []string

	uniqueLinks := make(map[string]struct{})

	for {
		tt := tokenizer.Next()

		switch tt {
		case html.ErrorToken:

			if tokenizer.Err() == io.EOF {
				return links
			}
			return links

		case html.StartTagToken:
			t := tokenizer.Token()

			if t.Data == "a" {
				for _, a := range t.Attr {
					if a.Key == "href" {
						link := a.Val

						resolvedURL, err := base.Parse(link)
						if err != nil {
							continue
						}

						resolvedURL.Fragment = ""

						absoluteLink := resolvedURL.String()

						if _, ok := uniqueLinks[absoluteLink]; !ok {
							uniqueLinks[absoluteLink] = struct{}{}
							links = append(links, absoluteLink)
						}
						break
					}
				}
			}
		}
	}
}
