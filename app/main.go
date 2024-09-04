package main

import (
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// CORSミドルウェアを追加
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},  // 許可するオリジンを指定
		AllowMethods:     []string{"POST", "GET", "OPTIONS"}, // 許可するHTTPメソッド
		AllowHeaders:     []string{"Content-Type"},           // 許可するヘッダー
		AllowCredentials: true,
	}))

	// ルートハンドラの定義
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	// 画像を受け取り、解像度を取得するハンドラ
	router.POST("/api/v1/upload", func(c *gin.Context) {
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(400, gin.H{"error": "画像の取得に失敗しました。"})
			return
		}

		src, err := file.Open()
		if err != nil {
			c.JSON(500, gin.H{"error": "画像ファイルを開くことができませんでした。"})
			return
		}
		defer src.Close()

		img, _, err := image.DecodeConfig(src)
		if err != nil {
			c.JSON(500, gin.H{"error": "画像情報を取得できませんでした。"})
			return
		}

		c.JSON(200, gin.H{
			"message":  "画像をアップロードしました。",
			"width":    img.Width,
			"height":   img.Height,
			"filename": file.Filename,
		})
	})

	router.Run(":1000")
}
