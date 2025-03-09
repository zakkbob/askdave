module github.com/ZakkBob/AskDave/backend

go 1.24.0

replace github.com/ZakkBob/AskDave/gocommon => ../gocommon

require github.com/jackc/pgx/v5 v5.7.2

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/text v0.21.0 // indirect
)
