module github.com/SamW94/chirpy

go 1.23.5

replace github.com/SamW94/chirpy/internal/database => ./internal/database

replace github.com/SamW94/chirpy/internal/auth => ./internal/auth

require (
	github.com/SamW94/chirpy/internal/auth v0.0.0
	github.com/SamW94/chirpy/internal/database v0.0.0
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	golang.org/x/crypto v0.38.0 // indirect
)
