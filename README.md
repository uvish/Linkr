<img src="ui/logo.png" style="width:150px">
<br>
# 🔗 Linkr

A self-hosted, minimal URL shortener built with Golang. Fast, simple, and Docker-ready.

## 🐳 Run with Docker

```bash
docker-compose up --build
```

- API : [http://localhost:8080](http://localhost:8080)  
- UI : [http://localhost:1234](http://localhost:1234)

## 📘 API Endpoints

### POST `/shorten`

Create a short URL.

**Request Body (JSON):**

```json
{
  "url": "https://example.com"
}
```

**Response:**

```json
{
  "shortCode": "abc123"
}
```

---

### GET `/:shortCode`

Redirect to the original URL for a given short code.

---

### GET `/urls`

Get all shortened URLs.

**Response:**

```json
[
  {
    "shortCode": "abc123",
    "originalURL": "https://example.com"
  }
]
```

---

### DELETE `/urls/:shortCode`

Delete a shortened URL by its short code.

**Response:**

```json
{
  "message": "Deleted successfully"
}
```

## 🛠️ Built Using

- Golang
- Gin
- PostgreSQL
- Docker & Docker Compose

## 📄 License

MIT License. See [LICENSE](./LICENSE) for details.
