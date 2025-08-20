package dynamostore

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
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
		return fmt.Errorf("dynamo: marshal user: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(user_id)"),
	})

	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) && ae.ErrorCode() == "ConditionalCheckFailedException" {
			return apperrors.ErrUserAlreadyExists
		}
		return fmt.Errorf("dynamo: put item: %w", err)
	}

	logger.Log.Info("user saved successfully in dynamo", "user_id", user.UserId)
	return nil
}

func (r *DynamoProfileRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	logger.Log.Info("querying dynamo for user by email",
		"email", email,
		"table", r.tableName,
	)

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
		var ae smithy.APIError
		if errors.As(err, &ae) {
			return nil, fmt.Errorf("dynamo: query [%s]: %s", ae.ErrorCode(), ae.ErrorMessage())
		}
		return nil, fmt.Errorf("dynamo: query: %w", err)
	}

	if len(out.Items) == 0 {
		return nil, apperrors.ErrUserNotFound
	}

	var user models.User
	if err := attributevalue.UnmarshalMap(out.Items[0], &user); err != nil {
		return nil, fmt.Errorf("dynamo: unmarshal user: %w", err)
	}

	logger.Log.Info("user found in dynamo", "user_id", user.UserId)
	return &user, nil
}
