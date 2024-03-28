# marketplace-api

## Содержание

1. [Запуск](#Запуск)
2. [Конфигурация](#Конфигурация)
2. [API](#API)
3. [Реализация](#Реализация)

## Структура проекта 

```
.
├── auth          // token provider utils
├── cmd           // entry points
│  └── app        // app entry point
├── config        // config loading utils
├── internal
│  ├── app        // main application package
│  ├── domain     // all business entities
│  ├── handler    // http handlers layer
│  ├── repository // database repository layer
│  └── service    // business logic services layer
└── migrations    // database migrations
 
```

## Конфигурация

В файле .env.local:
```
DB_USER=marketplace            # Пользователь базы данных
DB_PASSWORD=marketplace        # Пароль для базы данных
DB_HOST=db                     # Имя хоста базы данных
DB_PORT=5432                   # Порт базы данных
DB_NAME=marketplace            # Название базы данных
JWT_ACCESS_EXPIRATION_TIME=10  # Время действия access токена
JWT_REFRESH_EXPIRATION_TIME=60 # Время действия refresh токена
JWT_SECRET=jwt-secret 
REDIS_HOST=redis               # Имя хоста redis
REDIS_PORT=6379                # Порт redis
AD_PER_PAGE=2                  # Количество объявлений на одной странице
CHECK_IMAGE_IDLE_TIMEOUT=30    # Максимальное количество времени 
на получение изображения(с)
MIN_IMAGE_WIDTH=64             # Минимальная ширина изображения (пикс.)
MAX_IMAGE_WIDTH=4096           # Максимальная ширина изображения (пикс.)
MIN_IMAGE_HEIGHT=64            # Минимальная высота изображения (пикс.)
MAX_IMAGE_HEIGHT=4096          # Максимальная высота изображения (пикс.)
GIN_MODE=release               # Отключение режима отладки Gin
```

## Запуск

```shell
docker compose up
```

---
## API
### POST /user/signup
Регистрация пользователя посредством отправки логина и пароля. 
Логин может иметь от 4 до 64 символов. Пароль может иметь от 8 до 64 символов.

Запрос:

```
curl --location 'localhost:8080/user/signup' \
--header 'Content-Type: application/json' \
--data '{
    "login": "user",
    "password": "password"
}'
```

Ответ:

```json
{
    "id": 1,
    "login": "user",
    "password": "password"
}
```
---
### POST /user/login
Авторизация пользователя посредством отправки логина и пароля. Если 
данные имеют корректный формат и пользователь существует, возвращаются access и refresh токены.
Access токен передается в запросах в заголовке Authorization.

Запрос:

```
curl --location 'localhost:8080/user/login' \
--header 'Content-Type: application/json' \
--data '{
    "login": "user",
    "password": "password"
}'
```

Ответ:

```json
{
    "access_token":  "<access_token>",
    "refresh_token": "<refresh_token>"
}
```
---
### POST /user/logout
Выход пользователя посредством access токена

Запрос:

```
curl --location --request POST '185.25.51.102:8080/user/logout' \
--header 'Authorization: <access_token>'
```

Ответ:

```json
{
    "msg": "successful logout"
}
```

### GET /user/refresh
Обновление access и refresh токена пользователя по переданному refresh 
токену.

Запрос:

```
curl --location 'localhost:8080/user/refresh' \
--header 'Authorization: <refresh_token>'
```

Ответ:

```json
{
    "access_token": "<access_token>",
    "refresh_token": "<refresh_token>"
}
```
---
### POST /ad
Размещение объявления. Длина заголовка ограничена 255 символами. Длина 
текста объявления и URL на изображение ограничены 2048 символами. 
Принимаемый формат изображения: JPEG.

Запрос:

```
curl --location 'localhost:8080/ad' \
--header 'Authorization: <access_token>' \
--header 'Content-Type: application/json' \
--data '{
    "title": "Car",
    "body": "Mazda",
    "image_url": "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcStU7CfhDVInPhuGwmnZMkMQAIHEXwkEGUzhd6HYao5bw&s",
    "price": 150
}'
```

Ответ:

```json
{
    "id": 2,
    "user_id": 1,
    "title": "Car",
    "body": "Mazda",
    "image_url": "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcStU7CfhDVInPhuGwmnZMkMQAIHEXwkEGUzhd6HYao5bw&s",
    "price": 150,
    "created_at": "2024-03-28T10:54:33.984747Z"
}
```
---

### GET /ad?page={num}&sort={date,price}&dir={asc, desc}&min={num}&max={num}
Отображение ленты объвлений.

Для сортировки используется параметр запроса sort:
- sort=date - сортировка по дате;
- sort=price - сортировка по цене.

Для указания направления сортировки используется параметр запроса dir:
- dir=asc - по возрастанию;
- dir=desc - по убыванию.

Если параметр dir не указан, то используется сортировка по возрастанию.

Параметр min задает минимальную цену для отображения объявлений. Параметр max 
задает маскимальную цену, отображаемых объявлений.

Запрос:

```
curl --location 'localhost:8080/ad?sort=price&dir=desc' \
--header 'Authorization: <access_token>'
```

Ответ:

```json
[
    {
        "posted_by_you": true,
        "user_login": "user",
        "title": "iPhone",
        "body": "iPhone SE",
        "image_url": "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcStU7CfhDVInPhuGwmnZMkMQAIHEXwkEGUzhd6HYao5bw&s",
        "price": 200,
        "created_at": "2024-03-28T10:53:13.029279Z"
    },
    {
        "posted_by_you": true,
        "user_login": "user",
        "title": "Car",
        "body": "Mazda",
        "image_url": "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcStU7CfhDVInPhuGwmnZMkMQAIHEXwkEGUzhd6HYao5bw&s",
        "price": 150,
        "created_at": "2024-03-28T10:54:33.984747Z"
    }
]
```

## Реализация

- В качестве базы данных используется PostgreSQL;
- Для авторизации пользователей используются JWT токены. Access токен передается в заголовке Authorization запроса;
- Для хранение токенов используется хранилище Redis;
- Количество объявлений на одной странице и ограничения по размеру изображения задаются в файле .env.local