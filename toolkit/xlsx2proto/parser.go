package xlsx2proto

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/hechh/library/uerror"
	"github.com/hechh/library/util"
	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
)

type MsgParser struct {
	sts     map[string]*StructDescriptor
	data    map[string]*EnumDescriptor
	enums   []*EnumDescriptor
	structs []*StructDescriptor
	configs []*StructDescriptor
}

func NewMsgParser() *MsgParser {
	return &MsgParser{
		sts:  make(map[string]*StructDescriptor),
		data: make(map[string]*EnumDescriptor),
	}
}

// @config[:col]|sheet:MessageName
// @enum|sheet
func (d *MsgParser) ParseFile(filename string) error {
	fp, err := excelize.OpenFile(filename)
	if err != nil {
		return uerror.Err(-1, "打开文件(%s)失败：%v", filename, err)
	}
	defer fp.Close()

	// 读取生成表
	rows, err := fp.GetRows("生成表")
	if err != nil {
		return uerror.Err(-1, "文件(%s)生成表不存在：%v", filename, err)
	}

	for _, items := range rows {
		for _, val := range items {
			if !strings.HasPrefix(val, "@") {
				continue
			}
			strs := strings.Split(val, "|")
			switch strings.ToLower(strs[0]) {
			case "@config":
				names := strings.Split(strs[1], ":")
				rows, err := fp.GetRows(names[0])
				if err != nil {
					return uerror.Err(-1, "表格(%s)不存在", strs[1])
				}
				d.parseConfig(names[1], rows)
			case "@config:col":
				names := strings.Split(strs[1], ":")
				rows, err := fp.GetCols(names[0])
				if err != nil {
					return uerror.Err(-1, "表格(%s)不存在", strs[1])
				}
				d.parseConfig(names[1], rows)
			case "@enum":
				rows, err := fp.GetRows(strs[1])
				if err != nil {
					return uerror.Err(-1, "表格(%s)不存在", strs[1])
				}
				d.parseEnum(rows)
			case "@struct":
				rows, err := fp.GetRows(strs[1])
				if err != nil {
					return uerror.Err(-1, "表格(%s)不存在", strs[1])
				}
				d.parseStruct(rows)
			}
		}
	}
	return nil
}

func (d *MsgParser) parseConfig(name string, rows [][]string) {
	st := NewStructDescriptor(name)
	for i, item := range rows[1] {
		if len(item) <= 0 {
			continue
		}
		st.Put(int32(i)+1, rows[0][i], item, rows[2][i])
	}
	d.configs = append(d.configs, st)
}

// S|CoinReward|CoinType|CoinType|货币类型
func (d *MsgParser) parseStruct(rows [][]string) {
	for _, items := range rows {
		for _, val := range items {
			if !strings.HasPrefix(val, "S|") && !strings.HasPrefix(val, "s|") {
				continue
			}
			strs := strings.Split(val, "|")
			st, ok := d.sts[strs[1]]
			if !ok {
				st = NewStructDescriptor(strs[1])
				d.sts[st.Name] = st
				d.structs = append(d.structs, st)
			}
			st.Put(int32(len(st.List)+1), strs[2], strs[3], strs[4])
		}
	}
}

// E|游戏类型-德州NORMAL|GameType|Normal|1
func (d *MsgParser) parseEnum(rows [][]string) {
	for _, items := range rows {
		for _, val := range items {
			if !strings.HasPrefix(val, "E|") && !strings.HasPrefix(val, "e|") {
				continue
			}
			strs := strings.Split(val, "|")
			enum, ok := d.data[strs[2]]
			if !ok {
				enum = NewEnumDescriptor(strs[2])
				d.enums = append(d.enums, enum)
				d.data[strs[2]] = enum
			}
			enum.Put(cast.ToInt32(strs[4]), strs[3], strs[2], strs[1])
		}
	}
}

func (d *MsgParser) get() (rets []string) {
	flagEnum, flagStruct := 0, 0
	for _, item := range d.configs {
		for _, field := range item.List {
			if _, ok := d.data[field.Class]; ok && flagEnum == 0 {
				rets = append(rets, "enum.gen.proto")
				flagEnum++
			}
			if _, ok := d.sts[field.Class]; ok && flagStruct == 0 {
				rets = append(rets, "struct.gen.proto")
				flagStruct++
			}
			if flagEnum > 0 && flagStruct > 0 {
				return
			}
		}
	}
	return
}

func (d *MsgParser) Gen(pkgname, option, dst string) error {
	buf := bytes.NewBuffer(nil)
	pos, _ := buf.WriteString(fmt.Sprintf(headTempl, pkgname, option))
	// 枚举
	if len(d.enums) > 0 {
		for _, item := range d.enums {
			item.Sort()
		}
		sort.Slice(d.enums, func(i, j int) bool {
			return strings.Compare(d.enums[i].Name, d.enums[j].Name) <= 0
		})
		enumTpl, err := template.New("enum").Parse(enumTempl)
		if err != nil {
			return err
		}
		if err := enumTpl.Execute(buf, d.enums); err != nil {
			return err
		}
		if err := util.Save(dst, "enum.gen.proto", buf.Bytes()); err != nil {
			return err
		}
	}
	// 结构
	if len(d.structs) > 0 {
		for _, item := range d.structs {
			item.Sort()
		}
		sort.Slice(d.structs, func(i, j int) bool {
			return strings.Compare(d.structs[i].Name, d.structs[j].Name) <= 0
		})
		buf.Truncate(pos)
		structTpl, err := template.New("struct").Parse(structTempl)
		if err != nil {
			return err
		}
		if err := structTpl.Execute(buf, d.structs); err != nil {
			return err
		}
		if err := util.Save(dst, "struct.gen.proto", buf.Bytes()); err != nil {
			return err
		}
	}
	// 配置
	if len(d.configs) > 0 {
		for _, item := range d.configs {
			item.Sort()
		}
		sort.Slice(d.configs, func(i, j int) bool {
			return strings.Compare(d.configs[i].Name, d.configs[j].Name) <= 0
		})
		buf.Truncate(pos)
		importTpl, err := template.New("import").Parse(importTempl)
		if err != nil {
			return err
		}
		if err := importTpl.Execute(buf, d.get()); err != nil {
			return err
		}
		configTpl, err := template.New("config").Parse(configTempl)
		if err != nil {
			return err
		}
		if err := configTpl.Execute(buf, d.configs); err != nil {
			return err
		}
		if err := util.Save(dst, "table.gen.proto", buf.Bytes()); err != nil {
			return err
		}
	}
	return nil
}
