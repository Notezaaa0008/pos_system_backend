package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/wneessen/go-mail"
)

func SendResetPasswordEmail(toEmail string, token string) error {
	// 1. ดึงค่าคอนฟิกจาก Environment Variables (.env)
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	frontendURL := os.Getenv("FRONTEND_URL")

	if smtpHost == "" || smtpPortStr == "" || smtpUser == "" || smtpPass == "" || frontendURL == "" {
		log.Println("Error: SMTP HOST or PORT or USER or PASS or FRONTEND_URL is missing in .env")
		// ถ้าลืมประกาศตัวแปรนี้ใน .env เลย ให้เบรกระบบทันที
		return errors.New("internal server error: security configuration is missing")
	}

	// แปลง Port จาก string เป็น int
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		log.Println("Admin Warning: SMTP_PORT in .env must be a number:", err)
		return errors.New("internal server error: invalid security configuration format")
	}

	// 2. สร้างก้อนข้อความ (Message)
	m := mail.NewMsg()
	
	// ตั้งค่าผู้ส่ง และผู้รับ
	if err := m.From(smtpUser); err != nil {
		return fmt.Errorf("failed to set FROM: %w", err)
	}
	if err := m.To(toEmail); err != nil {
		return fmt.Errorf("failed to set TO: %w", err)
	}

	// ตั้งชื่อหัวข้ออีเมล (Subject)
	m.Subject("🔥 รีเซ็ตรหัสผ่านของคุณ (Reset Password)")

	// ประกอบร่างลิงก์ที่จะให้หน้าบ้านเอาไปใช้ต่อ
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)

	// เขียนเนื้อหาเป็น HTML สวยๆ
	htmlBody := fmt.Sprintf(`
		<div style="font-family: sans-serif; line-height: 1.6; color: #333;">
			<h2>ระบบรีเซ็ตรหัสผ่าน</h2>
			<p>คุณได้ทำการขอรีเซ็ตรหัสผ่านใหม่ กรุณาคลิกที่ปุ่มด้านล่างเพื่อดำเนินการต่อ:</p>
			<p style="margin: 24px 0;">
				<a href="%s" style="background-color: #4F46E5; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; font-weight: bold;">
					เปลี่ยนรหัสผ่านใหม่
				</a>
			</p>
			<p style="color: #666; font-size: 13px;">* ลิงก์นี้จะมีอายุการใช้งานแค่ 15 นาทีเท่านั้น</p>
			<hr style="border: none; border-top: 1px solid #eee; margin: 24px 0;" />
			<p style="color: #999; font-size: 12px;">หากคุณไม่ได้เป็นคนส่งคำขอนี้ สามารถปล่อยผ่านอีเมลฉบับนี้ไปได้เลยครับ</p>
		</div>
	`, resetLink)

	// ยัดเนื้อหา HTML ใส่ใน Message
	m.SetBodyString(mail.TypeTextHTML, htmlBody)

	// 3. สร้าง Client สำหรับเชื่อมต่อเซิร์ฟเวอร์ SMTP
	client, err := mail.NewClient(smtpHost,
		mail.WithPort(smtpPort),
		mail.WithSMTPAuth(mail.SMTPAuthPlain), // ใช้การยืนยันตัวตนแบบปกติ
		mail.WithUsername(smtpUser),
		mail.WithPassword(smtpPass),
		mail.WithTLSPolicy(mail.TLSMandatory), // บังคับเข้ารหัส TLS เพื่อความปลอดภัย (สำหรับ Port 587)
	)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}

	// 4. สั่งยิงเมลออกไปตรงๆ
	if err := client.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}