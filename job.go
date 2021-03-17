package main

type Job struct {
	Retry       bool          `json:retry`
	Queue       string        `json:queue`
	Class       string        `json:class`
	Args        []interface{} `json:args`
	Jid         string        `json:jid`
	Enqueued_at float64       `json:enqueued_at`
}
