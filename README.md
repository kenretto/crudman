# crudman
crudman completes common operations of adding, deleting, modifying and querying, making data operation more convenient

# Example
Because only gorm is implemented now, take gorm as an example

- Build a data structure
```go
    type user struct {
    	gorm.Model
    	Nickname string `json:"nickname" gorm:"column:nickname;type:varchar(16)" validate:"required"`
    	Age      uint8  `json:"age" gorm:"column:age;"`
    }
```

- Implement crudman.Tabler

```go
// TableName table name
func (user) TableName() string {
	return "members"
}
```

- Initialize the DB object
```go
db, err = gorm.Open(sqlite.Open("file:pager?mode=memory&cache=shared&_fk=1"), &gorm.Config{PrepareStmt: true, Logger: logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      false,
		},
	)})
```
You can initialize and migrate some data.

- func main
```go
	var crud = crudman.New() // get a crudman instance

    // register your custom struct validator
	driver.SetValidator(func(obj interface{}) interface{} {
		return validator.New().Struct(obj)
	})

    // register a data struct
	crud.Register(driver.NewGorm(db, "ID"), user{}, crudman.SetRoute("/crud"))

    // register http handler
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
```