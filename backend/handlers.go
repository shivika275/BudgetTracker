package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var registerData RegisterData
	if err := json.NewDecoder(r.Body).Decode(&registerData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Registering user as %v", registerData)

	err := CreateUserInCognito(registerData)
	if err != nil {
		log.Printf("Registering user err %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = CreateUserEntry(registerData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Registering user complete %v", registerData)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Registration successful"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var loginData LoginData

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Authenticating user as %v", loginData)
	result, err := AuthenticateUser(loginData)
	if err != nil {
		log.Printf("Authenticating user err %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userId, err := GetUserIdByUserName(loginData.Username)
	log.Printf("Got user as %v", userId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	result.UserId = userId

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// Income handlers
func AddIncomeHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var incomeItem IncomeItem
	err := json.NewDecoder(r.Body).Decode(&incomeItem)
	log.Printf("%v, %v", incomeItem, err)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Get the userId from the authenticated user's session
	// For now, we'll use the userId passed in the request

	// Validate the input
	if incomeItem.UserId == "" || incomeItem.IncomeItemName == "" || incomeItem.Month == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Add the income item to the database
	err = AddIncome(incomeItem)
	if err != nil {
		log.Printf("Failed to add income: " + err.Error())
		http.Error(w, "Failed to add income: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Added income")
	// Return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Income added successfully",
	})
}

func GetAllIncomeHandler(w http.ResponseWriter, r *http.Request) {
	// Get userId and month from query parameters
	userId := r.URL.Query().Get("userId")
	monthStr := r.URL.Query().Get("month")

	// TODO: In a real application, get the userId from the authenticated user's session
	// instead of from query parameters

	// Validate the input
	if userId == "" || monthStr == "" {
		http.Error(w, "Missing required query parameters: userId and month", http.StatusBadRequest)
		return
	}

	// Get the income items from the database
	incomeItems, err := GetAllIncome(userId, monthStr)
	if err != nil {
		http.Error(w, "Failed to get income items: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the income items
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(incomeItems)
}

func UpdateIncomeHandler(w http.ResponseWriter, r *http.Request) {
	// Get userId, month, and incomeItemName from URL parameters
	vars := mux.Vars(r)
	userId := vars["userId"]
	monthStr := vars["month"]
	incomeItemName := vars["incomeItemName"]

	// TODO:get the userId from the authenticated user's session
	// instead of from URL parameters

	// Validate the input
	if userId == "" || monthStr == "" || incomeItemName == "" {
		http.Error(w, "Missing required parameters: userId, month, and incomeItemName", http.StatusBadRequest)
		return
	}

	// Parse the request body to get the new value
	var updateRequest struct {
		NewValue float64 `json:"newValue"`
	}
	err := json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update the income item in the database
	err = UpdateIncome(userId, monthStr, incomeItemName, updateRequest.NewValue)
	if err != nil {
		http.Error(w, "Failed to update income item: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Income item updated successfully",
	})
}

func DeleteIncomeHandler(w http.ResponseWriter, r *http.Request) {
	// Get userId, month, and incomeItemName from URL parameters
	vars := mux.Vars(r)
	userId := vars["userId"]
	monthStr := vars["month"]
	incomeItemName := vars["incomeItemName"]

	// TODO: get the userId from the authenticated user's session
	// instead of from URL parameters

	// Validate the input
	if userId == "" || monthStr == "" || incomeItemName == "" {
		http.Error(w, "Missing required parameters: userId, month, and incomeItemName", http.StatusBadRequest)
		return
	}

	// Delete the income item from the database
	err := DeleteIncome(userId, monthStr, incomeItemName)
	if err != nil {
		http.Error(w, "Failed to delete income item: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Income item deleted successfully",
	})
}

// Budget handlers
func AddBudgetHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var budgetItem BudgetItem
	err := json.NewDecoder(r.Body).Decode(&budgetItem)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Get the userId from the authenticated user's session
	// For now, we'll use the userId passed in the request

	// Validate the input
	if budgetItem.UserID == "" || budgetItem.BudgetItemName == "" || budgetItem.Month == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Add the budget item to the database
	err = AddBudget(budgetItem)
	if err != nil {
		http.Error(w, "Failed to add budget: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Budget added successfully",
	})
}

func GetAllBudgetHandler(w http.ResponseWriter, r *http.Request) {
	// Get userId and month from query parameters
	userId := r.URL.Query().Get("userId")
	monthStr := r.URL.Query().Get("month")

	// TODO: get the userId from the authenticated user's session
	// instead of from query parameters

	// Validate the input
	if userId == "" || monthStr == "" {
		http.Error(w, "Missing required query parameters: userId and month", http.StatusBadRequest)
		return
	}

	// Get the budget items from the database
	budgetItems, err := GetAllBudget(userId, monthStr)
	if err != nil {
		http.Error(w, "Failed to get budget items: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the budget items
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(budgetItems)
}

func UpdateBudgetHandler(w http.ResponseWriter, r *http.Request) {
	// Get userId, month, and budgetItemName from URL parameters
	vars := mux.Vars(r)
	userId := vars["userId"]
	monthStr := vars["month"]
	budgetItemName := vars["budgetItemName"]

	// TODO:get the userId from the authenticated user's session
	// instead of from URL parameters

	// Validate the input
	if userId == "" || monthStr == "" || budgetItemName == "" {
		http.Error(w, "Missing required parameters: userId, month, and budgetItemName", http.StatusBadRequest)
		return
	}

	// Parse the request body to get the new value
	var updateRequest struct {
		NewValue float64 `json:"newValue"`
	}
	err := json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update the budget item in the database
	err = UpdateBudget(userId, monthStr, budgetItemName, updateRequest.NewValue)
	if err != nil {
		http.Error(w, "Failed to update budget item: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Budget item updated successfully",
	})
}

func DeleteBudgetHandler(w http.ResponseWriter, r *http.Request) {
	// Get userId, month, and budgetItemName from URL parameters
	vars := mux.Vars(r)
	userId := vars["userId"]
	monthStr := vars["month"]
	budgetItemName := vars["budgetItemName"]

	// TODO:get the userId from the authenticated user's session
	// instead of from URL parameters

	// Validate the input
	if userId == "" || monthStr == "" || budgetItemName == "" {
		http.Error(w, "Missing required parameters: userId, month, and budgetItemName", http.StatusBadRequest)
		return
	}

	// Delete the budget item from the database
	err := DeleteBudget(userId, monthStr, budgetItemName)
	if err != nil {
		http.Error(w, "Failed to delete budget item: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Budget item deleted successfully",
	})
}

func AddExpensesHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var requestBody struct {
		Expenses []ExpenseItem `json:"expenses"`
	}

	// Parse the request body
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if any expenses were provided
	if len(requestBody.Expenses) == 0 {
		http.Error(w, "No expense items provided", http.StatusBadRequest)
		return
	}

	// TODO: Get the userId from the authenticated user's session
	// For now, we'll use the userId passed in the first expense item

	// Validate the input
	for _, item := range requestBody.Expenses {
		if item.UserId == "" || item.ExpenseItemName == "" || item.Month == "" {
			http.Error(w, "Missing required fields in one or more expense items", http.StatusBadRequest)
			return
		}
	}

	// Add the expense items to the database
	err = AddExpenses(requestBody.Expenses)
	if err != nil {
		http.Error(w, "Failed to add expenses: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("%d expense(s) added successfully", len(requestBody.Expenses)),
	})
}

func GetAllExpenseHandler(w http.ResponseWriter, r *http.Request) {
	// Get userId and month from query parameters
	userId := r.URL.Query().Get("userId")
	monthStr := r.URL.Query().Get("month")

	// TODO: get the userId from the authenticated user's session
	// instead of from query parameters

	// Validate the input
	if userId == "" || monthStr == "" {
		http.Error(w, "Missing required query parameters: userId and month", http.StatusBadRequest)
		return
	}

	// Get the expense items from the database
	expenseItems, err := GetAllExpenses(userId, monthStr)
	if err != nil {
		http.Error(w, "Failed to get expense items: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the expense items
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenseItems)
}

func UpdateExpenseHandler(w http.ResponseWriter, r *http.Request) {
	// Get userId, month, and expenseItemName from URL parameters
	vars := mux.Vars(r)
	userId := vars["userId"]
	monthStr := vars["month"]
	expenseItemName := vars["expenseItemName"]

	// TODO: get the userId from the authenticated user's session
	// instead of from URL parameters

	// Validate the input
	if userId == "" || monthStr == "" || expenseItemName == "" {
		http.Error(w, "Missing required parameters: userId, month, and expenseItemName", http.StatusBadRequest)
		return
	}

	// Parse the request body to get the new values
	var updateRequest struct {
		NewValue float64  `json:"newValue"`
		NewTags  []string `json:"newTags"`
	}
	err := json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update the expense item in the database
	err = UpdateExpense(userId, monthStr, expenseItemName, updateRequest.NewValue, updateRequest.NewTags)
	if err != nil {
		http.Error(w, "Failed to update expense item: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Expense item updated successfully",
	})
}

func DeleteExpenseHandler(w http.ResponseWriter, r *http.Request) {
	// Get userId, month, and expenseItemName from URL parameters
	vars := mux.Vars(r)
	userId := vars["userId"]
	monthStr := vars["month"]
	expenseItemName := vars["expenseItemName"]

	// TODO: get the userId from the authenticated user's session
	// instead of from URL parameters

	// Validate the input
	if userId == "" || monthStr == "" || expenseItemName == "" {
		http.Error(w, "Missing required parameters: userId, month, and expenseItemName", http.StatusBadRequest)
		return
	}

	// Delete the expense item from the database
	err := DeleteExpense(userId, monthStr, expenseItemName)
	if err != nil {
		http.Error(w, "Failed to delete expense item: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Expense item deleted successfully",
	})
}
