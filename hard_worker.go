package main

import (
   	    "database/sql"
	      "fmt"
)

type Hardwork struct {
	*Job
	room_id float64
	user_id float64
	newMessage string
}

func (work *Hardwork) Perform() {
        database, _ :=  sql.Open("sql", "development.sqlite3")
        statement, _ := database.Prepare("INSERT INTO room_messages (room_id, user_id, message,created_at,updated_at) VALUES (?, ?,? , ?,?)")
        statement.Exec(work.room_id, work.user_id, work.newMessage,"2019-11-20 22:00:16.835008","2019-11-20 22:00:16.835008")

	fmt.Printf("Processed message: %s\n", work.newMessage)
}

func NewHardwork(job Job) Worker {
	room_id := job.Args[0].(float64)
	user_id := job.Args[1].(float64)
	newMessage := job.Args[2].(string)
   return &Hardwork{&job, room_id, user_id, newMessage}
}
