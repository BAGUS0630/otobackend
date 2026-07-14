package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

// FileUploadConfig berisi konfigurasi upload file
type FileUploadConfig struct {
	MaxFileSize int64
	AllowedExt  []string
}

// DefaultFileUploadConfig returns default configuration
func DefaultFileUploadConfig() FileUploadConfig {
	return FileUploadConfig{
		MaxFileSize: 5 * 1024 * 1024, // 5MB
		AllowedExt:  []string{"jpg", "jpeg", "png", "gif"},
	}
}

// ValidateFile memvalidasi file sebelum upload
func ValidateFile(file *multipart.FileHeader, config FileUploadConfig) error {
	// Cek ukuran file
	if file.Size > config.MaxFileSize {
		return fmt.Errorf("ukuran file terlalu besar, maksimal %.2f MB", float64(config.MaxFileSize)/(1024*1024))
	}

	// Cek ekstensi file
	parts := strings.Split(file.Filename, ".")
	if len(parts) < 2 {
		return fmt.Errorf("file harus memiliki ekstensi")
	}

	ext := strings.ToLower(parts[len(parts)-1])
	isAllowed := false
	for _, allowed := range config.AllowedExt {
		if ext == allowed {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return fmt.Errorf("format file tidak didukung. Gunakan: %v", config.AllowedExt)
	}

	return nil
}

// UploadToSupabase mengupload file ke Supabase Storage
func UploadToSupabase(file *multipart.FileHeader, bucket string, folder string) (string, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		return "", fmt.Errorf("Supabase credentials tidak dikonfigurasi")
	}

	// Validasi file
	config := DefaultFileUploadConfig()
	if err := ValidateFile(file, config); err != nil {
		return "", err
	}

	// Buat filename unik dengan timestamp
	parts := strings.Split(file.Filename, ".")
	ext := parts[len(parts)-1]
	filename := fmt.Sprintf("%s-%d.%s", strings.Join(parts[:len(parts)-1], "."), time.Now().Unix(), ext)
	filepath := folder + "/" + filename

	// Buka file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("gagal membuka file: %v", err)
	}
	defer src.Close()

	// Baca file ke memory
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, src); err != nil {
		return "", fmt.Errorf("gagal membaca file: %v", err)
	}

	// Upload ke Supabase Storage
	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", supabaseURL, bucket, filepath)

	req, err := http.NewRequest("POST", uploadURL, buf)
	if err != nil {
		return "", fmt.Errorf("gagal membuat request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+supabaseKey)
	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("gagal upload ke Supabase: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		
		// Coba parse pesan error dari JSON
		var errData map[string]interface{}
		if err := json.Unmarshal(body, &errData); err == nil {
			if msg, ok := errData["message"].(string); ok {
				// Cek pesan spesifik yang sering muncul
				if strings.Contains(msg, "Compact JWS") {
					return "", fmt.Errorf("API Key Supabase tidak valid (terpotong atau salah format)")
				}
				return "", fmt.Errorf("Supabase: %s", msg)
			}
			if errStr, ok := errData["error"].(string); ok {
				return "", fmt.Errorf("Supabase: %s", errStr)
			}
		}
		
		return "", fmt.Errorf("Gagal mengunggah file ke Supabase")
	}

	// Return URL file
	fileURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", supabaseURL, bucket, filepath)
	return fileURL, nil
}

// DeleteFromSupabase menghapus file dari Supabase Storage
func DeleteFromSupabase(fileURL string, bucket string) error {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		return fmt.Errorf("Supabase credentials tidak dikonfigurasi")
	}

	// Extract filepath dari URL
	baseURL := fmt.Sprintf("%s/storage/v1/object/public/%s/", supabaseURL, bucket)
	filepath := strings.TrimPrefix(fileURL, baseURL)

	if filepath == fileURL {
		// Coba format berbeda
		baseURL = fmt.Sprintf("%s/storage/v1/object/%s/", supabaseURL, bucket)
		filepath = strings.TrimPrefix(fileURL, baseURL)
	}

	deleteURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", supabaseURL, bucket, filepath)

	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		return fmt.Errorf("gagal membuat request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+supabaseKey)
	req.Header.Set("apikey", supabaseKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("gagal delete dari Supabase: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		
		var errData map[string]interface{}
		if err := json.Unmarshal(body, &errData); err == nil {
			if msg, ok := errData["message"].(string); ok {
				if strings.Contains(msg, "Compact JWS") {
					return fmt.Errorf("API Key Supabase tidak valid (terpotong atau salah format)")
				}
				return fmt.Errorf("Supabase: %s", msg)
			}
		}
		
		return fmt.Errorf("Gagal menghapus file dari Supabase")
	}

	return nil
}
