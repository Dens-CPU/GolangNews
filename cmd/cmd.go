package main

import (
	"GplangNews/pkg/api"
	"GplangNews/pkg/postgres"
	"GplangNews/pkg/rss"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// Структура конфиг файла
type Config struct {
	URLs           []string `json:"rss"`
	Request_period int      `json:"request_period"`
}

// Парсинга кофиг файла
func ParseConfig(fileName string) (Config, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Println("Ошибка прочтения конфиг файла")
		return Config{}, err
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Println("Ошибка депарсинга конфиг файла")
		return Config{}, err
	}
	return config, nil
}

const dsn = "postgres://postgres:12345@127.0.0.1:5432/GolangNews"

func main() {
	var wg sync.WaitGroup

	//Канал для постов
	postsChan := make(chan []rss.XML, 6)

	//Канал для ошибок
	errorChan := make(chan error, 5)

	// Подключение к базе данных
	db, err := postgres.New(dsn)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных")
	}
	log.Println("Подключение успешно")
	defer db.DB.Close()

	//Создание объекта API, использующего БД
	api := api.New(&db)

	//Парсинг конфиг файла
	configFile, err := ParseConfig("config.json")
	if err != nil {
		log.Println(err)
	}
	// Запуск RSS агрегатора
	rss := rss.New()

	//Прослушивание RSS ленты
	for _, url := range configFile.URLs {

		//Запуск горутины с соответвующим URL
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			for {
				channel, err := rss.GetPosts(url) //Получения XML-документа
				if err != nil {
					errorChan <- err //Запись ошибки в канал
				} else {
					postsChan <- channel.Item //Запись постов в канал
				}
				time.Sleep(time.Duration(configFile.Request_period) * time.Second) //Пауза
			}
		}(url)
	}

	//Запись XML-файлов в БД
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			db.AddPost(postsChan)
		}
	}()

	//Обработчик ошибок
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case err := <-errorChan:
				log.Println(err)
			default:
				<-time.After(5 * time.Second)
			}
		}
	}()
	//Запуск сетевой службу и HTTP-сервера на всех IP-адресах и порту 80
	err = http.ListenAndServe(":80", api.Router())
	if err != nil {
		log.Fatal(err)
	}
	wg.Wait()
}
