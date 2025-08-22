package rss

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// Структура данных публикации в соответсвии с XML-файлом
type XML struct {
	ID              int    //Номер записи
	Title           string `xml:"title" `       //Заголовок
	Description     string `xml:"description" ` //Содержание
	PublicationDate string `xml:"pubDate" `     //Дата публикации
	Link            string `xml:"link" `        //Ссылка на источник
}

// Структура RSS(XML) документа
type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Item []XML `xml:"item"`
}

// Конструктор для RSS
func New() *RSS {
	return &RSS{Channel: Channel{Item: make([]XML, 0)}}
}

// Получение XML-документа с RSS-ленты
func (rss *RSS) GetPosts(url string) (Channel, error) {
	resp, err := http.Get(url) //Запрос XML по адрессу
	if err != nil {
		log.Println("Ошибка получения XML документа", err)
		return Channel{}, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body) //Прочтение XML документа
	if err != nil {
		log.Println("Ошибка прочтения XML-файла", err)
	}
	var result RSS
	err = xml.Unmarshal(b, &result) //Депарсинг XML-документа
	if err != nil {
		log.Println("Ошибка парсинга XML:", err)
		return Channel{}, err
	}
	for i := range result.Channel.Item {
		result.Channel.Item[i].Description = extractText(result.Channel.Item[i].Description)
	}
	return result.Channel, nil
}

// HTML-парсер
// HTML-парсер
func extractText(str string) string {
	// Сначала заменяем escape-последовательности
	str = strings.ReplaceAll(str, `\u003c`, "<")
	str = strings.ReplaceAll(str, `\u003e`, ">")
	str = strings.ReplaceAll(str, `\u0026amp;`, "&")

	doc, err := html.Parse(strings.NewReader(str))
	if err != nil {
		return str
	}

	var f func(*html.Node) string
	f = func(n *html.Node) string {
		if n.Type == html.TextNode {
			return n.Data
		}
		if n.Type == html.ElementNode && (n.Data == "code" || n.Data == "script" || n.Data == "style") {
			return ""
		}
		var res string
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			res += f(c)
		}
		return res
	}

	return strings.TrimSpace(f(doc))
}
