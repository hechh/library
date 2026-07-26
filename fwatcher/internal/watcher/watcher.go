package watcher

import (
	"fmt"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/hechh/library/base/fileutil"
	"github.com/hechh/library/base/logic"
	"github.com/hechh/library/base/safe"
	"github.com/hechh/library/fwatcher/internal/registry"
	"github.com/hechh/library/mlog"
)

type Watcher struct {
	dataPath  string
	xlsxPath  string
	ext       string
	pattern   string
	fswatcher *fsnotify.Watcher
	exitCh    chan struct{}
}

func NewWatcher(dataPath, xlsxPath string, ext string) *Watcher {
	return &Watcher{
		dataPath: dataPath,
		xlsxPath: xlsxPath,
		ext:      ext,
		exitCh:   make(chan struct{}),
	}
}

func (d *Watcher) Init() error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	d.fswatcher = w

	// 创建目录
	abspath, err := filepath.Abs(d.dataPath)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}
	if err := fileutil.EnsureDir(abspath); err != nil {
		return err
	}

	// 监听目录 文件变更
	if err := d.fswatcher.Add(abspath); err != nil {
		return err
	}

	// 初始化加载
	d.pattern = fmt.Sprintf("%s/*%s", abspath, d.ext)
	files, err := registry.Glob(d.pattern)
	if err != nil {
		return err
	}
	if err := registry.Load(files); err != nil {
		return err
	}

	// 监听变更事件
	safe.SafeGo(mlog.Fatalf, d.watch)
	return nil
}

func (d *Watcher) Close() {
	close(d.exitCh)
}

func (d *Watcher) watch() {
	defer d.fswatcher.Close()
	for {
		select {
		case <-d.exitCh:
			return
		case event, ok := <-d.fswatcher.Events:
			if !ok {
				return
			}
			if logic.Has(event.Op, fsnotify.Write) || logic.Has(event.Op, fsnotify.Create) {
				files, err := registry.Glob(d.pattern)
				if err != nil {
					mlog.Errorf("搜索所有*.%s 文件失败 error=%v", d.ext, err)
					continue
				}

				if err := registry.Load(files); err != nil {
					mlog.Errorf("游戏配置加载失败 error=%v", err)
				}
			}
		}
	}
}
