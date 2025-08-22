package api

import (
	"GplangNews/pkg/postgres"
	"GplangNews/pkg/rss"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPI_posts(t *testing.T) {
	//Создание чистого объекта API для теста
	const dsn = "postgres://postgres:12345@127.0.0.1:5432/GolangNews"
	dbase, err := postgres.New(dsn)
	if err != nil {
		t.Errorf("Ошибка подключения к базе данных")
	}
	channel := make(chan []rss.XML, 1)
	channel <- []rss.XML{
		{
			Title:           "test",
			Description:     "desc",
			PublicationDate: "12",
			Link:            "http://example.com",
		},
	}
	dbase.AddPost(channel)

	api := New(&dbase)

	//Создание HTTP запроса
	req := httptest.NewRequest(http.MethodGet, "/news/10", nil)

	//Создание объекта для записи ответа обработкика
	rr := httptest.NewRecorder()

	//Вызов маршрутизатора.Маршрутизатор для пути и метода запроса вызовет обработчик
	//Обработчик запишет ответ в созданный объект

	api.r.ServeHTTP(rr, req)

	//Проверка кода ответа
	if rr.Code != http.StatusOK {
		t.Errorf("Код неверен: полкчили %d, а хотели %d", rr.Code, http.StatusOK)
	}

	// Прочтение тела ответа
	b, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("Не удалось раскодировать ответ сервера: %v", err)
	}

	//Декодирование JSON в массив заказов
	var data []postgres.Post
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("Ошибка декодирования: %v", err)
	}

	//Проверка длины массива
	const wantLen = 1
	if len(data) != wantLen {
		t.Fatalf("Получено %d записей, хотели %d", len(data), wantLen)
	}
	err = api.db.DeletePost()
	if err != nil {
		t.Fatalf("Ошибки очистки базы данных %v", err)
	}
}
