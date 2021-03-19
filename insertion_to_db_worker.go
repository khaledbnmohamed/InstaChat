package main

import (
	"database/sql"
	"net/http"
	"io/ioutil"
	"time"
	"os"
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

	baseUrl := os.Getenv("DEVELOPMENT_HOST")

	statement, err := db.Prepare("INSERT INTO messages (chat_id, number, text,created_at,updated_at) VALUES (?,?,?,?,?)")
	if err != nil { panic(err.Error()) }

	statement.Exec(work.chatID, work.number, work.newMessage, time.Now(), time.Now())

	fmt.Printf("Added message: %s\n", work.newMessage)

	resp, err := http.Get(baseUrl + "/api/v1/messages/reindex")
	if err != nil { panic(err.Error()) }

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil { panic(err.Error()) }

	sb := string(body)
	fmt.Printf("Indexed the message: %s\n", sb)

}

func NewInsetionToDBWorker(job Job) Worker {

	chatID := job.Args[0].(float64)
	newMessage := job.Args[1].(string)
	number := job.Args[2].(float64)

   return &InsetionToDBWorker{&job, chatID, newMessage,number}
}
