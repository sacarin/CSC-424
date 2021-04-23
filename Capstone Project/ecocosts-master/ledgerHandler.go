package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"
)

func ledgerHandler(w http.ResponseWriter, r *http.Request) {
	var c client
	if err := c.readCookie(w, r); err != nil {
		http.Redirect(w, r, "/auth", http.StatusFound)
		return
	}

	if r.Method == "POST" {
		var t transaction

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		amt, err := strconv.ParseFloat(r.FormValue("amount"), 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Amount = amt

		desc := r.FormValue("description")
		if desc == "" {
			t.Description = sql.NullString{"", false}
		} else {
			t.Description = sql.NullString{desc, true}
		}

		category, err := strconv.ParseInt(r.FormValue("category"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if category == 0 {
			t.Cat_id = sql.NullInt64{0, false}
		} else {
			t.Cat_id = sql.NullInt64{category, true}
		}

		t.Client_id = c.id
		t.Time = time.Now()
		if err := t.insert(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := c.updateTransactions(1000); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categories, err := getCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Client":     c,
		"Categories": categories,
	}

	if err := templates.ExecuteTemplate(w, "ledger.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
