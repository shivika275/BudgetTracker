package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Return a successful response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {

	//CreateTables()
	r := mux.NewRouter()

	r.HandleFunc("/", HealthCheckHandler).Methods("GET")

	// Auth routes
	r.HandleFunc("/api/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/api/login", LoginHandler).Methods("POST")

	// Protected routes
	api := r.PathPrefix("/api").Subrouter()
	//api.Use(AuthMiddleware)

	// Income routes
	api.HandleFunc("/income", AddIncomeHandler).Methods("POST")
	api.HandleFunc("/income", GetAllIncomeHandler).Methods("GET")
	api.HandleFunc("/income/{userId}/{month}/{incomeItemName}", UpdateIncomeHandler).Methods("PUT")
	api.HandleFunc("/income/{userId}/{month}/{incomeItemName}", DeleteIncomeHandler).Methods("DELETE")

	// Budget routes
	api.HandleFunc("/budget", AddBudgetHandler).Methods("POST")
	api.HandleFunc("/budget", GetAllBudgetHandler).Methods("GET")
	api.HandleFunc("/budget/{userId}/{month}/{budgetItemName}", UpdateBudgetHandler).Methods("PUT")
	api.HandleFunc("/budget/{userId}/{month}/{budgetItemName}", DeleteBudgetHandler).Methods("DELETE")

	// Expense routes
	api.HandleFunc("/expense", AddExpensesHandler).Methods("POST")
	api.HandleFunc("/expense", GetAllExpenseHandler).Methods("GET")
	api.HandleFunc("/expense/{userId}/{month}/{expenseItemName}", UpdateExpenseHandler).Methods("PUT")
	api.HandleFunc("/expense/{userId}/{month}/{expenseItemName}", DeleteExpenseHandler).Methods("DELETE")

	log.Println("Server starting on port 8080...")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Add your frontend URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Origin", "Accept", "*"},
		AllowCredentials: true,
	})

	// Wrap the router with the CORS middleware
	handler := c.Handler(r)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", handler))
}
