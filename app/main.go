package main

import (
	"image"
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

	"github.com/Yuki-w6/dpi/imageUtil"
)

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
		xDPI, yDPI := imageUtil.CalculateDPI(fileData)

		// ファイルパスを暗号化
		encryptedFilePath, err := imageUtil.Encrypt(filePath)
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
		filePath, err := imageUtil.Decrypt(encryptedFilePath)
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
			img = imageUtil.ConvertToGrayscale(img)
		}

		// リサイズのオプション
		widthStr := c.PostForm("width")
		heightStr := c.PostForm("height")
		if widthStr != "" && heightStr != "" {
			width, _ := strconv.Atoi(widthStr)
			height, _ := strconv.Atoi(heightStr)
			img = imageUtil.ResizeImage(img, uint(width), uint(height))
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
