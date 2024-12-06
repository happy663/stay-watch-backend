package model

import (
	"time"

	"gorm.io/gorm"
)

// DBで user_name ならここでは UserName こうしてDBの値を取得している

type User struct {
	gorm.Model

	UUID        string
	Name        string
	Email       string
	Role        int64
	BeaconId    int64
	CommunityId int64
}

type DeletedUser struct {
	gorm.Model

	UUID        string
	Name        string
	Email       string
	Role        int64
	BeaconId    int64
	CommunityId int64
	UserId      int64
}

type Log struct {
	gorm.Model
	RoomID  int64
	StartAt time.Time
	EndAt   time.Time
	UserID  int64
	Rssi    int64
}

type Room struct {
	gorm.Model
	Name        string
	BuildingID  int64
	CommunityID int64
	Polygon     string
	RecieverID  string
}

type Building struct {
	gorm.Model
	Name    string
	MapFile string
}

type Stayer struct {
	gorm.Model
	UserID int64
	RoomID int64
	Rssi   int64
}

type Tag struct {
	gorm.Model
	Name        string
	CommunityId int64
}

type TagMap struct {
	gorm.Model
	UserID int64
	TagID  int64
}

type Attendance struct {
	gorm.Model
	UserID int64
	Date   string
	Flag   bool
}

type AttendanceTmp struct {
	UserID int64 `json:"userID" xorm:"user_id"`
	Flag   int64 `json:"exit" xorm:"flag"`
}

type Beacon struct {
	gorm.Model
	Type         string
	UuidEditable bool
}

type Community struct {
	gorm.Model
	Name string
}

type EditedLog struct {
	gorm.Model
	UserId    int64
	Date      time.Time `gorm:"type:date"`
	Reporting time.Time `gorm:"type:datetime"`
	Leave     time.Time `gorm:"type:datetime"`
}

type Cluster struct {
	gorm.Model
	Date      time.Time `gorm:"type:date"`
	Reporting bool
	Average   float64
	Sd        float64
	Count     int64
	UserId    int64
}
