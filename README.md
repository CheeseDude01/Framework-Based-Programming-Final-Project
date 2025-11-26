# ThriftShop - Refactored

This is a refactored starter for a Thrift Shop web app using:
- Go + Gin
- GORM (SQLite)
- Simple user registration + login (bcrypt)
- Item CRUD + Buy flow (creates orders)

How to run:
1. go mod tidy
2. go run .
3. Open http://localhost:8080

Notes:
- Currently uses plain responses; add JWT middleware for real auth.
- Database file thriftshop.db will be created in project root.
