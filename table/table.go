package table

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/hechh/library/mlog"
	"github.com/hechh/library/util"
)

var (
	tableObj = &Watcher{parsers: make(map[string]*Parser)}
)

type Watcher struct {
	path    string
	watcher *fsnotify.Watcher
	parsers map[string]*Parser
}

func Register(sheet string, f ParseFunc) {
	item, ok := tableObj.parsers[sheet]
	if !ok {
		item = NewParser(sheet)
		tableObj.parsers[sheet] = item
	}
	item.Register(f)
}

func Listen(sheet string, fs ...ChangeFunc) {
	item, ok := tableObj.parsers[sheet]
	if !ok {
		item = NewParser(sheet)
		tableObj.parsers[sheet] = item
	}
	item.Listen(fs...)
}

func Init(path string) error {
	if fs, err := fsnotify.NewWatcher(); err != nil {
		return err
	} else {
		tableObj.path = path
		tableObj.watcher = fs
	}

	// 设置监听目录
	if err := tableObj.watcher.Add(path); err != nil {
		return err
	}

	// 加载配置
	if err := load(true); err != nil {
		return err
	}

	go watch()
	return nil
}

func load(isload bool) error {
	files, err := util.Glob(tableObj.path, ".*\\.conf", true)
	if err != nil {
		return err
	}
	for _, filename := range files {
		sheet := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
		item, ok := tableObj.parsers[sheet]
		if !ok {
			continue
		}

		// 加载配置
		buf, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		if err := item.Parse(isload, buf); err != nil {
			return err
		}
	}
	return nil
}

func watch() {
	for {
		select {
		case _, ok := <-tableObj.watcher.Events:
			if !ok {
				return
			}
			// 加载配置
			if err := load(false); err != nil {
				mlog.Errorf("加载配置失败：%v", err)
			}
		case err, ok := <-tableObj.watcher.Errors:
			if !ok {
				return
			}
			mlog.Errorf("配置监听过程中发生错误：%v", err)
		}
	}
}

func Close() {
	tableObj.watcher.Close()
}
