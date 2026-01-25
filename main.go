package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Desc string `json:"description"`
}

var category = []Category{
	{ID: 1, Name: "Horror", Desc: "Movies about ghost"},
	{ID: 2, Name: "Comedy", Desc: "Movies that brings laughter"},
	{ID: 3, Name: "Drama", Desc: "Movies that shake your heart"},
}

func addNewCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory Category
	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	newCategory.ID = len(category) + 1
	category = append(category, newCategory)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "Add Category Successful",
	})
}

func findCategoryById(id int, w http.ResponseWriter) {
	for _, p := range category {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "Category Not Found", http.StatusNotFound)
}

func updateCategory(id int, w http.ResponseWriter, r *http.Request) {
	var updateCategory Category
	err := json.NewDecoder(r.Body).Decode(&updateCategory)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	for i := range category {
		if category[i].ID == id {
			updateCategory.ID = id
			category[i] = updateCategory
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"status": "Update Successful",
			})
			return
		}
	}

	http.Error(w, "Category Not Found", http.StatusNotFound)
}

func deleteCategory(id int, w http.ResponseWriter) {
	for i, c := range category {
		if c.ID == id {
			category = append(category[:i], category[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"status": "Delete Successful",
			})
			return
		}
	}

	http.Error(w, "Category Not Found", http.StatusNotFound)
}

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(category)
		case "POST":
			addNewCategory(w, r)
		}
	})

	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case "GET":
			findCategoryById(id, w)
		case "PUT":
			updateCategory(id, w, r)
		case "DELETE":
			deleteCategory(id, w)
		}
	})

	fmt.Println("Server running di localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Failed running server")
	}
}
