package controller

import (
	"Stay_watch/model"
	"Stay_watch/service"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func Detail(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}

func CreateUser(c *gin.Context) {
	RegistrationUserForm := model.RegistrationUserForm{}
	c.Bind(&RegistrationUserForm)

	UserService := service.UserService{}
	//userIDがないなら新規登録
	if RegistrationUserForm.ID == 0 {
		user := model.User{
			Name:  RegistrationUserForm.Name,
			Email: RegistrationUserForm.Email,
			Role:  RegistrationUserForm.Role,
			UUID:  UserService.NewUUID(),
		}

		err := UserService.RegisterUser(&user)
		if err != nil {
			fmt.Printf("Cannnot register user: %v", err)
			c.String(http.StatusInternalServerError, "Server Error")
			return
		}
	}

	//userIDがあるなら更新
	if RegistrationUserForm.ID != 0 {
		//userNameが空なので、userIDからuserNameを取得する
		err := UserService.UpdateUser(
			int(RegistrationUserForm.ID),
			RegistrationUserForm.Email,
		)

		if err != nil {
			c.String(http.StatusInternalServerError, "Server Error")
			return
		}
	}

	if !strings.HasSuffix(os.Args[0], ".test") {
		mailService := service.MailService{}
		mailService.SendMail("滞在ウォッチユーザ登録の完了のお知らせ", "ユーザ登録が完了したので滞在ウォッチを閲覧することが可能になりました\n一度プロジェクトをリセットしたので再度ログインお願いします。\nアプリドメイン\nhttps://stay-watch-go.kajilab.tk/", RegistrationUserForm.Email)
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "ok",
	})
}

func UserList(c *gin.Context) {

	UserService := service.UserService{}
	communityId, _ := strconv.ParseInt(c.Param("communityId"), 10, 64) // string -> int64

	if c.Query("fields") == "admin" {
		// 編集画面のユーザの情報を返す
		fmt.Print("コミュニティID：")
		fmt.Println(communityId)
		fmt.Printf("%T\n", communityId)

		edit_users, err := UserService.GetEditUsersByCommunityId(communityId)
		if err != nil {
			c.String(http.StatusInternalServerError, "Server Error")
		}

		userEditorResponse := []model.UserEditorResponse{}

		for _, user := range edit_users {

			tags := make([]model.TagGetResponse, 0)
			tagsID, err := UserService.GetUserTagsID(int64(user.Model.ID))
			if err != nil {
				c.String(http.StatusInternalServerError, "Server Error")
				return
			}

			for _, tagID := range tagsID {
				//タグIDからタグ名を取得する
				tagName, err := UserService.GetTagName(tagID)
				if err != nil {
					c.String(http.StatusInternalServerError, "Server Error")
					return
				}
				tag := model.TagGetResponse{
					ID:   tagID,
					Name: tagName,
				}
				tags = append(tags, tag)
			}

			userEditorResponse = append(userEditorResponse, model.UserEditorResponse{
				ID:         int64(user.ID),
				Name:       user.Name,
				Uuid:       user.UUID,
				Email:      user.Email,
				Role:       user.Role,
				BeaconType: user.BeaconTypeId,
				BeaconName: "android",
				Tags:       tags,
			})
		}
		c.JSON(http.StatusOK, userEditorResponse)

	} else {
		// 一覧画面でのユーザ情報
		users, err := UserService.GetAllUser()
		if err != nil {
			c.String(http.StatusInternalServerError, "Server Error")
			return
		}

		userInformationGetResponse := []model.UserInformationGetResponse{}

		for _, user := range users {

			tags := make([]model.TagGetResponse, 0)
			tagsID, err := UserService.GetUserTagsID(int64(user.Model.ID))
			if err != nil {
				c.String(http.StatusInternalServerError, "Server Error")
				return
			}

			for _, tagID := range tagsID {
				//タグIDからタグ名を取得する
				tagName, err := UserService.GetTagName(tagID)
				if err != nil {
					c.String(http.StatusInternalServerError, "Server Error")
					return
				}
				tag := model.TagGetResponse{
					ID:   tagID,
					Name: tagName,
				}
				tags = append(tags, tag)
			}

			userInformationGetResponse = append(userInformationGetResponse, model.UserInformationGetResponse{
				ID:   int64(user.ID),
				Name: user.Name,
				Tags: tags,
			})
		}
		c.JSON(http.StatusOK, userInformationGetResponse)
	}
}

func ExtendedUserList(c *gin.Context) {

	UserService := service.UserService{}
	users, err := UserService.GetAllUser()
	if err != nil {
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}

	extendedUserInformationGetResponses := []model.ExtendedUserInformationGetResponse{}
	// userInformationGetResponse := []model.UserInformationGetResponse{}

	for _, user := range users {

		tags := make([]model.TagGetResponse, 0)
		tagsID, err := UserService.GetUserTagsID(int64(user.Model.ID))
		if err != nil {
			c.String(http.StatusInternalServerError, "Server Error")
			return
		}

		for _, tagID := range tagsID {
			//タグIDからタグ名を取得する
			tagName, err := UserService.GetTagName(tagID)
			if err != nil {
				c.String(http.StatusInternalServerError, "Server Error")
				return
			}
			tag := model.TagGetResponse{
				ID:   tagID,
				Name: tagName,
			}
			tags = append(tags, tag)
		}

		extendedUserInformationGetResponses = append(extendedUserInformationGetResponses, model.ExtendedUserInformationGetResponse{
			ID:   int64(user.ID),
			Name: user.Name,
			Tags: tags,
			Role: user.Role,
			Uuid: user.UUID,
		})
	}

	c.JSON(http.StatusOK, extendedUserInformationGetResponses)
}

// for _, user := range users {

// 	tags := make([]model.TagGetResponse, 0)
// 	tagsID, err := UserService.GetUserTagsID(int64(user.Model.ID))
// 	if err != nil {
// 		c.String(http.StatusInternalServerError, "Server Error")
// 		return
// 	}

// 	for _, tagID := range tagsID {
// 		//タグIDからタグ名を取得する
// 		tagName, err := UserService.GetTagName(tagID)
// 		if err != nil {
// 			c.String(http.StatusInternalServerError, "Server Error")
// 			return
// 		}
// 		tag := model.TagGetResponse{
// 			ID:   tagID,
// 			Name: tagName,
// 		}
// 		tags = append(tags, tag)
// 	}

// 	userInformationGetResponse = append(userInformationGetResponse, model.UserInformationGetResponse{
// 		ID:   int64(user.ID),
// 		Name: user.Name,
// 		Tags: tags,
// 	})
// }

// c.JSON(http.StatusOK, userInformationGetResponse)

func Attendance(c *gin.Context) {

	//構造体定義
	type Meeting struct {
		ID int64 `json:"meetingID"`
	}
	var meeting Meeting
	err := c.Bind(&meeting)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(meeting.ID)
	UserService := service.UserService{}
	//attendaance_tmpテーブルから全てのデータを取得する
	allAttendancesTmp, err := UserService.GetAllAttendancesTmp()
	if err != nil {
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}

	isExist := true
	flagCount := 0
	if meeting.ID == 2 {
		for i := 0; i < 16; i++ {
			if allAttendancesTmp[i].Flag == 0 {
				flagCount++
			}
		}
		if flagCount == 16 {
			isExist = false
		}
	}
	if meeting.ID == 1 {
		for i := 16; i < 28; i++ {
			if allAttendancesTmp[i].Flag == 0 {
				flagCount++
			}
		}
		if flagCount == 12 {
			isExist = false
		}
	}

	ExcelService := service.ExcelService{}
	if isExist {
		ExcelService.WriteExcel(allAttendancesTmp, meeting.ID)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// func SimultaneousStayUserList(c *gin.Context) {
// 	userID := c.Param("user_id")
// 	//int64に変換
// 	userIDInt64, err := strconv.ParseInt(userID, 10, 64)
// 	if err != nil {
// 		c.String(http.StatusInternalServerError, "Server Error")
// 		return
// 	}

// 	UserService := service.UserService{}
// 	RoomService := service.RoomService{}

// 	logs, err := RoomService.GetLogByUserAndDate(userIDInt64, 14)
// 	if err != nil {
// 		c.String(http.StatusInternalServerError, "Server Error")
// 		return
// 	}
// 	simultaneousStayUserGetResponses, err := UserService.GetSameTimeUser(logs)
// 	if err != nil {
// 		c.String(http.StatusInternalServerError, "Server Error")
// 		return
// 	}

// 	c.JSON(200, simultaneousStayUserGetResponses)
// }

func Check(c *gin.Context) {

	firebaseUserInfo, err := verifyCheck(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "invalid token",
		})
		return
	}

	UserService := service.UserService{}
	user, err := UserService.GetUserByEmail(firebaseUserInfo["Email"])
	if err != nil {
		fmt.Printf("Cannnot find user: %v", err)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}
	fmt.Println(user)

	//メールアドレスが存在しない場合はUserは存在しないのでリクエスト失敗
	if (user == model.User{}) {
		c.JSON(http.StatusForbidden, gin.H{
			"status": "権限がありません 管理者にユーザ追加を依頼してください",
		})
		return
	}

	userRole := model.UserRoleGetResponse{
		ID:   int64(user.ID),
		Role: user.Role,
	}

	c.JSON(http.StatusOK,
		userRole,
	)
}

func SignUp(c *gin.Context) {

	firebaseUserInfo, err := verifyCheck(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "invalid token",
		})
		return
	}

	UserService := service.UserService{}
	user, err := UserService.GetUserByEmail(firebaseUserInfo["Email"])
	if err != nil {
		fmt.Printf("Cannnot find user: %v", err)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}
	fmt.Println(user)

	//メールアドレスが存在しない場合はUserは存在しないのでリクエスト失敗
	if (user == model.User{}) {
		c.JSON(http.StatusForbidden, gin.H{
			"status": "権限がありません 管理者にユーザ追加を依頼してください",
		})
		return
	}

	// userRole := model.UserRoleGetResponse{
	// 	ID:   int64(user.ID),
	// 	Role: user.Role,
	// }

	c.JSON(http.StatusOK,
		user,
	)
}
