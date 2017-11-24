package duckduckgo

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/aurieh/ddg-ng/utils"
	"net/http"
	"net/url"
)

const DDGURL = "https://api.duckduckgo.com/?q=%s&o=json&kp=%d&no_redirect=1&no_html=1&d=%d&t=%s"

type Icon struct {
	URL    string
	Height interface{} `json:"Height,omitempty"`
	Width  interface{} `json:"Width,omitempty"`
}

type Topic struct {
	Result   string
	Icon     Icon
	FirstURL string
	Text     string
}

type Answer struct {
	Abstract       string
	AbstractText   string
	AbstractSource string
	AbstractURL    string
	Image          string
	Heading        string

	Answer     string
	AnswerType string

	Definition       string
	DefinitionSource string
	DefinitionURL    string

	RelatedTopics []Topic `json:"RelatedTopics,omitempty"`
	Results       []Topic `json:"Results,omitempty"`

	Redirect string
}

func Query(client *http.Client, query string, useragent string, safesearch bool, meanings bool) (*Answer, error) {
	answer := &Answer{}

	qQuery := url.QueryEscape(query)
	var qSafesearch, qMeanings int
	if safesearch {
		qSafesearch = 1
	} else {
		qSafesearch = -1
	}
	if meanings {
		qMeanings = 0
	} else {
		qMeanings = 1
	}

	url := fmt.Sprintf(DDGURL, qQuery, qSafesearch, qMeanings, url.QueryEscape(useragent))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", useragent)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		log.Errorln("duckduckgo answer status: " + res.Status)
		return nil, errors.New("unknown error")
	}
	err = utils.GetJson(res, answer)
	if err != nil {
		return nil, err
	}

	return answer, nil
}
