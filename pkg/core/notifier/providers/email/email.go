package email

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/wneessen/go-mail"

	"github.com/certimate-go/certimate/pkg/core/notifier"
	xtls "github.com/certimate-go/certimate/pkg/utils/tls"
)

type NotifierConfig struct {
	// SMTP 服务器地址。
	SmtpHost string `json:"smtpHost"`
	// SMTP 服务器端口。
	// 零值时根据是否启用 TLS 决定。
	SmtpPort int32 `json:"smtpPort"`
	// 是否启用 TLS。
	SmtpTls bool `json:"smtpTls"`
	// 用户名。
	Username string `json:"username"`
	// 密码。
	Password string `json:"password"`
	// 发件人邮箱。
	SenderAddress string `json:"senderAddress"`
	// 发件人显示名称。
	SenderName string `json:"senderName,omitempty"`
	// 收件人邮箱。
	ReceiverAddress string `json:"receiverAddress"`
	// 是否允许不安全的连接。
	AllowInsecureConnections bool `json:"allowInsecureConnections,omitempty"`
}

type Notifier struct {
	config *NotifierConfig
	logger *slog.Logger
}

var _ notifier.Provider = (*Notifier)(nil)

func NewNotifier(config *NotifierConfig) (*Notifier, error) {
	if config == nil {
		return nil, errors.New("the configuration of the notifier provider is nil")
	}

	return &Notifier{
		config: config,
		logger: slog.Default(),
	}, nil
}

func (n *Notifier) SetLogger(logger *slog.Logger) {
	if logger == nil {
		n.logger = slog.New(slog.DiscardHandler)
	} else {
		n.logger = logger
	}
}

func (n *Notifier) Notify(ctx context.Context, subject string, message string) (*notifier.NotifyResult, error) {
	clientOptions := []mail.Option{
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithUsername(n.config.Username),
		mail.WithPassword(n.config.Password),
	}

	if n.config.SmtpPort == 0 {
		if n.config.SmtpTls {
			clientOptions = append(clientOptions, mail.WithPort(mail.DefaultPortTLS))
		} else {
			clientOptions = append(clientOptions, mail.WithPort(mail.DefaultPort))
		}
	} else {
		clientOptions = append(clientOptions, mail.WithPort(int(n.config.SmtpPort)))
	}

	if n.config.SmtpTls {
		tlsConfig := xtls.NewCompatibleConfig()
		if n.config.AllowInsecureConnections {
			tlsConfig.InsecureSkipVerify = true
		} else {
			tlsConfig.ServerName = n.config.SmtpHost
		}

		mail.WithSSL()
		mail.WithSSLPort(true)
		mail.WithTLSConfig(tlsConfig)
	} else {
		mail.WithTLSPolicy(mail.TLSOpportunistic)
	}

	client, err := mail.NewClient(n.config.SmtpHost, clientOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to create smtp client: %w", err)
	}

	client.ErrorHandlerRegistry.RegisterHandler("smtp.qq.com", "QUIT", &wQQMailQuitErrorHandler{})
	defer client.Close()

	msg := mail.NewMsg()
	msg.Subject(subject)
	msg.SetBodyString(mail.TypeTextPlain, message)
	if n.config.SenderName == "" {
		msg.From(n.config.SenderAddress)
	} else {
		msg.FromFormat(n.config.SenderName, n.config.SenderAddress)
	}
	msg.To(n.config.ReceiverAddress)

	if err := client.DialAndSend(msg); err != nil {
		return nil, fmt.Errorf("failed to send mail: %w", err)
	}

	return &notifier.NotifyResult{}, nil
}
