package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

func (mgr *UploadMgr) uploadDir(localPath, remotePath string) error {
	// 最外层目录
	mgr.sftp.MkdirAll(remotePath)

	// 遍历本地目录
	err := filepath.WalkDir(localPath, func(filePath string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if _, ok := mgr.options.filter[info.Name()]; ok {
			if info.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			}
			//fmt.Printf("name:%s conitnue\n", info.Name())
		}

		localFilePath := filePath
		remoteFilePath := path.Join(remotePath, WindowsPathToLinuxPath(filePath))

		// 判断是否是文件,是文件直接上传.是文件夹,先远程创建文件夹,再递归复制内部文件
		if info.IsDir() {
			mgr.sftp.MkdirAll(remoteFilePath)
		} else {
			//fmt.Println(localFilePath, "   ", remoteFilePath)
			err = mgr.uploadFile(localFilePath, remoteFilePath)
			if err != nil {
				return err
			}
		}
		fmt.Printf("File %s uploaded successfully!\n", remoteFilePath)

		return err
	})
	return err
}

func (mgr *UploadMgr) uploadFile(localFilePath, remoteFilePath string) error {
	// 本地文件流
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		fmt.Printf("os.Open error : %s", localFilePath)
		return err
	}
	defer srcFile.Close()
	fileInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()

	dstFile, err := mgr.sftp.Create(remoteFilePath)
	if err != nil {
		fmt.Printf("sftpClient.OnCreate %s err: %s", remoteFilePath, err)
		return err
	}
	defer dstFile.Close()

	// 使用bufio.Reader来包装localFile，以便使用缓冲区读取数据
	reader := bufio.NewReader(srcFile)
	buffer := make([]byte, SegmentSize) // 设置一个合适大小的缓冲区

	// 循环读取本地文件并写入远程文件
	okNum := 0
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic("Failed to read local file: " + err.Error())
		}

		// 将读取到的数据写入远程文件
		_, err = dstFile.Write(buffer[:n])
		if err != nil {
			panic("Failed to write to remote file: " + err.Error())
		}

		okNum += n

		fmt.Printf("上传进度: %.2f%%\n", float64(okNum)/float64(fileSize)*100)
	}

	return nil
}

// WindowsPathToLinuxPath 将Windows风格的路径转换为Linux风格的路径
func WindowsPathToLinuxPath(windowsPath string) string {
	// 使用strings.ReplaceAll替换所有的反斜杠为正斜杠
	return strings.ReplaceAll(windowsPath, `\`, `/`)
}

func winBase(path string) string {
	if path == "" {
		return "."
	}
	// Strip trailing slashes.
	for len(path) > 0 && path[len(path)-1] == '\\' {
		path = path[0 : len(path)-1]
	}
	// Find the last element
	if i := strings.LastIndex(path, "\\"); i >= 0 {
		path = path[i+1:]
	}
	if i := strings.LastIndex(path, "/"); i >= 0 {
		path = path[i+1:]
	}
	// If empty now, it had only slashes.
	if path == "" {
		return "/"
	}
	return path
}
func (mgr *UploadMgr) uploadDirectory(localPath string, remotePath string) error {
	// 本地文件夹流
	localFiles, err := os.ReadDir(localPath)
	if err != nil {
		fmt.Printf("路径错误 %s", err.Error())
		return err
	}
	// 先创建最外层文件夹
	mgr.sftp.Mkdir(remotePath)
	// 遍历文件夹内容
	for _, backupDir := range localFiles {
		localFilePath := path.Join(localPath, backupDir.Name())
		remoteFilePath := path.Join(remotePath, backupDir.Name())
		// 判断是否是文件,是文件直接上传.是文件夹,先远程创建文件夹,再递归复制内部文件
		if backupDir.IsDir() {
			mgr.sftp.Mkdir(remoteFilePath)
			mgr.uploadDirectory(localFilePath, remoteFilePath)
		} else {
			mgr.uploadFile1(path.Join(localPath, backupDir.Name()), remotePath)
		}
	}
	fmt.Printf("File %s uploaded successfully!\n", remotePath)
	return nil
}
func (mgr *UploadMgr) uploadFile1(localFilePath string, remotePath string) error {
	// 本地文件流
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		fmt.Printf("os.Open error : %s", localFilePath)
		return err
	}
	defer srcFile.Close()
	fileInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()

	// 创建一个通道来接收已传输的字节数
	progressCh := make(chan int64)
	// 启动一个goroutine来打印进度
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		var lastReported int64
		for reported := range progressCh {
			if reported == fileSize {
				// 上传完成，打印最终进度
				fmt.Printf("上传进度: 100.00%%\n")
				break
			}
			progress := float64(reported) / float64(fileSize) * 100
			if reported != lastReported {
				// 仅当进度更新时打印
				fmt.Printf("上传进度: %.2f%%\n", progress)
				lastReported = reported
			}
			select {
			case <-ticker.C:
				// 定时更新，确保即使进度没有变化也打印一次
				fmt.Printf("上传进度: %.2f%%\n", progress)
			default:
				// 如果没有到定时器的时间，则不执行任何操作
			}
		}
	}()

	// 使用自定义的progressReader来包装本地文件
	progressReader := &ProgressReader{
		reader: srcFile,
		total:  fileSize,
		ch:     progressCh,
	}

	// 上传到远端服务器的文件名,与本地路径末尾相同
	var remoteFileName string
	if runtime.GOOS == "windows" {
		remoteFileName = winBase(localFilePath)
	} else {
		remoteFileName = path.Base(localFilePath)
	}

	// 远程文件流
	dstFile, err := mgr.sftp.Create(path.Join(remotePath, remoteFileName))
	if err != nil {
		fmt.Printf("sftpClient.OnCreate %s err: %s\n", path.Join(remotePath, remoteFileName), err)
		return err
	}
	defer dstFile.Close()

	// 使用io.Copy来复制文件，并使用progressReader来跟踪进度
	_, err = io.Copy(dstFile, progressReader)
	if err != nil {
		return err
	}

	fmt.Printf("File %s uploaded successfully!\n", path.Join(remotePath, remoteFileName))
	return nil
}

// progressReader 包装一个io.Reader来跟踪已读取的字节数
type ProgressReader struct {
	reader   io.Reader
	reported int64
	total    int64
	ch       chan<- int64
}

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.reader.Read(p)
	atomic.AddInt64(&pr.reported, int64(n))
	select {
	case pr.ch <- pr.reported:
	default: // 如果channel满了，我们就不发送进度更新了（防止阻塞）
	}
	return n, err
}
