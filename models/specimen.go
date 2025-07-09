package models

import (
	"time"

	"gorm.io/gorm"
)

type Specimen struct {
	gorm.Model
	// 採集・観察について
	Latitude          float64 // 緯度
	Longitude         float64 // 経度
	Date              *time.Time // 年月日
	Time              int    // 時間
	CollectionMethod  string // 採集方法
	Collector         string // 採集・観察者
	IndividualCount   int    // 個体数
	Weather           string // 天気
	Temperature       float64 // 気温
	Humidity          float64 // 湿度
	Precipitation     float64 // 降水量
	Environment       string // 環境
	CollectionRemarks string // 採集・観察備考

	// 標本について
	SpecimenCreator   string // 標本作成者
	SpecimenType      string // 標本種類
	SpecimenLocation  string // 標本所在地
	SpecimenBoxID     string // 標本箱、サンプル番号
	SpecimenCreationDate *time.Time // 標本作成日
	SpecimenRemarks   string // 標本備考

	// 共通
	Sex              string `gorm:"not null"` // 雌雄
	SpeciesName      string `gorm:"not null"` // 種名(学名)
	JapaneseName     string // 和名
	Age              string // 年齢・成体/幼体
	Identifier       string // 同定者
	ProjectName      string // プロジェクト名
	DataInputDate    *time.Time // データ入力日
	
	CommonRemarks    string // 共通備考
	Photos           []Photo `gorm:"foreignKey:SpecimenID"`
}
