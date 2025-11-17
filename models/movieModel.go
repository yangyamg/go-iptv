package models

type IptvMovie struct {
	ID    int64  `gorm:"column:id;primaryKey;autoIncrement" json:"-"`
	Name  string `gorm:"column:name" json:"name"`
	Api   string `gorm:"column:api" json:"api"`
	State int64  `gorm:"default:1;column:state" json:"-"`
}

func (IptvMovie) TableName() string {
	return "iptv_movie"
}
