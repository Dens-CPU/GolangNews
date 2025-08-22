package postgres

import (
	"GplangNews/pkg/rss"
	"testing"
)

func TestStore_GetPosts(t *testing.T) {
	const dsn = "postgres://postgres:12345@127.0.0.1:5432/GolangNews"
	db, _ := New(dsn)
	var post = make(chan []rss.XML, 1)
	post <- []rss.XML{
		{
			Title:           "1",
			Description:     "1",
			PublicationDate: "1",
			Link:            "1",
		},
		{
			Title:           "2",
			Description:     "2",
			PublicationDate: "2",
			Link:            "2",
		},
	}
	err := db.AddPost(post)
	if err != nil {
		t.Errorf(err.Error())
	}
	numbers := 2
	want := 2
	dates, _ := db.GetPosts(numbers)
	get := len(dates)
	if want != get {
		t.Errorf("Получили %d, ожидали %d\n", get, want)
	}

	err = db.DeletePost()
	if err != nil {
		t.Fatalf("Ошибка очитски БД: %v", err)
	}
}
