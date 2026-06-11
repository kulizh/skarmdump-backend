# Skarmdump Backend

HTTP-сервис для загрузки PNG-изображений. Принимает файл, вычисляет SHA-256 хеш, сохраняет локально или в S3, возвращает URL. Дубликаты определяются по хешу — повторная загрузка того же файла возвращает ту же ссылку.

---

## Запуск через Docker

### 1. Настройка окружения

Скопируйте `.env.example` в `.env` и заполните:

```env
PORT=8080
DOMAIN=http://localhost:8080
LOCAL_PATH=./img
S3_BUCKET=
S3_REGION=
S3_URL=
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
USER_AGENT=
```

| Переменная | Обязательно | По умолчанию | Описание |
|---|---|---|---|
| `PORT` | нет | `8080` | Порт внутри контейнера |
| `DOMAIN` | да | — | Базовый URL (возвращается в ответе) |
| `LOCAL_PATH` | нет | `./img` | Директория для локального хранения |
| `API_KEY` | нет | — | Ключ авторизации (см. раздел Авторизация) |
| `USER_AGENT` | нет | — | Разрешить запросы только с этого User-Agent, если не указан, принимаем всё|
| `HASH_LENGTH` | нет | `12` | Длина хеша в символах |
| `MAX_FILE_SIZE_MB` | нет | `10` | Максимальный размер файла в МБ |
| `S3_BUCKET` | для S3 | — | Имя S3-бакета |
| `S3_REGION` | для S3 | — | Регион S3 |
| `S3_URL` | для S3 | — | Публичный URL бакета |
| `AWS_ACCESS_KEY_ID` | для S3 | — | AWS-ключ |
| `AWS_SECRET_ACCESS_KEY` | для S3 | — | AWS-секрет |

### 2. Запуск

```bash
make compose-up
```

Или вручную:

```bash
docker compose up -d --build
```

Сервис будет доступен на `http://localhost:8080`.

Остановка:

```bash
make compose-down
```

---

## API

### Авторизация

Если в `.env` задан `API_KEY`, все запросы к `/upload` должны содержать заголовок:

```
Authorization: <значение API_KEY>
```

Если `API_KEY` пустой (не задан) — авторизация не требуется.

### Загрузить изображение

```
POST /upload
Content-Type: multipart/form-data
Authorization: <API_KEY>
```

**Поля формы:**

| Поле | Тип | Обязательно | Описание |
|---|---|---|---|
| `image` | file | да | PNG-файл |
| `s3` | string | нет | `true` — сохранить в S3, иначе локально |

**Успех (200):**

```json
{
  "success": true,
  "url": "http://localhost:8080/abc12345"
}
```

**Пример запроса:**

```bash
curl -X POST http://localhost:8080/upload \
  -F "image=@photo.png" \
  -F "s3=false"
```

### Embed-страница (Open Graph)

```
GET /<hash>
```

Возвращает HTML-страницу с OG-тегами для embed-превью при отправке ссылки в мессенджеры / соцсети.

**Пример:**

```
GET /abc12345
```

Страница содержит:
- `og:image` — ссылка на изображение (`{DOMAIN}/abc12345.png`)
- `og:title`, `og:description` — заголовок и описание
- `twitter:card` — `summary_large_image`
- meta refresh + тело с `<img>` для ручного открытия страницы

Само изображение доступно по `/<hash>.png` (например `/abc12345.png`) — эту ссылку nginx отдаёт напрямую без участия Go.

**Nginx:** путь `/<hash>.png` (с расширением) nginx отдаёт сам. Остальные пути (`/<hash>`, `/upload`) проксируются на Go-бэкенд.

---

## Коды ошибок

| HTTP | `error` | `message` | Причина |
|---|---|---|---|
| 401 | `auth_required` | `authorization header is required` | Не передан заголовок `Authorization` |
| 403 | `invalid_key` | `invalid API key` | Неверный API-ключ |
| 403 | `bad_user_agent` | `request from this User-Agent is not allowed` | User-Agent не совпадает с `USER_AGENT` |
| 400 | `missing_image` | `image file is required` | Не передан файл в поле `image` |
| 400 | `read_error` | `cannot read uploaded image` | Ошибка чтения файла |
| 400 | `invalid_image` | `only PNG files are accepted` | Файл не является PNG |
| 429 | `rate_limit` | `rate limit exceeded` | Слишком много запросов с одного IP |
| 500 | `storage_error` | `failed to save image` | Ошибка записи (диск / S3) |

**Пример ответа с ошибкой:**

```json
{
  "success": false,
  "error": "invalid_image",
  "message": "only PNG files are accepted"
}
```

---

## Ограничения

- Принимаются только PNG
- Нет аутентификации
- Нет удаления файлов
- Rate limit: 1 запрос/сек на upload, 3 запроса/сек на GET (по IP)