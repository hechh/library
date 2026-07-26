package registry

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hechh/library/base/fileutil"
	"github.com/hechh/library/fwatcher/domain"
	"github.com/hechh/library/fwatcher/internal/parser"
	"github.com/hechh/library/mlog"
	"google.golang.org/protobuf/encoding/prototext"
)

var (
	parsers = make(map[string]domain.IParser)
)

// Register 注册配置解析函数
func Register[T any](sheet string, parseFunc func(*T) error) {
	parsers[sheet] = parser.NewParser(sheet, parseFunc)
}

// RegisterChange 注册配置变更回调函数
func RegisterChange(sheet string, changeFunc func()) {
	if item, ok := parsers[sheet]; ok {
		item.RegisterChange(changeFunc)
	}
}

func Glob(pattern string) (map[string]string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, filename := range matches {
		sheet := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
		result[sheet] = filename
	}
	return result, nil
}

func Load(files map[string]string) error {
	hh := md5.New()
	for sheet, par := range parsers {
		filename, ok := files[sheet]
		if !ok {
			return fmt.Errorf("config sheet %q not found", sheet)
		}

		// 读取文件内容
		buf, err := os.ReadFile(filename)
		if err != nil {
			return err
		}

		// 解析配置
		if err := par.Parse(hh, buf); err != nil {
			mlog.Errorf("失败加载配置:%s，error:%v", filename, err)
			return err
		}
	}
	return nil
}

func Save(sheet, filename string, body []byte) error {
	parse, ok := parsers[sheet]
	if !ok {
		return fmt.Errorf("配置%s未编译注册", sheet)
	}

	ary, err := parse.New(body)
	if err != nil {
		return err
	}

	text, err := prototext.Marshal(ary)
	if err != nil {
		return err
	}

	return fileutil.Save(filename, text)
}
