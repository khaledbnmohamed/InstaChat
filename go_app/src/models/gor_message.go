

// Package models includes the functions on the model Message.
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

type Message struct {
	Id int64 `json:"id,omitempty" db:"id" valid:"-"`
Text string `json:"text,omitempty" db:"text" valid:"required"`
CreatedAt time.Time `json:"created_at,omitempty" db:"created_at" valid:"-"`
UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at" valid:"-"`
ChatId int64 `json:"chat_id,omitempty" db:"chat_id" valid:"-"`
Number string `json:"number,omitempty" db:"number" valid:"-"`
Chat Chat `json:"chat,omitempty" db:"chat" valid:"-"`
}

// DataStruct for the pagination
type MessagePage struct {
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

// Current get the current page of MessagePage object for pagination.
func (_p *MessagePage) Current() ([]Message, error) {
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
	messages, err := FindMessagesWhere(whereStr, whereParams...)
	if err != nil {
		return nil, err
	}
	if len(messages) != 0 {
		_p.FirstId, _p.LastId = messages[0].Id, messages[len(messages)-1].Id
	}
	return messages, nil
}

// Previous get the previous page of MessagePage object for pagination.
func (_p *MessagePage) Previous() ([]Message, error) {
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
	messages, err := FindMessagesWhere(whereStr, whereParams...)
	if err != nil {
		return nil, err
	}
	if len(messages) != 0 {
		_p.FirstId, _p.LastId = messages[0].Id, messages[len(messages)-1].Id
	}
	_p.PageNum -= 1
	return messages, nil
}

// Next get the next page of MessagePage object for pagination.
func (_p *MessagePage) Next() ([]Message, error) {
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
	messages, err := FindMessagesWhere(whereStr, whereParams...)
	if err != nil {
		return nil, err
	}
	if len(messages) != 0 {
		_p.FirstId, _p.LastId = messages[0].Id, messages[len(messages)-1].Id
	}
	_p.PageNum += 1
	return messages, nil
}

// GetPage is a helper function for the MessagePage object to return a corresponding page due to
// the parameter passed in, i.e. one of "previous, current or next".
func (_p *MessagePage) GetPage(direction string) (ps []Message, err error) {
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

// buildOrder is for MessagePage object to build a SQL ORDER BY clause.
func (_p *MessagePage) buildOrder() {
	tempList := []string{}
	for k, v := range _p.Order {
		tempList = append(tempList, fmt.Sprintf("%v %v", k, v))
	}
	_p.orderStr = " ORDER BY " + strings.Join(tempList, ", ")
}

// buildIdRestrict is for MessagePage object to build a SQL clause for ID restriction,
// implementing a simple keyset style pagination.
func (_p *MessagePage) buildIdRestrict(direction string) (idStr string, idParams []interface{}) {
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

// buildPageCount calculate the TotalItems/TotalPages for the MessagePage object.
func (_p *MessagePage) buildPageCount() error {
	count, err := MessageCountWhere(_p.WhereString, _p.WhereParams...)
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


// FindMessage find a single message by an ID.
func FindMessage(id int64) (*Message, error) {
	if id == 0 {
		return nil, errors.New("Invalid ID: it can't be zero")
	}
	_message := Message{}
	err := DB.Get(&_message, DB.Rebind(`SELECT COALESCE(messages.chat_id, 0) AS chat_id, COALESCE(messages.number, '') AS number, messages.id, messages.text, messages.created_at, messages.updated_at FROM messages WHERE messages.id = ? LIMIT 1`), id)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_message, nil
}

// FirstMessage find the first one message by ID ASC order.
func FirstMessage() (*Message, error) {
	_message := Message{}
	err := DB.Get(&_message, DB.Rebind(`SELECT COALESCE(messages.chat_id, 0) AS chat_id, COALESCE(messages.number, '') AS number, messages.id, messages.text, messages.created_at, messages.updated_at FROM messages ORDER BY messages.id ASC LIMIT 1`))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_message, nil
}

// FirstMessages find the first N messages by ID ASC order.
func FirstMessages(n uint32) ([]Message, error) {
	_messages := []Message{}
	sql := fmt.Sprintf("SELECT COALESCE(messages.chat_id, 0) AS chat_id, COALESCE(messages.number, '') AS number, messages.id, messages.text, messages.created_at, messages.updated_at FROM messages ORDER BY messages.id ASC LIMIT %v", n)
	err := DB.Select(&_messages, DB.Rebind(sql))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _messages, nil
}

// LastMessage find the last one message by ID DESC order.
func LastMessage() (*Message, error) {
	_message := Message{}
	err := DB.Get(&_message, DB.Rebind(`SELECT COALESCE(messages.chat_id, 0) AS chat_id, COALESCE(messages.number, '') AS number, messages.id, messages.text, messages.created_at, messages.updated_at FROM messages ORDER BY messages.id DESC LIMIT 1`))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_message, nil
}

// LastMessages find the last N messages by ID DESC order.
func LastMessages(n uint32) ([]Message, error) {
	_messages := []Message{}
	sql := fmt.Sprintf("SELECT COALESCE(messages.chat_id, 0) AS chat_id, COALESCE(messages.number, '') AS number, messages.id, messages.text, messages.created_at, messages.updated_at FROM messages ORDER BY messages.id DESC LIMIT %v", n)
	err := DB.Select(&_messages, DB.Rebind(sql))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _messages, nil
}

// FindMessages find one or more messages by the given ID(s).
func FindMessages(ids ...int64) ([]Message, error) {
	if len(ids) == 0 {
		msg := "At least one or more ids needed"
		log.Println(msg)
		return nil, errors.New(msg)
	}
	_messages := []Message{}
	idsHolder := strings.Repeat(",?", len(ids)-1)
	sql := DB.Rebind(fmt.Sprintf(`SELECT COALESCE(messages.chat_id, 0) AS chat_id, COALESCE(messages.number, '') AS number, messages.id, messages.text, messages.created_at, messages.updated_at FROM messages WHERE messages.id IN (?%s)`, idsHolder))
	idsT := []interface{}{}
	for _,id := range ids {
		idsT = append(idsT, interface{}(id))
	}
	err := DB.Select(&_messages, sql, idsT...)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _messages, nil
}

// FindMessageBy find a single message by a field name and a value.
func FindMessageBy(field string, val interface{}) (*Message, error) {
	_message := Message{}
	sqlFmt := `SELECT COALESCE(messages.chat_id, 0) AS chat_id, COALESCE(messages.number, '') AS number, messages.id, messages.text, messages.created_at, messages.updated_at FROM messages WHERE %s = ? LIMIT 1`
	sqlStr := fmt.Sprintf(sqlFmt, field)
	err := DB.Get(&_message, DB.Rebind(sqlStr), val)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_message, nil
}

// FindMessagesBy find all messages by a field name and a value.
func FindMessagesBy(field string, val interface{}) (_messages []Message, err error) {
	sqlFmt := `SELECT COALESCE(messages.chat_id, 0) AS chat_id, COALESCE(messages.number, '') AS number, messages.id, messages.text, messages.created_at, messages.updated_at FROM messages WHERE %s = ?`
	sqlStr := fmt.Sprintf(sqlFmt, field)
	err = DB.Select(&_messages, DB.Rebind(sqlStr), val)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _messages, nil
}

// AllMessages get all the Message records.
func AllMessages() (messages []Message, err error) {
	err = DB.Select(&messages, "SELECT COALESCE(messages.chat_id, 0) AS chat_id, COALESCE(messages.number, '') AS number, messages.id, messages.text, messages.created_at, messages.updated_at FROM messages")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return messages, nil
}

// MessageCount get the count of all the Message records.
func MessageCount() (c int64, err error) {
	err = DB.Get(&c, "SELECT count(*) FROM messages")
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return c, nil
}

// MessageCountWhere get the count of all the Message records with a where clause.
func MessageCountWhere(where string, args ...interface{}) (c int64, err error) {
	sql := "SELECT count(*) FROM messages"
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

// MessageIncludesWhere get the Message associated models records, currently it's not same as the corresponding "includes" function but "preload" instead in Ruby on Rails. It means that the "sql" should be restricted on Message model.
func MessageIncludesWhere(assocs []string, sql string, args ...interface{}) (_messages []Message, err error) {
	_messages, err = FindMessagesWhere(sql, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(assocs) == 0 {
		log.Println("No associated fields ard specified")
		return _messages, err
	}
	if len(_messages) <= 0 {
		return nil, errors.New("No results available")
	}
	ids := make([]interface{}, len(_messages))
	for _, v := range _messages {
		ids = append(ids, interface{}(v.Id))
	}
	return _messages, nil
}

// MessageIds get all the IDs of Message records.
func MessageIds() (ids []int64, err error) {
	err = DB.Select(&ids, "SELECT id FROM messages")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return ids, nil
}

// MessageIdsWhere get all the IDs of Message records by where restriction.
func MessageIdsWhere(where string, args ...interface{}) ([]int64, error) {
	ids, err := MessageIntCol("id", where, args...)
	return ids, err
}

// MessageIntCol get some int64 typed column of Message by where restriction.
func MessageIntCol(col, where string, args ...interface{}) (intColRecs []int64, err error) {
	sql := "SELECT " + col + " FROM messages"
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

// MessageStrCol get some string typed column of Message by where restriction.
func MessageStrCol(col, where string, args ...interface{}) (strColRecs []string, err error) {
	sql := "SELECT " + col + " FROM messages"
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

// FindMessagesWhere query use a partial SQL clause that usually following after WHERE
// with placeholders, eg: FindUsersWhere("first_name = ? AND age > ?", "John", 18)
// will return those records in the table "users" whose first_name is "John" and age elder than 18.
func FindMessagesWhere(where string, args ...interface{}) (messages []Message, err error) {
	sql := "SELECT COALESCE(messages.chat_id, 0) AS chat_id, COALESCE(messages.number, '') AS number, messages.id, messages.text, messages.created_at, messages.updated_at FROM messages"
	if len(where) > 0 {
		sql = sql + " WHERE " + where
	}
	stmt, err := DB.Preparex(DB.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = stmt.Select(&messages, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return messages, nil
}

// FindMessageBySql query use a complete SQL clause
// with placeholders, eg: FindUserBySql("SELECT * FROM users WHERE first_name = ? AND age > ? ORDER BY DESC LIMIT 1", "John", 18)
// will return only One record in the table "users" whose first_name is "John" and age elder than 18.
func FindMessageBySql(sql string, args ...interface{}) (*Message, error) {
	stmt, err := DB.Preparex(DB.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	_message := &Message{}
	err = stmt.Get(_message, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return _message, nil
}

// FindMessagesBySql query use a complete SQL clause
// with placeholders, eg: FindUsersBySql("SELECT * FROM users WHERE first_name = ? AND age > ?", "John", 18)
// will return those records in the table "users" whose first_name is "John" and age elder than 18.
func FindMessagesBySql(sql string, args ...interface{}) (messages []Message, err error) {
	stmt, err := DB.Preparex(DB.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = stmt.Select(&messages, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return messages, nil
}

// CreateMessage use a named params to create a single Message record.
// A named params is key-value map like map[string]interface{}{"first_name": "John", "age": 23} .
func CreateMessage(am map[string]interface{}) (int64, error) {
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
	sqlFmt := `INSERT INTO messages (%s) VALUES (%s)`
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

// Create is a method for Message to create a record.
func (_message *Message) Create() (int64, error) {
	ok, err := govalidator.ValidateStruct(_message)
	if !ok {
		errMsg := "Validate Message struct error: Unknown error"
		if err != nil {
			errMsg = "Validate Message struct error: " + err.Error()
		}
		log.Println(errMsg)
		return 0, errors.New(errMsg)
	}
	t := time.Now()
	_message.CreatedAt = t
	_message.UpdatedAt = t
    sql := `INSERT INTO messages (text,created_at,updated_at,chat_id,number) VALUES (:text,:created_at,:updated_at,:chat_id,:number)`
    result, err := DB.NamedExec(sql, _message)
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



// CreateChat is a method for a Message object to create an associated Chat record.
func (_message *Message) CreateChat(am map[string]interface{}) error {
	am["message_id"] = _message.Id
	_, err := CreateChat(am)
	return err
}

// Destroy is method used for a Message object to be destroyed.
func (_message *Message) Destroy() error {
	if _message.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := DestroyMessage(_message.Id)
	return err
}

// DestroyMessage will destroy a Message record specified by the id parameter.
func DestroyMessage(id int64) error {
	stmt, err := DB.Preparex(DB.Rebind(`DELETE FROM messages WHERE id = ?`))
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

// DestroyMessages will destroy Message records those specified by the ids parameters.
func DestroyMessages(ids ...int64) (int64, error) {
	if len(ids) == 0 {
		msg := "At least one or more ids needed"
		log.Println(msg)
		return 0, errors.New(msg)
	}
	idsHolder := strings.Repeat(",?", len(ids)-1)
	sql := fmt.Sprintf(`DELETE FROM messages WHERE id IN (?%s)`, idsHolder)
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

// DestroyMessagesWhere delete records by a where clause restriction.
// e.g. DestroyMessagesWhere("name = ?", "John")
// And this func will not call the association dependent action
func DestroyMessagesWhere(where string, args ...interface{}) (int64, error) {
	sql := `DELETE FROM messages WHERE `
	if len(where) > 0 {
		sql = sql + where
	} else {
		return 0, errors.New("No WHERE conditions provided")
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


// Save method is used for a Message object to update an existed record mainly.
// If no id provided a new record will be created. FIXME: A UPSERT action will be implemented further.
func (_message *Message) Save() error {
	ok, err := govalidator.ValidateStruct(_message)
	if !ok {
		errMsg := "Validate Message struct error: Unknown error"
		if err != nil {
			errMsg = "Validate Message struct error: " + err.Error()
		}
		log.Println(errMsg)
		return errors.New(errMsg)
	}
	if _message.Id == 0 {
		_, err = _message.Create()
		return err
	}
	_message.UpdatedAt = time.Now()
	sqlFmt := `UPDATE messages SET %s WHERE id = %v`
	sqlStr := fmt.Sprintf(sqlFmt, "text = :text, updated_at = :updated_at, chat_id = :chat_id, number = :number", _message.Id)
    _, err = DB.NamedExec(sqlStr, _message)
    return err
}

// UpdateMessage is used to update a record with a id and map[string]interface{} typed key-value parameters.
func UpdateMessage(id int64, am map[string]interface{}) error {
	if len(am) == 0 {
		return errors.New("Zero key in the attributes map!")
	}
	am["updated_at"] = time.Now()
	keys := allKeys(am)
	sqlFmt := `UPDATE messages SET %s WHERE id = %v`
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

// Update is a method used to update a Message record with the map[string]interface{} typed key-value parameters.
func (_message *Message) Update(am map[string]interface{}) error {
	if _message.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := UpdateMessage(_message.Id, am)
	return err
}

// UpdateAttributes method is supposed to be used to update Message records as corresponding update_attributes in Ruby on Rails.
func (_message *Message) UpdateAttributes(am map[string]interface{}) error {
	if _message.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := UpdateMessage(_message.Id, am)
	return err
}

// UpdateColumns method is supposed to be used to update Message records as corresponding update_columns in Ruby on Rails.
func (_message *Message) UpdateColumns(am map[string]interface{}) error {
	if _message.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := UpdateMessage(_message.Id, am)
	return err
}

// UpdateMessagesBySql is used to update Message records by a SQL clause
// using the '?' binding syntax.
func UpdateMessagesBySql(sql string, args ...interface{}) (int64, error) {
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
