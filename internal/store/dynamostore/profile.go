package dynamostore

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/barretot/ifkpass/internal/apperrors"
	"github.com/barretot/ifkpass/internal/config"
	"github.com/barretot/ifkpass/internal/logger"
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
		tableName: cfg.ProfilesTableName,
	}
}

func (r *DynamoProfileRepository) Save(ctx context.Context, user models.User) error {
	logger.Log.Info("saving user to dynamo",
		"user_id", user.UserId,
		"email", user.Email,
		"table", r.tableName,
	)

	item, err := attributevalue.MarshalMap(user)

	if err != nil {
		logger.Log.Error("failed to marshal user for dynamo", "err", err)
		return fmt.Errorf("marshal user: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	if err != nil {
		logger.Log.Error("failed to put item in dynamo", "user_id", user.UserId, "err", err)
		return fmt.Errorf("put item: %w", err)
	}

	logger.Log.Info("user saved successfully in dynamo", "user_id", user.UserId)
	return nil
}

func (r *DynamoProfileRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	logger.Log.Info("querying dynamo for user by email",
		"email", email,
		"table", r.tableName,
	)

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
		logger.Log.Error("failed to query dynamo", "email", email, "err", err)
		return nil, fmt.Errorf("query FindByEmail: %w", err)
	}

	if len(out.Items) == 0 {
		logger.Log.Warn("no user found in dynamo with given email", "email", email)
		return nil, apperrors.ErrUserNotFound
	}

	var user models.User
	if err := attributevalue.UnmarshalMap(out.Items[0], &user); err != nil {
		logger.Log.Error("failed to unmarshal dynamo item to user struct", "email", email, "err", err)
		return nil, fmt.Errorf("unmarshal user: %w", err)
	}

	logger.Log.Info("user found in dynamo", "user_id", user.UserId, "email", user.Email)
	return &user, nil
}
