package api

import (
	"GplangNews/pkg/postgres"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Структура API приложения
type API struct {
	r  *mux.Router     //Маршрутизатор
	db *postgres.Store //База данных
}

// Конструктор для API
func New(db *postgres.Store) *API {
	api := API{}
	api.r = mux.NewRouter()
	api.db = db
	api.endpoints()
	return &api
}

// Router возвращает маршрутизатор запросов
func (api *API) Router() *mux.Router {
	return api.r
}

// Регистрация методов в маршрутизаторе запросов
func (api *API) endpoints() {
	api.r.HandleFunc("/news/{n}", api.posts).Methods(http.MethodGet, http.MethodOptions) //Получение последних n новостей

	//Веб приложение
	api.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

func (api *API) posts(w http.ResponseWriter, r *http.Request) {
	//Считывание папраметра n из пути запроса
	// /news/10
	s := mux.Vars(r)["n"]
	n, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var news []postgres.Post
	news, err = api.db.GetPosts(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//Отправка данных клиенту в формате JSON
	json.NewEncoder(w).Encode(news)
}
