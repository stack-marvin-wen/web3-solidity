package main

import (
	"fmt"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()
	dirs := []string{"./logs"}
	for _, dir := range dirs {
		err = watcher.Add(dir)
		if err != nil {
			panic(err)
		}
	}
	go processEvents(watcher)
	fmt.Println("文件监控系统启动")
	select {}
}
func processEvents(watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			handleEvent(event)
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("error: %v", err)
		}
	}
}
func handleEvent(event fsnotify.Event) {
	filename := event.Name
	op := event.Op
	timestamp := time.Now().Format("2026-10-25 11:11:12")
	switch op {
	case fsnotify.Create:
		fmt.Printf("[%s] 创建文件：%s\n", timestamp, filename)
	case fsnotify.Write:
		fmt.Printf("[%s] 修改文件：%s\n", timestamp, filename)
	case fsnotify.Remove:
		fmt.Printf("[%s] 删除文件：%s\n", timestamp, filename)
	case fsnotify.Rename:
		fmt.Printf("[%s] 重命名文件：%s\n", timestamp, filename)
	case fsnotify.Chmod:
		fmt.Printf("[%s] 修改文件权限：%s\n", timestamp, filename)
	}
}
