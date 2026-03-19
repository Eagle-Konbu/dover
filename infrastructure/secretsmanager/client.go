package secretsmanager

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type Client struct {
	client *secretsmanager.Client
}

func New(ctx context.Context) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("secretsmanager: load aws config: %w", err)
	}
	return &Client{client: secretsmanager.NewFromConfig(cfg)}, nil
}

func (c *Client) GetSecretValue(ctx context.Context, name string) (string, error) {
	out, err := c.client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &name,
	})
	if err != nil {
		return "", fmt.Errorf("secretsmanager: get %q: %w", name, err)
	}
	if out.SecretString == nil {
		return "", fmt.Errorf("secretsmanager: %q has no string value", name)
	}
	return *out.SecretString, nil
}

func (c *Client) PutSecretValue(ctx context.Context, name string, value string) error {
	_, err := c.client.PutSecretValue(ctx, &secretsmanager.PutSecretValueInput{
		SecretId:     &name,
		SecretString: &value,
	})
	if err != nil {
		return fmt.Errorf("secretsmanager: put %q: %w", name, err)
	}
	return nil
}
