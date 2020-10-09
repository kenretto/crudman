package main

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/kenretto/crudman"
	"github.com/kenretto/crudman/driver"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/http"
	"os"
	"time"
)

type user struct {
	gorm.Model
	Nickname string `json:"nickname" gorm:"column:nickname;type:varchar(16)" validate:"required"`
	Age      uint8  `json:"age" gorm:"column:age;"`
}

type order struct {
	gorm.Model
	*driver.GormManager `gorm:"-"`
	Name                string `json:"name" gorm:"column:name;type:varchar(16)" validate:"required"`
}

// TableName table name
func (order) TableName() string {
	return "orders"
}

// List custom list data
func (order) List(r *http.Request) interface{} {
	return map[string]string{"hello": "world"}
}

// TableName table name
func (user) TableName() string {
	return "members"
}

var u user

var users = []user{
	{Nickname: "a", Age: 1},
	{Nickname: "b", Age: 2},
	{Nickname: "c", Age: 3},
	{Nickname: "d", Age: 4},
	{Nickname: "e", Age: 5},
	{Nickname: "f", Age: 6},
	{Nickname: "g", Age: 7},
	{Nickname: "h", Age: 8},
	{Nickname: "i", Age: 9},
	{Nickname: "j", Age: 10},
	{Nickname: "k", Age: 11},
	{Nickname: "l", Age: 12},
	{Nickname: "m", Age: 13},
	{Nickname: "n", Age: 14},
	{Nickname: "o", Age: 15},
	{Nickname: "p", Age: 16},
	{Nickname: "q", Age: 17},
	{Nickname: "r", Age: 18},
	{Nickname: "s", Age: 19},
	{Nickname: "t", Age: 20},
	{Nickname: "u", Age: 21},
	{Nickname: "v", Age: 22},
	{Nickname: "w", Age: 23},
	{Nickname: "x", Age: 24},
	{Nickname: "y", Age: 25},
	{Nickname: "z", Age: 26},
}

var (
	db  *gorm.DB
	err error
)

func init() {
	db, err = gorm.Open(sqlite.Open("file:pager?mode=memory&cache=shared&_fk=1"), &gorm.Config{PrepareStmt: true, Logger: logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      false,
		},
	)})
	if err != nil {
		log.Fatalln(err)
	}
	_ = db.AutoMigrate(&u)
	_ = db.AutoMigrate(&order{})
	db.Model(&u).Create(&users)
}

func main() {
	var crud = crudman.New()
	driver.SetValidator(func(obj interface{}) interface{} {
		return validator.New().Struct(obj)
	})
	crud.Register(driver.NewGorm(db, "ID"), user{}, crudman.SetRoute("/crud"))
	crud.Register(&order{GormManager: driver.NewGorm(db, "ID")}, order{})

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		var obj, err = crud.Handler(writer, request)
		var body = map[string]interface{}{
			"data": obj,
			"msg": func() string {
				if err != nil {
					return err.Error()
				}
				return ""
			}(),
		}

		data, _ := json.Marshal(body)
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write(data)
	})
	log.Fatalln(http.ListenAndServe(":3359", nil))
}
