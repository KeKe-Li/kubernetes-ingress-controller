package utils

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/chanxuehong/log"
)

type IFile interface {
	// GetPath 获取项目路径
	GetPath(ctx context.Context) string
	// IsDirExists 判断文件目录否存在
	IsDirExists(ctx context.Context, path string) bool
	// MkdirFile 创建文件
	MkdirFile(ctx context.Context, path string) error
	// Dir 目录返回的执行用户的主目录。
	Dir(ctx context.Context) (string, error)
	// Expand 扩展的目录
	Expand(ctx context.Context, path string) (string, error)
	// dirUnix unix的目录
	dirUnix(ctx context.Context) (string, error)
	//  dirWindows windows的目录
	dirWindows(ctx context.Context) (string, error)
	// DeleteFile 删除文件
	DeleteFile(ctx context.Context, filePath string) error
	// MkDir 创建文件夹,支持x/a/a  多层级
	MkDir(ctx context.Context, path string) error
}

type file struct{}

func NewFile() IFile {
	return &file{}
}

func (impl *file) GetPath(ctx context.Context) string {
	var apiroot string

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.ErrorContext(ctx, "GetPath failed", "error", err.Error())
		return apiroot
	}

	apiroot = strings.Replace(dir, "\\", "/", -1)
	return apiroot
}

func (impl *file) IsDirExists(ctx context.Context, path string) bool {
	file, err := os.Stat(path)
	if err != nil {
		log.ErrorContext(ctx, "IsDirExists filed", "error", err.Error())
		return false
	}
	return file.IsDir()
}

func (impl *file) MkdirFile(ctx context.Context, path string) error {
	//在当前目录下生成md目录
	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		log.ErrorContext(ctx, "MkdirFile failed", "error", err.Error())
		return err
	}
	return nil
}

func (impl *file) Dir(ctx context.Context) (string, error) {
	var DisableCache bool
	var homedirCache string
	var cacheLock sync.RWMutex

	if !DisableCache {
		cacheLock.RLock()
		cached := homedirCache
		cacheLock.RUnlock()
		if cached != "" {
			log.DebugContext(ctx, "cached is not nil", "cached", cached)
			return cached, nil
		}
	}
	cacheLock.Lock()
	defer cacheLock.Unlock()

	var result string
	var err error
	if runtime.GOOS == "windows" {
		result, err = impl.dirWindows(ctx)
	} else {
		// Unix-like system, so just assume Unix
		result, err = impl.dirUnix(ctx)
	}

	if err != nil {
		log.ErrorContext(ctx, "Dir failed", "error", err.Error())
		return "", err
	}
	homedirCache = result
	return result, nil
}

func (impl *file) Expand(ctx context.Context, path string) (string, error) {
	if len(path) == 0 {
		log.DebugContext(ctx, "path==0", "path", path)
		return path, nil
	}

	if path[0] != '~' {
		log.DebugContext(ctx, "path[0] != '~'", "path", path)
		return path, nil
	}

	if len(path) > 1 && path[1] != '/' && path[1] != '\\' {
		log.DebugContext(ctx, "path is not correctly", "path", path)
		return "", errors.New("cannot expand user-specific home dir")
	}

	dir, err := impl.Dir(ctx)
	if err != nil {
		log.ErrorContext(ctx, "path is not correctly", "path", path, "err", err.Error())
		return "", err
	}

	return filepath.Join(dir, path[1:]), nil
}

func (impl *file) dirUnix(ctx context.Context) (string, error) {
	homeEnv := "HOME"
	if runtime.GOOS == "plan9" {
		// On plan9, env vars are lowercase.
		homeEnv = "home"
	}

	// First prefer the HOME environmental variable
	if home := os.Getenv(homeEnv); home != "" {
		return home, nil
	}

	var stdout bytes.Buffer

	// If that fails, try OS specific commands
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("sh", "-c", `dscl -q . -read /Users/"$(whoami)" NFSHomeDirectory | sed 's/^[^ ]*: //'`)
		cmd.Stdout = &stdout
		if err := cmd.Run(); err == nil {
			result := strings.TrimSpace(stdout.String())
			if result != "" {
				log.DebugContext(ctx, "result is not correctly", "result", result)
				return result, nil
			}
		}
	} else {
		cmd := exec.Command("getent", "passwd", strconv.Itoa(os.Getuid()))
		cmd.Stdout = &stdout
		if err := cmd.Run(); err != nil {
			// If the error is ErrNotFound, we ignore it. Otherwise, return it.
			if err != exec.ErrNotFound {
				log.DebugContext(ctx, "exec.ErrNotFound", "err", err.Error())
				return "", err
			}
		} else {
			if passwd := strings.TrimSpace(stdout.String()); passwd != "" {
				// username:password:uid:gid:gecos:home:shell
				passwdParts := strings.SplitN(passwd, ":", 7)
				if len(passwdParts) > 5 {
					return passwdParts[5], nil
				}
			}
		}
	}

	// If all else fails, try the shell
	stdout.Reset()
	cmd := exec.Command("sh", "-c", "cd && pwd")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		log.ErrorContext(ctx, "cmd.Run failed", "err", err.Error())
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		log.DebugContext(ctx, "result is empty", "result", result)
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func (impl *file) dirWindows(ctx context.Context) (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		log.DebugContext(ctx, "Getenv home", "home", home)
		return home, nil
	}

	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		log.DebugContext(ctx, "home & path is empty", "drive", drive, "path", path)
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		log.DebugContext(ctx, "home is empty", "home", home)
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}

func (impl *file) DeleteFile(ctx context.Context, filePath string) error {
	log.DebugContext(ctx, "DeleteFile", "filePath", filePath)
	return os.RemoveAll(filePath)
}

// 创建文件夹,支持x/a/a  多层级
func (impl *file) MkDir(ctx context.Context, path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			//文件夹不存在，创建
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				log.ErrorContext(ctx, "os.MkdirAll failed", "error", err.Error())
				return err
			}
		} else {
			log.ErrorContext(ctx, "os.Stat failed", "error", err.Error())
			return err
		}
	}
	return nil
}
