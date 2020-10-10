package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kenretto/crudman"
	"github.com/kenretto/pager"
	"github.com/kenretto/pager/driver"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"reflect"
)

var validate func(obj interface{}) interface{}

// SetValidator set up default validator and make basic configuration for github.com/go-playground/validator
func SetValidator(valid func(obj interface{}) interface{}) {
	validate = valid
}

// GormManager gorm
type GormManager struct {
	TableTyp   reflect.Type
	Route      string
	Table      crudman.Tabler
	PrimaryKey string
	db         *gorm.DB
	validator  func(obj interface{}) interface{}
}

// NewGorm gorm driver
func NewGorm(db *gorm.DB, primaryKey string) *GormManager {
	return &GormManager{
		PrimaryKey: primaryKey,
		db:         db,
	}
}

// WithValidator custom validator for a manager
func (manager *GormManager) WithValidator(validator func(obj interface{}) interface{}) *GormManager {
	manager.validator = validator
	return manager
}

// GetRoute get route
func (manager *GormManager) GetRoute() string { return manager.Route }

// SetRoute setup route
func (manager *GormManager) SetRoute(route string) { manager.Route = route }

// SetTableTyp save table struct reflect.Type
func (manager *GormManager) SetTableTyp(typ reflect.Type) { manager.TableTyp = typ }

// GetTableTyp get table struct reflect.Type
func (manager *GormManager) GetTableTyp() reflect.Type { return manager.TableTyp }

// GetTable get table struct
func (manager *GormManager) GetTable() crudman.Tabler { return manager.Table }

// SetTable save table struct
func (manager *GormManager) SetTable(table crudman.Tabler) { manager.Table = table }

// List list data
func (manager *GormManager) List(r *http.Request) interface{} {
	result := pager.New(r, driver.NewGormDriver(manager.db)).SetIndex(manager.Table.TableName()).
		SetPaginationField(manager.PrimaryKey).Find(manager.Table).Result()
	return result
}

// Post create data
func (manager *GormManager) Post(r *http.Request) (interface{}, error) {
	var (
		newInstance = reflect.New(manager.TableTyp)
	)

	switch r.Header.Get("Content-Type") {
	case "application/json":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(body, newInstance.Interface())
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("only application/json is supported for the time being")
	}

	var validate = validate
	if manager.validator != nil {
		validate = manager.validator
	}
	if rs := validate(newInstance.Interface()); rs != nil {
		return rs, errors.New("params valid error")
	}

	err := manager.db.Create(newInstance.Interface()).Error
	if err != nil {
		return nil, err
	}

	return newInstance.Interface(), nil
}

// Put update data
func (manager *GormManager) Put(r *http.Request) (interface{}, error) {
	var newInstance = reflect.New(manager.TableTyp)

	switch r.Header.Get("Content-Type") {
	case "application/json":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(body, newInstance.Interface())
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("only application/json is supported for the time being")
	}

	if newInstance.Elem().FieldByName(manager.PrimaryKey).IsZero() {
		return nil, errors.New("put must set primary key")
	}

	var validate = validate
	if manager.validator != nil {
		validate = manager.validator
	}
	if rs := validate(newInstance.Interface()); rs != nil {
		return rs, errors.New("params valid error")
	}

	err := manager.db.Table(manager.GetTable().TableName()).
		Model(newInstance.Interface()).Where("id = ?", newInstance.Elem().FieldByName(manager.PrimaryKey).Interface()).Updates(newInstance.Interface()).Error
	if err != nil {
		return nil, err
	}

	return newInstance.Interface(), nil
}

// Delete delete data
func (manager *GormManager) Delete(r *http.Request) error {
	id := r.URL.Query().Get("id")
	if id == "" {
		return errors.New("operate id can not be null")
	}

	var newInstance = reflect.New(manager.TableTyp)

	var stmt = manager.db.Model(manager.Table).Statement
	err := stmt.Parse(newInstance.Interface())
	if err != nil {
		return err
	}
	primaryKey := stmt.Schema.LookUpField("ID")
	if primaryKey.PrimaryKey {
		err := manager.db.Table(manager.GetTable().TableName()).Where(fmt.Sprintf("%s = ?", primaryKey.DBName), id).Delete(newInstance.Interface()).Error
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("primary key error")
}
