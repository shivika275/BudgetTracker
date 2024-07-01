package main

type UserData struct {
	Username string `json:"username"`
	UserId   string `json:"userId"`
}

type RegisterData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResult struct {
	AccessToken  string `json:"accessToken"`
	IdToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
	TokenType    string `json:"tokenType"`
	UserId       string `json:"userId"`
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
