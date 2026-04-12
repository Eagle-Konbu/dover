package ports

import "context"

type SecretStore interface {
	GetSecretValue(ctx context.Context, name string) (string, error)
}
