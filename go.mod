module github.com/George-Anagnostou/countries

go 1.24

require (
	github.com/gorilla/sessions v1.3.0
	github.com/joho/godotenv v1.5.1
	github.com/labstack/echo-contrib v0.17.1
	github.com/labstack/echo/v4 v4.12.0
	github.com/mattn/go-sqlite3 v1.14.22
	golang.org/x/crypto v0.31.0
)

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/gorilla/context v1.1.2 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/time v0.5.0 // indirect
)

replace github.com/George-Anagnostou/countries/internal/routes => ./internal/routes

replace github.com/George-Anagnostou/countries/internal/utils => ./internal/utils

replace github.com/George-Anagnostou/countries/internal/templates => ./internal/templates
