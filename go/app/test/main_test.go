// main_test.go
package main

import (
	controller "Stay_watch/controller"
	"Stay_watch/model"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strconv"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPostStayer(t *testing.T) {
	response := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(response)

	beaconsRoom := model.BeaconRoom{
		Beacons: []*model.Beacon{
			{
				Uuid: "e7d61ea3f8dd49c88f2ff2484c07ac00",
				Rssi: -60,
			},
		},
		RoomID: 1,
	}
	//jsonに変換
	jsonBeaconsRoom, err := json.Marshal(beaconsRoom)
	if err != nil {
		t.Fatal(err)
	}

	// リクエスト情報をコンテキストに入れる
	ginContext.Request, _ = http.NewRequest(http.MethodPost, "/stayers", nil)
	ginContext.Request.Header.Set("Content-Type", "application/json")
	ginContext.Request.Body = ioutil.NopCloser(bytes.NewBuffer(jsonBeaconsRoom))
	controller.Beacon(ginContext)
	asserts := assert.New(t)
	// レスポンスのステータスコードの確認
	asserts.Equal(http.StatusCreated, response.Code)
}

func TestGetStayer(t *testing.T) {
	response := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(response)

	// リクエストの生成
	// 今回はmiddlewareのテストのためpathはなんでも可
	req, _ := http.NewRequest(http.MethodGet, "/stayers", nil)

	// リクエスト情報をコンテキストに入れる
	ginContext.Request = req
	controller.Stayer(ginContext)

	asserts := assert.New(t)

	// レスポンスのステータスコードの確認
	asserts.Equal(http.StatusOK, response.Code)
	// レスポンスのボディを構造体に変換
	var responseStayer []model.StayerGetResponse
	json.Unmarshal(response.Body.Bytes(), &responseStayer)
	// レスポンスのボディの確認
	asserts.Equal("kaji", responseStayer[0].Name)
	asserts.Equal("梶研", responseStayer[0].Tags[0].Name)
	asserts.Equal("梶研", responseStayer[0].Tags[0].Name)
	asserts.Equal(1, int(responseStayer[0].Tags[0].ID))
	asserts.Equal(1, int(responseStayer[0].ID))
	asserts.Equal(1, int(responseStayer[0].RoomID))
}

func TestGetLog(t *testing.T) {
	response := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(response)
	// リクエストの生成
	// 今回はmiddlewareのテストのためpathはなんでも可
	req, _ := http.NewRequest(http.MethodGet, "/logs", nil)
	// リクエスト情報をコンテキストに入れる
	ginContext.Request = req
	controller.Log(ginContext)
	asserts := assert.New(t)
	// レスポンスのステータスコードの確認
	asserts.Equal(http.StatusOK, response.Code)
	// レスポンスのボディを構造体に変換
	var responseLog []model.LogGetResponse
	json.Unmarshal(response.Body.Bytes(), &responseLog)
	// レスポンスのボディの確認
	asserts.Equal("kaji", responseLog[0].Name)
	asserts.Equal("梶研-学生部屋", responseLog[0].Room)
	// asserts.Equal(1, int(responseLog[0].ID))
	// asserts.Equal("2021-05-01 00:00:00", responseLog[0].StartAt)
	// asserts.Equal("2021-05-01 00:00:00", responseLog[0].EndAt)
}

func TestGetUser(t *testing.T) {
	response := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(response)
	// リクエストの生成
	// 今回はmiddlewareのテストのためpathはなんでも可
	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	// リクエスト情報をコンテキストに入れる
	ginContext.Request = req
	controller.UserList(ginContext)
	asserts := assert.New(t)
	// レスポンスのステータスコードの確認
	asserts.Equal(http.StatusOK, response.Code)
	// レスポンスのボディを構造体に変換
	var responseUser []model.UserInformationGetResponse
	json.Unmarshal(response.Body.Bytes(), &responseUser)
	// レスポンスのボディの確認
	//fmt.Println(responseUser)
	asserts.Equal("kaji", responseUser[0].Name)
	asserts.Equal("梶研", responseUser[0].Tags[0].Name)
	asserts.Equal(1, int(responseUser[0].ID))

}

// 管理者画面でのユーザ取得API
func TestGetEditorUser(t *testing.T) {

	router := gin.Default()
	router.GET("/api/v1/users/:communityId", controller.UserList)

	asserts := assert.New(t)

	lastSumUsers := 0
	isAllSumEqual := true

	for i := 0; i < 10; i++ {
		// HTTPリクエストの生成
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/"+strconv.Itoa(i)+"?fields=admin", nil)

		// レスポンスのレコーダーを作成
		res := httptest.NewRecorder()

		// リクエストをハンドル
		router.ServeHTTP(res, req)
		// レスポンスのステータスコードの確認
		asserts.Equal(http.StatusOK, res.Code)

		// レスポンスのボディを構造体に変換
		var responseUser []model.UserEditorResponse
		json.Unmarshal(res.Body.Bytes(), &responseUser)

		if lastSumUsers != len(responseUser) {
			isAllSumEqual = false
		}
	}
	if isAllSumEqual {
		// 全てユーザ数が同じ場合は正常なら存在しないため
		// community_idによる絞り込みができていないか、データがそもそも取れていないかなど
		t.Fatalf("All community users have the same count")
	}
}
