package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type rssFeed struct {
	XMLName xml.Name  `xml:"rss"`
	Items   []rssItem `xml:"channel>item"`
}

type rssItem struct {
	XMLName     xml.Name   `xml:"item"`
	Title       string     `xml:"title"`
	PostName    string     `xml:"post_name"`
	Link        string     `xml:"link"`
	Description string     `xml:"description"`
	PostType    string     `xml:"post_type"`
	Meta        []postMeta `xml:"postmeta"`
	Status      string     `xml:"status"`
	PubDate     string     `xml:"pubDate"`
	Category    []category `xml:"category"`
}

type category struct {
	Domain   string `xml:"domain,attr"`
	Nicename string `xml:"nicename,attr"`
	Value    string `xml:",chardata"`
}

type postMeta struct {
	MetaKey   string `xml:"meta_key"`
	MetaValue string `xml:"meta_value"`
}

type wpBlogPost struct {
	Title       string   `yaml:"title"`
	RedirectURL string   `yaml:"redirectURL"`
	PublishedOn string   `yaml:"date"`
	Tags        []string `yaml:"tags"`
	Category    []string `yaml:"categories"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var archiveLocation = flag.String("l", "archives", "location of archives")
	var file = flag.String("f", "mywp.xml", "xml file")
	flag.Parse()

	xmlFile, err := os.Open(*file)
	if err != nil {
		fmt.Println(err)
	}

	defer xmlFile.Close()

	var rssfeed rssFeed

	byteValue, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		fmt.Printf("Error reading file. %s\n", err)
	}
	xml.Unmarshal(byteValue, &rssfeed)

	if err := os.Mkdir(*archiveLocation, 0777); err != nil {
		fmt.Printf("Error creating archive directory. %v", err)
	}

	for _, item := range rssfeed.Items {
		if isPost(item.PostType) &&
			isPublished(item.Status) && !isReBlog(item.Meta) {
			wppost := *archiveLocation + "/" + item.PostName + ".md"
			if _, err := os.Stat(wppost); os.IsNotExist(err) {
				f, err := os.Create(wppost)
				check(err)
				defer f.Close()
			}

			var tagList []string
			var categoryList []string

			for _, cat := range item.Category {
				if cat.Domain == "post_tag" {
					if isUnique(categoryList, cat.Value) {
						tagList = append(tagList, cat.Value)
					}
				}
				if cat.Domain == "category" {
					if isUnique(categoryList, cat.Value) {
						categoryList = append(categoryList, cat.Value)
					}
				}
			}

			t, _ := time.Parse("Mon, 2 Jan 2006 15:04:05 +0000", item.PubDate)
			pubdate := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d-00:00",
				t.Year(), t.Month(), t.Day(),
				t.Hour(), t.Minute(), t.Second())

			post := wpBlogPost{item.Title, item.Link, pubdate, tagList, categoryList}

			ypost, err := yaml.Marshal(post)
			hyphen := []byte("---\n")
			content := append(append(hyphen, ypost...), hyphen...)
			err = ioutil.WriteFile(wppost, content, 0777)
			check(err)
		}
	}
}

func isUnique(l []string, key string) bool {
	for _, v := range l {
		if strings.ToUpper(v) == strings.ToUpper(key) {
			return false
		}
	}
	return true
}

func isPublished(status string) bool {
	if status == "publish" {
		return true
	}
	return false
}

func isPost(postType string) bool {
	if postType == "post" {
		return true
	}
	return false
}

func isReBlog(meta []postMeta) bool {
	if stringInArray("is_reblog", meta) {
		return true
	}
	return false
}

func stringInArray(key string, meta []postMeta) bool {
	for _, k := range meta {
		if k.MetaKey == key {
			return true
		}
	}
	return false
}
