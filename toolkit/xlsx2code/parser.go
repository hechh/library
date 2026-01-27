package xlsx2code

import (
	"bytes"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/hechh/library/uerror"
	"github.com/hechh/library/util"
	"github.com/iancoleman/strcase"
	"github.com/xuri/excelize/v2"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

type MsgParser struct {
	pkg   string
	files *protoregistry.Files
	data  map[string]*StructDescriptor
	list  []*StructDescriptor
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

	d := &MsgParser{data: make(map[string]*StructDescriptor)}
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
				if err := d.parseStruct(names[1], rows, strs[2:]...); err != nil {
					return err
				}
			case "@config:col":
				names := strings.Split(strs[1], ":")
				rows, err := fp.GetCols(names[0])
				if err != nil {
					return uerror.Err(-1, "表格(%s)不存在", strs[1])
				}
				if err := d.parseStruct(names[1], rows, strs[2:]...); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (d *MsgParser) parseStruct(name string, rows [][]string, rules ...string) error {
	cfgType, err := d.GetMessageType(name)
	if err != nil {
		return err
	}
	st := NewStructDescriptor(name, cfgType)
	for i, item := range rows[1] {
		if len(item) <= 0 {
			continue
		}
		st.Put(int32(i)+1, rows[0][i], item)
	}
	for _, rule := range rules {
		st.AddIndex(rule)
	}
	d.data[st.Name] = st
	d.list = append(d.list, st)
	return nil
}

func (d *MsgParser) Gen(dst string, tpl *template.Template) error {
	buf := bytes.NewBuffer(nil)
	for _, item := range d.list {
		pkgname := strcase.ToSnake(item.Name)
		if err := tpl.Execute(buf, item); err != nil {
			return err
		}

		if err := util.SaveGo(path.Join(dst, pkgname), item.Name+".gen.go", buf.Bytes()); err != nil {
			return err
		}
		buf.Reset()
	}
	return nil
}
