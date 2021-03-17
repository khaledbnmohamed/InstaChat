package main

import (
	"database/sql"
	"time"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type InsetionToDBWorker struct {
	*Job
	chat_id float64
	newMessage string
}

func (work *InsetionToDBWorker) Perform() {
	db, err := sql.Open("mysql", "instachat:instachat@/instachat_development")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	statement, _ := db.Prepare("INSERT INTO messages (chat_id, text,created_at,updated_at) VALUES (?,? , ?,?)")
	statement.Exec(work.chat_id, work.newMessage,"2019-11-20 22:00:16.835008","2019-11-20 22:00:16.835008")

	fmt.Printf("Processed message: %s\n", work.newMessage)
}

func NewInsetionToDBWorker(job Job) Worker {
	fmt.Println("khalodaaaaaaaaa for jobs...", job.Args)

	newMessage := "2323" 
	chat_id := 22.1
   return &InsetionToDBWorker{&job, chat_id, newMessage}
}
