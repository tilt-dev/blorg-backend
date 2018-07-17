package golink_test

import (
	"bytes"
	"database/sql"
	"io"
	"testing"

	"github.com/windmilleng/blorg-backend/golink"

	_ "github.com/proullon/ramsql/driver"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func TestLinkFromName(t *testing.T) {
	batch := []string{
		`CREATE TABLE links (name TEXT PRIMARY KEY, url TEXT);`,
		`INSERT INTO links (name, url) VALUES ('cat', 'meow.com');`,
	}

	db, err := sql.Open("ramsql", "TestLinkFromName")
	if err != nil {
		t.Fatalf("sql.Open : Error : %s\n", err)
	}
	defer db.Close()

	for _, b := range batch {
		_, err = db.Exec(b)
		if err != nil {
			t.Fatalf("sql.Exec: Error: %s\n", err)
		}
	}

	g := golink.NewGolink(db)

	// Test Case: link doesn't exist
	l, err := g.LinkFromName("aqdfhlkjasdf")
	if err != nil {
		t.Fatal(err)
	}
	if l != "" {
		t.Errorf("Expected empty string, got %s", l)
	}

	// Test Case: link exists
	l, err = g.LinkFromName("cat")
	if err != nil {
		t.Fatal(err)
	}
	if l != "meow.com" {
		t.Errorf("Expected meow.com, got %s", l)
	}
}

func TestParseParams(t *testing.T) {
	s := `{"name": "foo", "address": "http://neopets.com"}`

	body := nopCloser{bytes.NewBufferString(s)}
	expected := golink.Link{
		Name:    "foo",
		Address: "http://neopets.com",
	}

	actual, err := golink.ParseParams(body)
	if err != nil {
		t.Fatal(err)
	}

	if expected.Name != actual.Name || expected.Address != actual.Address {
		t.Errorf("Expected (%+v) does not equal actual (%+v)", expected, actual)
	}
}

// TODO add fail to parse test case

func TestWriteLink(t *testing.T) {
	batch := []string{
		`CREATE TABLE links (name TEXT PRIMARY KEY, url TEXT);`,
	}

	db, err := sql.Open("ramsql", "TestWriteLink")
	if err != nil {
		t.Fatalf("sql.Open : Error : %s\n", err)
	}
	defer db.Close()

	for _, b := range batch {
		_, err = db.Exec(b)
		if err != nil {
			t.Fatalf("sql.Exec: Error: %s\n", err)
		}
	}

	g := golink.NewGolink(db)

	// test case: write valid link
	p := &golink.Link{
		Address: "http://meow.com",
		Name:    "cat",
	}
	g.WriteLink(p)
	l, err := g.LinkFromName("cat")
	if err != nil {
		t.Fatal(err)
	}
	if l != "http://meow.com" {
		t.Errorf("Expected meow.com, got %s", l)
	}

}
