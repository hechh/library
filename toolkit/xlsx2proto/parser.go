package xlsx2proto

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/hechh/library/uerror"
	"github.com/hechh/library/util"
	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
)

type MsgParser struct {
	pkgname string
	option  string
	dst     string
	sts     map[string]*StructDescriptor
	data    map[string]*EnumDescriptor
	enums   []*EnumDescriptor
	structs []*StructDescriptor
	configs []*ConfigDescriptor
}

func NewMsgParser(p string, o string, dst string) *MsgParser {
	return &MsgParser{
		pkgname: p,
		option:  o,
		dst:     dst,
		sts:     make(map[string]*StructDescriptor),
		data:    make(map[string]*EnumDescriptor),
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
	st := NewConfigDescriptor(name)
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
				d.sts[st.name] = st
				d.structs = append(d.structs, st)
			}
			st.Put(int32(len(st.list)+1), strs[2], strs[3], strs[4])
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

func (d *MsgParser) Gen() error {
	buf := bytes.NewBuffer(nil)

	pos, _ := buf.WriteString(fmt.Sprintf(`
syntax = "proto3";

package %s;

option  go_package = "%s";
	
`, d.pkgname, d.option))

	if len(d.enums) > 0 {
		sort.Slice(d.enums, func(i, j int) bool {
			return strings.Compare(d.enums[i].name, d.enums[j].name) <= 0
		})
		for _, item := range d.enums {
			buf.WriteString(item.String())
		}
		if err := util.Save(d.dst, "enum.gen.proto", buf.Bytes()); err != nil {
			return err
		}
	}

	if len(d.structs) > 0 {
		sort.Slice(d.structs, func(i, j int) bool {
			return strings.Compare(d.structs[i].name, d.structs[j].name) <= 0
		})
		buf.Truncate(pos)
		if d.hasEnum() {
			buf.WriteString("import \"enum.gen.proto\";\n\n")
		}
		for _, item := range d.structs {
			buf.WriteString(item.String())
		}
		if err := util.Save(d.dst, "struct.gen.proto", buf.Bytes()); err != nil {
			return err
		}
	}

	if len(d.configs) > 0 {
		sort.Slice(d.configs, func(i, j int) bool {
			return strings.Compare(d.configs[i].name, d.configs[j].name) <= 0
		})
		buf.Truncate(pos)
		if d.hasEnum() {
			buf.WriteString("import \"enum.gen.proto\";\n\n")
		}
		if d.hasStruct() {
			buf.WriteString("import \"struct.gen.proto\";\n\n")
		}
		for _, item := range d.configs {
			buf.WriteString(item.String())
		}
		return util.Save(d.dst, "table.gen.proto", buf.Bytes())
	}
	return nil
}

func (d *MsgParser) hasEnum() bool {
	for _, item := range d.configs {
		for _, field := range item.list {
			if _, ok := d.data[field.class]; ok {
				return true
			}
		}
	}
	return false
}

func (d *MsgParser) hasStruct() bool {
	for _, item := range d.configs {
		for _, field := range item.list {
			if _, ok := d.sts[field.class]; ok {
				return true
			}
		}
	}
	return false
}
