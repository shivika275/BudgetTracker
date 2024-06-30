package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var db *dynamodb.DynamoDB

func init() {
	// // Initialize a session that the SDK will use to load
	// // credentials from the shared credentials file ~/.aws/credentials
	// // and region from the shared configuration file ~/.aws/config.
	// sess, err := session.NewSession(&aws.Config{
	// 	Region:      aws.String("us-west-2"),
	// 	Endpoint:    aws.String("http://localhost:8000"),
	// 	Credentials: credentials.NewStaticCredentials("dummy", "dummy", "dummy"),
	// })

	// if err != nil {
	// 	log.Fatalf("Failed to create session: %v", err)
	// }

	// // Create DynamoDB client
	// db = dynamodb.New(sess)

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	// Create an STS credentials provider
	//creds := stscreds.NewCredentials(sess, "arn:aws:iam::654654428853:role/BudgetBackendRole")

	log.Printf("Trying to creating new DB session")
	// Create DynamoDB client using the STS credentials
	db = dynamodb.New(sess)
	log.Printf("create DB %v", db)
}

func CreateTables() error {
	tables := []struct {
		Name       string
		Attributes []*dynamodb.AttributeDefinition
		KeySchema  []*dynamodb.KeySchemaElement
	}{
		{
			Name: "Users",
			Attributes: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("userName"),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("userName"),
					KeyType:       aws.String("HASH"),
				},
			},
		},
		{
			Name: "Income",
			Attributes: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("userId#month"),
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String("incomeItemName"),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("userId#month"),
					KeyType:       aws.String("HASH"),
				},
				{
					AttributeName: aws.String("incomeItemName"),
					KeyType:       aws.String("RANGE"),
				},
			},
		},
		{
			Name: "Budget",
			Attributes: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("userId#month"),
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String("budgetItemName"),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("userId#month"),
					KeyType:       aws.String("HASH"),
				},
				{
					AttributeName: aws.String("budgetItemName"),
					KeyType:       aws.String("RANGE"),
				},
			},
		},
		{
			Name: "Expenses",
			Attributes: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("userId#month"),
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String("expenseItemName"),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("userId#month"),
					KeyType:       aws.String("HASH"),
				},
				{
					AttributeName: aws.String("expenseItemName"),
					KeyType:       aws.String("RANGE"),
				},
			},
		},
	}

	for _, table := range tables {

		input := &dynamodb.CreateTableInput{
			AttributeDefinitions: table.Attributes,
			KeySchema:            table.KeySchema,
			BillingMode:          aws.String("PAY_PER_REQUEST"), // This enables on-demand capacity
			TableName:            aws.String(table.Name),
		}

		log.Println("Creating table %v", input)
		_, err := db.CreateTable(input)
		if err != nil {
			log.Println("failed %v", err)
			return fmt.Errorf("failed to create table %s: %v", table.Name, err)
		}
		fmt.Printf("Created table %s\n", table.Name)
	}

	return nil
}

// AddIncome adds a new income item to the Income table
func AddIncome(item IncomeItem) error {
	// Create the composite key
	userIdMonth := fmt.Sprintf("%s#%s", item.UserId, item.Month)

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal Income item: %v", err)
	}

	// Add the composite key to the item
	av["userId#month"] = &dynamodb.AttributeValue{S: aws.String(userIdMonth)}
	log.Printf("Adding : %v", av)
	input := &dynamodb.PutItemInput{
		TableName: aws.String("Income"),
		Item:      av,
	}

	_, err = db.PutItem(input)
	if err != nil {
		return fmt.Errorf("failed to add Income item: %v", err)
	}

	return nil
}

func GetAllIncome(userId string, month string) ([]IncomeItem, error) {
	// Create the composite key
	userIdMonth := fmt.Sprintf("%s#%s", userId, month)

	// Create the query input
	keyCond := expression.Key("userId#month").Equal(expression.Value(userIdMonth))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %v", err)
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String("Income"),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	// Execute the query
	result, err := db.Query(input)
	if err != nil {
		return nil, fmt.Errorf("failed to query Income items: %v", err)
	}

	// Unmarshal the results
	var incomeItems []IncomeItem
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &incomeItems)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal Income items: %v", err)
	}

	return incomeItems, nil
}

func UpdateIncome(userId string, month string, incomeItemName string, newValue float64) error {
	// Create the composite key
	userIdMonth := fmt.Sprintf("%s#%s", userId, month)

	// Create the update expression
	update := expression.Set(expression.Name("incomeItemValue"), expression.Value(newValue))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return fmt.Errorf("failed to build expression: %v", err)
	}

	// Create the update item input
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("Income"),
		Key: map[string]*dynamodb.AttributeValue{
			"userId#month":   {S: aws.String(userIdMonth)},
			"incomeItemName": {S: aws.String(incomeItemName)},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ReturnValues:              aws.String("UPDATED_NEW"),
	}

	// Execute the update
	_, err = db.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("failed to update Income item: %v", err)
	}

	return nil
}

// DeleteIncome removes an income item from the Income table
func DeleteIncome(userId string, month string, incomeItemName string) error {
	// Create the composite key
	userIdMonth := fmt.Sprintf("%s#%s", userId, month)

	// Create the delete item input
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("Income"),
		Key: map[string]*dynamodb.AttributeValue{
			"userId#month":   {S: aws.String(userIdMonth)},
			"incomeItemName": {S: aws.String(incomeItemName)},
		},
	}

	// Execute the delete operation
	_, err := db.DeleteItem(input)
	if err != nil {
		return fmt.Errorf("failed to delete Income item: %v", err)
	}

	return nil
}

// AddBudget adds a new budget item to the Budget table
func AddBudget(item BudgetItem) error {
	// Create the composite key
	userIdMonth := fmt.Sprintf("%s#%s", item.UserID, item.Month)

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal Budget item: %v", err)
	}

	// Add the composite key to the item
	av["userId#month"] = &dynamodb.AttributeValue{S: aws.String(userIdMonth)}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Budget"),
		Item:      av,
	}

	_, err = db.PutItem(input)
	if err != nil {
		return fmt.Errorf("failed to add Budget item: %v", err)
	}

	return nil
}

func GetAllBudget(userId string, month string) ([]BudgetItem, error) {
	// Create the composite key
	userIdMonth := fmt.Sprintf("%s#%s", userId, month)

	// Create the query input
	keyCond := expression.Key("userId#month").Equal(expression.Value(userIdMonth))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %v", err)
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String("Budget"),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	// Execute the query
	result, err := db.Query(input)
	if err != nil {
		return nil, fmt.Errorf("failed to query Budget items: %v", err)
	}

	// Unmarshal the results
	var budgetItems []BudgetItem
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &budgetItems)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal Budget items: %v", err)
	}

	return budgetItems, nil
}

func UpdateBudget(userId string, month string, budgetItemName string, newValue float64) error {
	// Create the composite key
	userIdMonth := fmt.Sprintf("%s#%s", userId, month)

	// Create the update expression
	update := expression.Set(expression.Name("budgetItemValue"), expression.Value(newValue))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return fmt.Errorf("failed to build expression: %v", err)
	}

	// Create the update item input
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("Budget"),
		Key: map[string]*dynamodb.AttributeValue{
			"userId#month":   {S: aws.String(userIdMonth)},
			"budgetItemName": {S: aws.String(budgetItemName)},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ReturnValues:              aws.String("UPDATED_NEW"),
	}

	// Execute the update
	_, err = db.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("failed to update Budget item: %v", err)
	}

	return nil
}

func DeleteBudget(userId string, month string, budgetItemName string) error {
	// Create the composite key
	userIdMonth := fmt.Sprintf("%s#%s", userId, month)

	// Create the delete item input
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("Budget"),
		Key: map[string]*dynamodb.AttributeValue{
			"userId#month":   {S: aws.String(userIdMonth)},
			"budgetItemName": {S: aws.String(budgetItemName)},
		},
	}

	// Execute the delete operation
	_, err := db.DeleteItem(input)
	if err != nil {
		return fmt.Errorf("failed to delete Budget item: %v", err)
	}

	return nil
}

func AddExpenses(items []ExpenseItem) error {
	// Create a list to hold the write requests
	var writeRequests []*dynamodb.WriteRequest

	for _, item := range items {
		// Create the composite key
		userIdMonth := fmt.Sprintf("%s#%s", item.UserId, item.Month)

		av, err := dynamodbattribute.MarshalMap(item)
		if err != nil {
			return fmt.Errorf("failed to marshal Expense item: %v", err)
		}

		// Add the composite key to the item
		av["userId#month"] = &dynamodb.AttributeValue{S: aws.String(userIdMonth)}

		// Create a PutRequest for each item
		writeRequests = append(writeRequests, &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: av,
			},
		})
	}

	// Split the write requests into batches of 25 (DynamoDB limit)
	for i := 0; i < len(writeRequests); i += 25 {
		end := i + 25
		if end > len(writeRequests) {
			end = len(writeRequests)
		}

		input := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]*dynamodb.WriteRequest{
				"Expenses": writeRequests[i:end],
			},
		}

		_, err := db.BatchWriteItem(input)
		if err != nil {
			return fmt.Errorf("failed to add Expense items: %v", err)
		}
	}

	return nil
}

func GetAllExpenses(userId string, month string) ([]ExpenseItem, error) {
	// Create the composite key
	userIdMonth := fmt.Sprintf("%s#%s", userId, month)

	// Create the query input
	keyCond := expression.Key("userId#month").Equal(expression.Value(userIdMonth))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %v", err)
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String("Expenses"),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	// Execute the query
	result, err := db.Query(input)
	if err != nil {
		return nil, fmt.Errorf("failed to query Expense items: %v", err)
	}

	// Unmarshal the results
	var expenseItems []ExpenseItem
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &expenseItems)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal Expense items: %v", err)
	}

	return expenseItems, nil
}

func UpdateExpense(userId string, month string, expenseItemName string, newValue float64, newTags []string) error {
	// Create the composite key
	userIdMonth := fmt.Sprintf("%s#%s", userId, month)

	// Create the update expression
	update := expression.
		Set(expression.Name("expenseItemValue"), expression.Value(newValue)).
		Set(expression.Name("expenseTags"), expression.Value(newTags))

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return fmt.Errorf("failed to build expression: %v", err)
	}

	// Create the update item input
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("Expenses"),
		Key: map[string]*dynamodb.AttributeValue{
			"userId#month":    {S: aws.String(userIdMonth)},
			"expenseItemName": {S: aws.String(expenseItemName)},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ReturnValues:              aws.String("UPDATED_NEW"),
	}

	// Execute the update
	_, err = db.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("failed to update Expense item: %v", err)
	}

	return nil
}

// DeleteExpense removes an expense item from the Expenses table
func DeleteExpense(userId string, month string, expenseItemName string) error {
	// Create the composite key
	userIdMonth := fmt.Sprintf("%s#%s", userId, month)

	// Create the delete item input
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("Expenses"),
		Key: map[string]*dynamodb.AttributeValue{
			"userId#month":    {S: aws.String(userIdMonth)},
			"expenseItemName": {S: aws.String(expenseItemName)},
		},
	}

	// Execute the delete operation
	_, err := db.DeleteItem(input)
	if err != nil {
		return fmt.Errorf("failed to delete Expense item: %v", err)
	}

	return nil
}
