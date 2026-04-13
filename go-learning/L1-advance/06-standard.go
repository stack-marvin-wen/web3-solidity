package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type Person struct {
	Age  int    `json:"age"`
	Name string `json:"name"`
}

func jsonDemo() {
	p := Person{
		Age:  18,
		Name: "Tom",
	}
	data, _ := json.Marshal(p)
	fmt.Println(string(data))
	p2 := Person{}
	json.Unmarshal(data, &p2)
	fmt.Printf("%+v\n", p2)
}

func fileDemo() {
	// 读取文件
	data, err := ioutil.ReadFile("file.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	// 写入文件
	err = ioutil.WriteFile("output.txt", []byte("Hello"), 0644)
	if err != nil {
		panic(err)
	}
}
func timeDemo() {
	now := time.Now()
	fmt.Println("当前时间:", now)
	fmt.Println("格式化:", now.Format("2006-01-02 15:04:05"))

	// 解析时间
	t, _ := time.Parse("2006-01-02", "2024-01-01")
	fmt.Println("解析的时间:", t)

	// 计算时间差
	duration := time.Now().Sub(t)
	fmt.Println("时间差:", duration)
}

// 加密哈希
func hashDemo() {
	data := "Hello World"

	// SHA256
	h := sha256.New()
	h.Write([]byte(data))
	fmt.Printf("SHA256: %x\n", h.Sum(nil))

	// MD5
	m := md5.New()
	m.Write([]byte(data))
	fmt.Printf("MD5: %x\n", m.Sum(nil))
}
func main() {
	jsonDemo()
	timeDemo()
	hashDemo()
}
