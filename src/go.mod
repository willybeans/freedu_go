module api/internal/src

go 1.21.3

require (
	api/internal/database v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi v1.5.5
	github.com/lib/pq v1.10.9
)

require (
	github.com/adrium/goheif v0.0.0-20230113233934-ca402e77a786 // indirect
	github.com/fatih/color v1.9.0 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/githubnemo/CompileDaemon v1.4.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.11 // indirect
	github.com/otiai10/gosseract/v2 v2.4.1 // indirect
	github.com/radovskyb/watcher v1.0.7 // indirect
	golang.org/x/sys v0.0.0-20191026070338-33540a1f6037 // indirect
)

replace api/internal/database => ../database
