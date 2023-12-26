module api/internal/src

go 1.21.3

require (
	api/internal/database v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi v1.5.5
	github.com/lib/pq v1.10.9
)

require (
	github.com/adrium/goheif v0.0.0-20230113233934-ca402e77a786 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/otiai10/gosseract/v2 v2.4.1 // indirect
)

replace api/internal/database => ../database
