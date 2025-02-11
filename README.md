# Сборка и запуск
1. Убедитесь, что порты 5432 и 9090 ничем не заняты. Если заняты - освободить.
2. В корне проекта создайте файл .env, пример содержания:</b>
```
DATABASE_URL=postgres://postgres:postgres@postgres:5432/testdb?sslmode=disable</br>
POSTGRES_DB=testdb</br>
POSTGRES_USER=postgres</br>
POSTGRES_PASSWORD=postgres</br>
```
## С помощью docker-compose
### Запуск с in-memory хранилищем
```
docker compose -f docker-compose.inmemory.yml up -d --build
```
### Остановка с in-memory хранилищем
```
docker compose -f docker-compose.inmemory.yml down --volumes
```
### Запуск с postgresql хранилищем
```
docker compose --profile postgres -f docker-compose.yml up -d --build
```
### Остановка с postgresql хранилищем
```
docker compose -f docker-compose.yml down --volumes
```

## С помощью Makefile
### Запуск с in-memory хранилищем
```
STORAGE=inmemory make -f MakeFile start
```
### Остановка с in-memory хранилищем
```
STORAGE=inmemory make -f MakeFile stop
```
### Запуск с postgresql хранилищем
```
STORAGE=postgres make -f MakeFile start
```

### Остановка с postgresql хранилищем
```
STORAGE=postgres make -f MakeFile stop
```

# Тестирование
## Запуск юнит-тестов
```
make -f MakeFile test
make -f MakeFile test-coverage
```

# API
### Важный момент: Запросы по HTTP по 8080 портуЮ одновременно по 9090 работает GRPC
> 1) Тело запроса/ответа - в формате JSON.
> 2) В случае ошибки возвращается необходимый HTTP код, в теле содержится описание ошибки (пример: ```{"error": "something went wrong"}```).

## POST /api/v1/create
Сокращение исходного URL.

- Параметры тела запроса:
    - OriginalURL - исходный URL.
- Тело ответа:
    - shortURL - уникальная сокращенная ссылка, присвоенная данному адресу URL.

**Пример HTTP**

Запрос:

```
curl --location 'http://localhost:8080/api/v1/create' \
--header 'Content-Type: application/json' \
--data '{
    "OriginalURL": "https://www.ozon.ru/product/chistaya-arhitektura-iskusstvo-razrabotki-programmnogo-obespecheniya-martin-robert-144499396/?at=XQtkZn4kvhDXMzoMurR1G6nc5pXQ77sq8vqlNfgVrLqo&avtc=1&avte=4&avts=1739295310&keywords=%D0%BA%D0%BD%D0%B8%D0%B3%D0%B0+%D1%80%D0%BE%D0%B1%D0%B5%D1%80%D1%82+%D0%BC%D0%B0%D1%80%D1%82%D0%B8%D0%BD"
}'
```

Ответ:

```
{
    "shortURL": "NW3l_ulICj"
}
```

**Пример GRPC**
Запрос: не разобрался как сохранить коллекцию из postman...

```
{"original_url": "https://www.ozon.ru/product/chistaya-arhitektura-iskusstvo-razrabotki-programmnogo-obespecheniya-martin-robert-144499396/?at=XQtkZn4kvhDXMzoMurR1G6nc5pXQ77sq8vqlNfgVrLqo&avtc=1&avte=4&avts=1739295310&keywords=%D0%BA%D0%BD%D0%B8%D0%B3%D0%B0+%D1%80%D0%BE%D0%B1%D0%B5%D1%80%D1%82+%D0%BC%D0%B0%D1%80%D1%82%D0%B8%D0%BD"}
```

Ответ:

```
{
    "short_url": "n6hZ_tlF0B"
}
```

## GET /api/v1/get
Восстановление оригинального URL.

- Тело ответа:
    - originalURL - оригинальный URL.

**Пример HTTP**

Запрос:

```
curl --location --request GET 'http://localhost:8080/api/v1/get' \
--header 'Content-Type: application/json' \
--data '{
    "shortURL": "NW3l_ulICj"
}'
```

Ответ:

```
{
    "originalURL": "https://www.ozon.ru/product/chistaya-arhitektura-iskusstvo-razrabotki-programmnogo-obespecheniya-martin-robert-144499396/?at=XQtkZn4kvhDXMzoMurR1G6nc5pXQ77sq8vqlNfgVrLqo&avtc=1&avte=4&avts=1739295310&keywords=%D0%BA%D0%BD%D0%B8%D0%B3%D0%B0+%D1%80%D0%BE%D0%B1%D0%B5%D1%80%D1%82+%D0%BC%D0%B0%D1%80%D1%82%D0%B8%D0%BD"
}
```

**Пример GRPC**

Запрос:

```
{"short_url": "n6hZ_tlF0B"}
```

Ответ:

```
{
    "original_url": "https://www.ozon.ru/product/chistaya-arhitektura-iskusstvo-razrabotki-programmnogo-obespecheniya-martin-robert-144499396/?at=XQtkZn4kvhDXMzoMurR1G6nc5pXQ77sq8vqlNfgVrLqo&avtc=1&avte=4&avts=1739295310&keywords=%D0%BA%D0%BD%D0%B8%D0%B3%D0%B0+%D1%80%D0%BE%D0%B1%D0%B5%D1%80%D1%82+%D0%BC%D0%B0%D1%80%D1%82%D0%B8%D0%BD"
}
```

## TODO
- Сделать нормальные тесты
- Улучшить или заменить алгоритм сокращения ссылок
- Сделать логирование
- Сделать конфиг
- swagger docs
- Метрики

## Что выполнено:
- Реализована работа по http и grpc
- Работа через postgres или inmemeory хранилище
- Рабочий конфиг для запуска приложения в докере
- Реализовано API, указанное в задании
- Graceful shutdown
- Написана конфигурация запуска линтеров