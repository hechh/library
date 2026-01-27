package xlsx2data

import (
	"os"
	"strings"

	"github.com/hechh/library/convertor"
	"github.com/hechh/library/uerror"
	"github.com/hechh/library/util"
	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

type MsgParser struct {
	pkg     string
	files   *protoregistry.Files
	data    map[string]*EnumDescriptor
	enums   []*EnumDescriptor
	configs []*StructDescriptor
}

func NewMsgParser(pkg string, filename string) (*MsgParser, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, uerror.Err(-1, " 加载文件(%s)失败: %v", filename, err)
	}

	fds := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(data, fds); err != nil {
		return nil, err
	}

	d := &MsgParser{data: make(map[string]*EnumDescriptor)}
	if files, err := protodesc.NewFiles(fds); err != nil {
		return nil, err
	} else {
		d.pkg = pkg
		d.files = files
	}
	return d, nil
}

func (d *MsgParser) GetFullName(name string) protoreflect.FullName {
	return protoreflect.FullName(d.pkg + "." + name)
}

func (d *MsgParser) GetMessageType(name string) (protoreflect.MessageDescriptor, error) {
	msgType, err := d.files.FindDescriptorByName(d.GetFullName(name))
	if err == protoregistry.NotFound {
		err = nil
	}
	return msgType.(protoreflect.MessageDescriptor), err
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
				if err := d.parseStruct(names[1], rows); err != nil {
					return err
				}
			case "@config:col":
				names := strings.Split(strs[1], ":")
				rows, err := fp.GetCols(names[0])
				if err != nil {
					return uerror.Err(-1, "表格(%s)不存在", strs[1])
				}
				if err := d.parseStruct(names[1], rows); err != nil {
					return err
				}
			case "@enum":
				rows, err := fp.GetRows(strs[1])
				if err != nil {
					return uerror.Err(-1, "表格(%s)不存在", strs[1])
				}
				d.parseEnum(rows)
			}
		}
	}
	return nil
}

func (d *MsgParser) parseStruct(name string, rows [][]string) error {
	aryType, err := d.GetMessageType(name + "Ary")
	if err != nil {
		return err
	}
	cfgType, err := d.GetMessageType(name)
	if err != nil {
		return err
	}
	st := NewStructDescriptor(name, aryType, cfgType, rows[3:])
	for i, item := range rows[1] {
		if len(item) <= 0 {
			continue
		}
		st.Put(int32(i)+1, rows[0][i], item)
	}
	d.configs = append(d.configs, st)
	return nil
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
				convertor.Register(func(val string) any { return enum.ToInt32(val) }, "int32", strs[2])
				d.enums = append(d.enums, enum)
				d.data[strs[2]] = enum
			}
			enum.Put(cast.ToInt32(strs[4]), strs[3], strs[2], strs[1])
		}
	}
}

func (d *MsgParser) Gen(dst string) error {
	for _, st := range d.configs {
		buf, err := st.Marshal()
		if err != nil {
			return err
		}
		if err := util.Save(dst, st.name+".conf", buf); err != nil {
			return err
		}
	}
	return nil
}
