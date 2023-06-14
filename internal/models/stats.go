package models

type Stats struct {
	CountUsers int `db:"count_users" json:"users"`
	CountURLs  int `db:"count_urls" json:"urls"`
}
