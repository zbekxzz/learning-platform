package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(fullName, email, password string) error {

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &User{
		FullName:     fullName,
		Email:        email,
		PasswordHash: string(hash),
		Role:         "student",
	}

	return CreateUser(user)
}

func Login(email, password string) (string, error) {

	user, err := GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secret))
}

func RequestRestoreCode(email string) error {
	// Verify user exists
	_, err := GetUserByEmail(email)
	if err != nil {
		return errors.New("пользователь с таким email не найден")
	}

	// Generate 6-digit random code
	codeBytes := make([]byte, 6)
	if _, err := rand.Read(codeBytes); err != nil {
		return errors.New("ошибка генерации кода")
	}
	for i := range codeBytes {
		codeBytes[i] = '0' + (codeBytes[i] % 10)
	}
	code := string(codeBytes)

	expiresAt := time.Now().UTC().Add(10 * time.Minute)

	err = SaveRestoreCode(email, code, expiresAt)
	if err != nil {
		return err
	}

	log.Printf("[RESTORE CODE] Email: %s, Code: %s", email, code)

	// Send code via email
	err = sendEmail(email, code)
	if err != nil {
		return errors.New("ошибка при отправке письма: " + err.Error())
	}

	return nil
}

func sendEmail(to string, code string) error {
	from := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	if from == "" || pass == "" || host == "" || port == "" {
		log.Println("WARNING: SMTP credentials not fully configured. Code logged only.")
		return nil
	}

	subject := "Subject: Код восстановления пароля - EduPlatform\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`<html>
<body style="font-family: Arial, sans-serif; background-color: #f3f4f6; padding: 20px; color: #1f2937;">
  <div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; padding: 30px; border-radius: 8px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
    <h2 style="color: #2563eb; margin-bottom: 20px;">Восстановление пароля</h2>
    <p>Здравствуйте!</p>
    <p>Ваш код подтверждения для восстановления пароля на EduPlatform:</p>
    <div style="background-color: #f3f4f6; border: 1px solid #e5e7eb; padding: 15px; border-radius: 6px; text-align: center; margin: 20px 0;">
      <span style="font-size: 28px; font-weight: bold; color: #2563eb; letter-spacing: 4px;">%s</span>
    </div>
    <p>Код действителен в течение 10 минут. Если вы не запрашивали восстановление пароля, просто проигнорируйте это письмо.</p>
  </div>
</body>
</html>`, code)

	msg := []byte(subject + mime + body)

	auth := smtp.PlainAuth("", from, pass, host)

	addr := host + ":" + port
	err := smtp.SendMail(addr, auth, from, []string{to}, msg)
	if err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return err
	}

	log.Printf("Email with restore code successfully sent to %s", to)
	return nil
}

func VerifyRestoreCode(email, code string) error {
	dbCode, expiresAt, err := GetRestoreCode(email)
	if err != nil {
		return errors.New("код не найден")
	}

	if dbCode != code {
		return errors.New("неверный код")
	}

	if time.Now().UTC().After(expiresAt) {
		return errors.New("время действия кода истекло")
	}

	return nil
}

func ResetPassword(email, code, password, confirmPassword string) error {
	if password != confirmPassword {
		return errors.New("пароли не совпадают")
	}

	err := VerifyRestoreCode(email, code)
	if err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("ошибка хеширования пароля")
	}

	err = UpdateUserPassword(email, string(hash))
	if err != nil {
		return err
	}

	// Code was verified and used successfully, delete it to prevent reuse
	_ = DeleteRestoreCode(email)

	return nil
}

