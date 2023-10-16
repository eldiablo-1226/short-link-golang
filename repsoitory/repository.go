package repsoitory

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

var connection *pgxpool.Pool

func Init(connectionString string) error {

	conn, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return err
	}

	connection = conn

	return nil
}

func Close() {
	connection.Close()
}

func GetCodeByUrl(url string) (string, error) {
	sql := "select su.code from short_urls as su where su.original_url = $1 limit 1"

	var code string
	err := connection.QueryRow(context.Background(), sql, url).Scan(&code)

	return code, err
}

func GetUrlByCode(code string) (string, error) {
	sql := "select su.original_url from short_urls as su where su.code = $1 limit 1"

	var url string
	err := connection.QueryRow(context.Background(), sql, code).Scan(&url)

	return url, err
}

func InsertShortLink(code string, url string, tag string) error {
	sql := "insert into short_urls (code, original_url, tag) VALUES ($1, $2, $3)"

	_, err := connection.Exec(context.Background(), sql, code, url, tag)

	return err
}
