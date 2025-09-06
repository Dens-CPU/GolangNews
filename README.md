# GolangNews
GolangNews-это приложение,написанное на языке Golang для получения последних новостей из RSS-лент новостей.

## Установка
### Клонирование объекта

```bash
git clone https://github.com/Dens-CPU/GolangNews
cd GolangNews
```
### Установка зависимостей
После кланирования пакета с приложением требуется установить или обновить зависимости

```bash
go mod tidy
```

### Настройка приложения
В приложении используется переменная окружения `DB_DSN`, находящаяся в файле `.env` в пакете `cmd`. Данная переменная отвечает за подключение к базе данных. Требуется указать свой `dsn` для подключения к БД PostgreSQL.
```bash
DB_DSN=postgres://postgres:******@host.docker.internal:5432/GolangNews
```
В базе данных требуется создать таблицу со структурой в файле `shema.sql`
```sql
DROP TABLE IF EXISTS posts;

CREATE TABLE posts (
id SERIAL PRIMARY KEY,
title TEXT NOT NULL,
content TEXT NOT NULL,
pubTime BIGINT NOT NULL,
link TEXT NOT NULL
);
```
### Локальный запуск
Запуск через командную строку
```bash
go run .\cmd\cmd.go
```
Сборка бинарного файла
```bash
go build .\cmd\cmd.go
```
После запуска сервер будет доступен по:
```arduino
http://localhost:8080
```
Или по прямому пути
```arduino
http://localhost:8080/news/n
```
n - количество последних новостей



