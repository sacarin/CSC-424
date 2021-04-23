package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type stock struct {
	Symbol   string
	Quantity int
	Price    quote
}

type quote struct {
	CurrPrice float64
	PrevClose float64
}

func (s *stock) getPrice() error {
	if err := s.Price.update(s.Symbol); err != nil {
		return err
	}

	return nil
}

// this overrides the current stock in the database; ideally you would probably
// want to add to the current quantity of stock present. no time.
func (s *stock) insert(client_id int) error {
	exist, err := rowExist(`SELECT quantity FROM stock WHERE client_id = $1
	AND symbol = $2`, client_id, s.Symbol)
	if exist {
		s.purge(client_id)
	}
	if err != nil {
		return err
	}

	// implement better error checking in production to see if stock exist.
	if err := s.getPrice(); err != nil {
		return fmt.Errorf("Stock exist?: %v", err)
	}

	_, err = db.Exec(`INSERT INTO stock (client_id, symbol, quantity)
		VALUES ($1, $2, $3)`, client_id, s.Symbol, s.Quantity)
	if err != nil {
		return err
	}

	return nil
}

func (s *stock) purge(client_id int) error {
	_, err := db.Exec("DELETE FROM stock WHERE client_id = $1 AND symbol = $2",
		client_id, s.Symbol)
	if err != nil {
		return err
	}
	return nil
}

func (q *quote) fetch(url string) ([]byte, error) {
	// http client with timeout
	client := http.Client{Timeout: 10 * time.Second}

	// fetch URL
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// convert to byte slice
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (q *quote) update(symbol string) error {
	resp, err := q.fetch("https://query2.finance.yahoo.com/v10/finance/quoteSummary/" +
		symbol + "?formatted=false&modules=price")
	if err != nil {
		return err
	}

	// get current stock price
	re := regexp.MustCompile(`regularMarketPrice\":[0-9]*\.[0-9]+`)
	currPrice := string(re.Find(resp))
	currPrice = strings.TrimPrefix(currPrice, "regularMarketPrice\":")
	q.CurrPrice, err = strconv.ParseFloat(currPrice, 64)
	if err != nil {
		return err
	}

	// get previous close
	re = regexp.MustCompile(`regularMarketPreviousClose\":[0-9]*\.[0-9]+`)
	prevClose := string(re.Find(resp))
	prevClose = strings.TrimPrefix(prevClose, "regularMarketPreviousClose\":")
	q.PrevClose, err = strconv.ParseFloat(prevClose, 64)
	if err != nil {
		return err
	}

	return nil
}
