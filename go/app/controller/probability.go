package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Stay_watch/model"
	"Stay_watch/service"
)

func GetProbability(c *gin.Context) {
	status := c.Param("status") // "reporting" or "leave"
	before := c.Param("before") // "before" or "after"
	user_id := c.Query("user_id")
	str_date := c.Query("date")
	str_time := c.Query("time")

	if user_id == "" || str_date == "" || str_time == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date or time is empty"})
		return
	}

	UserService := service.UserService{}
	uid, err := strconv.Atoi(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is not number"})
		return
	}
	user, err := UserService.GetUserNameByUserID(int64(uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user name"})
		return
	}

	url := "https://stay-estimate.kajilab.dev/app/probability/" + status + "/" + before + "?user_id=" + user_id + "&date=" + str_date + "&time=" + str_time
	req, _ := http.NewRequest("GET", url, nil)
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to access the processing server"})
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		c.JSON(resp.StatusCode, gin.H{"error": "Failed to get probability"})
		return
	}

	body, _ := io.ReadAll(resp.Body)
	var probability model.ProbabilityStayingResponse
	if err := json.Unmarshal(body, &probability); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal probability"})
		return
	}

	probability.UserName = user

	c.JSON(http.StatusOK, probability)
}
