package service

import (
	"Stay_watch/model"
	"Stay_watch/util"
	"fmt"

	// "log"
	"time"
)

type RoomService struct{}

func (RoomService) CreateLog(Log *model.Log) error {

	DbEngine := connect()
	closer, err := DbEngine.DB()
	if err != nil {
		return err
	}
	defer closer.Close()

	result := DbEngine.Create(&Log)
	if result.Error != nil {
		return fmt.Errorf(" failed to create log: %w", result.Error)
	}
	return nil
}

// 該当ユーザが存在するか確認
func (RoomService) GetStayer(userID int64) (error, bool) {
	DbEngine := connect()
	closer, err := DbEngine.DB()
	if err != nil {
		return err, false
	}
	defer closer.Close()

	//レコードの数を取得する
	var count int64
	result := DbEngine.Table("stayers").Where("user_id=?", userID).Count(&count)
	if result.Error != nil {
		return fmt.Errorf(" failed to get count: %w", result.Error), false
	}
	fmt.Println("count=", count)

	if count == 0 {
		return nil, false
	}
	return nil, true
}

// 滞在者全体を取得する
func (RoomService) GetAllStayer() ([]model.Stayer, error) {
	DbEngine := connect()
	closer, err := DbEngine.DB()
	if err != nil {
		return nil, err
	}
	defer closer.Close()
	stayers := make([]model.Stayer, 0)
	result := DbEngine.Table("stayers").Find(&stayers)
	if result.Error != nil {
		return nil, fmt.Errorf(" failed to get all stayer: %w", result.Error)
	}

	return stayers, nil
}

// 部屋のID、建物の名前、建物のID、部屋の範囲をデータベースへ保存
func (RoomService) UpdateRoom(roomID int, room_name string, buildingID int, polygon string) error {
	DbEngine := connect()
	closer, err := DbEngine.DB()
	if err != nil {
		return err
	}
	defer closer.Close()
	result := DbEngine.Model(&model.Room{}).Where("id = ?", roomID).Updates(model.Room{Name: room_name, Polygon: polygon, BuildingID: int64(buildingID)}) // 今は部屋名と範囲だけ
	if result.Error != nil {
		fmt.Printf("ユーザ更新失敗 %v", result.Error)
		return result.Error
	}
	return nil
}

func (RoomService) GetAllRooms() ([]model.Room, error) {
	DbEngine := connect()
	closer, err := DbEngine.DB()
	if err != nil {
		return nil, err
	}
	defer closer.Close()
	rooms := make([]model.Room, 0)
	result := DbEngine.Table("rooms").Find(&rooms)
	if result.Error != nil {
		return nil, fmt.Errorf(" failed to get all stayer: %w", result.Error)
	}

	return rooms, nil
}

//滞在者の一部を取得する
// func (RoomService) GetStayerByRoomID(roomID int64) ([]model.Stayer, error) {
// 	DbEngine := connect()
// 	closer, err := DbEngine.DB()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer closer.Close()
// 	stayers := make([]model.Stayer, 0)
// 	err := DbEngine.Table("stayer").Where("room_id=?", roomID).Find(&stayers)
// 	if err != nil {
// 		return nil, fmt.Errorf(" failed: %w", err)
// 	}
// 	return stayers, nil
// }

func (RoomService) CreateStayer(stayer *model.Stayer) error {
	DbEngine := connect()
	closer, err := DbEngine.DB()
	if err != nil {
		return err
	}
	defer closer.Close()

	result := DbEngine.Create(&stayer)
	if result.Error != nil {
		return fmt.Errorf(" failed: to create stayer: %w", result.Error)
	}
	return nil
}

func (RoomService) UpdateStayer(stayer *model.Stayer) error {
	DbEngine := connect()
	closer, err := DbEngine.DB()
	if err != nil {
		return err
	}
	defer closer.Close()

	result := DbEngine.Model(&model.Stayer{}).Where("user_id=?", stayer.UserID).Updates(&stayer)

	if result.Error != nil {
		return fmt.Errorf(" failed to update stayer: %w", result.Error)
	}

	return nil
}

func (RoomService) DeleteStayer(userID int64) error {
	DbEngine := connect()
	closer, err := DbEngine.DB()
	if err != nil {
		return err
	}
	defer closer.Close()

	result := DbEngine.Unscoped().Where("user_id=?", userID).Delete(&model.Stayer{})
	if result.Error != nil {
		return fmt.Errorf(" failed to delete stayer: %w", result.Error)
	}
	return nil
}

func (RoomService) UpdateEndAt(userID int64, endAt time.Time) error {
	DbEngine := connect()
	closer, err := DbEngine.DB()
	if err != nil {
		return err
	}
	defer closer.Close()

	result := DbEngine.Where("user_id=? ", userID).Last(&model.Log{}).Update("end_at", endAt)
	if result.Error != nil {
		return fmt.Errorf(" failed to update end_at: %w", result.Error)
	}

	return nil
}

// 指定したuserと現在の日付から指定した日付以内のログを取得する
func (RoomService) GetLogByUserAndDate(userID int64, date int64) ([]model.Log, error) {
	// currentTime := time.Now()
	logs := make([]model.Log, 0)
	// err := DbEngine.Table("log").Asc("start_at").Where("user_id=?", userID).And("start_at>=?", currentTime.AddDate(0, 0, -int(date)).Format("2006-01-02 15:04:05")).Find(&logs)
	// if err != nil {
	// 	return nil, fmt.Errorf(" failed: %w", err)
	// }
	return logs, nil
}

func (RoomService) GetLogWithinDate(date int64) ([]model.Log, error) {

	DbEngine := connect()
	closer, err := DbEngine.DB()
	if err != nil {
		return nil, err
	}
	currentTime := time.Now()
	logs := make([]model.Log, 0)

	// 現在の日付から指定した日付以内のログを取得する
	// err = DbEngine.Table("log").Asc("start_at").Where("start_at>=?", currentTime.AddDate(0, 0, -int(date)).Format("2006-01-02 15:04:05")).Find(&logs)
	result := DbEngine.Order("start_at asc").Where("start_at>=?", currentTime.AddDate(0, 0, -int(date)).Format("2006-01-02 15:04:05")).Find(&logs)
	if result.Error != nil {
		return nil, fmt.Errorf(" failed: %w", result.Error)
	}

	if err != nil {
		return nil, fmt.Errorf(" failed: %w", err)
	}

	defer closer.Close()

	return logs, nil
}

func (RoomService) GetGanttLog() ([]model.SimulataneousStayLogGetResponse, error) {

	DbEngine := connect()
	closer, err := DbEngine.DB()
	if err != nil {
		return nil, err
	}
	defer closer.Close()

	RoomService := RoomService{}
	UserService := UserService{}
	//1週間以内のログを取得
	logs, err := RoomService.GetLogWithinDate(7)
	fmt.Println(logs)
	if err != nil {
		return nil, err
	}

	dates := make([]string, 0)
	roomIDs := make([]int64, 0)
	for _, log := range logs {
		dates = append(dates,
			log.StartAt.Format("2006-01-02"))

		roomIDs = append(roomIDs, log.RoomID)
	}
	util := util.Util{}
	uniqueDates := util.SliceUniqueString(dates)
	uniqueRoomIDs := util.SliceUniqueNumber(roomIDs)

	simulataneousStayLogsGetResponse := make([]model.SimulataneousStayLogGetResponse, 0)

	for index, dates := range uniqueDates {
		simulataneousStayLogGetResponse := model.SimulataneousStayLogGetResponse{}
		simulataneousStayLogGetResponse.Date = dates
		simulataneousStayLogGetResponse.ID = int64(index)
		roomsGetResponse := make([]model.RoomGetResponse, 0)
		for _, roomID := range uniqueRoomIDs {

			roomGetResponse := model.RoomGetResponse{}
			roomGetResponse.ID = roomID
			roomGetResponse.Name, err = RoomService.GetRoomNameByRoomID(roomID)
			if err != nil {
				return nil, err
			}
			sortlogs := make([]model.Log, 0)

			result := DbEngine.Where("room_id=?", roomID).Where("start_at like ?", dates+"%").Find(&sortlogs)
			if result.Error != nil {
				return nil, fmt.Errorf(" failed: %w", err)
			}

			stayTimes := make([]model.StayTime, 0)
			for _, log := range sortlogs {
				stayTime := model.StayTime{}
				locationTime, err := util.ConvertDatetimeToLocationTime(log.StartAt.Format("2006-01-02 15:04:05"), "Asia/Tokyo")
				if err != nil {
					return nil, err
				}
				unixMilli := util.TimeToUnixMilli(locationTime)
				stayTime.StartAt = unixMilli

				// 終了時間が初期値　2016-01-01 00:00:00の場合は現在時刻を入れる
				if log.EndAt.Format("2006-01-01 15:04:05") == "2016-01-01 00:00:00" {
					locationTime, err = util.ConvertDatetimeToLocationTime(time.Now().Format("2006-01-02 15:04:05"), "Asia/Tokyo")
					if err != nil {
						return nil, err
					}
					unixMilli = util.TimeToUnixMilli(locationTime)
				} else {
					locationTime, err = util.ConvertDatetimeToLocationTime(log.EndAt.Format("2006-01-02 15:04:05"), "Asia/Tokyo")
					if err != nil {
						return nil, err
					}
					unixMilli = util.TimeToUnixMilli(locationTime)
				}

				fmt.Println(unixMilli)
				stayTime.EndAt = unixMilli

				userName, err := UserService.GetUserNameByUserID(log.UserID)
				if err != nil {
					return nil, err
				}
				stayTime.UserName = userName
				stayTime.Color = "green"
				stayTime.ID = int64(log.ID)
				stayTimes = append(stayTimes, stayTime)
			}
			roomGetResponse.StayTimes = stayTimes
			roomsGetResponse = append(roomsGetResponse, roomGetResponse)
		}
		simulataneousStayLogGetResponse.Rooms = roomsGetResponse
		simulataneousStayLogsGetResponse = append(simulataneousStayLogsGetResponse, simulataneousStayLogGetResponse)
	}

	fmt.Println(simulataneousStayLogsGetResponse)

	return simulataneousStayLogsGetResponse, nil

}

func (RoomService) GetRefinementSearchLogs(userID int64, limit int64, offset int64) ([]model.Log, error) {
	DbEngine := connect()
	closer, err := DbEngine.DB()
	if err != nil {
		return nil, err
	}
	defer closer.Close()

	//ログデータ初期化
	logs := make([]model.Log, 0)

	result := DbEngine.Table("logs")
	if userID != 0 {
		result.Where("user_id=?", userID)
	}

	result.Order("id DESC").Limit(int(limit)).Offset(int(offset)).Find(&logs)
	if result.Error != nil {
		return nil, result.Error
	}
	return logs, nil
}

// func (RoomService) GetSimultaneousList(userID int64) ([]model.SimulataneousStayLogGetResponse, error) {
// 	RoomService := RoomService{}
// 	logs, err := RoomService.GetLogByUserAndDate(userID, 14)
// 	if err != nil {
// 		return nil, fmt.Errorf(" failed: %w", err)
// 	}

// 	dates := make([]string, 0)
// 	roomIDs := make([]int64, 0)
// 	for _, log := range logs {
// 		dates = append(dates, log.StartAt.Format("2006-01-02"))
// 		roomIDs = append(roomIDs, log.RoomID)
// 	}

// 	util := util.Util{}

// 	//滞在した日付
// 	uniqueDates := util.SliceUniqueString(dates)
// 	//滞在した部屋
// 	uniqueRoomIDs := util.SliceUniqueNumber(roomIDs)

// 	dateSql := ""
// 	for index, uniqueDate := range uniqueDates {
// 		dateSql += "start_at like '" + uniqueDate + "%' or "
// 		if index == len(uniqueDates)-1 {
// 			dateSql = dateSql[:len(dateSql)-4]
// 		}
// 	}

// 	roomSql := ""
// 	for index, uniqueRoomID := range uniqueRoomIDs {
// 		roomSql += "room_id=" + fmt.Sprintf("%d", uniqueRoomID) + " or "
// 		if index == len(uniqueRoomIDs)-1 {
// 			roomSql = roomSql[:len(roomSql)-4]
// 		}
// 	}

// 	sameDayAndRoomlogs := make([]model.Log, 0)
// 	err = DbEngine.Table("log").Where(dateSql).And(roomSql).OrderBy("date_format(start_at,'%Y-%m-%d') ,room_id ").Find(&sameDayAndRoomlogs)
// 	if err != nil {
// 		return nil, fmt.Errorf(" failed: %w", err)
// 	}

// 	UserService := UserService{}

// 	simulataneousStayLogsGetResponse := make([]model.SimulataneousStayLogGetResponse, 0)
// 	roomsGetResponse := make([]model.RoomGetResponse, 0)
// 	stayTimes := make([]model.StayTime, 0)
// 	simulataneousStayLogGetResponse := model.SimulataneousStayLogGetResponse{}
// 	dateCount := 1

// 	for index, sameDayAndRoomlog := range sameDayAndRoomlogs {
// 		//Ascで昇順にしているため、違うuserIDになるまでループする
// 		roomGetResponse := model.RoomGetResponse{}
// 		stayTime := model.StayTime{}

// 		//最後のindexの時
// 		if index == len(sameDayAndRoomlogs)-1 {
// 			stayTime.ID = int64(sameDayAndRoomlog.ID)
// 			userName, err := UserService.GetUserNameByUserID(sameDayAndRoomlog.UserID)
// 			if err != nil {
// 				return nil, fmt.Errorf(" failed: %w", err)
// 			}
// 			stayTime.UserName = userName

// 			locationTime, err := util.ConvertDatetimeToLocationTime(sameDayAndRoomlog.StartAt.Format("2006-01-02 15:04:05"), "Asia/Tokyo")
// 			if err != nil {
// 				return nil, fmt.Errorf(" failed: %w", err)
// 			}
// 			unixMilli := util.TimeToUnixMilli(locationTime)
// 			stayTime.StartAt = unixMilli

// 			locationTime, err = util.ConvertDatetimeToLocationTime(sameDayAndRoomlog.EndAt.Format("2006-01-02 15:04:05"), "Asia/Tokyo")

// 			if err != nil {
// 				log.Fatal(err.Error())
// 				return nil, err
// 			}
// 			unixMilli = util.TimeToUnixMilli(locationTime)
// 			stayTime.EndAt = unixMilli

// 			//検索対象者は赤色にする
// 			if userID == sameDayAndRoomlog.UserID {
// 				stayTime.Color = "red"
// 			} else {
// 				stayTime.Color = "green"
// 			}
// 			stayTimes = append(stayTimes, stayTime)

// 			roomGetResponse.ID = sameDayAndRoomlog.RoomID
// 			roomName, err := RoomService.GetRoomNameByRoomID(sameDayAndRoomlog.RoomID)
// 			if err != nil {
// 				return nil, fmt.Errorf("failed: %w", err)
// 			}
// 			roomGetResponse.Name = roomName
// 			roomGetResponse.StayTimes = stayTimes
// 			roomsGetResponse = append(roomsGetResponse, roomGetResponse)
// 			stayTimes = nil

// 			//後でuniqueDateのindexに置き換えるかも
// 			simulataneousStayLogGetResponse.ID = int64(dateCount)
// 			dateCount++
// 			simulataneousStayLogGetResponse.Date = sameDayAndRoomlog.StartAt.Format("2006-01-02")
// 			simulataneousStayLogGetResponse.Rooms = roomsGetResponse
// 			simulataneousStayLogsGetResponse = append(simulataneousStayLogsGetResponse, simulataneousStayLogGetResponse)
// 			roomsGetResponse = nil
// 		} else {
// 			if sameDayAndRoomlog.StartAt.Format("2006-01-02") == sameDayAndRoomlogs[index+1].StartAt.Format("2006-01-02") {
// 				if sameDayAndRoomlog.RoomID == sameDayAndRoomlogs[index+1].RoomID {
// 					stayTime.ID = int64(sameDayAndRoomlog.ID)
// 					userName, err := UserService.GetUserNameByUserID(sameDayAndRoomlog.UserID)
// 					if err != nil {
// 						log.Fatal(err.Error())
// 						return nil, err
// 					}
// 					stayTime.UserName = userName

// 					locationTime, err := util.ConvertDatetimeToLocationTime(sameDayAndRoomlog.StartAt.Format("2006-01-02"), "Asia/Tokyo")
// 					if err != nil {
// 						log.Fatal(err.Error())
// 						return nil, err
// 					}
// 					unixMilli := util.TimeToUnixMilli(locationTime)
// 					stayTime.StartAt = unixMilli

// 					locationTime, err = util.ConvertDatetimeToLocationTime(sameDayAndRoomlog.StartAt.Format("2006-01-02"), "Asia/Tokyo")
// 					if err != nil {
// 						log.Fatal(err.Error())
// 						return nil, err
// 					}
// 					unixMilli = util.TimeToUnixMilli(locationTime)
// 					stayTime.EndAt = unixMilli

// 					//検索対象者は赤色にする
// 					if userID == sameDayAndRoomlog.UserID {
// 						stayTime.Color = "red"
// 					} else {
// 						stayTime.Color = "green"
// 					}
// 					stayTimes = append(stayTimes, stayTime)
// 				} else {
// 					stayTime.ID = int64(sameDayAndRoomlog.ID)
// 					userName, err := UserService.GetUserNameByUserID(sameDayAndRoomlog.UserID)
// 					if err != nil {
// 						log.Fatal(err.Error())
// 						return nil, err
// 					}
// 					stayTime.UserName = userName

// 					locationTime, err := util.ConvertDatetimeToLocationTime(sameDayAndRoomlog.StartAt.Format("2006-01-02"), "Asia/Tokyo")
// 					if err != nil {
// 						log.Fatal(err.Error())
// 						return nil, err
// 					}
// 					unixMilli := util.TimeToUnixMilli(locationTime)
// 					stayTime.StartAt = unixMilli

// 					locationTime, err = util.ConvertDatetimeToLocationTime(sameDayAndRoomlog.StartAt.Format("2006-01-02"), "Asia/Tokyo")
// 					if err != nil {
// 						log.Fatal(err.Error())
// 						return nil, err
// 					}
// 					unixMilli = util.TimeToUnixMilli(locationTime)
// 					stayTime.EndAt = unixMilli

// 					//検索対象者は赤色にする
// 					if userID == sameDayAndRoomlog.UserID {
// 						stayTime.Color = "red"
// 					} else {
// 						stayTime.Color = "green"
// 					}
// 					stayTimes = append(stayTimes, stayTime)

// 					roomGetResponse.ID = sameDayAndRoomlog.RoomID
// 					roomName, err := RoomService.GetRoomNameByRoomID(sameDayAndRoomlog.RoomID)
// 					if err != nil {
// 						log.Fatal(err.Error())
// 						return nil, err
// 					}
// 					roomGetResponse.Name = roomName
// 					roomGetResponse.StayTimes = stayTimes
// 					roomsGetResponse = append(roomsGetResponse, roomGetResponse)
// 					stayTimes = nil
// 				}
// 			} else {
// 				stayTime.ID = int64(sameDayAndRoomlog.ID)
// 				userName, err := UserService.GetUserNameByUserID(sameDayAndRoomlog.UserID)
// 				if err != nil {
// 					log.Fatal(err.Error())
// 					return nil, err
// 				}
// 				stayTime.UserName = userName

// 				locationTime, err := util.ConvertDatetimeToLocationTime(sameDayAndRoomlog.StartAt.Format("2006-01-02"), "Asia/Tokyo")
// 				if err != nil {
// 					log.Fatal(err.Error())
// 					return nil, err
// 				}
// 				unixMilli := util.TimeToUnixMilli(locationTime)
// 				stayTime.StartAt = unixMilli

// 				locationTime, err = util.ConvertDatetimeToLocationTime(sameDayAndRoomlog.StartAt.Format("2006-01-02"), "Asia/Tokyo")
// 				if err != nil {
// 					log.Fatal(err.Error())
// 					return nil, err
// 				}
// 				unixMilli = util.TimeToUnixMilli(locationTime)
// 				stayTime.EndAt = unixMilli

// 				//検索対象者は赤色にする
// 				if userID == sameDayAndRoomlog.UserID {
// 					stayTime.Color = "red"
// 				} else {
// 					stayTime.Color = "green"
// 				}
// 				stayTimes = append(stayTimes, stayTime)

// 				roomGetResponse.ID = sameDayAndRoomlog.RoomID
// 				roomName, err := RoomService.GetRoomNameByRoomID(sameDayAndRoomlog.RoomID)
// 				if err != nil {
// 					log.Fatal(err.Error())
// 					return nil, err
// 				}
// 				roomGetResponse.Name = roomName
// 				roomGetResponse.StayTimes = stayTimes
// 				roomsGetResponse = append(roomsGetResponse, roomGetResponse)
// 				stayTimes = nil

// 				//後でuniqueDateのindexに置き換えるかも
// 				simulataneousStayLogGetResponse.ID = int64(dateCount)
// 				dateCount++
// 				simulataneousStayLogGetResponse.Date = sameDayAndRoomlog.StartAt.Format("2006-01-02")
// 				simulataneousStayLogGetResponse.Rooms = roomsGetResponse
// 				simulataneousStayLogsGetResponse = append(simulataneousStayLogsGetResponse, simulataneousStayLogGetResponse)
// 				roomsGetResponse = nil
// 			}
// 		}
// 	}

// 	fmt.Println(simulataneousStayLogsGetResponse)

// 	return simulataneousStayLogsGetResponse, nil
// }

func (RoomService) GetTimesFromStartAtAndEntAt(startAt string, endAt string) ([]string, error) {
	times := []string{}
	startAtTime, err := time.Parse("2006-01-02 15:04:05", startAt)
	if err != nil {
		return nil, err
	}
	endAtTime, err := time.Parse("2006-01-02 15:04:05", endAt)
	if err != nil {
		return nil, err
	}
	for startAtTime.Before(endAtTime) {
		times = append(times, startAtTime.Format("15:04"))
		startAtTime = startAtTime.Add(time.Minute * 15)
	}
	return times, nil
}

// ルームIDからルームの名前を取得する
func (RoomService) GetRoomNameByRoomID(roomID int64) (string, error) {
	DbEngine := connect()
	closer, err := DbEngine.DB()
	if err != nil {
		return "", err
	}
	defer closer.Close()

	room := model.Room{}
	result := DbEngine.Take(&room, roomID)
	if result.Error != nil {
		fmt.Printf("Cannot get room: %v", result.Error)
		return "", result.Error
	}
	return room.Name, nil
}

// pageごとに30件のログを取得する
func (RoomService) GetLogsByPage(page int) ([]model.Log, error) {
	DbEngine := connect()
	closer, err := DbEngine.DB()
	if err != nil {
		return nil, err
	}
	defer closer.Close()

	if page == 1 {
		page = 0
	} else {
		page = (page - 1) * 30
	}
	logs := make([]model.Log, 0)
	result := DbEngine.Order("id desc").Limit(30).Offset(page).Find(&logs)
	if result.Error != nil {
		fmt.Printf("Cannot get logs: %v", result.Error)
		return nil, result.Error
	}

	return logs, nil
}

// 指定した時間のログを取得する
func (RoomService) GetLogsFromStartAtAndEntAt(startAt string, endAt string) ([]model.Log, error) {
	DbEngine := connect()
	closer, err := DbEngine.DB()
	if err != nil {
		return nil, err
	}
	defer closer.Close()
	logs := make([]model.Log, 0)
	startAtTime, err := time.Parse("2006-01-02 15:04:05", startAt)
	if err != nil {
		return nil, err
	}
	endAtTime, err := time.Parse("2006-01-02 15:04:05", endAt)
	if err != nil {
		return nil, err
	}
	// err = DbEngine.Where("start_at>=? and start_at<=?", startAtTime, endAtTime).Find(&logs)
	result := DbEngine.Where("start_at>=? and start_at<=?", startAtTime, endAtTime).Find(&logs)
	if result.Error != nil {
		fmt.Printf("failed to get logs from startAt and endAt: %v", result.Error)
		return nil, result.Error
	}

	return logs, nil
}
