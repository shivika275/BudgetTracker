package main

type User struct {
	UserID       string `json:"userId"`
	Username     string `json:"username"`
	PasswordHash string `json:"password,omitempty"`
}

type IncomeItem struct {
	UserId          string  `json:"userId"`
	IncomeItemName  string  `json:"incomeItemName"`
	Month           string  `json:"month"`
	IncomeItemValue float64 `json:"incomeItemValue"`
}

type BudgetItem struct {
	UserID          string  `json:"userId"`
	Month           string  `json:"month"`
	BudgetItemName  string  `json:"budgetItemName"`
	BudgetItemValue float64 `json:"budgetItemValue"`
}

type ExpenseItem struct {
	UserId          string   `json:"userId"`
	ExpenseItemName string   `json:"expenseItemName"`
	Month           string   `json:"month"`
	ExpenseValue    float64  `json:"expenseItemValue"`
	ExpenseTags     []string `json:"expenseTags"`
}
