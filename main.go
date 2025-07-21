package main

// I learn go
// https://go.dev/ref/spec

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func Hello(name string) (string, error) {
	if name == "" {
		return "", errors.New("name is empty")
	}
	message := fmt.Sprintf("Hello, %v", name)
	return message, nil
}

func InPlace(n1 int, n2 int) (int, int) {
	// xor
	n1 = n1 ^ n2
	n2 = n2 ^ n1
	n1 = n1 ^ n2

	return n1, n2
}

func Months() [12]string {
	var dt = time.Date(2000, time.January, 1, 1, 0, 0, 0, time.UTC)
	var months [12]string
	for i := 0; i < 12; i++ {
		months[i] = dt.AddDate(0, i, 0).Month().String()
	}
	return months
}

type Vec2 struct {
	x float32
	y float32
}

// v1 is a pointer (if not we get a copy of v1)
func (v1 *Vec2) add(v2 Vec2) {
	v1.x += v2.x
	v1.y += v2.y
}

// UseOp takes a function which takes two ints and returns an int
func UseOP(op func(n1 int, n2 int) int, n1 int, n2 int) int {
	return op(n1, n2)
}

func Minus(n1 int, n2 int) int {
	return n1 - n2
}

func main() {
	// args are passed in os.Args
	fmt.Println(len(os.Args), os.Args)
	// command is the last arg
	command := os.Args[len(os.Args)-1]

	switch command {
	case "hello":
		message, err := Hello("World")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(message)
	case "inplace":
		var n1, n2 int = 5, 9
		fmt.Println("n1=", strconv.Itoa(n1), "n2=", strconv.Itoa(n2))
		n1, n2 = InPlace(n1, n2)
		fmt.Println("n1=", strconv.Itoa(n1), "n2=", strconv.Itoa(n2))
	case "months":
		now := time.Now()
		// use _ to ignore the index
		for _, m := range Months() {
			if now.Month().String() == m {
				fmt.Println("This month is", m)
			}
		}
	case "vec2":
		v1 := Vec2{x: 10, y: 5}
		fmt.Println("x=", v1.x, "y=", v1.y)
		v2 := Vec2{x: -5, y: 5}
		v1.add(v2)
		fmt.Println("x=", v1.x, "y=", v1.y)
	case "higher":
		n := UseOP(Minus, 10, 8)
		fmt.Println("n=", n)

		n = UseOP(func(n1 int, n2 int) int { return n1 % n2 }, 12, 8)
		fmt.Println("n=", n)
	case "json":
		marshalled, err := GetColors("json")
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		fmt.Println(marshalled)
		path, err := os.Getwd()
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		err = os.WriteFile(path+"/colors.json", []byte(marshalled), 0777)
		if err != nil {
			fmt.Println(err.Error())
			break
		}

		bytes, err := os.ReadFile(path + "/colors.json")
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		colors := []Color{}
		json.Unmarshal(bytes, &colors)
		for _, c := range colors {
			fmt.Println(c.Name)
		}

	case "xml":
		fmt.Println(GetColors("xml"))
	case "csv":
		fmt.Println(GetColors("csv"))
	default:
		fmt.Println("Command not recognized")
	}
}
