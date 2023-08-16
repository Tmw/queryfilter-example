package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func TestFilters(t *testing.T) {
	db, err := setupDatabase()
	if err != nil {
		t.Errorf("error while setting up database: %s", err)
	}

	defer db.Close()

	// make mock server, passing in the database connection
	s := httptest.NewServer(handleGetShirts(db))
	e := httpexpect.Default(t, s.URL)

	body := e.GET("/shirts").
		WithQuery("price_min", "5").
		WithQuery("price_max", "100").
		WithQuery("color", "yellow").
		WithQuery("size", "L").
		Expect().
		Status(http.StatusOK).
		JSON()

	body.Path("$[*].id").Array().IsEqual([]int{4})
}
