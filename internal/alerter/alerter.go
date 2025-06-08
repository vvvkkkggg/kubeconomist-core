package alerter

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Messenger interface {
	SendMessage(ctx context.Context, text string) error
}

type Alerter struct {
	messenger Messenger
	prefix    string
	timeout   time.Duration
}

type Option func(*Alerter)

func WithTimeout(timeout time.Duration) Option {
	return func(a *Alerter) {
		a.timeout = timeout
	}
}

// DailyLimitMessenger — ограничивает отправку сообщений до 1 раза в сутки
type DailyLimitMessenger struct {
	messenger Messenger
	lastSend  time.Time
	mu        sync.Mutex
}

func NewDailyLimitMessenger(m Messenger) *DailyLimitMessenger {
	return &DailyLimitMessenger{
		messenger: m,
		lastSend:  time.Now().Add(-24*time.Hour - 1*time.Second),
	}
}

func (d *DailyLimitMessenger) SendMessage(ctx context.Context, text string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()
	if now.Sub(d.lastSend) < 24*time.Hour {
		return fmt.Errorf("daily limit exceeded: only one message allowed per 24 hours")
	}

	err := d.messenger.SendMessage(ctx, text)
	if err == nil {
		d.lastSend = now
	}

	return err
}

// WithDailyLimit — добавляет ограничение на отправку: максимум одно сообщение в сутки
func WithDailyLimit() Option {
	return func(a *Alerter) {
		a.messenger = NewDailyLimitMessenger(a.messenger)
	}
}

func New(messenger Messenger, opts ...Option) *Alerter {
	a := &Alerter{
		messenger: messenger,
		timeout:   10 * time.Second,
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

func (a *Alerter) Alert(ctx context.Context, message string) error {
	if a.prefix != "" {
		message = fmt.Sprintf("[%s] %s", a.prefix, message)
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	return a.messenger.SendMessage(ctxWithTimeout, message)
}
