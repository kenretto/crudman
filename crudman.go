// Package crudman curd api
package crudman

import (
	"net/http"
	"reflect"
)

type (
	// Managers all managers
	Managers struct {
		container map[string]ManagerInterface
	}
	// ManagerInterface manger interface
	ManagerInterface interface {
		List(r *http.Request) interface{}
		Post(r *http.Request) (interface{}, error)
		Put(r *http.Request) (interface{}, error)
		Delete(r *http.Request) error
		GetRoute() string
		SetRoute(route string)
		SetTableTyp(typ reflect.Type)
		GetTableTyp() reflect.Type
		GetTable() Tabler
		SetTable(table Tabler)
	}

	// Setup setting
	Setup interface {
		set(managerInterface ManagerInterface)
	}
	// Route router
	Route struct {
		route string
	}
)

// Tabler any structure that uses this library needs to implement this interface,
// for example, in mysql, this is identified as the table name
type Tabler interface {
	TableName() string
}

// New get a new Managers object
func New() *Managers {
	return &Managers{
		container: make(map[string]ManagerInterface, 0),
	}
}

// SetRoute setup custom route
func SetRoute(r string) *Route                        { return &Route{r} }
func (r Route) set(managerInterface ManagerInterface) { managerInterface.SetRoute(r.route) }

// NotFound 404
func (managers *Managers) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte(`404 not found`))
}

// Forbidden 403
func (managers *Managers) Forbidden(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write([]byte(`404 not found`))
}

// Handler http handler, support GET->list data, POST->create data, PUT->update data, DELETE->delete data
func (managers *Managers) Handler(w http.ResponseWriter, r *http.Request) (obj interface{}, err error) {
	var route = r.URL.Path
	var manager, ok = managers.Get(route)
	if !ok {
		managers.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		obj = manager.List(r)
	case http.MethodPost:
		obj, err = manager.Post(r)
	case http.MethodPut:
		obj, err = manager.Put(r)
	case http.MethodDelete:
		err = manager.Delete(r)
	default:
		managers.Forbidden(w, r)
		return
	}

	return
}

// Get get exist manager
func (managers *Managers) Get(route string) (ManagerInterface, bool) {
	manager, ok := managers.container[route]
	return manager, ok
}

// Register register manager
// you can inherit GormManager or MongoManager and then override the method to implement custom operations
func (managers *Managers) Register(manager ManagerInterface, entity Tabler, setups ...Setup) *Managers {
	manager.SetRoute("/" + entity.TableName())
	manager.SetTableTyp(reflect.TypeOf(entity))
	manager.SetTable(entity)
	for _, set := range setups {
		set.set(manager)
	}
	managers.container[manager.GetRoute()] = manager
	return managers
}
