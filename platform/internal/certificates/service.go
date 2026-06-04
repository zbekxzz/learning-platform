package certificates

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"platform/internal/database"
	"platform/internal/storage"

	"github.com/minio/minio-go/v7"
)

type Certificate struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	CourseID       int64     `json:"course_id"`
	CourseTitle    string    `json:"course_title"`
	CertificateURL string    `json:"certificate_url"`
	CreatedAt      time.Time `json:"created_at"`
}

func GenerateCertificate(userID, courseID int64) (string, error) {
	// 1. Fetch user's name
	var fullName string
	err := database.DB.QueryRow(context.Background(),
		"SELECT full_name FROM users WHERE id = $1", userID).Scan(&fullName)
	if err != nil {
		return "", fmt.Errorf("failed to fetch user name: %w", err)
	}

	// 2. Fetch course's title
	var courseTitle string
	err = database.DB.QueryRow(context.Background(),
		"SELECT title FROM courses WHERE id = $1", courseID).Scan(&courseTitle)
	if err != nil {
		return "", fmt.Errorf("failed to fetch course title: %w", err)
	}

	// 3. Format SVG certificate
	currentDate := time.Now().Format("02.01.2006")
	svgTemplate := `<?xml version="1.0" encoding="utf-8"?>
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 800 600" width="100%%" height="100%%">
  <defs>
    <linearGradient id="bg-grad" x1="0%%" y1="0%%" x2="100%%" y2="100%%">
      <stop offset="0%%" stop-color="#f3f4f6"/>
      <stop offset="100%%" stop-color="#e5e7eb"/>
    </linearGradient>
  </defs>
  
  <rect width="800" height="600" fill="url(#bg-grad)" />
  <rect x="20" y="20" width="760" height="560" fill="#ffffff" rx="12" stroke="#2563eb" stroke-width="6" />
  <rect x="32" y="32" width="736" height="536" fill="none" stroke="#d1d5db" stroke-width="2" stroke-dasharray="10 5" />
  
  <text x="400" y="110" font-family="'Helvetica Neue', Helvetica, Arial, sans-serif" font-size="38" font-weight="bold" fill="#1e3a8a" text-anchor="middle" letter-spacing="4">СЕРТИФИКАТ</text>
  <text x="400" y="150" font-family="'Helvetica Neue', Helvetica, Arial, sans-serif" font-size="16" font-style="italic" fill="#4b5563" text-anchor="middle">КУРСТЫ СӘТТІ АЯҚТАҒАНЫ ТУРАЛЫ</text>
  
  <text x="400" y="220" font-family="'Helvetica Neue', Helvetica, Arial, sans-serif" font-size="14" fill="#6b7280" text-anchor="middle">Осы сертификат иесіне берілді:</text>
  <text x="400" y="270" font-family="'Helvetica Neue', Helvetica, Arial, sans-serif" font-size="30" font-weight="bold" fill="#111827" text-anchor="middle">%s</text>
  <line x1="150" y1="290" x2="650" y2="290" stroke="#2563eb" stroke-width="2" />
  
  <text x="400" y="350" font-family="'Helvetica Neue', Helvetica, Arial, sans-serif" font-size="14" fill="#6b7280" text-anchor="middle">Тақырыбы бойынша курсты сәтті өткені үшін:</text>
  <text x="400" y="395" font-family="'Helvetica Neue', Helvetica, Arial, sans-serif" font-size="24" font-weight="bold" fill="#1e3a8a" text-anchor="middle">«%s»</text>
  
  <text x="400" y="470" font-family="'Helvetica Neue', Helvetica, Arial, sans-serif" font-size="14" fill="#9ca3af" text-anchor="middle">Күні: %s</text>
  
  <!-- Stepik-style ribbon / badge -->
  <g transform="translate(370, 495)">
    <circle cx="30" cy="30" r="28" fill="#2563eb" />
    <circle cx="30" cy="30" r="24" fill="none" stroke="#ffffff" stroke-width="2" />
    <polygon points="30,12 36,25 50,25 39,33 43,47 30,38 17,47 21,33 10,25 24,25" fill="#f59e0b" />
  </g>
  
  <text x="400" y="580" font-family="'Helvetica Neue', Helvetica, Arial, sans-serif" font-size="11" fill="#9ca3af" text-anchor="middle">EduPlatform - Онлайн Оқыту Платформасы</text>
</svg>
`
	svgContent := fmt.Sprintf(svgTemplate, fullName, courseTitle, currentDate)

	// 4. Create a temp file inside workspace
	tempDir := os.TempDir()
	tempFile, err := os.CreateTemp(tempDir, "certificate-*.svg")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tempFilePath := tempFile.Name()
	defer os.Remove(tempFilePath)
	defer tempFile.Close()

	if _, err := tempFile.WriteString(svgContent); err != nil {
		return "", fmt.Errorf("failed to write svg content: %w", err)
	}
	tempFile.Close()

	// 5. Ensure bucket exists in MinIO
	bucketName := "certificates"
	ctx := context.Background()
	exists, err := storage.Client.BucketExists(ctx, bucketName)
	if err == nil && !exists {
		err = storage.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Println("Warning: failed to make bucket certificates in MinIO:", err)
		}
	}

	// Set bucket policy to public read-only so certificates can be downloaded without signature
	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, bucketName)
	err = storage.Client.SetBucketPolicy(ctx, bucketName, policy)
	if err != nil {
		log.Println("Warning: failed to set public read-only policy on certificates bucket:", err)
	}

	// 6. Upload file to MinIO
	objectName := fmt.Sprintf("cert_%d_%d.svg", userID, courseID)
	url, err := storage.UploadFile(bucketName, objectName, tempFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to upload certificate to MinIO: %w", err)
	}

	// 7. Save to database
	err = SaveCertificate(userID, courseID, url)
	if err != nil {
		return "", fmt.Errorf("failed to save certificate to DB: %w", err)
	}

	return url, nil
}

func SaveCertificate(userID, courseID int64, url string) error {
	_, err := database.DB.Exec(context.Background(),
		`INSERT INTO certificates (user_id, course_id, certificate_url)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (user_id, course_id)
		 DO UPDATE SET certificate_url = $3`,
		userID, courseID, url)
	return err
}

func GetMyCertificates(userID int64) ([]Certificate, error) {
	rows, err := database.DB.Query(context.Background(),
		`SELECT c.id, c.user_id, c.course_id, co.title, c.certificate_url, c.created_at
		 FROM certificates c
		 JOIN courses co ON co.id = c.course_id
		 WHERE c.user_id = $1
		 ORDER BY c.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Certificate
	for rows.Next() {
		var c Certificate
		err = rows.Scan(&c.ID, &c.UserID, &c.CourseID, &c.CourseTitle, &c.CertificateURL, &c.CreatedAt)
		if err == nil {
			result = append(result, c)
		}
	}
	if result == nil {
		result = []Certificate{}
	}
	return result, nil
}

func HasUserEarnedCertificate(userID, courseID int64) (bool, string, error) {
	var url string
	err := database.DB.QueryRow(context.Background(),
		`SELECT certificate_url FROM certificates WHERE user_id = $1 AND course_id = $2`,
		userID, courseID).Scan(&url)
	if err != nil {
		return false, "", nil
	}
	return true, url, nil
}
