package security

import (
	"context"
	"errors"
	"fmt"
	"time"

	"math/rand"

	"github.com/abdulazizax/yelp/config"
	"github.com/abdulazizax/yelp/internal/adapter/redis/cache"
	gomail "gopkg.in/gomail.v2"
)

type emailService struct {
	cache cache.RedisCache
	cfg   *config.Config
}

type EmailService interface {
	VerifyEmail(ctx context.Context, email string, code int) error
	SendVerificationCode(ctx context.Context, email string) (time.Duration, error)
	GenerateVerificationCode() int
}

func NewEmailService(cache cache.RedisCache, cfg *config.Config) EmailService {
	return &emailService{
		cache: cache,
		cfg:   cfg,
	}
}

func (s *emailService) GenerateVerificationCode() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(90000) + 10000
}

func (s *emailService) SendVerificationCode(ctx context.Context, email string) (time.Duration, error) {
	code := s.GenerateVerificationCode()

	m := gomail.NewMessage()
	m.SetHeader("From", "Yelp")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Your Verification Code")
	m.SetBody("text/plain", fmt.Sprintf("Your verification code is: %d", code))

	d := gomail.NewDialer(s.cfg.Email.SmtpHost, s.cfg.Email.SmtpPort, s.cfg.Email.SmtpUser, s.cfg.Email.SmtpPass)

	if err := d.DialAndSend(m); err != nil {
		return 0, err
	}

	err := s.cache.StoreEmailAndCode(ctx, email, code, 2)
	if err != nil {
		return 0, err
	}

	return 2 * time.Minute, nil
}
func (s *emailService) VerifyEmail(ctx context.Context, email string, code int) error {
	c, err := s.cache.GetCodeByEmail(ctx, email)
	if err != nil {
		return err
	}
	if c != code {
		return errors.New("invalide code")
	}
	return nil
}
