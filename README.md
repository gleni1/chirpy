# 🌟 Chirpy 🐦

**Chirpy** is a lightweight 🧱 HTTP server written in Go that provides a range of tools for managing users and "chirps"️(posts). It exposes a set of RESTful 🌀 API endpoints to make user management, content creation ✍️, and admin 🛡️ operations easier.

## Features ✨


- 👥🔐 User authentication & token management 🔑
- 🆕🛠️ User creation & updating ✍️
- 📄 CRUD operations for chirps 🐤
- 🌐 Webhook 🌊 integration
- 📈 Monitoring 📊 & metrics for admins 🧑‍💻

---

## API Endpoints 🌐

### Public 🌏

#### **`GET /app/`** 🏡
- 📜 Displays the home 🖼️ page.
- Middleware 🚦 used for 📊 metrics tracking.

#### **`GET /api/healthz`** 💓
- Provides ⚕️ health check.
- ✅ Verifies server is 🔛 & working ⚙️.

---

### Admin 🔒

#### **`GET /admin/metrics`** 📊
- Shows server 📈 metrics.
- Needs admin privileges 🧑‍💼.

#### **`DELETE /admin/reset`** 🔁
- ♻️ Resets request counters for metrics 📏.

---

### Users 👤

#### **`POST /api/users`** 🆕
- Creates 👤.
- Request 📄 includes 📨 email & 🔑 password.

#### **`PUT /api/users`** ✏️
- Updates 👤.
- Needs authentication 🔒.

#### **`POST /api/login`** 🔑
- 👤 Login & receive 🛡️ JWT token.

#### **`POST /api/refresh`** 🔄
- Renews 🛡️ JWT token.
- Needs valid refresh 🔑.

#### **`POST /api/revoke`** ⛔
- Revokes 🔑 token.

---

### Chirps 🐤

#### **`POST /api/chirps`** 🆕🐦
- Adds chirp 🗨️.
- Needs authentication 🔐.

#### **`GET /api/chirps`** 🗃️
- 📜 List chirps.
- 🕵️ Filter/sort supported.

#### **`GET /api/chirps/{chirpID}`** 🔍🐦
- 📜 Specific chirp by ID 🆔.

#### **`DELETE /api/chirps/{chirpID}`** 🗑️🐦
- 🔥 Removes chirp.
- Needs authentication 🔒.

---

### Webhooks 🌊

#### **`POST /api/polka/webhooks`** 📩
- Handles webhook from Polka 🎵.

---

## Setup ⚙️

1. 📥 Clone repo:
   ```bash
   git clone https://github.com/your-username/chirpy.git
   cd chirpy

2. 📚 Install 📦 dependencies:
    ```bash
    go mod tidy
    
3. 🔨 Build:
    ```bash
    go build -o chirpy
4. 🏃 Run:
    ```bash
    ./chirpy
    
## 🤝 Contributing
- Love feedback 💬! Report bugs 🐛 or requests ⭐ via issue 📌 or PR ✍️.
    
