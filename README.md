# Forum

A web forum application built with **Go** and **SQLite** that enables users to share ideas, interact through posts and comments, and engage with the community using likes, dislikes, and category-based discussions.

The project implements user authentication, session management, and a clean forum experience without relying on frontend frameworks.

## ✨ Features

* 🔐 User registration and login
* 🍪 Session management using cookies
* 📝 Create and manage posts
* 💬 Comment on posts
* 👍 Like and 👎 dislike posts and comments
* 🏷️ Organize posts with categories
* 🔎 Filter posts by:

  * Categories
  * Your own posts
  * Posts you've liked
* ⚠️ Custom error pages and HTTP status handling
* 🗄️ Persistent data storage with SQLite
* 🐳 Docker support for easy deployment

## 🛠️ Tech Stack

* **Backend:** Go
* **Database:** SQLite
* **Frontend:** HTML, CSS, JavaScript
* **Authentication:** Cookies & Sessions
* **Password Hashing:** bcrypt
* **Containerization:** Docker

## 📁 Project Structure

```text
forum/
├── config/
├── database/
├── handlers/
├── middleware/
├── models/
├── routes/
├── static/
├── views/
├── Dockerfile
├── go.mod
├── go.sum
├── main.go
└── forum.db
```

## 🚀 Getting Started

### Clone the repository

```bash
git clone https://github.com/<your-username>/forum.git
cd forum
```

### Install dependencies

```bash
go mod tidy
```

### Run the application

```bash
go run .
```

The server will be available at:

```text
http://localhost:8080
```

### Run with Docker

```bash
docker build -t forum .
docker run -p 8080:8080 forum
```

## 🔒 Authentication

Registered users can:

* Create posts
* Write comments
* Like and dislike posts and comments
* Filter their own posts
* View posts they've liked

Visitors can browse all public posts and comments without creating an account.

## 🗃️ Database

SQLite is used to store:

* Users
* Sessions
* Posts
* Comments
* Categories
* Likes and dislikes

## 🤝 Contributors

* Kaoutar Hammaoui — https://github.com/KaoutarHammaoui
* Oumaima Talhaoui — https://github.com/oumaimatlh
* Mohamed barrah   — https://github.com/ku-uhaku
* Houssam zmarrou  — https://github.com/frigoDz

## 📄 License

This project was developed for educational purposes and to practice backend web development using Go, SQLite, and Docker.
