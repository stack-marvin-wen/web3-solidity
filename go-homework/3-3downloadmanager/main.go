package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type DownLoadTask struct {
	URL      string
	FilePath string
	Size     int64
	Progress int
	Done     bool // 新增：标记是否完成
	Error    error
	mu       sync.Mutex
}
type DownLoaderManager struct {
	tasks      []*DownLoadTask
	concurrent int
	wg         sync.WaitGroup
	mu         sync.Mutex
	progress   chan *DownLoadTask
}

func NewDownLoaderManager(concurrent int) *DownLoaderManager {
	return &DownLoaderManager{
		tasks:      make([]*DownLoadTask, 0),
		concurrent: concurrent,
		progress:   make(chan *DownLoadTask, 100),
	}
}

func (d *DownLoaderManager) addTask(url, filePath string) *DownLoadTask {
	d.mu.Lock()
	defer d.mu.Unlock()
	task := &DownLoadTask{
		URL:      url,
		FilePath: filePath,
		Progress: 0,
	}
	d.tasks = append(d.tasks, task)
	return task
}

func (d *DownLoaderManager) start() {
	for _, task := range d.tasks {
		if task.Done {
			continue
		}
		size, err := getFileSize(task.URL)
		if err != nil {
			task.Error = err
			continue
		}
		task.Size = size
	}
	semaphore := make(chan struct{}, d.concurrent)
	for _, task := range d.tasks {
		if task.Done || task.Error != nil {
			continue
		}
		d.wg.Add(1)
		go func(task *DownLoadTask) {
			defer d.wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			d.downloadTask(task)
		}(task)
	}
	go d.monitorProgress()
	d.wg.Wait()
	close(d.progress)
}
func (d *DownLoaderManager) downloadTask(task *DownLoadTask) {
	dir := filepath.Dir(task.FilePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		task.Error = err
		task.Done = true
		task.Progress = 100
		return
	}
	file, err := os.Create(task.FilePath)
	if err != nil {
		task.Error = err
		task.Done = true
		task.Progress = 100
		return
	}
	defer file.Close()
	resp, err := http.Get(task.URL)
	if err != nil {
		task.Error = err
		task.Done = true
		task.Progress = 100
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		task.Error = errors.New("http status code error")
		task.Done = true
		task.Progress = 100
		return
	}
	if task.Size == 0 {
		if contentLength := resp.Header.Get("Content-Length"); contentLength != "" {
			if size, err := strconv.ParseInt(contentLength, 10, 64); err == nil {
				task.Size = size
			}
		}
	}

	var downloaded int64
	buf := make([]byte, 32*1024) // 32KB缓冲区
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			// 写入文件
			if _, writeErr := file.Write(buf[:n]); writeErr != nil {
				task.Error = writeErr
				break
			}

			// 更新进度
			downloaded += int64(n)
			if task.Size > 0 {
				progress := int(float64(downloaded) / float64(task.Size) * 100)
				if progress > 100 {
					progress = 100
				}

				task.mu.Lock()
				task.Progress = progress
				task.mu.Unlock()
			}

			// 发送进度更新
			d.progress <- task
		}

		if err != nil {
			if err != io.EOF {
				task.Error = err
			}
			break
		}
	}
	task.mu.Lock()
	if task.Error == nil {
		task.Progress = 100
	}
	task.Done = true
	task.mu.Unlock()

	// 发送最终进度
	d.progress <- task
}
func (dm *DownLoaderManager) monitorProgress() {
	for task := range dm.progress {
		task.mu.Lock()
		if task.Error != nil {
			fmt.Printf("❌ 下载失败 %s: %v\n", filepath.Base(task.FilePath), task.Error)
		} else if task.Done {
			fmt.Printf("✅ 下载完成 %s: 100%%\n", filepath.Base(task.FilePath))
		} else {
			fmt.Printf("⏳ 下载中 %s: %d%%\n", filepath.Base(task.FilePath), task.Progress)
		}
		task.mu.Unlock()
	}
}
func getFileSize(url string) (int64, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	if contentLength := resp.Header.Get("Content-Length"); contentLength != "" {
		size, err := strconv.ParseInt(contentLength, 10, 64)
		if err != nil {
			return 0, err
		}
		return size, nil
	}

	return 0, nil
}
func (dm *DownLoaderManager) Wait() {
	dm.wg.Wait()
}
func (dm *DownLoaderManager) GetTasks() []*DownLoadTask {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	return dm.tasks
}
func main() {
	dm := NewDownLoaderManager(3)
	// 添加下载任务
	dm.addTask("https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf", "./downloads/dummy.pdf")
	dm.addTask("https://file-examples.com/wp-content/uploads/2017/02/file-sample_100kB.doc", "./downloads/sample.doc")
	dm.addTask("https://file-examples.com/wp-content/uploads/2017/10/file_example_JPG_100kB.jpg", "./downloads/sample.jpg")
	dm.addTask("https://file-examples.com/wp-content/uploads/2017/02/file_example_XLS_10.xls", "./downloads/sample.xls")
	fmt.Println("开始下载...")

	// 开始所有下载
	dm.start()

	fmt.Println("所有下载任务完成")

	// 打印最终状态
	fmt.Println("\n最终状态:")
	for i, task := range dm.GetTasks() {
		task.mu.Lock()
		status := "进行中"
		if task.Done {
			if task.Error != nil {
				status = "失败"
			} else {
				status = "完成"
			}
		}

		sizeStr := "未知"
		if task.Size > 0 {
			sizeStr = formatBytes(task.Size)
		}

		fmt.Printf("任务 %d: %s - %s (%s) - %d%% - 错误: %v\n",
			i+1,
			filepath.Base(task.FilePath),
			status,
			sizeStr,
			task.Progress,
			task.Error)
		task.mu.Unlock()
	}

}

// 格式化字节大小
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%dB", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB",
		float64(bytes)/float64(div), "KMGTPE"[exp])
}
