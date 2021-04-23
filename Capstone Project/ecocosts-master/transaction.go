package main

import (
	"database/sql"
	"errors"
	"time"
)

type transaction struct {
	ID          int
	Client_id   int
	Cat_id      sql.NullInt64
	Amount      float64
	Balance     float64
	Description sql.NullString
	Time        time.Time
	Cat_name    string
}

// Inserts values into the database with incremented ID. Check for unset values.
func (t *transaction) insert() error {
	if t.Client_id == 0 {
		return errors.New("Transaction Client_id is not set")
	} else if t.Amount == 0 {
		return errors.New("Transaction Amount is zero")
	} else if t.Time.IsZero() {
		return errors.New("Transaction Time is not set")
	}

	_, err := db.Exec(`
		INSERT INTO transaction
		(client_id, cat_id, amount, balance, description, time) 
		VALUES ($1, $2, $3, $4, $5, $6)`, t.Client_id, t.Cat_id, t.Amount,
		t.Balance, t.Description, t.Time)
	if err != nil {
		return err
	}

	return nil
}

func (t *transaction) updateCategory() error {
	if !t.Cat_id.Valid {
		t.Cat_name = "None"
		return nil
	}

	categories, err := getCategories()
	if err != nil {
		return err
	}
	for _, c := range categories {
		if t.Cat_id.Int64 == int64(c.ID) {
			t.Cat_name = c.Description
		}
	}

	return nil
}
