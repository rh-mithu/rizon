package mail

import (
	"context"
	"fmt"
	"github.com/rh-mithu/rizon/backend/config"
	"github.com/rh-mithu/rizon/backend/pkg/random"
	"net/smtp"
)

type EmailRepository struct {
	cfg *config.Config
}

func ProvideEmailRepository(cfg *config.Config) *EmailRepository {
	return &EmailRepository{cfg: cfg}
}

func (e *EmailRepository) SendEmail(ctx context.Context, toEmail string) error {
	auth := smtp.PlainAuth("Rizon Auth", e.cfg.SmtpUser, e.cfg.SmtpPassword, e.cfg.SmtpHost)
	token, err := random.Bytes(10)
	if err != nil {
		return err
	}
	message := []byte(fmt.Sprintf(
		"Subject: Your Login Link\r\n\r\nClick here to login:\n%s",
		token,
	))

	return smtp.SendMail(
		e.cfg.SmtpHost+":"+e.cfg.SmtpPort,
		auth,
		e.cfg.SmtpUser,
		[]string{toEmail},
		message,
	)
}
