module api/internal/handlers

go 1.21.3

require (
	api/internal/database v0.0.0-00010101000000-000000000000
	github.com/go-rod/rod v0.114.5
)

require (
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/ysmood/fetchup v0.2.3 // indirect
	github.com/ysmood/goob v0.4.0 // indirect
	github.com/ysmood/got v0.34.1 // indirect
	github.com/ysmood/gson v0.7.3 // indirect
	github.com/ysmood/leakless v0.8.0 // indirect
)

replace api/internal/database => ../database
