# ✏️ Welcome to Chirpy!

This is Chirpy, a fully functioning back end social media platform where users can create and edit accounts, make chirps (text based messages) and see other users chirps!

# ❓ But Why?

Chirpy was an important part of my growth and learning as it helped me understand how real world applications and web servers functioned. I build my own production style HTTP server in Go with no framework, using JSON, headers, and Status codes (way more fun than i expected) to communicate with clients via a RESTful API. I built a system for authentication and authorization for the users, and imitated the funciton of a webhook with API keys. On top of that, all user information was securely placed in a Postgres database.

# ⚙️ Installation

Prerequisites:
1. Go toolchain 1.22+ installed
2. PostgresSQL Installed and Running
3. Dependencies are managed with Go modules: run
   go mod tidy
in terminal and it should fetch them.

Inside a go module:
go get github.com/Rhyster42/Chirpy

How to run:
1. set env vars - (DB_HOST, DB_USER, etc.)
2. run
  go run ./cmd/server
or wherever the entry point you put is.

*You can also run go build then run the binary

3.Have fun!
