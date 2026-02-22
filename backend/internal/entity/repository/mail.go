package repository

import "context"

type EmailRepository interface {
	SendEmail(ctx context.Context, email string) error
}
