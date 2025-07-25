package main

import "time"

var contents = []Content{}

type Content struct {
	UserId   int
	UserName string
	Title    string
	Hash     string
	Created  time.Time
	Updated  time.Time
	Status   int
}
