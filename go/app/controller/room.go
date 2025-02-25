package controller

import (
	"Stay_watch/model"
	"Stay_watch/service"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// 例："100,101-200,201" -> [[100,101],[200,201]]
func ParseStringToIntSlice(str string) [][]int64 {
	polygonStringArray := strings.Split(str, "-")
	polygonIntArray := [][]int64{}
	for _, pointString := range polygonStringArray {
		parts := strings.Split(pointString, ",") // parts: ["100","101"]
		pointIntArray := []int64{}
		for _, part := range parts {
			tmp, _ := strconv.Atoi(part)
			pointIntArray = append(pointIntArray, int64(tmp))
		}
		polygonIntArray = append(polygonIntArray, pointIntArray)
	}
	return polygonIntArray
}

func UpdateRoom(c *gin.Context) {

	RoomForm := model.RoomEditorForm{}
	err := c.Bind(&RoomForm)

	if err != nil {
		fmt.Println(err)
		return
	}

	// [2][2]int64 -> string
	storePolygon := strconv.FormatInt(RoomForm.Polygon[0][0], 10) + "," + strconv.FormatInt(RoomForm.Polygon[0][1], 10) + "-" + strconv.FormatInt(RoomForm.Polygon[1][0], 10) + "," + strconv.FormatInt(RoomForm.Polygon[1][1], 10)

	RoomService := service.RoomService{}
	RoomService.UpdateRoom(int(RoomForm.RoomID), RoomForm.RoomName, int(RoomForm.BuildingID), storePolygon)

	c.JSON(http.StatusCreated, gin.H{
		"status": "ok",
	})
}

func GetRoomsByCommunityID(c *gin.Context) {
	communityID, err := strconv.ParseInt(c.Param("communityID"), 10, 64) // string -> int64
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "Type is not number")
	}

	RoomService := service.RoomService{}
	rooms, err := RoomService.GetAllRooms()
	if err != nil {
		fmt.Printf("failed: Cannnot get stayer %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get rooms"})
		return
	}
	roomsGetResponse := []model.RoomsGetResponse{}

	//---communityIDからコミュニティの名前を調べる機能実装予定---

	//----------------------------------------------------

	for _, room := range rooms {
		// fmt.Print("コミュニティID: ")
		// fmt.Println(room.CommunityID)
		// fmt.Print("ポリゴン: ")
		// fmt.Println(room.Polygon)

		// コミュニティIDが一致した部屋情報だけ返す
		if room.CommunityID == communityID {
			roomName := room.Name
			roomID := int64(room.ID)

			//---roomIDから建物の名前,IDを調べる機能実装予定---

			//----------------------------------------
			roomsGetResponse = append(roomsGetResponse, model.RoomsGetResponse{
				RoomID:        roomID,
				Name:          roomName,
				CommunityName: "梶研究室",
				BuildingName:  "4号館",
				Polygon:       ParseStringToIntSlice(room.Polygon),
				BuildingId:    room.BuildingID,
			})
		}
	}
	c.JSON(http.StatusOK, roomsGetResponse)

}

func Stayer(c *gin.Context) {

	RoomService := service.RoomService{}
	UserService := service.UserService{}

	//Stayerテーブルから全てのデータを取得する
	allStayer, err := RoomService.GetAllStayer()
	if err != nil {

		fmt.Printf("failed: Cannnot get stayer %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stayer"})
		return
	}

	stayerGetResponse := []model.StayerGetResponse{}

	for _, stayer := range allStayer {

		userName, err := UserService.GetUserNameByUserID(stayer.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user name"})
			return
		}
		roomName, err := RoomService.GetRoomNameByRoomID(stayer.RoomID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get room name"})
			return
		}
		tagsGetResponse := make([]model.TagGetResponse, 0)

		tagsID, err := UserService.GetUserTagsID(stayer.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user tags"})
			return
		}

		for _, tagID := range tagsID {
			//タグIDからタグ名を取得する
			tagName, err := UserService.GetTagName(tagID)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tag name"})
				return
			}
			tag := model.TagGetResponse{
				ID:   tagID,
				Name: tagName,
			}
			tagsGetResponse = append(tagsGetResponse, tag)
		}

		stayerGetResponse = append(stayerGetResponse, model.StayerGetResponse{
			ID:     stayer.UserID,
			Name:   userName,
			Room:   roomName,
			RoomID: int(stayer.RoomID),
			Tags:   tagsGetResponse,
		})
	}
	c.JSON(http.StatusOK, stayerGetResponse)
}

func Log(c *gin.Context) {
	RoomService := service.RoomService{}
	UserService := service.UserService{}

	//ページング処理
	page := c.Query("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}

	//ページごとにLogテーブルからデータを取得する
	pageLog, err := RoomService.GetLogsByPage(pageInt)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get logs"})
		return
	}

	logGetResponse := []model.LogGetResponse{}

	for _, log := range pageLog {

		userName, err := UserService.GetUserNameByUserID(log.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user name"})
			return
		}
		roomName, err := RoomService.GetRoomNameByRoomID(log.RoomID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get room name"})
			return
		}

		logGetResponse = append(logGetResponse, model.LogGetResponse{
			ID:      int64(log.ID),
			Name:    userName,
			Room:    roomName,
			StartAt: log.StartAt.Format("2006-01-02 15:04:05"),
			EndAt:   log.EndAt.Format("2006-01-02 15:04:05"),
		})
	}
	c.JSON(http.StatusOK, logGetResponse)
}

// func SimultaneousList(c *gin.Context) {
// 	userID := c.Param("user_id")

// 	RoomService := service.RoomService{}

// 	//userIDをint64に変換
// 	userIDInt, err := strconv.ParseInt(userID, 10, 64)
// 	if err != nil {
// 		c.String(http.StatusBadRequest, "Bad Request")
// 		return
// 	}

// 	SimultaneousList, err := RoomService.GetSimultaneousList(userIDInt)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to "})
// 		return
// 	}

// 	c.JSON(http.StatusOK, SimultaneousList)
// }

func LogGantt(c *gin.Context) {

	RoomService := service.RoomService{}
	GanttLogs, err := RoomService.GetGanttLog()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get gantt log"})
		return
	}
	c.JSON(http.StatusOK, GanttLogs)
}

func LogRefinementSearch(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Query("user-id"), 10, 64)
	if err != nil {
		// c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user_id Bad Request"})
		userID = 0
	}

	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil || limit == 0 {
		limit = 30 //デフォルト値
	}
	offset, err := strconv.ParseInt(c.Query("offset"), 10, 64)
	if err != nil {
		offset = 0 //デフォルト値
	}

	RoomService := service.RoomService{}
	UserService := service.UserService{}

	pageLog, err := RoomService.GetRefinementSearchLogs(userID, limit, offset)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user log"})
		return
	}

	SpecificUserResponseLog := []model.LogGetResponse{}

	for _, log := range pageLog {

		userName, err := UserService.GetUserNameByUserID(log.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user name"})
			return
		}
		roomName, err := RoomService.GetRoomNameByRoomID(log.RoomID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get room name"})
			return
		}

		SpecificUserResponseLog = append(SpecificUserResponseLog, model.LogGetResponse{
			ID:      int64(log.ID),
			Name:    userName,
			Room:    roomName,
			StartAt: log.StartAt.Format("2006-01-02 15:04:05"),
			EndAt:   log.EndAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, SpecificUserResponseLog)
}
