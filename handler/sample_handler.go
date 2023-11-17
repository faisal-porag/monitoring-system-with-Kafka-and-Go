package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"real_time_monitoring_system/producer"
	"real_time_monitoring_system/utils"
	"time"
)

type SendMessageDataPayload struct {
	Message string `json:"message"`
}

func SendMessageHandler() func(c *gin.Context) {
	type Success struct {
		Status string `json:"status"`
	}
	type MessageDataResponse struct {
		Message     string    `json:"message"`
		MessageTime time.Time `json:"message_time"`
	}
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json; charset=utf-8")

		var reqPayload SendMessageDataPayload
		if err := c.ShouldBindJSON(&reqPayload); err != nil {
			c.JSON(
				http.StatusBadRequest,
				map[string]interface{}{
					"message": "The provided information is invalid. Please recheck and try again",
				},
			)
			return
		}

		sendMessageData := MessageDataResponse{
			Message:     reqPayload.Message,
			MessageTime: time.Now(),
		}

		err := producer.ProduceMessage(sendMessageData, utils.NotificationBulk)
		if err != nil {
			log.Println(err)
		}

		c.JSON(http.StatusOK, Success{
			Status: "Ok",
		})

		return
	}
}
