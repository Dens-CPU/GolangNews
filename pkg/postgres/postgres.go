package postgres

import (
	"GplangNews/pkg/rss"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Хранилище данных
type Store struct {
	DB *pgxpool.Pool
}

// Структура данных публикации
type Post struct {
	ID      int    //Номер записи
	Title   string //Заголовок
	Content string //Содержание
	PubTime int64  //Дата публикации
	Link    string //Ссылка на источник
}

// Конструктор
func New(dsn string) (Store, error) {
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Println("Ошибка подключения к базе данных", err)
		return Store{nil}, err
	}
	return Store{DB: db}, nil
}

// Мапа для уникальности добавляемых постов
var CheckMap = make(map[string]bool)
var CheckMutex sync.Mutex

// Добавление публикаций
func (s *Store) AddPost(postChan chan []rss.XML) error {

	//Передача данных из канала
	select {
	case posts := <-postChan:

		//Регистрация транзакции
		tx, err := s.DB.Begin(context.Background())
		if err != nil {
			fmt.Println("Ошибка регистрация транзакии")
			return err
		}
		for _, post := range posts {

			//Проверка на уникальность статьи
			CheckMutex.Lock()
			if CheckMap[post.Title] {
				CheckMutex.Unlock()
				continue
			}
			CheckMap[post.Title] = true
			CheckMutex.Unlock()

			//Парсинг даты публикации в формат int64
			t, err := time.Parse(time.RFC1123Z, post.PublicationDate)
			if err != nil {
				t, err = time.Parse(time.RFC1123, post.PublicationDate) // запасной вариант
			}
			pubTime := t.Unix()

			//Запрос на вставку данных
			_, err = tx.Exec(context.Background(), `
		INSERT INTO posts (title,content,pubTime,link)
		VALUES ($1,$2,$3,$4);
		`, post.Title, post.Description, pubTime, post.Link)

			if err != nil {
				tx.Rollback(context.Background()) //Откат при ошибке
				log.Println("Ошибка вставки элекмента в базу данных")
				return err
			}
		}
		//Фиксация транзакции
		tx.Commit(context.Background())

	default:
		<-time.After(2 * time.Second)
	}
	return nil
}

// Прочтение публикаций из БД
func (s *Store) GetPosts(numbers int) ([]Post, error) {
	var posts []Post

	rows, err := s.DB.Query(context.Background(), `
	SELECT * FROM posts 
	ORDER BY id DESC 
	LIMIT $1
	;
	`, numbers)
	if err != nil {
		log.Println("Ошибка прочтения публикаций из БД:", err)
		return nil, err
	}
	for rows.Next() {
		var p Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			log.Println("Ошибка сканирования полученных результатов:", err)
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

// Удаление заказов из БД
func (s *Store) DeletePost() error {
	_, err := s.DB.Query(context.Background(), `
	DELETE FROM posts;
	`)
	if err != nil {
		return err
	}
	return nil
}
