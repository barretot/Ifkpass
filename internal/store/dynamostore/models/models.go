package models

type User struct {
	ID    string `dynamodbav:"id"`
	Name  string `dynamodbav:"name"`
	Email string `dynamodbav:"email"`
}
