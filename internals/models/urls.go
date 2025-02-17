 package models

 import (
	"database/sql"
 )

 type ShortenerData struct {
	OriginalURL, ShortenedURL string
	Clicks			int
 }

 type ShortenerDataModel struct {
 	DB *sql.DB
 }

 func (m *ShortenerDataModel) Latest() ([]*ShortenerData, error) {
	stmt := `SELECT original_url, shortened_url, clicks FROM urls`
	rows, err := m.DB.Query(stmt)
	if err != nil {
	 return nil, err
	}
	defer rows.Close()

	urls := []*ShortenerData{}
	for rows.Next() {
	 url := &ShortenerData{}
	 err := rows.Scan(&url.OriginalURL, &url.ShortenedURL, &url.Clicks)
	 if err != nil {
		return nil, err
	 }
	 urls = append(urls, url)
	}

	if err = rows.Err(); err != nil {
	 return nil, err
	}
	return urls, nil
 }
