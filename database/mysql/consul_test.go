package mysql

import (
	"fmt"
	"testing"

	"github.com/JREAMLU/j-core/constant"
	"github.com/jinzhu/gorm"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	consulAddr = "10.200.202.35:8500"
	cronTable  = "Cron"
)

func TestLoadConfig(t *testing.T) {
	Convey("load mysql test", t, func() {
		Convey("load by name", func() {
			dbs, err := Load(consulAddr, "BGCrawler")
			So(err, ShouldBeNil)
			So(len(dbs), ShouldBeGreaterThan, 0)
			for _, db := range dbs {
				err := db.Close()
				So(err, ShouldBeNil)
			}
		})

		Convey("load all", func() {
			dbs, err := Load(consulAddr)
			So(err, ShouldBeNil)
			So(len(dbs), ShouldBeGreaterThan, 0)
			for _, db := range dbs {
				err := db.Close()
				So(err, ShouldBeNil)
			}
		})
	})
}

func TestSQL(t *testing.T) {
	load()
	Convey("sql test", t, func() {
		Convey("insert", func() {
			cron := Cron{
				Name: "jream",
			}
			id, err := insert(cron)
			fmt.Println("++++++++++++: ", id, err)
		})
	})
}

var gx map[string]*gorm.DB

func load() {
	var err error
	gx, err = Load(consulAddr, "BGCrawler")
	if err != nil {
		panic(err)
	}
}

func db(isWrite bool) *gorm.DB {
	if isWrite {
		return gx["BGCrawler"]
	}

	return gx[GetReadOnly("BGCrawler")]
}

type Cron struct {
	ID   uint64 `gorm:"column:ID;primary_key"`
	Name string `gorm:"column:Name"`
}

func (cron Cron) TableName() string {
	return cronTable
}

func insert(cron Cron) (uint64, error) {
	result := db(true).Create(&cron)
	if result.Error != nil {
		return constant.ZeroUint64, result.Error
	}

	return cron.ID, nil
}

/*
CREATE TABLE `Crawler` (
	`ID` INT(11) NOT NULL AUTO_INCREMENT COMMENT 'ID，自增长',
	`Name` varchar(30) NOT NULL DEFAULT '' COMMENT '姓名',
	PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Crawler';
*/
