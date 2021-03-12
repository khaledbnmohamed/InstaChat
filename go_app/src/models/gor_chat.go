

// Package models includes the functions on the model Chat.
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

type Chat struct {
	Id int64 `json:"id,omitempty" db:"id" valid:"-"`
ApplicationId int64 `json:"application_id,omitempty" db:"application_id" valid:"-"`
Number string `json:"number,omitempty" db:"number" valid:"-"`
MessagesCount int64 `json:"messages_count,omitempty" db:"messages_count" valid:"-"`
CreatedAt time.Time `json:"created_at,omitempty" db:"created_at" valid:"-"`
UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at" valid:"-"`
Application Application `json:"application,omitempty" db:"application" valid:"-"`
Messages []Message `json:"messages,omitempty" db:"messages" valid:"-"`
}

// DataStruct for the pagination
type ChatPage struct {
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

// Current get the current page of ChatPage object for pagination.
func (_p *ChatPage) Current() ([]Chat, error) {
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
	chats, err := FindChatsWhere(whereStr, whereParams...)
	if err != nil {
		return nil, err
	}
	if len(chats) != 0 {
		_p.FirstId, _p.LastId = chats[0].Id, chats[len(chats)-1].Id
	}
	return chats, nil
}

// Previous get the previous page of ChatPage object for pagination.
func (_p *ChatPage) Previous() ([]Chat, error) {
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
	chats, err := FindChatsWhere(whereStr, whereParams...)
	if err != nil {
		return nil, err
	}
	if len(chats) != 0 {
		_p.FirstId, _p.LastId = chats[0].Id, chats[len(chats)-1].Id
	}
	_p.PageNum -= 1
	return chats, nil
}

// Next get the next page of ChatPage object for pagination.
func (_p *ChatPage) Next() ([]Chat, error) {
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
	chats, err := FindChatsWhere(whereStr, whereParams...)
	if err != nil {
		return nil, err
	}
	if len(chats) != 0 {
		_p.FirstId, _p.LastId = chats[0].Id, chats[len(chats)-1].Id
	}
	_p.PageNum += 1
	return chats, nil
}

// GetPage is a helper function for the ChatPage object to return a corresponding page due to
// the parameter passed in, i.e. one of "previous, current or next".
func (_p *ChatPage) GetPage(direction string) (ps []Chat, err error) {
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

// buildOrder is for ChatPage object to build a SQL ORDER BY clause.
func (_p *ChatPage) buildOrder() {
	tempList := []string{}
	for k, v := range _p.Order {
		tempList = append(tempList, fmt.Sprintf("%v %v", k, v))
	}
	_p.orderStr = " ORDER BY " + strings.Join(tempList, ", ")
}

// buildIdRestrict is for ChatPage object to build a SQL clause for ID restriction,
// implementing a simple keyset style pagination.
func (_p *ChatPage) buildIdRestrict(direction string) (idStr string, idParams []interface{}) {
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

// buildPageCount calculate the TotalItems/TotalPages for the ChatPage object.
func (_p *ChatPage) buildPageCount() error {
	count, err := ChatCountWhere(_p.WhereString, _p.WhereParams...)
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


// FindChat find a single chat by an ID.
func FindChat(id int64) (*Chat, error) {
	if id == 0 {
		return nil, errors.New("Invalid ID: it can't be zero")
	}
	_chat := Chat{}
	err := DB.Get(&_chat, DB.Rebind(`SELECT COALESCE(chats.number, '') AS number, COALESCE(chats.messages_count, 0) AS messages_count, chats.id, chats.application_id, chats.created_at, chats.updated_at FROM chats WHERE chats.id = ? LIMIT 1`), id)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_chat, nil
}

// FirstChat find the first one chat by ID ASC order.
func FirstChat() (*Chat, error) {
	_chat := Chat{}
	err := DB.Get(&_chat, DB.Rebind(`SELECT COALESCE(chats.number, '') AS number, COALESCE(chats.messages_count, 0) AS messages_count, chats.id, chats.application_id, chats.created_at, chats.updated_at FROM chats ORDER BY chats.id ASC LIMIT 1`))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_chat, nil
}

// FirstChats find the first N chats by ID ASC order.
func FirstChats(n uint32) ([]Chat, error) {
	_chats := []Chat{}
	sql := fmt.Sprintf("SELECT COALESCE(chats.number, '') AS number, COALESCE(chats.messages_count, 0) AS messages_count, chats.id, chats.application_id, chats.created_at, chats.updated_at FROM chats ORDER BY chats.id ASC LIMIT %v", n)
	err := DB.Select(&_chats, DB.Rebind(sql))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _chats, nil
}

// LastChat find the last one chat by ID DESC order.
func LastChat() (*Chat, error) {
	_chat := Chat{}
	err := DB.Get(&_chat, DB.Rebind(`SELECT COALESCE(chats.number, '') AS number, COALESCE(chats.messages_count, 0) AS messages_count, chats.id, chats.application_id, chats.created_at, chats.updated_at FROM chats ORDER BY chats.id DESC LIMIT 1`))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_chat, nil
}

// LastChats find the last N chats by ID DESC order.
func LastChats(n uint32) ([]Chat, error) {
	_chats := []Chat{}
	sql := fmt.Sprintf("SELECT COALESCE(chats.number, '') AS number, COALESCE(chats.messages_count, 0) AS messages_count, chats.id, chats.application_id, chats.created_at, chats.updated_at FROM chats ORDER BY chats.id DESC LIMIT %v", n)
	err := DB.Select(&_chats, DB.Rebind(sql))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _chats, nil
}

// FindChats find one or more chats by the given ID(s).
func FindChats(ids ...int64) ([]Chat, error) {
	if len(ids) == 0 {
		msg := "At least one or more ids needed"
		log.Println(msg)
		return nil, errors.New(msg)
	}
	_chats := []Chat{}
	idsHolder := strings.Repeat(",?", len(ids)-1)
	sql := DB.Rebind(fmt.Sprintf(`SELECT COALESCE(chats.number, '') AS number, COALESCE(chats.messages_count, 0) AS messages_count, chats.id, chats.application_id, chats.created_at, chats.updated_at FROM chats WHERE chats.id IN (?%s)`, idsHolder))
	idsT := []interface{}{}
	for _,id := range ids {
		idsT = append(idsT, interface{}(id))
	}
	err := DB.Select(&_chats, sql, idsT...)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _chats, nil
}

// FindChatBy find a single chat by a field name and a value.
func FindChatBy(field string, val interface{}) (*Chat, error) {
	_chat := Chat{}
	sqlFmt := `SELECT COALESCE(chats.number, '') AS number, COALESCE(chats.messages_count, 0) AS messages_count, chats.id, chats.application_id, chats.created_at, chats.updated_at FROM chats WHERE %s = ? LIMIT 1`
	sqlStr := fmt.Sprintf(sqlFmt, field)
	err := DB.Get(&_chat, DB.Rebind(sqlStr), val)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_chat, nil
}

// FindChatsBy find all chats by a field name and a value.
func FindChatsBy(field string, val interface{}) (_chats []Chat, err error) {
	sqlFmt := `SELECT COALESCE(chats.number, '') AS number, COALESCE(chats.messages_count, 0) AS messages_count, chats.id, chats.application_id, chats.created_at, chats.updated_at FROM chats WHERE %s = ?`
	sqlStr := fmt.Sprintf(sqlFmt, field)
	err = DB.Select(&_chats, DB.Rebind(sqlStr), val)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _chats, nil
}

// AllChats get all the Chat records.
func AllChats() (chats []Chat, err error) {
	err = DB.Select(&chats, "SELECT COALESCE(chats.number, '') AS number, COALESCE(chats.messages_count, 0) AS messages_count, chats.id, chats.application_id, chats.created_at, chats.updated_at FROM chats")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return chats, nil
}

// ChatCount get the count of all the Chat records.
func ChatCount() (c int64, err error) {
	err = DB.Get(&c, "SELECT count(*) FROM chats")
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return c, nil
}

// ChatCountWhere get the count of all the Chat records with a where clause.
func ChatCountWhere(where string, args ...interface{}) (c int64, err error) {
	sql := "SELECT count(*) FROM chats"
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

// ChatIncludesWhere get the Chat associated models records, currently it's not same as the corresponding "includes" function but "preload" instead in Ruby on Rails. It means that the "sql" should be restricted on Chat model.
func ChatIncludesWhere(assocs []string, sql string, args ...interface{}) (_chats []Chat, err error) {
	_chats, err = FindChatsWhere(sql, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(assocs) == 0 {
		log.Println("No associated fields ard specified")
		return _chats, err
	}
	if len(_chats) <= 0 {
		return nil, errors.New("No results available")
	}
	ids := make([]interface{}, len(_chats))
	for _, v := range _chats {
		ids = append(ids, interface{}(v.Id))
	}
	idsHolder := strings.Repeat(",?", len(ids)-1)
	for _, assoc := range assocs {
		switch assoc {
				case "messages":
							where := fmt.Sprintf("chat_id IN (?%s)", idsHolder)
						_messages, err := FindMessagesWhere(where, ids...)
						if err != nil {
							log.Printf("Error when query associated objects: %v\n", assoc)
							continue
						}
						for _, vv := range _messages {
							for i, vvv := range  _chats {
									if vv.ChatId == vvv.Id {
										vvv.Messages = append(vvv.Messages, vv)
									}
								_chats[i].Messages = vvv.Messages
						    }
					    }
		}
	}
	return _chats, nil
}

// ChatIds get all the IDs of Chat records.
func ChatIds() (ids []int64, err error) {
	err = DB.Select(&ids, "SELECT id FROM chats")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return ids, nil
}

// ChatIdsWhere get all the IDs of Chat records by where restriction.
func ChatIdsWhere(where string, args ...interface{}) ([]int64, error) {
	ids, err := ChatIntCol("id", where, args...)
	return ids, err
}

// ChatIntCol get some int64 typed column of Chat by where restriction.
func ChatIntCol(col, where string, args ...interface{}) (intColRecs []int64, err error) {
	sql := "SELECT " + col + " FROM chats"
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

// ChatStrCol get some string typed column of Chat by where restriction.
func ChatStrCol(col, where string, args ...interface{}) (strColRecs []string, err error) {
	sql := "SELECT " + col + " FROM chats"
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

// FindChatsWhere query use a partial SQL clause that usually following after WHERE
// with placeholders, eg: FindUsersWhere("first_name = ? AND age > ?", "John", 18)
// will return those records in the table "users" whose first_name is "John" and age elder than 18.
func FindChatsWhere(where string, args ...interface{}) (chats []Chat, err error) {
	sql := "SELECT COALESCE(chats.number, '') AS number, COALESCE(chats.messages_count, 0) AS messages_count, chats.id, chats.application_id, chats.created_at, chats.updated_at FROM chats"
	if len(where) > 0 {
		sql = sql + " WHERE " + where
	}
	stmt, err := DB.Preparex(DB.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = stmt.Select(&chats, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return chats, nil
}

// FindChatBySql query use a complete SQL clause
// with placeholders, eg: FindUserBySql("SELECT * FROM users WHERE first_name = ? AND age > ? ORDER BY DESC LIMIT 1", "John", 18)
// will return only One record in the table "users" whose first_name is "John" and age elder than 18.
func FindChatBySql(sql string, args ...interface{}) (*Chat, error) {
	stmt, err := DB.Preparex(DB.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	_chat := &Chat{}
	err = stmt.Get(_chat, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return _chat, nil
}

// FindChatsBySql query use a complete SQL clause
// with placeholders, eg: FindUsersBySql("SELECT * FROM users WHERE first_name = ? AND age > ?", "John", 18)
// will return those records in the table "users" whose first_name is "John" and age elder than 18.
func FindChatsBySql(sql string, args ...interface{}) (chats []Chat, err error) {
	stmt, err := DB.Preparex(DB.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = stmt.Select(&chats, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return chats, nil
}

// CreateChat use a named params to create a single Chat record.
// A named params is key-value map like map[string]interface{}{"first_name": "John", "age": 23} .
func CreateChat(am map[string]interface{}) (int64, error) {
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
	sqlFmt := `INSERT INTO chats (%s) VALUES (%s)`
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

// Create is a method for Chat to create a record.
func (_chat *Chat) Create() (int64, error) {
	ok, err := govalidator.ValidateStruct(_chat)
	if !ok {
		errMsg := "Validate Chat struct error: Unknown error"
		if err != nil {
			errMsg = "Validate Chat struct error: " + err.Error()
		}
		log.Println(errMsg)
		return 0, errors.New(errMsg)
	}
	t := time.Now()
	_chat.CreatedAt = t
	_chat.UpdatedAt = t
    sql := `INSERT INTO chats (application_id,number,messages_count,created_at,updated_at) VALUES (:application_id,:number,:messages_count,:created_at,:updated_at)`
    result, err := DB.NamedExec(sql, _chat)
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

// MessagesCreate is used for Chat to create the associated objects Messages
func (_chat *Chat) MessagesCreate(am map[string]interface{}) error {
			am["chat_id"] = _chat.Id
		_, err := CreateMessage(am)
	return err
}

// GetMessages is used for Chat to get associated objects Messages
// Say you have a Chat object named chat, when you call chat.GetMessages(),
// the object will get the associated Messages attributes evaluated in the struct.
func (_chat *Chat) GetMessages() error {
	_messages, err := ChatGetMessages(_chat.Id)
	if err == nil {
		_chat.Messages = _messages
    }
    return err
}

// ChatGetMessages a helper fuction used to get associated objects for ChatIncludesWhere().
func ChatGetMessages(id int64) ([]Message, error) {
			_messages, err := FindMessagesBy("chat_id", id)
	return _messages, err
}



// CreateApplication is a method for a Chat object to create an associated Application record.
func (_chat *Chat) CreateApplication(am map[string]interface{}) error {
	am["chat_id"] = _chat.Id
	_, err := CreateApplication(am)
	return err
}

// Destroy is method used for a Chat object to be destroyed.
func (_chat *Chat) Destroy() error {
	if _chat.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := DestroyChat(_chat.Id)
	return err
}

// DestroyChat will destroy a Chat record specified by the id parameter.
func DestroyChat(id int64) error {
	// Destroy association objects at first
	// Not care if exec properly temporarily
	destroyChatAssociations(id)
	stmt, err := DB.Preparex(DB.Rebind(`DELETE FROM chats WHERE id = ?`))
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

// DestroyChats will destroy Chat records those specified by the ids parameters.
func DestroyChats(ids ...int64) (int64, error) {
	if len(ids) == 0 {
		msg := "At least one or more ids needed"
		log.Println(msg)
		return 0, errors.New(msg)
	}
	// Destroy association objects at first
	// Not care if exec properly temporarily
	destroyChatAssociations(ids...)
	idsHolder := strings.Repeat(",?", len(ids)-1)
	sql := fmt.Sprintf(`DELETE FROM chats WHERE id IN (?%s)`, idsHolder)
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

// DestroyChatsWhere delete records by a where clause restriction.
// e.g. DestroyChatsWhere("name = ?", "John")
// And this func will not call the association dependent action
func DestroyChatsWhere(where string, args ...interface{}) (int64, error) {
	sql := `DELETE FROM chats WHERE `
	if len(where) > 0 {
		sql = sql + where
	} else {
		return 0, errors.New("No WHERE conditions provided")
	}
	ids, x_err := ChatIdsWhere(where, args...)
	if x_err != nil {
		log.Printf("Delete associated objects error: %v\n", x_err)
	} else {
		destroyChatAssociations(ids...)
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

// destroyChatAssociations is a private function used to destroy a Chat record's associated objects.
// The func not return err temporarily.
func destroyChatAssociations(ids ...int64) {
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

// Save method is used for a Chat object to update an existed record mainly.
// If no id provided a new record will be created. FIXME: A UPSERT action will be implemented further.
func (_chat *Chat) Save() error {
	ok, err := govalidator.ValidateStruct(_chat)
	if !ok {
		errMsg := "Validate Chat struct error: Unknown error"
		if err != nil {
			errMsg = "Validate Chat struct error: " + err.Error()
		}
		log.Println(errMsg)
		return errors.New(errMsg)
	}
	if _chat.Id == 0 {
		_, err = _chat.Create()
		return err
	}
	_chat.UpdatedAt = time.Now()
	sqlFmt := `UPDATE chats SET %s WHERE id = %v`
	sqlStr := fmt.Sprintf(sqlFmt, "application_id = :application_id, number = :number, messages_count = :messages_count, updated_at = :updated_at", _chat.Id)
    _, err = DB.NamedExec(sqlStr, _chat)
    return err
}

// UpdateChat is used to update a record with a id and map[string]interface{} typed key-value parameters.
func UpdateChat(id int64, am map[string]interface{}) error {
	if len(am) == 0 {
		return errors.New("Zero key in the attributes map!")
	}
	am["updated_at"] = time.Now()
	keys := allKeys(am)
	sqlFmt := `UPDATE chats SET %s WHERE id = %v`
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

// Update is a method used to update a Chat record with the map[string]interface{} typed key-value parameters.
func (_chat *Chat) Update(am map[string]interface{}) error {
	if _chat.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := UpdateChat(_chat.Id, am)
	return err
}

// UpdateAttributes method is supposed to be used to update Chat records as corresponding update_attributes in Ruby on Rails.
func (_chat *Chat) UpdateAttributes(am map[string]interface{}) error {
	if _chat.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := UpdateChat(_chat.Id, am)
	return err
}

// UpdateColumns method is supposed to be used to update Chat records as corresponding update_columns in Ruby on Rails.
func (_chat *Chat) UpdateColumns(am map[string]interface{}) error {
	if _chat.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := UpdateChat(_chat.Id, am)
	return err
}

// UpdateChatsBySql is used to update Chat records by a SQL clause
// using the '?' binding syntax.
func UpdateChatsBySql(sql string, args ...interface{}) (int64, error) {
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
