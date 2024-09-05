package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nfnt/resize"
	"github.com/rwcarlsen/goexif/exif"
)

const (
	physChunkLength = 9 // 'pHYs'チャンクのデータ長は9バイト
)

// 画像の解像度（dpi）を計算する関数
func calculateDPI(file *os.File) (int, int) {
	// ファイルの形式を確認
	_, format, err := image.DecodeConfig(file)
	if err != nil {
		log.Println("画像フォーマットの取得に失敗しました:", err)
		return 72, 72 // デフォルトのDPIを返す
	}

	// JPEGの場合、EXIFデータからDPIを取得
	if format == "jpeg" {
		file.Seek(0, 0) // ファイルポインタをリセット
		x, err := exif.Decode(file)
		if err != nil {
			log.Println("Exifデータのデコードに失敗しました:", err)
			return 72, 72 // デフォルトのDPIを返す
		}

		// X方向のDPIを取得
		xDPI := getDPI(x, exif.XResolution)
		// Y方向のDPIを取得
		yDPI := getDPI(x, exif.YResolution)

		return xDPI, yDPI
	}

	// PNGの場合はデフォルトのDPIを返す
	if format == "png" {
		xDPI, yDPI := readPNGResolution(file)
		return xDPI, yDPI
	}

	// 他のフォーマットは非対応
	log.Println("対応していない画像フォーマット:", format)
	return 72, 72
}

// DPIを取得するヘルパー関数
func getDPI(x *exif.Exif, tag exif.FieldName) int {
	dpi, err := x.Get(tag)
	if err != nil {
		log.Println("解像度の取得に失敗しました:", err)
		return 72 // デフォルトのDPIを返す
	}

	// dpi.Rat2が3つの戻り値を返すため、変数を調整
	numerator, denominator, err := dpi.Rat2(0)
	if err != nil {
		log.Println("DPIの変換に失敗しました:", err)
		return 72 // デフォルトのDPIを返す
	}

	// 分母が0でない場合にのみDPIを計算
	if denominator != 0 {
		return int(numerator / denominator)
	}
	return 72 // デフォルトのDPIを返す
}

func readPNGResolution(file *os.File) (int, int) {
	// ファイルの先頭に移動
	file.Seek(0, 0)

	// PNGファイルの署名（8バイト）をスキップ
	signature := make([]byte, 8)
	if _, err := file.Read(signature); err != nil {
		log.Println("PNG署名の読み取りに失敗しました:", err)
		return 72, 72
	}

	// 'pHYs'チャンクが見つかるまでチャンクを読み込む
	for {
		// チャンクの長さ（4バイト）とチャンクのタイプ（4バイト）を読み込む
		chunkHeader := make([]byte, 8)
		if _, err := file.Read(chunkHeader); err != nil {
			log.Println("チャンクヘッダーの読み取りに失敗しました:", err)
			return 72, 72
		}

		// チャンクのデータ長を取得
		chunkLength := binary.BigEndian.Uint32(chunkHeader[:4])
		chunkType := string(chunkHeader[4:8])

		// 'pHYs'チャンクを見つけた場合
		if chunkType == "pHYs" {
			// データ長が9バイトであることを確認
			if chunkLength != physChunkLength {
				log.Println("pHYsチャンクの長さが不正です:")
				return 0, 0
			}

			// pHYsチャンクのデータを読み込む
			chunkData := make([]byte, physChunkLength)
			if _, err := file.Read(chunkData); err != nil {
				log.Println("pHYsチャンクのデータ読み取りに失敗しました:", err)
				return 0, 0
			}

			// 水平方向と垂直方向のピクセル密度（メートルあたりのピクセル数）を取得
			xPixelsPerUnit := binary.BigEndian.Uint32(chunkData[:4])
			yPixelsPerUnit := binary.BigEndian.Uint32(chunkData[4:8])

			// 解像度の単位（1 = メートル、0 = 不明）
			unitSpecifier := chunkData[8]

			// 単位がメートルである場合、DPI（インチあたりのドット数）に変換
			if unitSpecifier == 1 {
				// 1メートル = 39.3701インチ、ピクセル密度をDPIに変換
				xDPI := int(float64(xPixelsPerUnit) / 39.3701)
				yDPI := int(float64(yPixelsPerUnit) / 39.3701)
				return xDPI, yDPI
			} else {
				// 単位が不明の場合
				log.Println("pHYsチャンクにはDPIの情報が含まれていません（単位が不明）:")
				return 0, 0
			}
		}

		// チャンクのデータ部分とCRC（データ長+4バイト）をスキップ
		file.Seek(int64(chunkLength+4), os.SEEK_CUR)
	}
}

// 暗号化キー（32バイト）
var encryptionKey = []byte("a very very very very secret key") // 32バイトのキーを使用

// ファイルパスを暗号化する関数
func encrypt(filePath string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(filePath), nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// ファイルパスを複合化する関数
func decrypt(encryptedFilePath string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(encryptedFilePath)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func convertToGrayscale(img image.Image) image.Image {
	grayImg := image.NewGray(img.Bounds())
	draw.Draw(grayImg, grayImg.Bounds(), img, image.Point{}, draw.Src)
	return grayImg
}

func resizeImage(img image.Image, width, height uint) image.Image {
	return resize.Resize(width, height, img, resize.Lanczos3)
}

func main() {
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file:", err)
		}
	}
	baseUrl := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")
	clientPort := os.Getenv("CLIENT_PORT")

	router := gin.Default()

	// CORSミドルウェアを追加
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{baseUrl + ":" + clientPort}, // 許可するオリジンを指定
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},   // 許可するHTTPメソッド
		AllowHeaders:     []string{"Content-Type"},             // 許可するヘッダー
		AllowCredentials: true,
	}))

	// 画像をアップロードするエンドポイント
	router.POST("/api/v1/upload", func(c *gin.Context) {
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(400, gin.H{"error": "画像の取得に失敗しました。"})
			return
		}

		// 現在の日時を取得してフォルダ名を作成
		timestamp := time.Now().Format("20060102150405")
		savePath := filepath.Join("./uploads/", timestamp)
		if err := os.MkdirAll(savePath, os.ModePerm); err != nil {
			c.JSON(500, gin.H{"error": "フォルダの作成に失敗しました。"})
			return
		}

		filePath := filepath.Join(savePath, file.Filename)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(500, gin.H{"error": "画像のアップロードに失敗しました。"})
			return
		}

		// 画像の幅と高さを取得
		fileData, err := os.Open(filePath)
		if err != nil {
			c.JSON(500, gin.H{"error": "画像ファイルを開くことができませんでした。"})
			return
		}
		defer fileData.Close()

		img, _, err := image.Decode(fileData)
		if err != nil {
			c.JSON(500, gin.H{"error": "画像のデコードに失敗しました。"})
			return
		}

		width := img.Bounds().Dx()
		height := img.Bounds().Dy()

		// 画像の解像度（dpi）を計算
		xDPI, yDPI := calculateDPI(fileData)

		// ファイルパスを暗号化
		encryptedFilePath, err := encrypt(filePath)
		if err != nil {
			c.JSON(500, gin.H{"error": "ファイルパスの暗号化に失敗しました。"})
			return
		}

		c.JSON(200, gin.H{
			"message":  "画像をアップロードしました。",
			"filename": file.Filename,
			"path":     encryptedFilePath,
			"width":    width,
			"height":   height,
			"xDPI":     xDPI,
			"yDPI":     yDPI,
		})
	})

	// 画像をリサイズ、グレースケール処理するエンドポイント
	router.POST("/api/v1/process", func(c *gin.Context) {
		// ファイル名を取得
		encryptedFilePath := c.PostForm("path")

		// ファイルパスを複合化
		filePath, err := decrypt(encryptedFilePath)
		if err != nil {
			c.JSON(500, gin.H{"error": "ファイルパスの複合化に失敗しました。"})
			return
		}

		// ファイルが存在するか確認
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(404, gin.H{"error": "画像が見つかりません。"})
			return
		}

		// 画像ファイルを開く
		fileData, err := os.Open(filePath)
		if err != nil {
			c.JSON(500, gin.H{"error": "画像ファイルを開くことができませんでした。"})
			return
		}
		defer fileData.Close()

		// 画像をデコード
		img, _, err := image.Decode(fileData)
		if err != nil {
			c.JSON(500, gin.H{"error": "画像のデコードに失敗しました。"})
			return
		}

		// グレースケール変換のオプション
		if c.PostForm("grayscale") == "true" {
			img = convertToGrayscale(img)
		}

		// リサイズのオプション
		widthStr := c.PostForm("width")
		heightStr := c.PostForm("height")
		if widthStr != "" && heightStr != "" {
			width, _ := strconv.Atoi(widthStr)
			height, _ := strconv.Atoi(heightStr)
			img = resizeImage(img, uint(width), uint(height))
		}

		// 画像をクライアントに返す
		c.Header("Content-Type", "image/jpeg")
		err = jpeg.Encode(c.Writer, img, nil)
		if err != nil {
			c.JSON(500, gin.H{"error": "画像のエンコードに失敗しました。"})
			return
		}
	})

	router.Run(":" + port)
}
