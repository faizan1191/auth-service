package email

import "fmt"

func ResetPasswordTemplate(resetURL string) string {
	return fmt.Sprintf(`
		<h2>Password Reset</h2>
		<p>You requested a password reset.</p>
		<p>
			<a href="%s">Reset your password</a>
		</p>
		<p>This link expires in 15 minutes.</p>
	`, resetURL)
}
