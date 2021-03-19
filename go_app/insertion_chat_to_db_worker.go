package main

import (
	"database/sql"
	"time"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type InsetionChatToDBWorker struct {
	*Job
	applicaitonID float64
	number float64
}

func (work *InsetionChatToDBWorker) Perform() {
	db, err := sql.Open("mysql", "instachat:instachat@tcp(instachat_database_development)/instachat_development")
	if err != nil { panic(err.Error()) }


	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	statement, err := db.Prepare("INSERT INTO chats (application_id, number,created_at,updated_at) VALUES (?,?,?,?)")
	if err != nil { panic(err.Error()) }

	statement.Exec(work.applicaitonID, work.number, time.Now(), time.Now())

	fmt.Printf("Added chat: %f\n", work.number)
}

func NewInsetionChatToDBWorker(job Job) Worker {

	applicaitonID := job.Args[0].(float64)
	number := job.Args[1].(float64)

   return &InsetionChatToDBWorker{&job, applicaitonID,number}
}
