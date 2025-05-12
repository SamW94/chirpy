module github.com/SamW94/chirpy

go 1.23.5

replace github.com/SamW94/chirpy/internal/database => ./internal/database

require (
	github.com/SamW94/chirpy/internal/database v0.0.0
	github.com/google/uuid v1.6.0 // indirect
	github.com/joho/godotenv v1.5.1 
	github.com/lib/pq v1.10.9
)
