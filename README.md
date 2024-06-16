# D0SL_organizer
### Запуск: 
Поднимаем milvus:

`docker compose up -d`

запускаем сервер на порте 8001:

`go run cmd/main/main.go`

### Ручки:

GET `http://localhost:8001/similar-videos` - схожие видео. Тело запроса: `{"message" : "1,2,3,4, ..."}`

POST `http://localhost:8001/video` - добавить новое видео. Тело запроса: `{{"id":1, "link":"sdc", "description":"dsc", "vector":[1,1, 1, 1, ...]}}`

Размерность векторов: `768`
