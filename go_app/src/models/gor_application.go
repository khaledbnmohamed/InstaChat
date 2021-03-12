

// Package models includes the functions on the model Application.
package models

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
)

// set flags to output more detailed log
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type Application struct {
	Id int64 `json:"id,omitempty" db:"id" valid:"-"`
Name string `json:"name,omitempty" db:"name" valid:"required"`
Number string `json:"number,omitempty" db:"number" valid:"-"`
ChatsCount int64 `json:"chats_count,omitempty" db:"chats_count" valid:"-"`
CreatedAt time.Time `json:"created_at,omitempty" db:"created_at" valid:"-"`
UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at" valid:"-"`
Chats []Chat `json:"chats,omitempty" db:"chats" valid:"-"`
}

// DataStruct for the pagination
type ApplicationPage struct {
	WhereString string
	WhereParams []interface{}
	Order       map[string]string
	FirstId     int64
	LastId      int64
	PageNum     int
	PerPage     int
	TotalPages  int
	TotalItems  int64
	orderStr    string
}

// Current get the current page of ApplicationPage object for pagination.
func (_p *ApplicationPage) Current() ([]Application, error) {
	if _, exist := _p.Order["id"]; !exist {
		return nil, errors.New("No id order specified in Order map")
	}
	err := _p.buildPageCount()
	if err != nil {
		return nil, fmt.Errorf("Calculate page count error: %v", err)
	}
	if _p.orderStr == "" {
		_p.buildOrder()
	}
	idStr, idParams := _p.buildIdRestrict("current")
	whereStr := fmt.Sprintf("%s %s %s LIMIT %v", _p.WhereString, idStr, _p.orderStr, _p.PerPage)
	whereParams := []interface{}{}
	whereParams = append(append(whereParams, _p.WhereParams...), idParams...)
	applications, err := FindApplicationsWhere(whereStr, whereParams...)
	if err != nil {
		return nil, err
	}
	if len(applications) != 0 {
		_p.FirstId, _p.LastId = applications[0].Id, applications[len(applications)-1].Id
	}
	return applications, nil
}

// Previous get the previous page of ApplicationPage object for pagination.
func (_p *ApplicationPage) Previous() ([]Application, error) {
	if _p.PageNum == 0 {
		return nil, errors.New("This's the first page, no previous page yet")
	}
	if _, exist := _p.Order["id"]; !exist {
		return nil, errors.New("No id order specified in Order map")
	}
	err := _p.buildPageCount()
	if err != nil {
		return nil, fmt.Errorf("Calculate page count error: %v", err)
	}
	if _p.orderStr == "" {
		_p.buildOrder()
	}
	idStr, idParams := _p.buildIdRestrict("previous")
	whereStr := fmt.Sprintf("%s %s %s LIMIT %v", _p.WhereString, idStr, _p.orderStr, _p.PerPage)
	whereParams := []interface{}{}
	whereParams = append(append(whereParams, _p.WhereParams...), idParams...)
	applications, err := FindApplicationsWhere(whereStr, whereParams...)
	if err != nil {
		return nil, err
	}
	if len(applications) != 0 {
		_p.FirstId, _p.LastId = applications[0].Id, applications[len(applications)-1].Id
	}
	_p.PageNum -= 1
	return applications, nil
}

// Next get the next page of ApplicationPage object for pagination.
func (_p *ApplicationPage) Next() ([]Application, error) {
	if _p.PageNum == _p.TotalPages-1 {
		return nil, errors.New("This's the last page, no next page yet")
	}
	if _, exist := _p.Order["id"]; !exist {
		return nil, errors.New("No id order specified in Order map")
	}
	err := _p.buildPageCount()
	if err != nil {
		return nil, fmt.Errorf("Calculate page count error: %v", err)
	}
	if _p.orderStr == "" {
		_p.buildOrder()
	}
	idStr, idParams := _p.buildIdRestrict("next")
	whereStr := fmt.Sprintf("%s %s %s LIMIT %v", _p.WhereString, idStr, _p.orderStr, _p.PerPage)
	whereParams := []interface{}{}
	whereParams = append(append(whereParams, _p.WhereParams...), idParams...)
	applications, err := FindApplicationsWhere(whereStr, whereParams...)
	if err != nil {
		return nil, err
	}
	if len(applications) != 0 {
		_p.FirstId, _p.LastId = applications[0].Id, applications[len(applications)-1].Id
	}
	_p.PageNum += 1
	return applications, nil
}

// GetPage is a helper function for the ApplicationPage object to return a corresponding page due to
// the parameter passed in, i.e. one of "previous, current or next".
func (_p *ApplicationPage) GetPage(direction string) (ps []Application, err error) {
	switch direction {
	case "previous":
		ps, _ = _p.Previous()
	case "next":
		ps, _ = _p.Next()
	case "current":
		ps, _ = _p.Current()
	default:
		return nil, errors.New("Error: wrong dircetion! None of previous, current or next!")
	}
	return
}

// buildOrder is for ApplicationPage object to build a SQL ORDER BY clause.
func (_p *ApplicationPage) buildOrder() {
	tempList := []string{}
	for k, v := range _p.Order {
		tempList = append(tempList, fmt.Sprintf("%v %v", k, v))
	}
	_p.orderStr = " ORDER BY " + strings.Join(tempList, ", ")
}

// buildIdRestrict is for ApplicationPage object to build a SQL clause for ID restriction,
// implementing a simple keyset style pagination.
func (_p *ApplicationPage) buildIdRestrict(direction string) (idStr string, idParams []interface{}) {
	switch direction {
	case "previous":
		if strings.ToLower(_p.Order["id"]) == "desc" {
			idStr += "id > ? "
			idParams = append(idParams, _p.FirstId)
		} else {
			idStr += "id < ? "
			idParams = append(idParams, _p.FirstId)
		}
	case "current":
		// trick to make Where function work
		if _p.PageNum == 0 && _p.FirstId == 0 && _p.LastId == 0 {
			idStr += "id > ? "
			idParams = append(idParams, 0)
		} else {
			if strings.ToLower(_p.Order["id"]) == "desc" {
				idStr += "id <= ? AND id >= ? "
				idParams = append(idParams, _p.FirstId, _p.LastId)
			} else {
				idStr += "id >= ? AND id <= ? "
				idParams = append(idParams, _p.FirstId, _p.LastId)
			}
		}
	case "next":
		if strings.ToLower(_p.Order["id"]) == "desc" {
			idStr += "id < ? "
			idParams = append(idParams, _p.LastId)
		} else {
			idStr += "id > ? "
			idParams = append(idParams, _p.LastId)
		}
	}
	if _p.WhereString != "" {
		idStr = " AND " + idStr
	}
	return
}

// buildPageCount calculate the TotalItems/TotalPages for the ApplicationPage object.
func (_p *ApplicationPage) buildPageCount() error {
	count, err := ApplicationCountWhere(_p.WhereString, _p.WhereParams...)
	if err != nil {
		return err
	}
	_p.TotalItems = count
	if _p.PerPage == 0 {
		_p.PerPage = 10
	}
	_p.TotalPages = int(math.Ceil(float64(_p.TotalItems) / float64(_p.PerPage)))
	return nil
}


// FindApplication find a single application by an ID.
func FindApplication(id int64) (*Application, error) {
	if id == 0 {
		return nil, errors.New("Invalid ID: it can't be zero")
	}
	_application := Application{}
	err := DB.Get(&_application, DB.Rebind(`SELECT COALESCE(applications.chats_count, 0) AS chats_count, applications.id, applications.name, applications.number, applications.created_at, applications.updated_at FROM applications WHERE applications.id = ? LIMIT 1`), id)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_application, nil
}

// FirstApplication find the first one application by ID ASC order.
func FirstApplication() (*Application, error) {
	_application := Application{}
	err := DB.Get(&_application, DB.Rebind(`SELECT COALESCE(applications.chats_count, 0) AS chats_count, applications.id, applications.name, applications.number, applications.created_at, applications.updated_at FROM applications ORDER BY applications.id ASC LIMIT 1`))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_application, nil
}

// FirstApplications find the first N applications by ID ASC order.
func FirstApplications(n uint32) ([]Application, error) {
	_applications := []Application{}
	sql := fmt.Sprintf("SELECT COALESCE(applications.chats_count, 0) AS chats_count, applications.id, applications.name, applications.number, applications.created_at, applications.updated_at FROM applications ORDER BY applications.id ASC LIMIT %v", n)
	err := DB.Select(&_applications, DB.Rebind(sql))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _applications, nil
}

// LastApplication find the last one application by ID DESC order.
func LastApplication() (*Application, error) {
	_application := Application{}
	err := DB.Get(&_application, DB.Rebind(`SELECT COALESCE(applications.chats_count, 0) AS chats_count, applications.id, applications.name, applications.number, applications.created_at, applications.updated_at FROM applications ORDER BY applications.id DESC LIMIT 1`))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_application, nil
}

// LastApplications find the last N applications by ID DESC order.
func LastApplications(n uint32) ([]Application, error) {
	_applications := []Application{}
	sql := fmt.Sprintf("SELECT COALESCE(applications.chats_count, 0) AS chats_count, applications.id, applications.name, applications.number, applications.created_at, applications.updated_at FROM applications ORDER BY applications.id DESC LIMIT %v", n)
	err := DB.Select(&_applications, DB.Rebind(sql))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _applications, nil
}

// FindApplications find one or more applications by the given ID(s).
func FindApplications(ids ...int64) ([]Application, error) {
	if len(ids) == 0 {
		msg := "At least one or more ids needed"
		log.Println(msg)
		return nil, errors.New(msg)
	}
	_applications := []Application{}
	idsHolder := strings.Repeat(",?", len(ids)-1)
	sql := DB.Rebind(fmt.Sprintf(`SELECT COALESCE(applications.chats_count, 0) AS chats_count, applications.id, applications.name, applications.number, applications.created_at, applications.updated_at FROM applications WHERE applications.id IN (?%s)`, idsHolder))
	idsT := []interface{}{}
	for _,id := range ids {
		idsT = append(idsT, interface{}(id))
	}
	err := DB.Select(&_applications, sql, idsT...)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _applications, nil
}

// FindApplicationBy find a single application by a field name and a value.
func FindApplicationBy(field string, val interface{}) (*Application, error) {
	_application := Application{}
	sqlFmt := `SELECT COALESCE(applications.chats_count, 0) AS chats_count, applications.id, applications.name, applications.number, applications.created_at, applications.updated_at FROM applications WHERE %s = ? LIMIT 1`
	sqlStr := fmt.Sprintf(sqlFmt, field)
	err := DB.Get(&_application, DB.Rebind(sqlStr), val)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_application, nil
}

// FindApplicationsBy find all applications by a field name and a value.
func FindApplicationsBy(field string, val interface{}) (_applications []Application, err error) {
	sqlFmt := `SELECT COALESCE(applications.chats_count, 0) AS chats_count, applications.id, applications.name, applications.number, applications.created_at, applications.updated_at FROM applications WHERE %s = ?`
	sqlStr := fmt.Sprintf(sqlFmt, field)
	err = DB.Select(&_applications, DB.Rebind(sqlStr), val)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _applications, nil
}

// AllApplications get all the Application records.
func AllApplications() (applications []Application, err error) {
	err = DB.Select(&applications, "SELECT COALESCE(applications.chats_count, 0) AS chats_count, applications.id, applications.name, applications.number, applications.created_at, applications.updated_at FROM applications")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return applications, nil
}

// ApplicationCount get the count of all the Application records.
func ApplicationCount() (c int64, err error) {
	err = DB.Get(&c, "SELECT count(*) FROM applications")
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return c, nil
}

// ApplicationCountWhere get the count of all the Application records with a where clause.
func ApplicationCountWhere(where string, args ...interface{}) (c int64, err error) {
	sql := "SELECT count(*) FROM applications"
	if len(where) > 0 {
		sql = sql + " WHERE " + where
	}
	stmt, err := DB.Preparex(DB.Rebind(sql))
	if err != nil {
		log.Println(err)
		return 0, err
	}
	err = stmt.Get(&c, args...)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return c, nil
}

// ApplicationIncludesWhere get the Application associated models records, currently it's not same as the corresponding "includes" function but "preload" instead in Ruby on Rails. It means that the "sql" should be restricted on Application model.
func ApplicationIncludesWhere(assocs []string, sql string, args ...interface{}) (_applications []Application, err error) {
	_applications, err = FindApplicationsWhere(sql, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(assocs) == 0 {
		log.Println("No associated fields ard specified")
		return _applications, err
	}
	if len(_applications) <= 0 {
		return nil, errors.New("No results available")
	}
	ids := make([]interface{}, len(_applications))
	for _, v := range _applications {
		ids = append(ids, interface{}(v.Id))
	}
	idsHolder := strings.Repeat(",?", len(ids)-1)
	for _, assoc := range assocs {
		switch assoc {
				case "chats":
							where := fmt.Sprintf("application_id IN (?%s)", idsHolder)
						_chats, err := FindChatsWhere(where, ids...)
						if err != nil {
							log.Printf("Error when query associated objects: %v\n", assoc)
							continue
						}
						for _, vv := range _chats {
							for i, vvv := range  _applications {
									if vv.ApplicationId == vvv.Id {
										vvv.Chats = append(vvv.Chats, vv)
									}
								_applications[i].Chats = vvv.Chats
						    }
					    }
		}
	}
	return _applications, nil
}

// ApplicationIds get all the IDs of Application records.
func ApplicationIds() (ids []int64, err error) {
	err = DB.Select(&ids, "SELECT id FROM applications")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return ids, nil
}

// ApplicationIdsWhere get all the IDs of Application records by where restriction.
func ApplicationIdsWhere(where string, args ...interface{}) ([]int64, error) {
	ids, err := ApplicationIntCol("id", where, args...)
	return ids, err
}

// ApplicationIntCol get some int64 typed column of Application by where restriction.
func ApplicationIntCol(col, where string, args ...interface{}) (intColRecs []int64, err error) {
	sql := "SELECT " + col + " FROM applications"
	if len(where) > 0 {
		sql = sql + " WHERE " + where
	}
	stmt, err := DB.Preparex(DB.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = stmt.Select(&intColRecs, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return intColRecs, nil
}

// ApplicationStrCol get some string typed column of Application by where restriction.
func ApplicationStrCol(col, where string, args ...interface{}) (strColRecs []string, err error) {
	sql := "SELECT " + col + " FROM applications"
	if len(where) > 0 {
		sql = sql + " WHERE " + where
	}
	stmt, err := DB.Preparex(DB.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = stmt.Select(&strColRecs, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return strColRecs, nil
}

// FindApplicationsWhere query use a partial SQL clause that usually following after WHERE
// with placeholders, eg: FindUsersWhere("first_name = ? AND age > ?", "John", 18)
// will return those records in the table "users" whose first_name is "John" and age elder than 18.
func FindApplicationsWhere(where string, args ...interface{}) (applications []Application, err error) {
	sql := "SELECT COALESCE(applications.chats_count, 0) AS chats_count, applications.id, applications.name, applications.number, applications.created_at, applications.updated_at FROM applications"
	if len(where) > 0 {
		sql = sql + " WHERE " + where
	}
	stmt, err := DB.Preparex(DB.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = stmt.Select(&applications, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return applications, nil
}

// FindApplicationBySql query use a complete SQL clause
// with placeholders, eg: FindUserBySql("SELECT * FROM users WHERE first_name = ? AND age > ? ORDER BY DESC LIMIT 1", "John", 18)
// will return only One record in the table "users" whose first_name is "John" and age elder than 18.
func FindApplicationBySql(sql string, args ...interface{}) (*Application, error) {
	stmt, err := DB.Preparex(DB.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	_application := &Application{}
	err = stmt.Get(_application, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return _application, nil
}

// FindApplicationsBySql query use a complete SQL clause
// with placeholders, eg: FindUsersBySql("SELECT * FROM users WHERE first_name = ? AND age > ?", "John", 18)
// will return those records in the table "users" whose first_name is "John" and age elder than 18.
func FindApplicationsBySql(sql string, args ...interface{}) (applications []Application, err error) {
	stmt, err := DB.Preparex(DB.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = stmt.Select(&applications, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return applications, nil
}

// CreateApplication use a named params to create a single Application record.
// A named params is key-value map like map[string]interface{}{"first_name": "John", "age": 23} .
func CreateApplication(am map[string]interface{}) (int64, error) {
	if len(am) == 0 {
		return 0, fmt.Errorf("Zero key in the attributes map!")
	}
	t := time.Now()
	for _, v := range []string{"created_at", "updated_at"} {
		if am[v] == nil {
			am[v] = t
		}
	}
	keys := allKeys(am)
	sqlFmt := `INSERT INTO applications (%s) VALUES (%s)`
	sql := fmt.Sprintf(sqlFmt, strings.Join(keys, ","), ":"+strings.Join(keys, ",:"))
	result, err := DB.NamedExec(sql, am)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return lastId, nil
}

// Create is a method for Application to create a record.
func (_application *Application) Create() (int64, error) {
	ok, err := govalidator.ValidateStruct(_application)
	if !ok {
		errMsg := "Validate Application struct error: Unknown error"
		if err != nil {
			errMsg = "Validate Application struct error: " + err.Error()
		}
		log.Println(errMsg)
		return 0, errors.New(errMsg)
	}
	t := time.Now()
	_application.CreatedAt = t
	_application.UpdatedAt = t
    sql := `INSERT INTO applications (name,number,chats_count,created_at,updated_at) VALUES (:name,:number,:chats_count,:created_at,:updated_at)`
    result, err := DB.NamedExec(sql, _application)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return lastId, nil
}

// ChatsCreate is used for Application to create the associated objects Chats
func (_application *Application) ChatsCreate(am map[string]interface{}) error {
			am["application_id"] = _application.Id
		_, err := CreateChat(am)
	return err
}

// GetChats is used for Application to get associated objects Chats
// Say you have a Application object named application, when you call application.GetChats(),
// the object will get the associated Chats attributes evaluated in the struct.
func (_application *Application) GetChats() error {
	_chats, err := ApplicationGetChats(_application.Id)
	if err == nil {
		_application.Chats = _chats
    }
    return err
}

// ApplicationGetChats a helper fuction used to get associated objects for ApplicationIncludesWhere().
func ApplicationGetChats(id int64) ([]Chat, error) {
			_chats, err := FindChatsBy("application_id", id)
	return _chats, err
}




// Destroy is method used for a Application object to be destroyed.
func (_application *Application) Destroy() error {
	if _application.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := DestroyApplication(_application.Id)
	return err
}

// DestroyApplication will destroy a Application record specified by the id parameter.
func DestroyApplication(id int64) error {
	// Destroy association objects at first
	// Not care if exec properly temporarily
	destroyApplicationAssociations(id)
	stmt, err := DB.Preparex(DB.Rebind(`DELETE FROM applications WHERE id = ?`))
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

// DestroyApplications will destroy Application records those specified by the ids parameters.
func DestroyApplications(ids ...int64) (int64, error) {
	if len(ids) == 0 {
		msg := "At least one or more ids needed"
		log.Println(msg)
		return 0, errors.New(msg)
	}
	// Destroy association objects at first
	// Not care if exec properly temporarily
	destroyApplicationAssociations(ids...)
	idsHolder := strings.Repeat(",?", len(ids)-1)
	sql := fmt.Sprintf(`DELETE FROM applications WHERE id IN (?%s)`, idsHolder)
	idsT := []interface{}{}
	for _,id := range ids {
		idsT = append(idsT, interface{}(id))
	}
	stmt, err := DB.Preparex(DB.Rebind(sql))
	result, err := stmt.Exec(idsT...)
	if err != nil {
		return 0, err
	}
	cnt, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// DestroyApplicationsWhere delete records by a where clause restriction.
// e.g. DestroyApplicationsWhere("name = ?", "John")
// And this func will not call the association dependent action
func DestroyApplicationsWhere(where string, args ...interface{}) (int64, error) {
	sql := `DELETE FROM applications WHERE `
	if len(where) > 0 {
		sql = sql + where
	} else {
		return 0, errors.New("No WHERE conditions provided")
	}
	ids, x_err := ApplicationIdsWhere(where, args...)
	if x_err != nil {
		log.Printf("Delete associated objects error: %v\n", x_err)
	} else {
		destroyApplicationAssociations(ids...)
	}
	stmt, err := DB.Preparex(DB.Rebind(sql))
	result, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	cnt, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// destroyApplicationAssociations is a private function used to destroy a Application record's associated objects.
// The func not return err temporarily.
func destroyApplicationAssociations(ids ...int64) {
	idsHolder := ""
	if len(ids) > 1 {
		idsHolder = strings.Repeat(",?", len(ids)-1)
	}
	idsT := []interface{}{}
	for _, id := range ids {
		idsT = append(idsT, interface{}(id))
	}
	var err error
	// make sure no declared-and-not-used exception
	_, _, _ = idsHolder, idsT, err
}

// Save method is used for a Application object to update an existed record mainly.
// If no id provided a new record will be created. FIXME: A UPSERT action will be implemented further.
func (_application *Application) Save() error {
	ok, err := govalidator.ValidateStruct(_application)
	if !ok {
		errMsg := "Validate Application struct error: Unknown error"
		if err != nil {
			errMsg = "Validate Application struct error: " + err.Error()
		}
		log.Println(errMsg)
		return errors.New(errMsg)
	}
	if _application.Id == 0 {
		_, err = _application.Create()
		return err
	}
	_application.UpdatedAt = time.Now()
	sqlFmt := `UPDATE applications SET %s WHERE id = %v`
	sqlStr := fmt.Sprintf(sqlFmt, "name = :name, number = :number, chats_count = :chats_count, updated_at = :updated_at", _application.Id)
    _, err = DB.NamedExec(sqlStr, _application)
    return err
}

// UpdateApplication is used to update a record with a id and map[string]interface{} typed key-value parameters.
func UpdateApplication(id int64, am map[string]interface{}) error {
	if len(am) == 0 {
		return errors.New("Zero key in the attributes map!")
	}
	am["updated_at"] = time.Now()
	keys := allKeys(am)
	sqlFmt := `UPDATE applications SET %s WHERE id = %v`
	setKeysArr := []string{}
	for _,v := range keys {
		s := fmt.Sprintf(" %s = :%s", v, v)
		setKeysArr = append(setKeysArr, s)
	}
	sqlStr := fmt.Sprintf(sqlFmt, strings.Join(setKeysArr, ", "), id)
	_, err := DB.NamedExec(sqlStr, am)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Update is a method used to update a Application record with the map[string]interface{} typed key-value parameters.
func (_application *Application) Update(am map[string]interface{}) error {
	if _application.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := UpdateApplication(_application.Id, am)
	return err
}

// UpdateAttributes method is supposed to be used to update Application records as corresponding update_attributes in Ruby on Rails.
func (_application *Application) UpdateAttributes(am map[string]interface{}) error {
	if _application.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := UpdateApplication(_application.Id, am)
	return err
}

// UpdateColumns method is supposed to be used to update Application records as corresponding update_columns in Ruby on Rails.
func (_application *Application) UpdateColumns(am map[string]interface{}) error {
	if _application.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := UpdateApplication(_application.Id, am)
	return err
}

// UpdateApplicationsBySql is used to update Application records by a SQL clause
// using the '?' binding syntax.
func UpdateApplicationsBySql(sql string, args ...interface{}) (int64, error) {
	if sql == "" {
		return 0, errors.New("A blank SQL clause")
	}
	stmt, err := DB.Preparex(DB.Rebind(sql))
	result, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	cnt, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return cnt, nil
}
