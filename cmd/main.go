package main

// I learn go
// https://go.dev/ref/spec

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
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
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
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
		bytes := []byte("SETYB NO ROX")
		fmt.Println(string(bytes), len(bytes))
		for j := 0; j < len(bytes)/2; j++ {
			bytes[j] = bytes[j] ^ bytes[len(bytes)-j-1]
			bytes[len(bytes)-j-1] = bytes[len(bytes)-j-1] ^ bytes[j]
			bytes[j] = bytes[j] ^ bytes[len(bytes)-j-1]

			fmt.Println(string(bytes), len(bytes))
		}
	case "async":
		waitgroup := sync.WaitGroup{}
		waitgroup.Add(4)
		now := time.Now()
		mu := sync.Mutex{}
		var log string = "started " + now.String()
		for i := 2; i < 10; i += 2 {
			go func() {
				var exp time.Duration = (time.Duration(i) * time.Second)
				defer waitgroup.Done()
				time.Sleep(exp)
				mu.Lock()
				log += "\n" + strconv.Itoa(i) + " seconds elapsed"
				mu.Unlock()
			}()
		}
		waitgroup.Wait()
		fmt.Println(log)
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
		err = os.WriteFile(cwd+"/colors.json", []byte(marshalled), 0777)
		if err != nil {
			fmt.Println(err.Error())
			break
		}

		jsonBytes, err := os.ReadFile(cwd + "/colors.json")
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		colors := []Color{}
		json.Unmarshal(jsonBytes, &colors)
		for _, c := range colors {
			fmt.Println(c.Name)
		}
	case "xml":
		marshalled, err := GetColors("xml")
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		fmt.Println(marshalled)
		err = os.WriteFile(cwd+"/colors.xml", []byte(marshalled), 0777)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		xmlBytes, err := os.ReadFile(cwd + "/colors.xml")
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		root := Root{}
		xml.Unmarshal(xmlBytes, &root)
		for _, c := range root.Colors {
			fmt.Println(c.Name)
		}
	case "csv":
		marshalled, err := GetColors("csv")
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		fmt.Println(marshalled)
		err = os.WriteFile(cwd+"/colors.csv", []byte(marshalled), 0777)
		if err != nil {
			fmt.Println(err.Error())
			break
		}

		file, err := os.Open(cwd + "/colors.csv")
		if err != nil {
			fmt.Println(err.Error())
			break
		}

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			fmt.Println(err.Error())
			break
		}

		i := 0
		colors := []Color{}
		for _, r := range records {
			if i > 0 {
				color := Color{Name: r[0], Hash: r[1]}
				colors = append(colors, color)
			}
			i++
		}

		for _, c := range colors {
			fmt.Println(c.Name)
		}

	default:
		fmt.Println("Command not recognized")
	}
}
