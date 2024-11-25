# Файлы для итогового задания

`cmd/scheduler` - лежит main.go  
`config` - конфигурация приложения  
`internal` - код необходимый для работы приложения  
`storage` - директория для базы  

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.

# Задачи со звёздочкой

- [x] TODO_PORT
- [x] TODO_DBFILE
- [ ] NextDate
  - [x] w
  - [ ] m
- [ ] Search
- [ ] Auth
- [x] Docker

# docker

Для запуска в docker нужен файл .env

```
TODO_PORT=7540
TODO_DBFILE="./storage/scheduler.db"
```

Команда для запуска:

```
docker compose up
```
