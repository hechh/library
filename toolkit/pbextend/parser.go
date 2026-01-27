package pbextend

import (
	"bytes"
	"go/ast"
	"strings"
	"text/template"

	"github.com/hechh/library/util"
)

type StructDescriptor struct {
	Name string
	TypeDescriptor
}

type EnumDescriptor struct {
	Name string
	TypeDescriptor
}

type Parser struct {
	pkgName string
	list    []*StructDescriptor
	enums   []*EnumDescriptor
}

func (p *Parser) Visit(n ast.Node) ast.Visitor {
	switch vv := n.(type) {
	case *ast.File:
		p.pkgName = vv.Name.Name
		return p
	case *ast.GenDecl:
		return p
	case *ast.TypeSpec:
		switch vv.Type.(type) {
		case *ast.StructType:
			item := ParseType(vv.Type)
			p.list = append(p.list, &StructDescriptor{
				Name:           vv.Name.Name,
				TypeDescriptor: item,
			})
		case *ast.Ident:
			item := ParseType(vv.Type)
			p.enums = append(p.enums, &EnumDescriptor{
				Name:           vv.Name.Name,
				TypeDescriptor: item,
			})
		}
		return nil
	}
	return nil
}

func (p *Parser) GetPkgName() string {
	return p.pkgName
}

func (p *Parser) GetAllEnum() []*EnumDescriptor {
	return p.enums
}

func (p *Parser) GetAllStruct() (rets []*StructDescriptor) {
	for _, item := range p.list {
		if strings.HasSuffix(item.Name, "Rsp") ||
			strings.HasSuffix(item.Name, "Req") ||
			strings.HasSuffix(item.Name, "Config") ||
			strings.HasSuffix(item.Name, "ConfigAry") {
			continue
		}
		rets = append(rets, item)
	}
	return
}

func (p *Parser) GetAllRsp() (rets []*StructDescriptor) {
	for _, item := range p.list {
		if strings.HasSuffix(item.Name, "Rsp") {
			rets = append(rets, item)
		}
	}
	return
}

func (p *Parser) Gen(dst string) error {
	tplObj := template.Must(template.New("pb").Parse(templ))
	buf := bytes.NewBuffer(nil)
	tplObj.Execute(buf, p)
	return util.SaveGo(dst, "common.gen.pb.go", buf.Bytes())
}
