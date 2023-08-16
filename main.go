package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	qf "github.com/tmw/queryfilter"
)

type Size string

const (
	SizeS  Size = "S"
	SizeM       = "M"
	SizeL       = "L"
	SizeXL      = "XL"
)

type Filter struct {
	Colors   []string `filter:"color,op=in"`
	Sizes    []Size   `filter:"size,op=in"`
	PriceMin int      `filter:"price,op=gte"`
	PriceMax int      `filter:"price,op=lte"`
}

func setupDatabase() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		return nil, fmt.Errorf("error setting up database: %w", err)
	}

	// setup minimal table
	db.MustExec(`
		CREATE TABLE tshirts (
		  id INTEGER PRIMARY KEY,
		  color TEXT,
		  size  TEXT,
		  price INT
		);
	`)

	// seed some tshirts
	db.MustExec(`
		INSERT INTO tshirts (color, size, price) VALUES
		("red", "S", 12),
		("red", "L", 12),
		("red", "XL", 12),
		("yellow", "L", 10),
		("yellow", "XL", 12),
		("yellow", "M", 8),
		("pink", "S", 25),
		("pink", "M", 30),
		("blue", "L", 17),
		("green", "M", 8),
		("green", "L", 9),
		("navy", "XL", 22)
	`)

	return db, nil
}

type Tshirt struct {
	ID    int    `db:"id" json:"id"`
	Color string `db:"color" json:"color"`
	Size  Size   `db:"size"  json:"size"`
	Price int    `db:"price" json:"price"`
}

func main() {
	db, err := setupDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/shirts", handleGetShirts(db))

	log.Println("Server listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleGetShirts(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sizes, colors, minPrice, maxPrice := parseQueryString(r.URL.Query())

		// setting up the filter
		f := Filter{
			PriceMin: minPrice,
			PriceMax: maxPrice,
			Sizes:    sizes,
			Colors:   colors,
		}

		sql, args, err := qf.ToSQL(f)
		if err != nil {
			log.Fatal(err)
		}

		sql = fmt.Sprintf("SELECT * FROM tshirts WHERE %s", sql)

		var results []Tshirt
		err = db.Select(&results, sql, args...)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Add("content-type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}

func parseQueryString(query url.Values) ([]Size, []string, int, int) {
	sizes := make([]Size, len(query["size"]))
	for _, size := range query["size"] {
		sizes = append(sizes, Size(size))
	}

	minPrice := 0
	if price := query.Get("min_price"); price != "" {
		if price, err := strconv.Atoi(price); err == nil {
			minPrice = price
		}
	}

	maxPrice := math.MaxInt
	if price := query.Get("max_price"); price != "" {
		if price, err := strconv.Atoi(price); err == nil {
			maxPrice = price
		}
	}

	return sizes, query["color"], minPrice, maxPrice
}
