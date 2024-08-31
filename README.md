# notes-app
Сервис создания заметок.

Используемые технологии и пакеты:  
- Хранилище - PostgreSQL (database/sql, pq, migrate)
- Docker
- Сервер - net/http
- Маршрутизатор - gorilla/mux
- Логирование - log/slog
- Конфигурации - .env, yaml (cleanenv, godotenv)
- Аутентификация/авторизация - JWT-токены
  
Сервис был написан с использованием чистой архитектуры, был реализован Graceful Shutdown для корректного завершения работы. Валидация орфографических ошибок происходит путем добавления результата проверки в тело ответа на запрос добавления заметки.

# start app  
- Перед запуском установить необходимые конфиги (создать .env файл. Шаблон env конфига в файле .env.example)
- Запуск ```docker-compose up --build```

# curl templates
## Регистрация  
```
curl --location --request POST 'localhost:YOUR-PORT/api/auth/registration' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "YOUR-EMAIL",
    "password": "YOUR-PASSWORD-GTE8"
}'
```
## Логин  
```
curl --location --request POST 'localhost:YOUR-PORT/api/auth/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "YOUR-EMAIL",
    "password": "YOUR-PASSWORD-GTE8"
}'
```
## Обновление токена  
```
curl --location --request PUT 'localhost:YOUR-PORT/api/auth/refresh' \
--header 'Cookie: refresh-token=YOUR-REFRESH-TOKEN'
```
## Добавление заметки  
```
curl --location --request POST 'localhost:YOUR-PORT/api/notes' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer YOUR-ACCESS-TOKEN' \
--data '{
    "note": "YOUR-NOTE"
}'
```
## Получение заметок  
```
curl --location --request GET 'localhost:YOUR-PORT/api/notes' \
--header 'Authorization: Bearer YOUR-ACCESS-TOKEN'
```
## Получение заметок с пагинацией  
```
curl --location --request GET 'localhost:YOUR-PORT/api/notes?page=PAGE-NUMBER&limit=LIMIT-COUNT' \
--header 'Authorization: Bearer YOUR-ACCESS-TOKEN'
```

# examples
## Регистрация  
- Запрос
```
curl --location --request POST 'localhost:12021/api/auth/registration' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "email@gmail.com",
    "password": "my-good-password"
}'
```
- Ответ
```
{
    "id": 2
}
```
## Логин  
- Запрос
```
curl --location --request POST 'localhost:12021/api/auth/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "email@gmail.com",
    "password": "my-good-password"
}'
```
- Ответ
```
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6MiwiZXhwIjoxNzI1MTA4MDgwfQ.jU-CUU9HHCusRbnf6_EBd_6ZaB9aJKUHVXpSWXDT7Bg"
}
+ Set-Cookie: refresh-token=afcef4d1-b207-4adf-90c3-4ae79dd2a317; Path=/api/auth; Domain=localhost; HttpOnly; SameSite=Strict
```
## Обновление токена  
- Запрос
```
curl --location --request PUT 'localhost:12021/api/auth/refresh' \
--header 'Cookie: refresh-token=afcef4d1-b207-4adf-90c3-4ae79dd2a317'
```
- Ответ
```
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6MiwiZXhwIjoxNzI1MTA4MzIwfQ.S_ozyoDeW9n0G08huwunfiD1cXRVCqUMEPGa9unzQEg"
}
+ Set-Cookie: refresh-token=12b81c5f-4d7a-4027-8c7f-b8f39cbdf5f5; Path=/api/auth; Domain=localhost; HttpOnly; SameSite=Strict
```
## Добавление заметки  
- Запрос (заметка без ошибок)
```
curl --location --request POST 'localhost:12021/api/notes' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6MiwiZXhwIjoxNzI1MTA4MzIwfQ.S_ozyoDeW9n0G08huwunfiD1cXRVCqUMEPGa9unzQEg' \
--data '{
    "note": "пример заметки без ошибок"
}'
```
- Ответ
```
{
    "id": 2,
    "spellingErrors": []
}
```
- Запрос (заметка с ошибкой)
```
curl --location --request POST 'localhost:12021/api/notes' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6MiwiZXhwIjoxNzI1MTA4MzIwfQ.S_ozyoDeW9n0G08huwunfiD1cXRVCqUMEPGa9unzQEg' \
--data '{
    "note": "пример заметки с ашибкой"
}'
```
- Ответ
```
{
    "id": 3,
    "spellingErrors": [
        {
            "code": 1,
            "pos": 17,
            "word": "ашибкой",
            "s": [
                "ошибкой"
            ]
        }
    ]
}
```
## Получение заметок  
- Запрос
```
curl --location --request GET 'localhost:12021/api/notes' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6MiwiZXhwIjoxNzI1MTA4MzIwfQ.S_ozyoDeW9n0G08huwunfiD1cXRVCqUMEPGa9unzQEg'
```
- Ответ
```
{
    "notes": [
        {
            "id": 2,
            "note": "пример заметки без ошибок",
            "createdAt": "2024-08-31T12:33:01.847906Z"
        },
        {
            "id": 3,
            "note": "пример заметки с ашибкой",
            "createdAt": "2024-08-31T12:34:47.713226Z"
        }
    ]
}
```
## Получение заметок с пагинацией  
- Список всех заметок пользователя
```
{
    "notes": [
        {
            "id": 2,
            "note": "пример заметки без ошибок",
            "createdAt": "2024-08-31T12:33:01.847906Z"
        },
        {
            "id": 3,
            "note": "пример заметки с ашибкой",
            "createdAt": "2024-08-31T12:34:47.713226Z"
        },
        {
            "id": 4,
            "note": "заметка 3",
            "createdAt": "2024-08-31T12:37:45.34428Z"
        },
        {
            "id": 5,
            "note": "заметка 4",
            "createdAt": "2024-08-31T12:37:49.37496Z"
        },
        {
            "id": 6,
            "note": "заметка 5",
            "createdAt": "2024-08-31T12:37:52.172993Z"
        },
        {
            "id": 7,
            "note": "заметка 6",
            "createdAt": "2024-08-31T12:37:55.063469Z"
        },
        {
            "id": 8,
            "note": "заметка 7",
            "createdAt": "2024-08-31T12:37:58.242772Z"
        },
        {
            "id": 9,
            "note": "заметка 8",
            "createdAt": "2024-08-31T12:38:03.944704Z"
        }
    ]
}
```
- Запрос (страница 1, лимит 2)
```
curl --location --request GET 'localhost:12021/api/notes?page=1&limit=2' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6MiwiZXhwIjoxNzI1MTA4MzIwfQ.S_ozyoDeW9n0G08huwunfiD1cXRVCqUMEPGa9unzQEg'
```
- Ответ
```
{
    "notes": [
        {
            "id": 2,
            "note": "пример заметки без ошибок",
            "createdAt": "2024-08-31T12:33:01.847906Z"
        },
        {
            "id": 3,
            "note": "пример заметки с ашибкой",
            "createdAt": "2024-08-31T12:34:47.713226Z"
        }
    ]
}
```
- Запрос (страница 2, лимит 4)
```
curl --location --request GET 'localhost:12021/api/notes?page=2&limit=4' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6MiwiZXhwIjoxNzI1MTA4MzIwfQ.S_ozyoDeW9n0G08huwunfiD1cXRVCqUMEPGa9unzQEg'
```
- Ответ
```
{
    "notes": [
        {
            "id": 6,
            "note": "заметка 5",
            "createdAt": "2024-08-31T12:37:52.172993Z"
        },
        {
            "id": 7,
            "note": "заметка 6",
            "createdAt": "2024-08-31T12:37:55.063469Z"
        },
        {
            "id": 8,
            "note": "заметка 7",
            "createdAt": "2024-08-31T12:37:58.242772Z"
        },
        {
            "id": 9,
            "note": "заметка 8",
            "createdAt": "2024-08-31T12:38:03.944704Z"
        }
    ]
}
```
