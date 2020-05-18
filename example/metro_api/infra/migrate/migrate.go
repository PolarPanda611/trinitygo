package migrate

import (
	"metro_api/infra/db"
	"metro_api/infra/util"
	"strconv"

	"metro_api/domain/model"

	"github.com/prometheus/common/log"
)

// Migrate sync db and default value
func Migrate() {
	db.DB.AutoMigrate(&model.User{})
	db.DB.AutoMigrate(&model.Group{})
	db.DB.AutoMigrate(&model.Line{})
	db.DB.AutoMigrate(&model.Metro{})
	db.DB.AutoMigrate(&model.Station{})
	db.DB.AutoMigrate(&model.RunHistory{})

	initMetro()
	initLineData()
	initStationData()
	initGroup()
	initUser()

}
func initGroup() {
	defaultGroup := []model.Group{
		model.Group{
			Code: "G9999",
			Name: "系统管理员",
		},
		model.Group{
			Code: "G0001",
			Name: "站长",
		},
		model.Group{
			Code: "G0002",
			Name: "值班员",
		},
	}
	for _, v := range defaultGroup {
		db.DB.Where("code = ?", v.Code).FirstOrCreate(&v)
	}
}
func initLineData() {
	defaultLine := []model.Line{
		model.Line{
			Code: "1",
			Name: "一号线",
		},
	}
	for _, v := range defaultLine {
		db.DB.Where("code = ?", v.Code).FirstOrCreate(&v)
	}
}

func initStationData() {
	var lineOne model.Line
	if err := db.DB.Where("code = ?", "1").First(&lineOne).Error; err != nil {
		log.Fatalf("find line one err : %v ", err)
	}
	defaultStation := []string{
		"莘庄",
		"外环路",
		"莲花百路",
		"锦江乐园",
		"上海南站",
		"漕宝路",
		"上海体育馆",
		"徐家汇",
		"衡山知路",
		"常熟路",
		"陕西南路",
		"黄陂南路",
		"人民广场",
		"新闸路",
		"汉中路",
		"上海火车站",
		"中山道北路",
		"延长路",
		"上海马戏城",
		"汶水路",
		"彭浦新村",
		"共康路",
		"通河新村",
		"呼兰路",
		"共富新村",
		"宝安公路",
		"友谊西路",
		"富锦路",
	}
	var nextStationID int64
	for k, v := range defaultStation {
		station := model.Station{
			Code:          util.GetCodeUtil("S", "00000", strconv.Itoa(k+1), 5),
			Name:          v,
			LineID:        lineOne.ID,
			NextStationID: nextStationID,
			Direction:     "Up",
		}
		if err := db.DB.Where("name = ?", v).FirstOrCreate(&station).Error; err != nil {
			log.Fatal(err)
		}
		nextStationID = station.ID
	}
}

func initMetro() {

	defaultMetro := []string{"M1111", "M2222", "M3333"}
	for _, v := range defaultMetro {
		metro := model.Metro{
			Code: v,
		}
		if err := db.DB.Where("code = ?", v).FirstOrCreate(&metro).Error; err != nil {
			log.Fatal(err)
		}
	}

}

func initUser() {
	defaultUser := []string{"SuperAdmin", "dtan11"}
	for _, v := range defaultUser {
		user := model.User{
			UserName: v,
		}
		if err := db.DB.Where("user_name = ?", v).FirstOrCreate(&user).Error; err != nil {
			log.Fatal(err)
		}
	}

}
