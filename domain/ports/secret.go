package ports

import "context"

type SecretStore interface {
	GetSecretValue(ctx context.Context, name string) (string, error)
	PutSecretValue(ctx context.Context, name string, value string) error
}
