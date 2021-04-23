package main

import (
	"net/http"
	"strconv"
)

func budgetHandler(w http.ResponseWriter, r *http.Request) {
	var c client
	if err := c.readCookie(w, r); err != nil {
		http.Redirect(w, r, "/auth", http.StatusFound)
		return
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		amt, err := strconv.ParseFloat(r.FormValue("amount"), 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cat_id, err := strconv.Atoi(r.FormValue("category"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		b := budget{
			Amount:   amt,
			Category: category{ID: cat_id},
		}

		if err := b.insert(c.id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := c.updateBudgets(); err != nil {
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

	if err := templates.ExecuteTemplate(w, "budget.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
