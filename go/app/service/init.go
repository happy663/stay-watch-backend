package service

import (
	"Stay_watch/model"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	db := connect()
	closer, err := db.DB()
	if err != nil {
		return
	}
	defer closer.Close()
	db.AutoMigrate(&model.User{}, &model.Log{}, &model.Room{}, &model.Stayer{}, &model.Tag{}, &model.TagMap{})

	var count int64
	db.Model(&model.User{}).Count(&count)
	if count == 0 {
		//複数のユーザーを作成
		users := []model.User{
			{
				Name:  "kaji",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac00",
			},
			{
				Name:  "ogane",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac01",
			},
			{
				Name:  "miyagawa",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac02",
			},
			{
				Name:  "ayato",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac03",
			},
			{
				Name:  "ken",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac04",
			},
			{
				Name:  "suzaki",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac05",
			},
			{
				Name:  "akito",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac06",
			},
			{
				Name:  "fueta",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac07",
			},
			{
				Name:  "kameda",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac08",
			},
			{
				Name:  "maruyama",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac09",
			},
			{
				Name:  "ohashi",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac0a",
			},
			{
				Name:  "rui",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac0b",
			},
			{
				Name:  "shamo",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac0c",
			},
			{
				Name:  "terada",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac0d",
			},
			{
				Name:  "toyama",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac0e",
			},
			{
				Name:  "ukai",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac0f",
			},
			{
				Name:  "isiguro",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac10",
			},
			{
				Name:  "ao",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac11",
			},
			{
				Name:  "fuma",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac12",
			},
			{
				Name:  "ueji",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac13",
			},
			{
				Name:  "oiwa",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac14",
			},
			{
				Name:  "togawa",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac15",
			},
			{
				Name:  "yada",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac16",
			},
			{
				Name:  "yokoyama",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac17",
			},
			{
				Name:  "kazuo",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac18",
			},
			{
				Name:  "sakai",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac19",
			},
			{
				Name:  "iwaguti",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac1a",
			},
			{
				Name:  "makino",
				Email: "",
				Role:  1,
				UUID:  "e7d61ea3f8dd49c88f2ff2484c07ac1b",
			},
		}
		db.Create(&users)

	}

	db.Model(&model.Room{}).Count(&count)
	if count == 0 {
		rooms := []model.Room{
			{
				Name: "梶研-学生部屋",
			},
			{
				Name: "梶研-スマートルーム",
			},
			{
				Name: "梶研-院生部屋",
			},
			{
				Name: "梶研-FA部屋",
			},
			{
				Name: "梶研-先生部屋",
			},
		}
		db.Create(&rooms)
	}

	db.Model(&model.Tag{}).Count(&count)
	if count == 0 {
		tags := []model.Tag{
			{
				Name: "梶研",
			},
			{
				Name: "ロケーション",
			},
			{
				Name: "インタラクション",
			},
			{
				Name: "センシング",
			},
			{
				Name: "B1",
			},
			{
				Name: "B2",
			},
			{
				Name: "B3",
			},
			{
				Name: "B4",
			},
			{
				Name: "M1",
			},
			{
				Name: "M2",
			},
			{
				Name: "Professor",
			},
		}
		db.Create(&tags)
	}

	db.Model(&model.TagMap{}).Count(&count)
	if count == 0 {
		tagMaps := []model.TagMap{
			{
				UserID: 1,
				TagID:  1,
			},
			{
				UserID: 1,
				TagID:  11,
			},
			{
				UserID: 2,
				TagID:  1,
			},
			{
				UserID: 2,
				TagID:  2,
			},
			{
				UserID: 2,
				TagID:  10,
			},
			{
				UserID: 3,
				TagID:  1,
			},
			{
				UserID: 3,
				TagID:  2,
			},
			{
				UserID: 3,
				TagID:  10,
			},
			{
				UserID: 4,
				TagID:  1,
			},
			{
				UserID: 4,
				TagID:  4,
			},
			{
				UserID: 4,
				TagID:  9,
			},
			{
				UserID: 5,
				TagID:  1,
			},
			{
				UserID: 5,
				TagID:  4,
			},
			{
				UserID: 5,
				TagID:  9,
			},
			{
				UserID: 6,
				TagID:  1,
			},
			{
				UserID: 6,
				TagID:  2,
			},
			{
				UserID: 6,
				TagID:  9,
			},
			{
				UserID: 7,
				TagID:  1,
			},
			{
				UserID: 7,
				TagID:  2,
			},
			{
				UserID: 7,
				TagID:  8,
			},
			{
				UserID: 8,
				TagID:  1,
			},
			{
				UserID: 8,
				TagID:  3,
			},
			{
				UserID: 8,
				TagID:  8,
			},
			{
				UserID: 9,
				TagID:  1,
			},
			{
				UserID: 9,
				TagID:  2,
			},
			{
				UserID: 9,
				TagID:  8,
			},
			{
				UserID: 10,
				TagID:  1,
			},
			{
				UserID: 10,
				TagID:  4,
			},
			{
				UserID: 10,
				TagID:  8,
			},
			{
				UserID: 11,
				TagID:  1,
			},
			{
				UserID: 11,
				TagID:  2,
			},
			{
				UserID: 11,
				TagID:  8,
			},
			{
				UserID: 12,
				TagID:  1,
			},
			{
				UserID: 12,
				TagID:  4,
			},
			{
				UserID: 12,
				TagID:  8,
			},
			{
				UserID: 13,
				TagID:  1,
			},
			{
				UserID: 13,
				TagID:  3,
			},
			{
				UserID: 13,
				TagID:  8,
			},
			{
				UserID: 14,
				TagID:  1,
			},
			{
				UserID: 14,
				TagID:  2,
			},
			{
				UserID: 15,
				TagID:  1,
			},
			{
				UserID: 15,
				TagID:  8,
			},
			{
				UserID: 16,
				TagID:  1,
			},
			{
				UserID: 16,
				TagID:  2,
			},
			{
				UserID: 16,
				TagID:  8,
			},
			{
				UserID: 17,
				TagID:  1,
			},
			{
				UserID: 17,
				TagID:  7,
			},
			{
				UserID: 18,
				TagID:  1,
			},
			{
				UserID: 18,
				TagID:  7,
			},
			{
				UserID: 19,
				TagID:  1,
			},
			{
				UserID: 19,
				TagID:  7,
			},
			{
				UserID: 20,
				TagID:  1,
			},
			{
				UserID: 20,
				TagID:  7,
			},
			{
				UserID: 21,
				TagID:  1,
			},
			{
				UserID: 21,
				TagID:  7,
			},
			{
				UserID: 22,
				TagID:  1,
			},
			{
				UserID: 22,
				TagID:  7,
			},
			{
				UserID: 23,
				TagID:  1,
			},
			{
				UserID: 23,
				TagID:  7,
			},
			{
				UserID: 24,
				TagID:  1,
			},
			{
				UserID: 24,
				TagID:  7,
			},
			{
				UserID: 25,
				TagID:  1,
			},
			{
				UserID: 25,
				TagID:  7,
			},
			{
				UserID: 26,
				TagID:  1,
			},
			{
				UserID: 26,
				TagID:  7,
			},
			{
				UserID: 27,
				TagID:  1,
			},
			{
				UserID: 27,
				TagID:  7,
			},
			{
				UserID: 28,
				TagID:  1,
			},
			{
				UserID: 28,
				TagID:  7,
			},
		}
		db.Create(&tagMaps)
	}

	// db.Model(&model.User{}).Count()

}

func connect() *gorm.DB {

	envPath := "../../.env"
	//test実行時はtest用のenvを読み込むJ
	if strings.HasSuffix(os.Args[0], ".test") {
		envPath = "../../../.env"
	}
	fmt.Println(envPath)
	err := godotenv.Load(
		envPath,
	)
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	dsn := "root:root@tcp(vol_mysql:3306)/app?charset=utf8mb4&parseTime=true&Asia%2FTokyo"
	if os.Getenv("ENVIRONMENT") == "dev" {
		dsn = "root:root@tcp(localhost:33066)/app?charset=utf8mb4&parseTime=true&loc=Asia%2FTokyo"
	}
	fmt.Println(dsn)

	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to the database")
	return gormDB
}
