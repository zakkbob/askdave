module github.com/ZakkBob/AskDave/backend

go 1.24.0

replace github.com/ZakkBob/AskDave/gocommon => ../gocommon

require (
	github.com/ZakkBob/AskDave/gocommon v0.0.0
	github.com/jackc/pgx/v5 v5.7.2
	github.com/pashagolub/pgxmock/v4 v4.5.0
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
