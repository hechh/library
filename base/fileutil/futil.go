package fileutil

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"

	"golang.org/x/tools/imports"
)

// EnsureDir 判断目录是否存在，如果不存在则创建目录
func EnsureDir(dir string) error {
	_, err := os.Stat(dir)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}
	return os.MkdirAll(dir, os.FileMode(0o755))
}

// CreateFile 创建文件
func CreateFile(fileName string, flag int) (fb *os.File, err error) {
	// 判断路径是否存在
	pp := filepath.Dir(fileName)
	if err := os.MkdirAll(pp, os.FileMode(0o755)); err != nil {
		return nil, err
	}
	// 创建文件
	if fb, err = os.OpenFile(fileName, flag, os.FileMode(0o644)); err != nil {
		return nil, err
	}
	return
}

// IsSameFile 文件是否相同
func IsSameFile(fb *os.File, name string) bool {
	st2, _ := os.Stat(name)
	st1, _ := fb.Stat()
	return os.SameFile(st1, st2)
}

// Save 保存文件（自动创建目录，.go 文件自动格式化 import）
func Save(fileName string, buf []byte) error {
	dir := filepath.Dir(fileName)
	if err := os.MkdirAll(dir, os.FileMode(0o755)); err != nil {
		return err
	}
	if ext := filepath.Ext(fileName); ext == ".go" {
		var err error
		if buf, err = imports.Process(fileName, buf, nil); err != nil {
			return err
		}
	}
	return os.WriteFile(fileName, buf, os.FileMode(0o644))
}

// ParseFiles 解析go文件
func ParseFiles(v ast.Visitor, files ...string) error {
	fset := token.NewFileSet()
	for _, filename := range files {
		// 解析语法树
		fs, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		// 遍历语法树
		ast.Walk(v, fs)
		/*
			buf := bytes.NewBuffer(nil)
			ast.Fprint(buf, fset, fs, nil)
			os.WriteFile(fmt.Sprintf("%s.ini", filename), buf.Bytes(), 0644)
		*/
	}
	return nil
}

// Glob 遍历目录所有文件
func Glob(dir, pattern string, recursive bool) (rets []string, err error) {
	pre, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// 不深度迭代
		if !recursive && info.IsDir() && dir != path {
			return filepath.SkipDir
		}
		// 过滤目录
		if info.IsDir() {
			return nil
		}
		// 是否配置
		if pre.MatchString(path) {
			rets = append(rets, path)
		}
		return nil
	})
	return
}

func SearchFile(filename string, depth int) string {
	if _, err := os.Stat(filename); err == nil {
		abs, _ := filepath.Abs(filename)
		return abs
	}
	if depth <= 0 {
		return ""
	}
	return SearchFile(filepath.Join("..", filename), depth-1)
}
