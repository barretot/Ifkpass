package dynamostore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/barretot/ifkpass/internal/config"
	"github.com/barretot/ifkpass/internal/repo"
	"github.com/barretot/ifkpass/internal/store/dynamostore/models"
)

type DynamoProfileRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoProfileRepository(cfg config.AppConfig) repo.ProfileRepository {
	awsCfg, _ := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(cfg.Region),
	)

	return &DynamoProfileRepository{
		client:    dynamodb.NewFromConfig(awsCfg),
		tableName: cfg.TableName,
	}
}

func (r *DynamoProfileRepository) Save(ctx context.Context, user models.User) error {
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	return err
}

var ErrUserNotFound = errors.New("user not found")

func (r *DynamoProfileRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("email-index"),
		KeyConditionExpression: aws.String("email = :email"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: email},
		},
		Limit: aws.Int32(1),
	})

	if err != nil {
		return nil, fmt.Errorf("query FindByEmail: %w", err)
	}

	if len(out.Items) == 0 {
		return nil, ErrUserNotFound
	}

	var user models.User
	if err := attributevalue.UnmarshalMap(out.Items[0], &user); err != nil {
		return nil, fmt.Errorf("unmarshal user: %w", err)
	}
	return &user, nil
}
