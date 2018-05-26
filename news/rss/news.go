package rss

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/anovakovic01/tap-backend/news"
	"github.com/mmcdole/gofeed"
	xpath "gopkg.in/xmlpath.v2"
)

var _ news.Collector = (*rssCollector)(nil)

type rssCollector struct {
	parser *gofeed.Parser
}

// NewCollector instantiates new collector instance.
func NewCollector(parser *gofeed.Parser) news.Collector {
	return rssCollector{parser}
}

func (rc rssCollector) Collect(src string) ([]news.News, error) {
	res, err := http.Get(src)
	if err != nil {
		return []news.News{}, err
	}
	defer res.Body.Close()

	data, _ := ioutil.ReadAll(res.Body)

	imageTitles := []string{}
	imageLinks := []string{}

	titlePath := xpath.MustCompile("//image/@title")
	linkPath := xpath.MustCompile("//image/@src")

	root, err := xpath.Parse(bytes.NewReader(data))
	if err != nil {
		return []news.News{}, err
	}

	titleIter := titlePath.Iter(root)
	for titleIter.Next() {
		imageTitles = append(imageTitles, titleIter.Node().String())
	}

	linkIter := linkPath.Iter(root)
	for linkIter.Next() {
		imageLinks = append(imageLinks, linkIter.Node().String())
	}

	feed, err := rc.parser.Parse(bytes.NewReader(data))
	if err != nil {
		return []news.News{}, err
	}

	items := []news.News{}
	for i, item := range feed.Items {
		n := news.News{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			ImageTitle:  imageTitles[i],
			Image:       imageLinks[i],
			PubDate:     *item.PublishedParsed,
		}
		items = append(items, n)
	}

	return items, nil
}
