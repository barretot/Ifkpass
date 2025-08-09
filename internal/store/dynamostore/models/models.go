package models

type User struct {
	UserId   string `dynamodbav:"userId"`
	Name     string `dynamodbav:"name"`
	LastName string `dynamodbav:"lastname"`
	Email    string `dynamodbav:"email"`
}
