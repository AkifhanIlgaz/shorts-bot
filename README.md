# YouTube Video Upload Bot

Bu uygulama, YouTube API kullanarak video yüklemenizi sağlar.

## Gereksinimler

1. **YouTube API Key & OAuth2 Credentials**: Google Cloud Console'dan YouTube Data API v3'ü etkinleştirin ve OAuth2 kimlik bilgilerini oluşturun.

## Kurulum

1. Bağımlılıkları yükleyin:

```bash
go mod download
```

2. Uygulamayı derleyin:

```bash
go build -o youtube-uploader main.go
```

## YouTube API Ayarları

1. [Google Cloud Console](https://console.cloud.google.com/)'a gidin
2. Yeni bir proje oluşturun veya mevcut projeyi seçin
3. **APIs & Services > Library** bölümünden "YouTube Data API v3"ü etkinleştirin
4. **APIs & Services > Credentials** bölümünden "OAuth 2.0 Client IDs" oluşturun
5. Application type olarak "Desktop application" seçin
6. Client ID ve Client Secret bilgilerini not edin

## Kullanım

```bash
./youtube-uploader \
  -filename="video.mp4" \
  -title="Video Başlığı" \
  -description="Video açıklaması" \
  -clientid="YOUR_CLIENT_ID" \
  -clientsecret="YOUR_CLIENT_SECRET" \
  -privacy="unlisted"
```

### Parametreler

- `-filename`: Yüklenecek video dosyasının yolu (zorunlu)
- `-title`: Video başlığı (varsayılan: "Test Title")
- `-description`: Video açıklaması (varsayılan: "Test Description")
- `-category`: Video kategorisi (varsayılan: "22")
- `-keywords`: Virgülle ayrılmış anahtar kelimeler
- `-privacy`: Gizlilik durumu (unlisted, private, public) (varsayılan: "unlisted")
- `-clientid`: OAuth2 Client ID (zorunlu)
- `-clientsecret`: OAuth2 Client Secret (zorunlu)

### Örnek Kullanım

```bash
./youtube-uploader \
  -filename="my-video.mp4" \
  -title="Harika Video" \
  -description="Bu çok güzel bir video" \
  -keywords="eğlence,müzik,komedi" \
  -clientid="123456789-abcdefghijklmnop.apps.googleusercontent.com" \
  -clientsecret="GOCSPX-abcdefghijklmnopqrstuvwxyz" \
  -privacy="public"
```

## İlk Çalıştırma

İlk kez çalıştırdığınızda:

1. Uygulama size bir OAuth2 URL'si verecek
2. Bu URL'yi tarayıcınızda açın
3. Google hesabınızla giriş yapın ve uygulamaya izin verin
4. Verilen authorization code'u terminale girin
5. Video yükleme işlemi başlayacak

## Desteklenen Video Formatları

- MP4
- MOV
- AVI
- WMV
- FLV
- WebM

## Notlar

- Video dosyası boyutu sınırları YouTube'un belirlediği limitlere tabidir
- İlk yükleme sırasında OAuth2 yetkilendirmesi gereklidir
- Uygulamayı kullanabilmek için YouTube Data API v3'ün etkin olması gerekir
# shorts-bot
