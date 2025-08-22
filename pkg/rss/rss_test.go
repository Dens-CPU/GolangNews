package rss

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRSS_GetPosts(t *testing.T) {
	//Заглушка XML-документа
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
	<rss version="2.0">
		<channel>
			<item>
				<title>Test Title</title>
				<link>http://example.com</link>
				<description>Test Description</description>
				<pubDate>Tue, 19 Aug 2025 14:15:49 GMT</pubDate>
			</item>
		</channel>
	</rss>`

	//Создание тестового сервера, возращающий XML-документ
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(xmlData))
	}))
	defer server.Close()

	rss := New()
	//Вызов функции получения XML-документа
	channel, err := rss.GetPosts(server.URL)
	if err != nil {
		t.Errorf("Ожидалось nil, получили %v", err)
	}
	if len(channel.Item) != 1 {
		t.Errorf("получили %d, ожидалось %d", len(channel.Item), 1)
	}

}
