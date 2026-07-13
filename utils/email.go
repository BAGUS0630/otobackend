package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/smtp"
	"os"
	"time"
)

// EmailVerification struct untuk menyimpan data verifikasi email
type EmailVerification struct {
	Email      string
	Code       string
	ExpiresAt  time.Time
	IsVerified bool
}

// GenerateVerificationCode menghasilkan kode verifikasi 6 digit
func GenerateVerificationCode() string {
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		return "000000"
	}
	code := int(b[0])<<16 | int(b[1])<<8 | int(b[2])
	return fmt.Sprintf("%06d", code%1000000)
}

// GenerateResetToken menghasilkan random token untuk reset password
func GenerateResetToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

// SendVerificationEmail mengirim email verifikasi
func SendVerificationEmail(recipientEmail, verificationCode string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("SMTP_EMAIL")
	senderPassword := os.Getenv("SMTP_PASSWORD")

	// Jika SMTP tidak dikonfigurasi, skip
	if smtpHost == "" {
		return nil
	}

	subject := "Verifikasi Email OtoMeet"
	body := fmt.Sprintf(`
	<html>
		<body style="font-family: Arial, sans-serif;">
			<h2>Verifikasi Email OtoMeet</h2>
			<p>Halo,</p>
			<p>Terima kasih telah mendaftar di OtoMeet. Silakan gunakan kode verifikasi berikut untuk mengkonfirmasi email Anda:</p>
			<h3 style="background-color: #f0f0f0; padding: 10px; text-align: center;">%s</h3>
			<p>Kode ini berlaku selama 15 menit.</p>
			<p>Jika Anda tidak mendaftar di OtoMeet, abaikan email ini.</p>
			<p>Salam,<br/>Tim OtoMeet</p>
		</body>
	</html>
	`, verificationCode)

	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	header := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n",
		senderEmail, recipientEmail, subject)

	fullMessage := header + body

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{recipientEmail}, []byte(fullMessage))
	return err
}

// SendPasswordResetEmail mengirim email reset password
func SendPasswordResetEmail(recipientEmail, resetToken string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("SMTP_EMAIL")
	senderPassword := os.Getenv("SMTP_PASSWORD")

	// Jika SMTP tidak dikonfigurasi, skip
	if smtpHost == "" {
		return nil
	}

	resetLink := fmt.Sprintf("http://localhost:5173/reset-password?token=%s", resetToken)

	subject := "Reset Password OtoMeet"
	body := fmt.Sprintf(`
	<html>
		<body style="font-family: Arial, sans-serif;">
			<h2>Reset Password OtoMeet</h2>
			<p>Halo,</p>
			<p>Kami menerima permintaan untuk mereset password akun Anda. Klik link berikut untuk mereset password:</p>
			<p><a href="%s">Reset Password</a></p>
			<p>Link ini berlaku selama 1 jam.</p>
			<p>Jika Anda tidak meminta reset password, abaikan email ini.</p>
			<p>Salam,<br/>Tim OtoMeet</p>
		</body>
	</html>
	`, resetLink)

	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	header := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n",
		senderEmail, recipientEmail, subject)

	fullMessage := header + body

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{recipientEmail}, []byte(fullMessage))
	return err
}

// SendWelcomeEmail mengirim email sambutan setelah registrasi
func SendWelcomeEmail(recipientEmail, username string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("SMTP_EMAIL")
	senderPassword := os.Getenv("SMTP_PASSWORD")

	// Jika SMTP tidak dikonfigurasi, skip
	if smtpHost == "" {
		return nil
	}

	subject := "Sambutan dari OtoMeet"
	body := fmt.Sprintf(`
	<html>
		<body style="font-family: Arial, sans-serif;">
			<h2>Selamat Datang di OtoMeet</h2>
			<p>Halo %s,</p>
			<p>Akun Anda telah berhasil dibuat. Anda sekarang dapat mengikuti berbagai agenda touring komunitas motor terbesar di Indonesia.</p>
			<p>Fitur yang dapat Anda nikmati:</p>
			<ul>
				<li>Lihat daftar agenda touring</li>
				<li>Daftar sebagai peserta touring</li>
				<li>Diskusi dengan komunitas</li>
				<li>Kelola profil pribadi</li>
			</ul>
			<p>Terima kasih telah bergabung dengan OtoMeet!</p>
			<p>Salam,<br/>Tim OtoMeet</p>
		</body>
	</html>
	`, username)

	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	header := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n",
		senderEmail, recipientEmail, subject)

	fullMessage := header + body

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{recipientEmail}, []byte(fullMessage))
	return err
}
