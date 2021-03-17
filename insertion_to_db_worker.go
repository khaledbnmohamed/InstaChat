package main

import (
	"database/sql"
	"time"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type InsetionToDBWorker struct {
	*Job
	chatID float64
	newMessage string
	number float64
}

func (work *InsetionToDBWorker) Perform() {
	db, err := sql.Open("mysql", "instachat:instachat@/instachat_development")
	if err != nil { panic(err.Error()) }


	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	statement, err := db.Prepare("INSERT INTO messages (chat_id, number, text,created_at,updated_at) VALUES (?,?,?,?,?)")
	if err != nil { panic(err.Error()) }

	statement.Exec(work.chatID, work.number, work.newMessage, time.Now(), time.Now())

	fmt.Printf("Processed message: %s\n", work.newMessage)
}

func NewInsetionToDBWorker(job Job) Worker {

	chatID := job.Args[0].(float64)
	newMessage := job.Args[1].(string)
	number := job.Args[2].(float64)

   return &InsetionToDBWorker{&job, chatID, newMessage,number}
}
