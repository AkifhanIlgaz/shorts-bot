// main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// TokenManager struct to handle token operations
type TokenManager struct {
	config *oauth2.Config
}

// NewTokenManager creates a new token manager
func NewTokenManager() *TokenManager {
	config := &oauth2.Config{
		ClientID:     os.Getenv("YOUTUBE_CLIENT_ID"),
		ClientSecret: os.Getenv("YOUTUBE_CLIENT_SECRET"),
		Scopes:       []string{youtube.YoutubeUploadScope},
		Endpoint:     google.Endpoint,
	}

	return &TokenManager{config: config}
}

// GetValidToken returns a valid access token, refreshing if necessary
func (tm *TokenManager) GetValidToken(ctx context.Context) (*oauth2.Token, error) {
	refreshToken := os.Getenv("YOUTUBE_REFRESH_TOKEN")
	if refreshToken == "" {
		return nil, fmt.Errorf("refresh token bulunamadı")
	}

	// Refresh token ile yeni access token al
	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	// Token source oluştur (otomatik refresh yapar)
	tokenSource := tm.config.TokenSource(ctx, token)

	// Fresh token al
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("token yenilenemedi: %v", err)
	}

	return newToken, nil
}

// VideoUploader handles YouTube video uploads
type VideoUploader struct {
	service *youtube.Service
	ctx     context.Context
}

// NewVideoUploader creates a new video uploader
func NewVideoUploader(ctx context.Context, token *oauth2.Token) (*VideoUploader, error) {
	tm := NewTokenManager()
	client := tm.config.Client(ctx, token)

	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("YouTube servisi başlatılamadı: %v", err)
	}

	return &VideoUploader{
		service: service,
		ctx:     ctx,
	}, nil
}

// VideoMetadata contains video information
type VideoMetadata struct {
	Title         string
	Description   string
	Tags          []string
	CategoryID    string
	PrivacyStatus string
	Language      string
}

// UploadVideo uploads a video to YouTube
func (vu *VideoUploader) UploadVideo(videoPath string, metadata VideoMetadata) (*youtube.Video, error) {
	// Video dosyası kontrol
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("video dosyası bulunamadı: %s", videoPath)
	}

	// Dosya boyutu bilgisi
	fileInfo, err := os.Stat(videoPath)
	if err != nil {
		return nil, fmt.Errorf("dosya bilgisi alınamadı: %v", err)
	}
	fileSizeMB := float64(fileInfo.Size()) / (1024 * 1024)
	fmt.Printf("📁 Video dosyası: %s (%.2f MB)\n", videoPath, fileSizeMB)

	// Video dosyası aç
	file, err := os.Open(videoPath)
	if err != nil {
		return nil, fmt.Errorf("video dosyası açılamadı: %v", err)
	}
	defer file.Close()

	// Video metadata oluştur
	video := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:                metadata.Title,
			Description:          metadata.Description,
			Tags:                 metadata.Tags,
			CategoryId:           metadata.CategoryID,
			DefaultLanguage:      metadata.Language,
			DefaultAudioLanguage: metadata.Language,
		},
		Status: &youtube.VideoStatus{
			PrivacyStatus:           metadata.PrivacyStatus,
			SelfDeclaredMadeForKids: false,
		},
	}

	// Video yükleme
	fmt.Println("📤 Video YouTube'a yükleniyor...")
	call := vu.service.Videos.Insert([]string{"snippet", "status"}, video)
	response, err := call.Media(file).Do()
	if err != nil {
		return nil, fmt.Errorf("video yüklenemedi: %v", err)
	}

	return response, nil
}

// GenerateDefaultMetadata creates default video metadata
func GenerateDefaultMetadata() VideoMetadata {
	// 1 Ağustos 2025 başlangıç tarihi
	startDate, _ := time.ParseInLocation("2006-01-02", "2025-08-01", time.Local)

	// Türkiye saati
	loc, _ := time.LoadLocation("Europe/Istanbul")
	now := time.Now().In(loc)
	dateStr := now.Format("2006-01-02 15:04")

	// 1 Ağustos 2025'ten bugüne kadar geçen gün sayısı
	daysSinceStart := int(now.Sub(startDate).Hours() / 24)

	// Eğer henüz 1 Ağustos 2025 gelmemişse, 0 gün göster
	if daysSinceStart < 0 {
		daysSinceStart = 1
	}

	description := fmt.Sprintf(`Bu video %s tarihinde otomatik olarak yüklendi.

🤖 Otomatik yükleme sistemi ile oluşturulmuştur.
📅 Yüklenme tarihi: %s
📊 1 Ağustos 2025'ten beri geçen gün: %d
⚙️  GitHub Actions ile otomatik yükleme
🔄 Günlük içerik programı

#OtomatikYükleme #GithubActions #Bot #Günlük #Automation #Gün%d`, dateStr, dateStr, daysSinceStart, daysSinceStart)

	return VideoMetadata{
		Title:         fmt.Sprintf("Gün %d - Kel Ahmet", daysSinceStart),
		Description:   description,
		Tags:          []string{"ahmet", "ten", "verdi", "verdi", fmt.Sprintf("gün%d", daysSinceStart)},
		CategoryID:    "22", // People & Blogs
		PrivacyStatus: "public",
		Language:      "tr",
	}
}

// WriteGitHubOutput writes output for GitHub Actions
func WriteGitHubOutput(videoID, videoTitle string) {
	githubOutput := os.Getenv("GITHUB_OUTPUT")
	if githubOutput == "" {
		return
	}

	f, err := os.OpenFile(githubOutput, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("GitHub output yazılamadı: %v\n", err)
		return
	}
	defer f.Close()

	videoURL := fmt.Sprintf("https://youtube.com/watch?v=%s", videoID)
	fmt.Fprintf(f, "video_id=%s\n", videoID)
	fmt.Fprintf(f, "video_url=%s\n", videoURL)
	fmt.Fprintf(f, "video_title=%s\n", videoTitle)
}

func main() {
	ctx := context.Background()

	fmt.Println("🚀 YouTube video yükleme başlatılıyor...")

	// Environment variables kontrol
	requiredEnvs := []string{"YOUTUBE_CLIENT_ID", "YOUTUBE_CLIENT_SECRET", "YOUTUBE_REFRESH_TOKEN"}
	for _, env := range requiredEnvs {
		if os.Getenv(env) == "" {
			log.Fatalf("❌ Gerekli environment variable eksik: %s", env)
		}
	}

	// Token manager oluştur
	tokenManager := NewTokenManager()

	// Valid token al
	token, err := tokenManager.GetValidToken(ctx)
	if err != nil {
		log.Fatalf("❌ Token alınamadı: %v", err)
	}

	fmt.Println("✅ Token başarıyla alındı")

	// Video uploader oluştur
	uploader, err := NewVideoUploader(ctx, token)
	if err != nil {
		log.Fatalf("❌ Uploader oluşturulamadı: %v", err)
	}

	// Video metadata oluştur
	metadata := GenerateDefaultMetadata()

	// Video yükle
	videoPath := "video.mp4"
	response, err := uploader.UploadVideo(videoPath, metadata)
	if err != nil {
		log.Fatalf("❌ Video yükleme hatası: %v", err)
	}

	// Başarı mesajları
	loc, _ := time.LoadLocation("Europe/Istanbul")
	now := time.Now().In(loc)

	fmt.Println("✅ Video başarıyla yüklendi!")
	fmt.Printf("🔗 Video ID: %s\n", response.Id)
	fmt.Printf("🌐 URL: https://youtube.com/watch?v=%s\n", response.Id)
	fmt.Printf("📝 Başlık: %s\n", metadata.Title)
	fmt.Printf("🔒 Gizlilik: %s\n", metadata.PrivacyStatus)
	fmt.Printf("⏰ Yüklenme zamanı: %s\n", now.Format("2006-01-02 15:04:05"))

	// GitHub Actions output
	WriteGitHubOutput(response.Id, metadata.Title)

	fmt.Println("🎉 İşlem tamamlandı!")
}
