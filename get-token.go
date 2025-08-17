// package main

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/exec"
// 	"runtime"
// 	"time"

// 	"golang.org/x/oauth2"
// 	"golang.org/x/oauth2/google"
// 	"google.golang.org/api/option"
// 	"google.golang.org/api/youtube/v3"
// )

// // Konfigürasyon değerleri - buraya kendi bilgilerinizi girin
// var (
// 	filename     = "video.mp4"                 // Yüklenecek video dosyası
// 	title        = "Test Video"                // Video başlığı
// 	description  = "Bu video API ile yüklendi" // Video açıklaması
// 	category     = "22"                        // Video kategorisi (22 = People & Blogs)
// 	keywords     = "test,api,youtube"          // Anahtar kelimeler (virgülle ayrılmış)
// 	privacy      = "unlisted"                  // Gizlilik durumu (public, unlisted, private)
// 	clientID     = "YOUR_CLIENT_ID_HERE"       // OAuth2 Client ID
// 	clientSecret = "YOUR_CLIENT_SECRET_HERE"   // OAuth2 Client Secret
// )

// func getClient(config *oauth2.Config) *http.Client {
// 	tokenFile := "token.json"
// 	tok, err := tokenFromFile(tokenFile)
// 	if err != nil {
// 		tok = getTokenFromWeb(config)
// 		saveToken(tokenFile, tok)
// 	}
// 	return config.Client(context.Background(), tok)
// }

// func tokenFromFile(file string) (*oauth2.Token, error) {
// 	f, err := os.Open(file)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()
// 	tok := &oauth2.Token{}
// 	err = json.NewDecoder(f).Decode(tok)
// 	return tok, err
// }

// func saveToken(path string, token *oauth2.Token) {
// 	fmt.Printf("Token kaydedildi: %s\n", path)
// 	f, err := os.Create(path)
// 	if err != nil {
// 		log.Fatalf("Token kaydedilemedi: %v", err)
// 	}
// 	defer f.Close()
// 	json.NewEncoder(f).Encode(token)
// }

// // openBrowser tarayıcıda URL açar
// func openBrowser(url string) {
// 	var err error
// 	switch runtime.GOOS {
// 	case "linux":
// 		err = exec.Command("xdg-open", url).Start()
// 	case "windows":
// 		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
// 	case "darwin":
// 		err = exec.Command("open", url).Start()
// 	default:
// 		err = fmt.Errorf("unsupported platform")
// 	}
// 	if err != nil {
// 		fmt.Printf("Tarayıcı açılamadı, manuel olarak açın: %s\n", url)
// 	}
// }

// func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
// 	// HTTP server oluştur
// 	codeCh := make(chan string)
// 	server := &http.Server{Addr: ":8080"}

// 	// OAuth2 callback endpoint
// 	http.HandleFunc("/oauth2callback", func(w http.ResponseWriter, r *http.Request) {
// 		code := r.URL.Query().Get("code")
// 		if code == "" {
// 			http.Error(w, "Authorization code not found", http.StatusBadRequest)
// 			return
// 		}

// 		// Başarı sayfası göster
// 		w.Header().Set("Content-Type", "text/html")
// 		fmt.Fprint(w, `
// 			<html>
// 			<head><title>Yetkilendirme Tamamlandı</title></head>
// 			<body>
// 				<h1>✅ Yetkilendirme Başarılı!</h1>
// 				<p>Bu pencereyi kapatabilirsiniz.</p>
// 				<script>setTimeout(function(){window.close();}, 3000);</script>
// 			</body>
// 			</html>
// 		`)

// 		codeCh <- code
// 	})

// 	// Redirect URL'yi ayarla
// 	config.RedirectURL = "http://localhost:8080/oauth2callback"

// 	// Authorization URL oluştur
// 	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

// 	fmt.Println("🌐 Tarayıcıda açılıyor...")
// 	fmt.Println("Manuel açmak için: ", authURL)

// 	// Tarayıcıyı aç
// 	openBrowser(authURL)

// 	// Server'ı başlat
// 	go func() {
// 		if err := server.ListenAndServe(); err != http.ErrServerClosed {
// 			log.Printf("HTTP server hatası: %v", err)
// 		}
// 	}()

// 	fmt.Println("📡 OAuth2 callback bekleniyor...")

// 	// Authorization code'u bekle
// 	code := <-codeCh

// 	// Server'ı kapat
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	server.Shutdown(ctx)

// 	// Token al
// 	tok, err := config.Exchange(context.Background(), code)
// 	if err != nil {
// 		log.Fatalf("Token exchange hatası: %v", err)
// 	}

// 	return tok
// }
// func main() {
// 	ctx := context.Background()

// 	// Google API istemci bilgileri
// 	b, err := os.ReadFile("client_secret.json")
// 	if err != nil {
// 		log.Fatalf("client_secret.json dosyası açılamadı: %v", err)
// 	}

// 	config, err := google.ConfigFromJSON(b, youtube.YoutubeUploadScope)
// 	if err != nil {
// 		log.Fatalf("Config oluşturulamadı: %v", err)
// 	}

// 	client := getClient(config)

// 	// YouTube API servisini başlat
// 	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
// 	if err != nil {
// 		log.Fatalf("YouTube servisi başlatılamadı: %v", err)
// 	}

// 	today := time.Now().Format("2006-01-02")

// 	// Video bilgileri
// 	video := &youtube.Video{
// 		Snippet: &youtube.VideoSnippet{
// 			Title:       fmt.Sprintf("Enayiiii", today),
// 			Description: "Enayi Ahmet",
// 			Tags:        []string{"go", "youtube", "bot"},
// 			CategoryId:  "22", // People & Blogs
// 		},
// 		Status: &youtube.VideoStatus{
// 			PrivacyStatus: "public",
// 		},
// 	}

// 	// Video dosyası
// 	file, err := os.Open("video.mp4")
// 	if err != nil {
// 		log.Fatalf("Video dosyası açılamadı: %v", err)
// 	}
// 	defer file.Close()

// 	// Yükleme işlemi
// 	call := service.Videos.Insert([]string{"snippet", "status"}, video)
// 	response, err := call.Media(file).Do()
// 	if err != nil {
// 		log.Fatalf("Video yüklenemedi: %v", err)
// 	}

// 	fmt.Printf("✅ Video yüklendi! Video ID: %v\n", response.Id)
// }
