package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	pushnotifications "github.com/pusher/push-notifications-go"
	"github.com/pusher/pusher-http-go"
)

const (
	instanceId = "1038202c-9ec2-498d-b25b-a0b3d71c4c52"
	secretKey  = "FA07044558EDB6307EF4E4D15227D5A1F8E2D20ECF1F2E4C65A8781C0FBA4CB3"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// chat app
	pusherClient := pusher.Client{
		AppID:   "1353465",
		Key:     "caf06d53f3267feb766a",
		Secret:  "827422f36923313a8d5f",
		Cluster: "ap1",
		Secure:  true,
	}

	beamsClient, err := pushnotifications.New(instanceId, secretKey)
	if err != nil {
		fmt.Println("Could not create Beams Client:", err.Error())
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// chat app
	r.GET("/chat", func(c *gin.Context) {
		data := map[string]string{"message": "hello world"}
		pusherClient.Trigger("chat", "test", data)

		c.JSON(200, gin.H{"message": "ok"})
	})

	r.GET("/noti", func(c *gin.Context) {
		publishRequest := map[string]interface{}{
			"apns": map[string]interface{}{
				"aps": map[string]interface{}{
					"alert": map[string]interface{}{
						"title": "Hello",
						"body":  "Hello, world",
					},
				},
			},
			"fcm": map[string]interface{}{
				"notification": map[string]interface{}{
					"title": "Hello",
					"body":  "Hello, world",
				},
			},
			"web": map[string]interface{}{
				"notification": map[string]interface{}{
					"title": "Hello",
					"body":  "Hello, world",
				},
			},
		}

		pubId, err := beamsClient.PublishToInterests([]string{"hello"}, publishRequest)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Publish Id:", pubId)
		}

		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	r.GET("/auth", func(c *gin.Context) {
		userID := c.GetHeader("Authorization")
		// fmt.Println(userID)

		userIDInQueryParam := c.Query("user_id")
		// fmt.Println(userIDInQueryParam)

		if userID != userIDInQueryParam {
			c.JSON(401, gin.H{
				"msg": "error",
			})
			return
		}

		beamsToken, err := beamsClient.GenerateToken(userID)
		if err != nil {
			c.JSON(400, gin.H{
				"msg": "error",
			})
			return
		}

		beamsTokenJson, err := json.Marshal(beamsToken)
		fmt.Println(beamsToken)
		fmt.Println(beamsTokenJson)
		if err != nil {
			c.JSON(400, gin.H{
				"msg": "error",
			})
			return
		}

		c.JSON(200, gin.H{
			"token": beamsToken["token"],
		})
	})

	r.GET("/noti/user", func(c *gin.Context) {
		publishRequest := map[string]interface{}{
			"apns": map[string]interface{}{
				"aps": map[string]interface{}{
					"alert": map[string]interface{}{
						"title": "Hello",
						"body":  "Hello, world",
					},
				},
			},
			"fcm": map[string]interface{}{
				"notification": map[string]interface{}{
					"title": "Hello",
					"body":  "Hello, world",
				},
			},
			"web": map[string]interface{}{
				"notification": map[string]interface{}{
					"title": "Hello2",
					"body":  "Hello, world2",
				},
			},
		}

		pubId, err := beamsClient.PublishToUsers([]string{"user-001", "user-002"}, publishRequest)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Publish Id:", pubId)
		}

		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
