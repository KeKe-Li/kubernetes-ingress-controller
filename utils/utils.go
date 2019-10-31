package utils

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/chanxuehong/log"
)

type IUtils interface {
	// GetDate 返回当前时间
	GetDate(ctx context.Context) string

	// ToString 时间转字符串
	ToString(ctx context.Context, t time.Time) string

	// ToStringWithFormat 时间转字符串,带格式
	ToStringWithFormat(ctx context.Context, t time.Time, format string) string

	// GetRunTime 获取当前系统环境
	GetRunTime(ctx context.Context) string

	// GetMD5 MD5 加密字符串
	GetMD5(ctx context.Context, plainText string) string

	// GetMd5 计算文件的md5，适用于本地文件计算
	GetMd5(ctx context.Context, path string) (string, error)

	// GetMd52 从流中直接读取数据计算md5 并返回流的副本
	GetMd52(ctx context.Context, file io.Reader) (io.Reader, string, error)

	// 解压
	DeCompress(ctx context.Context, zipFile, dest string) error

	// getDir 获取文件的目录
	getDir(ctx context.Context, path string) string

	// subString 截取字符串
	subString(ctx context.Context, str string, start, end int) string

	// UploadFile 上传文件
	UploadFile(ctx context.Context, file *multipart.FileHeader, path string) (string, error)

	// GetFileSize 获取文件的大小
	GetFileSize(ctx context.Context, filePath string) (int64, error)

	// CheckFileIsExist 判断文件是否存在
	CheckFileIsExist(ctx context.Context, filename string) bool
}

type utils struct {
	File IFile
}

func NewUtils(file IFile) IUtils {
	return &utils{
		File: file,
	}
}

// 返回当前时间
func (impl *utils) GetDate(ctx context.Context) string {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02 03:04:05")
}

// 时间转字符串
func (impl *utils) ToString(ctx context.Context, t time.Time) string {
	return impl.ToStringWithFormat(ctx, t, "2006-01-02 15:04:05")
}

// 时间转字符串,带格式
func (impl *utils) ToStringWithFormat(ctx context.Context, t time.Time, format string) string {
	return t.In(time.Local).Format(format)
}

// 获取当前系统环境
func (impl *utils) GetRunTime(ctx context.Context) string {
	//获取系统环境变量
	runTime := os.Getenv("RUN_TIME")
	if runTime == "" {
		log.DebugContext(ctx, "runTime is empty")
		fmt.Println("No RUN_TIME Can't start")
	}
	return runTime
}

// MD5 加密字符串
func (impl *utils) GetMD5(ctx context.Context, plainText string) string {
	h := md5.New()
	h.Write([]byte(plainText))
	return hex.EncodeToString(h.Sum(nil))
}

// 计算文件的md5，适用于本地文件计算
func (impl *utils) GetMd5(ctx context.Context, path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		log.ErrorContext(ctx, "io.Copy failed", "error", err.Error())
		return "", err
	}
	return hex.EncodeToString(md5hash.Sum(nil)), nil
}

// 从流中直接读取数据计算md5 并返回流的副本，不能用于计算大文件流否则内存占用很大
func (impl *utils) GetMd52(ctx context.Context, file io.Reader) (io.Reader, string, error) {
	var b bytes.Buffer
	md5hash := md5.New()
	if _, err := io.Copy(&b, io.TeeReader(file, md5hash)); err != nil {
		return nil, "", err
	}
	return &b, hex.EncodeToString(md5hash.Sum(nil)), nil
}

// 解压
func (impl *utils) DeCompress(ctx context.Context, zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		log.ErrorContext(ctx, "zip.OpenReader failed", "error", err.Error())
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			continue
		}
		rc, err := file.Open()
		if err != nil {
			log.ErrorContext(ctx, "file.Open failed", "error", err.Error())
			return err
		}
		defer rc.Close()
		filename := dest + file.Name
		err = os.MkdirAll(impl.getDir(ctx, filename), 0755)
		if err != nil {
			log.ErrorContext(ctx, "os.MkdirAll failed", "error", err.Error())
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			log.ErrorContext(ctx, "os.Create failed", "error", err.Error())
			return err
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			log.ErrorContext(ctx, "io.Copy failed", "error", err.Error())
			return err
		}
	}
	return nil
}

// 获取文件的目录
func (impl *utils) getDir(ctx context.Context, path string) string {
	return impl.subString(ctx, path, 0, strings.LastIndex(path, "/"))
}

// 截取字符串
func (impl *utils) subString(ctx context.Context, str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		log.DebugContext(ctx, "start is not enough", "start", start)
		return ""
	}

	if end < start || end > length {
		log.DebugContext(ctx, "end is not enough", "start", start)
		return ""
	}

	return string(rs[start:end])
}

// 上传文件
func (impl *utils) UploadFile(ctx context.Context, file *multipart.FileHeader, path string) (string, error) {
	if reflect.ValueOf(file).IsNil() || !reflect.ValueOf(file).IsValid() {
		log.DebugContext(ctx, "reflect.ValueOf failed", "file", file)
		return "", errors.New("invalid memory address or nil pointer dereference")
	}
	src, err := file.Open()
	defer src.Close()
	if err != nil {
		log.ErrorContext(ctx, "file.Open failed", "error", err.Error())
		return "", err
	}
	err = impl.File.MkDir(ctx, path)
	if err != nil {
		log.ErrorContext(ctx, "file.Open failed", "error", err.Error())
		return "", err
	}
	// Destination 去除空格
	filename := strings.Replace(file.Filename, " ", "", -1)
	// 去除换行符
	filename = strings.Replace(filename, "\n", "", -1)

	dst, err := os.Create(path + filename)
	if err != nil {
		log.ErrorContext(ctx, "os.Create failed", "error", err.Error())
		return "", err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		log.ErrorContext(ctx, "io.Copy failed", "error", err.Error())
		return "", err
	}
	return filename, nil
}

// 获取文件的大小
func (impl *utils) GetFileSize(ctx context.Context, filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.ErrorContext(ctx, "GetFileSize failed", "error", err.Error())
		return 0, err
	}
	//文件大小
	fsize := fileInfo.Size()
	return fsize, nil
}

// 判断文件是否存在
func (impl *utils) CheckFileIsExist(ctx context.Context, filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
