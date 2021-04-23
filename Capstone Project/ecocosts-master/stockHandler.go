package main

import (
	"net/http"
	"strconv"
	"regexp"
)

func stockHandler(w http.ResponseWriter, r *http.Request) {
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

		symbol := r.FormValue("symbol")
		qty, err := strconv.Atoi(r.FormValue("quantity"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s := stock{
			Symbol: symbol,
			Quantity: qty,
		}

		if err := s.insert(c.id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := c.updateStocks(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := templates.ExecuteTemplate(w, "stock.html", c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func stockDeleteHandler(w http.ResponseWriter, r *http.Request) {
	var c client
	if err := c.readCookie(w, r); err != nil {
		http.Redirect(w, r, "/auth", http.StatusFound)
		return
	}

	validPath := regexp.MustCompile("^/stock/delete/([a-zA-Z0-9]+)$")
	symbol := validPath.FindStringSubmatch(r.URL.Path)
	if symbol == nil {
		http.Error(w, "no symbol set", http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		s := stock{
			Symbol: symbol[1],
		}

		if err := s.purge(c.id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/stock", http.StatusFound)
	}

	data := map[string]interface{}{
		"Client": c,
		"Symbol": symbol[1],
	}

	if err := templates.ExecuteTemplate(w, "stockDelete.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
