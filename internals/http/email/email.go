package email

import (
	"fmt"

	"github.com/omniflare/go_starter/internals/config"
	"github.com/resendlabs/resend-go"
)

var client *resend.Client

func init() {
	apiKey := config.NewEnvConfig().ResendAPI
	if apiKey == "" {
		panic("RESEND_API_KEY is not set in environment variables")
	}
	client = resend.NewClient(apiKey)
}

const (
	verificationEmailTemplate = `
<!DOCTYPE html>
<html>
<body>
    <h1>Verify Your Email Address</h1>
    <p>Your verification code is:</p>
    <h2 style="font-size: 32px; letter-spacing: 5px; text-align: center; padding: 20px; background: #f5f5f5; border-radius: 8px;">%s</h2>
    <p>This code will expire in 24 hours.</p>
    <p>If you didn't request this code, please ignore this email.</p>
</body>
</html>
`

	welcomeEmailTemplate = `
    <!DOCTYPE html>
    <html>
    <body>
        <h1>Welcome to our platform, %s!</h1>
        <p>We're excited to have you on board.</p>
        <p>Start exploring our services and let us know if you need any help!</p>
    </body>
    </html>
    `
)

func SendVerificationEmail(email, verificationCode string) error {
	params := &resend.SendEmailRequest{
		From:    "Acme <onboarding@resend.dev>",
		To:      []string{email},
		Subject: "Verify Your Email Address Now",
		Html:    fmt.Sprintf(verificationEmailTemplate, verificationCode),
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

func SendWelcomeEmail(email, name string) error {
	params := &resend.SendEmailRequest{
		From:    "Acme <onboarding@resend.dev>",
		To:      []string{email},
		Subject: "Welcome to our company",
		Html:    fmt.Sprintf(welcomeEmailTemplate, name),
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send welcome email: %w", err)
	}

	return nil
}

func SendPasswordResetEmail(email, resetURL string) error {
	params := &resend.SendEmailRequest{
		From:    "Acme <onboarding@resend.dev>",
		To:      []string{email},
		Subject: "Reset Your Password",
		Html: fmt.Sprintf(`
            <!DOCTYPE html>
            <html>
            <body>
                <h1>Password Reset Request</h1>
                <p>Click <a href="%s">here</a> to reset your password</p>
                <p>If you didn't request this, please ignore this email.</p>
            </body>
            </html>
        `, resetURL),
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	return nil
}

func SendResetSuccessEmail(email string) error {
	params := &resend.SendEmailRequest{
		From:    "Acme <onboarding@resend.dev>",
		To:      []string{email},
		Subject: "Password Reset Was Successful",
		Html: `
            <!DOCTYPE html>
            <html>
            <body>
                <h1>Password Reset Successful</h1>
                <p>Your password was reset successfully.</p>
                <p>If you didn't make this change, please contact support immediately.</p>
            </body>
            </html>
        `,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send reset success email: %w", err)
	}

	return nil
}
