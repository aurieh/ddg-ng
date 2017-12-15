package htmlmeta

import (
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"net/http"
	"strings"
)

// New create new MetaParser
func New(client *http.Client, url string) (*MetaParser, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	return &MetaParser{
		Root: root,
	}, nil
}

// MetaParser html meta info parser
type MetaParser struct {
	Root *html.Node
}

// GetTitle get website head title
func (p *MetaParser) GetTitle() string {
	matcher := func(n *html.Node) bool {
		return n.DataAtom == atom.Title
	}
	titles := scrape.FindAll(p.Root, matcher)
	if len(titles) == 0 {
		return ""
	}
	return scrape.Text(titles[0])
}

// GetMeta get head meta tag
func (p *MetaParser) GetMeta(name string, attr string) string {
	matcher := func(n *html.Node) bool {
		return n.DataAtom == atom.Meta && scrape.Attr(n, "name") == name
	}
	tags := scrape.FindAll(p.Root, matcher)
	if len(tags) == 0 {
		return ""
	}
	return scrape.Attr(tags[0], attr)
}

// GetOGPMeta get ogp meta tag
func (p *MetaParser) GetOGPMeta(prop string) string {
	if !strings.HasPrefix(prop, "og:") {
		prop = "og:" + prop
	}
	matcher := func(n *html.Node) bool {
		return n.DataAtom == atom.Meta && scrape.Attr(n, "property") == prop
	}
	tags := scrape.FindAll(p.Root, matcher)
	if len(tags) == 0 {
		return ""
	}
	return scrape.Attr(tags[0], "content")
}
