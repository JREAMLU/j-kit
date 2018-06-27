package mysql

import (
	"testing"

	"github.com/JREAMLU/j-kit/constant"
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
			err := Load(consulAddr, false, "BGCrawler")
			So(err, ShouldBeNil)
			dbs := GetAllDB()
			So(len(dbs), ShouldBeGreaterThan, 0)
			for _, db := range dbs {
				err := db.Close()
				So(err, ShouldBeNil)
			}
		})

		Convey("load all", func() {
			err := Load(consulAddr, false)
			So(err, ShouldBeNil)
			dbs := GetAllDB()
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
	var uid int64
	Convey("sql test", t, func() {
		Convey("insert", func() {
			cron := Cron{
				Name: "jream",
			}
			id, err := insert(cron)
			So(err, ShouldBeNil)
			So(id, ShouldBeGreaterThan, constant.ZeroInt64)
			uid = id
		})

		Convey("query", func() {
			crons, err := query([]int64{uid})
			So(err, ShouldBeNil)
			So(len(crons), ShouldBeGreaterThan, 0)
		})

		Convey("raw query", func() {
			crons, err := queryRaw([]int64{uid})
			So(err, ShouldBeNil)
			So(len(crons), ShouldBeGreaterThan, 0)
		})

		Convey("update", func() {
			cron, err := update(uid)
			So(err, ShouldBeNil)
			So(cron, ShouldNotBeEmpty)
		})

		Convey("updates", func() {
			cron, err := updates(uid)
			So(err, ShouldBeNil)
			So(cron, ShouldNotBeEmpty)
		})

		Convey("deletes", func() {
			err := deletes([]int64{uid})
			So(err, ShouldBeNil)
		})

		Convey("transact", func() {
			err := transact()
			So(err, ShouldBeNil)
		})
	})
}

func load() {
	err := Load(consulAddr, true, "BGCrawler")
	if err != nil {
		panic(err)
	}
}

func db(isWrite bool) *gorm.DB {
	if isWrite {
		return GetDB("BGCrawler")
	}

	return GetReadOnlyDB("BGCrawler")
}

type Cron struct {
	ID   int64  `gorm:"column:ID;primary_key"`
	Name string `gorm:"column:Name"`
}

func (cron Cron) TableName() string {
	return cronTable
}

func insert(cron Cron) (int64, error) {
	result := db(true).Create(&cron)
	if result.Error != nil {
		return constant.ZeroInt64, result.Error
	}

	return cron.ID, nil
}

func query(ids []int64) (crons []Cron, err error) {
	result := db(false).Where("ID IN (?)", ids).Find(&crons)
	if result.Error != nil {
		return crons, result.Error
	}

	return crons, nil
}

func queryRaw(ids []int64) (crons []Cron, err error) {
	sql := `
SELECT  ID, Name
FROM cron
WHERE ID IN (?)
`
	result := db(false).Raw(sql, ids).Scan(&crons)
	if result.Error != nil {
		return crons, result.Error
	}

	return crons, nil
}

func update(id int64) (cron Cron, err error) {
	result := db(true).Model(&cron).Where("ID = ?", id).Update("Name", "LUj").Find(&cron)
	if result.Error != nil {
		return cron, result.Error
	}

	return cron, nil
}

func updates(id int64) (cron Cron, err error) {
	result := db(true).Model(&cron).Where("ID = ?", id).Updates(&Cron{Name: "JREAM"}).Find(&cron)
	if result.Error != nil {
		return cron, result.Error
	}

	return cron, nil
}

func deletes(ids []int64) error {
	result := db(true).Delete(Cron{}, "ID IN (?)", ids)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func transact() error {
	tx := db(true).Begin()
	cron := Cron{
		Name: "jream-tx",
	}
	id, err := insert(cron)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = deletes([]int64{id})
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

/*
CREATE TABLE `Crawler` (
	`ID` INT(11) NOT NULL AUTO_INCREMENT COMMENT 'INCREMENT ID',
	`Name` varchar(30) NOT NULL DEFAULT '' COMMENT 'Name',
	PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Crawler';
*/
