package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/manticoresoftware/go-sdk/manticore"
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

var ParseMode,
	showConsole bool
var outputPath,
	matchMode,
	columns string

func main() {

	flag.BoolVarP(&ParseMode, "parse", "p", false, "Парсинг страниц с наименованиями книг в файл csv")
	flag.StringVarP(&outputPath, "file", "f", "./list.csv", "csv файл с наименованиями книг для проверки")
	flag.StringVarP(&matchMode, "matchMode", "m", "query_string", "режим поиска, query_string, match_phrase, match")
	flag.StringVarP(&columns, "columns", "c", "all", "колонки csv (author, title), которые будут соединены в строку запроса")
	flag.BoolVarP(&showConsole, "showConsole", "s", false, "вывод результатов в консоль без сохранения в файл")
	flag.Parse()

	if ParseMode {
		parseMode()
		return
	}

	var filename string
	if !showConsole {
		currentTime := time.Now()
		filename = fmt.Sprintf("./%v_result_log.txt", currentTime.Format("15-04-05_02012006"))
		f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()

		log.SetOutput(f)
	}

	// Читаем файл со списком книг
	records := readCsvFile(outputPath)

	// Создаем клиент Manticore
	cl := manticore.NewClient()
	cl.SetServer("localhost", 9308)
	_, err := cl.Open()
	if err != nil {
		fmt.Printf("Conn: %v", err)
	}

	for n, item := range records {
		if n == 0 {
			continue
		}

		var b strings.Builder

		if columns == "author" {
			fmt.Fprintf(&b, "%v ", item[0])
		}
		if columns == "title" {
			fmt.Fprintf(&b, "%v ", item[1])
		}
		if columns != "author" && columns != "title" {
			fmt.Fprintf(&b, "%v ", item[0])
			b.WriteString(item[1])
		}
		query := fmt.Sprintln(b.String())

		var currentMatchMode EMatchMode

		switch matchMode {
		case "query_string":
			currentMatchMode = MatchAll
		case "match_phrase":
			currentMatchMode = MatchPhrase
		case "match":
			currentMatchMode = MatchAny
		default:
			currentMatchMode = MatchAll
		}

		search := NewSearch(query, "minjust_list")
		search.LogMessage = fmt.Sprintf("Строка: %v\r\nЗапрос(%v): %v", item, matchMode, query)
		search.MatchMode = currentMatchMode

		manticoreHttpJson(search)
	}

	if filename != "" {
		fmt.Printf("Обработка завершена. Создан файл: %v\r\n", filename)
	} else {
		fmt.Println("Обработка завершена")
	}
}

func manticoreHttpJson(search Search) {

	if search.Query == "" {
		return
	}

	var query string
	switch search.MatchMode {
	case MatchAll:
		query = fmt.Sprintf("\"query\": {  \"query_string\": \"%v\" },", search.Query)
	case MatchPhrase:
		query = fmt.Sprintf("\"query\": {  \"match_phrase\": {\"*\": \"%v\"} },", search.Query)
	case MatchAny:
		query = fmt.Sprintf("\"query\": {  \"match\": {\"*\": \"%v\"} },", search.Query)
	}
	var jsonData = []byte(fmt.Sprintf(`{
		"index": "%v",
		%v
		"highlight": {"limit": 0}
	}`, search.Index, query))

	resp, err := http.Post("http://localhost:9308/search", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Errorf("%v", err)
	}
	defer resp.Body.Close()
	//body, err := io.ReadAll(resp.Body)

	var searchResult SearchResult
	err = json.NewDecoder(resp.Body).Decode(&searchResult)
	if err != nil {
		panic(err)
	}
	if searchResult.Hits.Total > 0 {

		var b strings.Builder
		fmt.Fprintf(&b, "%v", search.LogMessage)

		for n, hit := range searchResult.Hits.Hits {
			for _, item := range hit.Highlight.Name {
				fmt.Fprintf(&b, "%d. %v\r\n", n+1, item)
			}
		}

		log.Println(b.String())
	}
	//log.Printf("Результат: %v\r\n\r\n", string(body))
}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}
	return records
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
