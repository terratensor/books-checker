package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/avast/retry-go"
	flag "github.com/spf13/pflag"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type List []Item

type Item struct {
	Author   string
	Title    string
	BookYear string
}

type Letters []string

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

var ParseMode bool
var outputPath string

func main() {

	flag.BoolVarP(&ParseMode, "parse", "p", false, "Парсинг страниц с наименованиями книг в файл csv")
	flag.StringVarP(&outputPath, "output", "o", "./parsed-files", "путь сохранения файлов")

	if ParseMode {
		parseMode()
	}

}

func parseMode() {
	letters := Letters{"a", "b", "v", "g", "d", "e", "zh", "z", "i", "k", "l", "m", "n", "o", "p", "r", "s", "t", "u", "f", "h", "ts", "ch", "sh", "sch", "ea", "yu", "ya"}
	list := List{}
	for _, letter := range letters {

		url := fmt.Sprintf("http://militera.lib.ru/1/cats/all/%v/index.html", letter)
		doc, err := getTopicBody(url)
		if err != nil {
			log.Printf("%v\n", err)
		}
		log.Printf("parse %v\n", url)

		list = parseDoc(doc, list)
	}
	path := "./csv"
	if err := ensureDir(path); err != nil {
		fmt.Println("Directory creation failed with error: " + err.Error())
		os.Exit(1)
	}

	currentTime := time.Now()
	file := fmt.Sprintf("%v/%v_list.csv", path, currentTime.Format("150405_02012006"))
	writeCSVFile(list, file)

	log.Printf("file %v was successful writing\n", file)
}

func parseDoc(n *html.Node, list List) List {

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && nodeHasRequiredCssClass("item_list", n) {
			// проходим по узлу с атрибутом class block item_list}
			list = append(list, parseItem(n))
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return list
}

func parseItem(node *html.Node) Item {
	//var nAnchor *html.Node

	var bufInnerHtml bytes.Buffer
	w := io.Writer(&bufInnerHtml)

	item := Item{}

	exit := false

	var f func(*html.Node)
	f = func(n *html.Node) {

		if n.Type == html.ElementNode && nodeHasRequiredCssClass("author_list", n) {
			item.Author = getInnerText(n)
			err := html.Render(w, n)
			if err != nil {
				return
			}
			bufInnerHtml.Reset()
		}

		if n.Type == html.ElementNode && nodeHasRequiredCssClass("title_list", n) {
			item.Title = getInnerText(n)
			err := html.Render(w, n)
			if err != nil {
				return
			}
			bufInnerHtml.Reset()
		}

		if n.Type == html.ElementNode && nodeHasRequiredCssClass("book_year", n) {
			item.BookYear = getInnerText(n)
			err := html.Render(w, n)
			if err != nil {
				return
			}
			bufInnerHtml.Reset()
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if exit {
				break
			}
			f(c)
		}
	}
	f(node)

	return item
}

func getTopicBody(url string) (*html.Node, error) {
	// Для того чтобы не следовать автоматическим перенаправлениям,
	// создадим свой экземпляр http.Client с методом проверки CheckRedirect.
	// Это поможет нам возвращать код состояния и адрес до перенаправления.
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}
	//resp, err := client.Get(url)
	resp, err := fetchDataWithRetries(client, url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	return doc, nil
}

// fetchDataWithRetries is your wrapped retrieval.
// It works with a static configuration for the retries,
// but obviously, you can generalize this function further.
func fetchDataWithRetries(client *http.Client, url string) (r *http.Response, err error) {
	retry.Do(
		// The actual function that does "stuff"
		func() error {
			log.Printf("Retrieving data from '%s'", url)
			r, err = client.Get(url)
			return err
		},
		// A function to decide whether you actually want to
		// retry or not. In this case, it would make sense
		// to actually stop retrying, since the host does not exist.
		// Return true if you want to retry, false if not.
		retry.RetryIf(
			func(error) bool {
				log.Printf("Retrieving data: %s", err)
				log.Printf("Deciding whether to retry")
				return true
			}),
		retry.OnRetry(func(try uint, orig error) {
			log.Printf("Retrying to fetch data. Try: %d", try+2)
		}),
		retry.Attempts(3),
		// Basically, we are setting up a delay
		// which randoms between 2 and 4 seconds.
		retry.Delay(4*time.Second),
		retry.MaxJitter(1*time.Second),
	)

	return
}

func writeCSVFile(list List, outputPath string) {
	// Define header row
	headerRow := []string{
		"author", "title", "book_year",
	}

	// Data array to write to CSV
	data := [][]string{
		headerRow,
	}

	for _, item := range list {
		data = append(data, []string{
			strings.TrimSpace(item.Author),
			strings.TrimSpace(item.Title),
			strings.TrimSpace(item.BookYear),
		})
	}

	// Create file
	file, err := os.Create(outputPath)
	checkError("Cannot create file", err)
	defer file.Close()

	// Create writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write rows into file
	for _, value := range data {
		err = writer.Write(value)
		checkError("Cannot write to file", err)
	}
}

func getInnerText(node *html.Node) string {
	for el := node.FirstChild; el != nil; el = el.NextSibling {
		if el.Type == html.TextNode {
			return el.Data
		}
	}
	return ""
}

func getRequiredDataAttr(rda string, n *html.Node) string {
	for _, attr := range n.Attr {
		if attr.Key == rda {
			return attr.Val
		}
	}
	return ""
}

// Перебирает аттрибуты токена в цикле и возвращает bool
// если в html token найден переданный css class
func nodeHasRequiredCssClass(rcc string, n *html.Node) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			classes := strings.Split(attr.Val, " ")
			for _, class := range classes {
				if class == rcc {
					return true
				}
			}
		}
	}
	return false
}

func ensureDir(dirName string) error {
	err := os.Mkdir(dirName, os.ModePerm)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		// check that the existing path is a directory
		info, err := os.Stat(dirName)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return errors.New("path exists but is not a directory")
		}
		return nil
	}
	return err
}
