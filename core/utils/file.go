package utils

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"zrWorker/core/slog"

	"github.com/pkg/errors"
)

// 换行的数据
func ReadLineData(userDict string) (users []string, err error) {
	file, err := os.Open(userDict)
	if err != nil {
		slog.Println(slog.DEBUG, "ReadLineData:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		user := strings.TrimSpace(scanner.Text())
		if user != "" {
			users = append(users, user)
		}
	}
	return users, err
}

func Write(path, str string) {
	f, err := os.OpenFile(path, os.O_CREATE+os.O_RDWR+os.O_TRUNC, 0764)
	if err != nil {
		slog.Println(slog.DEBUG, err)
	}

	//jsonBuf := append([]byte(result),[]byte("\r\n")...)
	f.WriteString(str)
}

func WriteAppend(path, str string) {
	f, err := os.OpenFile(path, os.O_CREATE+os.O_RDWR+os.O_APPEND, 0764)
	if err != nil {
		slog.Println(slog.DEBUG, err)
	}

	//jsonBuf := append([]byte(result),[]byte("\r\n")...)
	str += "\n"
	f.WriteString(str)
}

func Read(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		slog.Println(slog.DEBUG, err)
	}
	return string(content)
}

func WritePng(path string, buf []byte) {
	//slog.Println(slog.WARN, path)
	_, err := os.Stat(GetScreenPath())
	if err != nil {
		os.MkdirAll(GetScreenPath(), 0777)
	}

	f, err := os.OpenFile(path, os.O_CREATE+os.O_RDWR, 0664)
	if err != nil {
		slog.Println(slog.WARN, err)
		return
	}

	f.Write(buf)

	//slog.Println(slog.DEBUG, "图片写入完成")
}

// 压缩为zip格式
// source为要压缩的文件或文件夹, 绝对路径和相对路径都可以
// target是目标文件
// filter是过滤正则(Golang 的 包 path.Match)
func ZipFile(source, target, filter string) error {
	var err error
	if isAbs := filepath.IsAbs(source); !isAbs {
		source, err = filepath.Abs(source) // 将传入路径直接转化为绝对路径
		if err != nil {
			return errors.WithStack(err)
		}
	}
	//创建zip包文件
	zipfile, err := os.Create(target)
	if err != nil {
		return errors.WithStack(err)
	}

	defer func() {
		if err := zipfile.Close(); err != nil {
			slog.Println(slog.DEBUG, err.Error(), zipfile.Name())
		}
	}()

	//创建zip.Writer
	zw := zip.NewWriter(zipfile)

	defer func() {
		if err := zw.Close(); err != nil {
			slog.Println(slog.DEBUG, err.Error())
		}
	}()

	info, err := os.Stat(source)
	if err != nil {
		return errors.WithStack(err)
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return errors.WithStack(err)
		}

		//将遍历到的路径与pattern进行匹配
		ism, err := filepath.Match(filter, info.Name())

		if err != nil {
			return errors.WithStack(err)
		}
		//如果匹配就忽略
		if !ism {
			return nil
		}
		//创建文件头
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return errors.WithStack(err)
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		//写入文件头信息
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return errors.WithStack(err)
		}

		if info.IsDir() {
			return nil
		}
		//写入文件内容
		file, err := os.Open(path)
		if err != nil {
			return errors.WithStack(err)
		}

		defer func() {
			if err := file.Close(); err != nil {
				slog.Println(slog.DEBUG, err.Error(), file.Name())
			}
		}()
		_, err = io.Copy(writer, file)

		return errors.WithStack(err)
	})

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func ByteCountIEC(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

func ReadList(fs io.Reader) (users []string, err error) {
	sourceBuf, err := io.ReadAll(fs)
	if err != nil {
		return nil, err
	}
	source := strings.Split(string(sourceBuf), "\n")

	return source, nil
}
