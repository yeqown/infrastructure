package gormic

import (
	"fmt"
	"testing"

	"github.com/yeqown/infrastructure/types"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func (m UserModel) TableName() string {
	return "users"
}

type UserModel struct {
	gorm.Model
	Name     string       `gorm:"column:name"`
	SchoolID uint         `gorm:"column:school_id"`
	GeoID    uint         `gorm:"column:geo_id"`
	School   *SchoolModel `gorm:"column:_;foreignkey:school_id;"`
	Geo      *GeoModel    `gorm:"column:_;foreignkey:geo_id;"`
}

func (m SchoolModel) TableName() string {
	return "schools"
}

type SchoolModel struct {
	gorm.Model
	Name string `gorm:"column:name"`
}

func (m GeoModel) TableName() string {
	return "geos"
}

type GeoModel struct {
	gorm.Model
	Name string `gorm:"column:name"`
}

func prepareTestdata() (*gorm.DB, error) {
	db, err := ConnectSqlite3(&types.SQLite3Config{
		Name: "sqlite3.db",
	})
	if err != nil {
		return nil, err
	}

	var (
		um = new(UserModel)
		sm = new(SchoolModel)
		gm = new(GeoModel)
	)

	db.DropTableIfExists(um)
	db.DropTableIfExists(sm)
	db.DropTableIfExists(gm)

	if err := db.
		CreateTable(um).
		CreateTable(sm).
		CreateTable(gm).Error; err != nil {
		// pass
	}

	// create 3 school and 3 geo
	for i := 0; i < 3; i++ {
		sm = &SchoolModel{
			Name: fmt.Sprintf("school_%d", i+1),
		}

		gm = &GeoModel{
			Name: fmt.Sprintf("geo_%d", i+1),
		}

		if err := db.Create(sm).Error; err != nil {
			return nil, err
		}

		if err := db.Create(gm).Error; err != nil {
			return nil, err
		}
	}

	randName := func(i int) string {
		switch i % 3 {
		case 0:
			return fmt.Sprintf("test_%d", i)
		case 1:
			return fmt.Sprintf("mock_%d", i)
		case 2:
			return fmt.Sprintf("simu_%d", i)
		}

		return fmt.Sprintf("null_%d", i)
	}

	for i := 0; i < 21; i++ {
		um := &UserModel{
			Name:     randName(i),
			SchoolID: uint(i%3 + 1),
			GeoID:    uint(i%3 + 1),
		}

		if err := db.Create(um).Error; err != nil {
			return nil, err
		}

	}

	return db, nil
}

func Test_PagingUsers(t *testing.T) {
	db, err := prepareTestdata()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	p := &PagingParam{
		DB:     db,
		Limit:  10,
		Offset: 0,
		// Select:   "",
		// OrderBy:  "",
		// GroupBy:  "",
		// Joins:    "",
		Wheres: []Where{
			{
				Key:   "name LIKE ?",
				Value: "test%",
			},
		},
		ORs: []Where{
			{
				Key:   "name LIKE ?",
				Value: "mock%",
			},
		},
		// Not: "",
		Preloads: []string{
			"School", "Geo",
		},
		ShowSQL: true,
	}
	out := make([]UserModel, 0)
	paginator := Pagging(p, &out)
	if paginator.Error != nil {
		t.Error(paginator.Error)
		t.FailNow()
	}

	t.Log(paginator.Total) // want 14
	for _, v := range *paginator.Records.(*[]UserModel) {
		t.Logf("name=%s, school_id=%d, geo_id=%d, school=(%d: %s), geo=(%d: %s)",
			v.Name, v.SchoolID, v.GeoID, v.School.ID, v.School.Name, v.Geo.ID, v.Geo.Name)
	}
}
