package main

import (
	"encoding/json"
	"os"
	"time"
)

var loaded = false
var contents = []Content{}

type Content struct {
	ContentId int
	UserId    int
	UserName  string
	Title     string
	Body      string
	Created   time.Time
	Updated   time.Time
	Status    int // 2 = published, 1 = draft  0 = deleted
	Access    bool
}

func PostContent(c *Content) ([]Content, error) {
	if c.ContentId == 0 {
		c.ContentId = len(contents) + 1
		contents = append(contents, *c)
	} else {
		contents[c.ContentId-1] = *c
	}
	err := writeContent()

	return GetUserContent(c.UserId), err
}

func GetContent(id int) Content {
	if !loaded {
		readContent()
	}
	return contents[id-1]
}

func GetUserContent(user int) []Content {
	if !loaded {
		readContent()
	}
	userContent := []Content{}
	for _, c := range contents {
		c.Access = (c.UserId == user)
		if (c.UserId == user && c.Status == 1) || c.Status == 2 {
			userContent = append(userContent, c)
		}
	}

	return userContent
}

func readContent() error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	bytes, err := os.ReadFile(path + "/contents.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &contents)

	if err != nil {
		return err
	}

	loaded = true
	return nil
}

func writeContent() error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(contents)
	if err != nil {
		return err
	}
	err = os.WriteFile(path+"/contents.json", bytes, 0777)
	if err != nil {
		return err
	}
	return nil
}
