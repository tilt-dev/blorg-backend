package golink

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
)

type Link struct {
	Name    string
	Address string
}

type Golink struct {
	db *sql.DB
}

func NewGolink(db *sql.DB) *Golink {
	return &Golink{
		db: db,
	}
}

func (g *Golink) LinkFromName(name string) (string, error) {
	query := "SELECT links.url FROM links WHERE links.name = $1;"

	rows, err := g.db.Query(query, name)
	if err != nil {
		return "", err
	}

	var url string
	for rows.Next() {
		if err := rows.Scan(&url); err != nil {
			return "", err
		}
		// just get the first result
		break
	}

	return url, nil
}

func (*Golink) WriteLink(payload *Link) error {
	return nil
}

func ParseParams(body io.ReadCloser) (*Link, error) {
	l := Link{}
	buf := new(bytes.Buffer)
	// NOTE(dmiller): this is inefficient. If someone sent us a giant url
	// we would load it in to memory and may become sad :(
	buf.ReadFrom(body)
	s := buf.String()
	err := json.Unmarshal([]byte(s), &l)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func LinkAsJSON(name string, link string) (string, error) {
	l := &Link{
		Address: link,
		Name:    "",
	}
	j, err := json.Marshal(l)
	if err != nil {
		return "", err
	}

	s := string(j)

	return s, nil
}
