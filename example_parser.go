package gowsdl

import (
	"fmt"
    "sort"
	"go/parser"
	"go/token"
	"go/ast"
    "strings"
)

const (
    SOAP_METHOD_FLAG = "@SoapMethod"
    SOAP_PARAM_FLAG  = "@SoapParam"
    SOAP_RET_FLAG    = "@SoapReturn"
)

type ExampleSoapCommentGroup struct {
    FuncName string
    Method   SoapCommentInterface
    Ret      []SoapCommentInterface
    Params   []SoapCommentInterface
}

type ExampleSoapComment struct {
    Type     int
    ValType  string
    ValName  string
}

func (s ExampleSoapCommentGroup) GetFuncName() string {
    return s.FuncName
}

func (s ExampleSoapCommentGroup) GetMethod() SoapCommentInterface {
    return s.Method
}

func (s ExampleSoapCommentGroup) GetParams() []SoapCommentInterface {
    return s.Params
}

func (s ExampleSoapCommentGroup) GetRet() []SoapCommentInterface {
    return s.Ret
}

func (s ExampleSoapComment) GetType() int {
    return s.Type
}

func (s ExampleSoapComment) GetValType() string {
    return s.ValType
}

func (s ExampleSoapComment) GetValName() string {
    return s.ValName
}

type ExampleSoapCommentGroupSlice []ExampleSoapCommentGroup

func (s ExampleSoapCommentGroupSlice) Len() int {
    return len(s)
}

func (s ExampleSoapCommentGroupSlice) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

func (s ExampleSoapCommentGroupSlice) Less(i, j int) bool {
    return s[i].Method.GetValName() < s[j].Method.GetValName()
}

func ParseComments(dir string) (commentGroups []SoapCommentGroupInterface, err error) {
    var pkgs map[string]*ast.Package
	fset := token.NewFileSet()
	pkgs, err = parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		return
	}

    var tmps ExampleSoapCommentGroupSlice
    for _, pkg := range pkgs {
        for _, file := range pkg.Files {
            for _, o := range file.Scope.Objects {
                if o.Kind != ast.Fun {
                    continue
                }

                decl, ok := o.Decl.(*ast.FuncDecl)
                if !ok {
                    continue
                }

                if decl.Doc == nil {
                    continue
                }

                commentGroup := ExampleSoapCommentGroup{FuncName: o.Name}

                for _, t := range decl.Doc.List {
                    parseComment(&commentGroup, trimComment(t.Text))
                }

                if commentGroup.Method != nil && len(commentGroup.Ret) > 0 {
                    tmps = append(tmps, commentGroup)
                }
            }
        }
    }
    sort.Sort(tmps)
    for _, t := range tmps {
        commentGroups = append(commentGroups, t)
    }
    return
}

func trimComment(txt string) string {
    if len(txt) < 2 {
        return txt
    }

    switch txt[1] {
        case '/':
            //-style comment (no newline at the end)
            txt = txt[2:]
            // strip first space - required for Example tests
            if len(txt) > 0 && txt[0] == ' ' {
                txt = txt[1:]
            }
        case '*':
            /*-style comment */
            txt = txt[2 : len(txt)-2]
    }
    return txt
}

/*
  Parse Golang Func Comments like this:
      // @SoapMethod string v1_active_client
      // @SoapParam  string tbg
      // @SoapParam  string ted
      // @SoapReturn string return
*/
func parseComment(group *ExampleSoapCommentGroup, txt string) (err error) {
    comment := ExampleSoapComment{}
    ary := trimStrSlice(strings.Split(txt, " "))
    if len(ary) < 3 {
        err = fmt.Errorf(`syntax invalid "%s"`, txt)
        return
    }
    comment.ValName = strings.TrimSpace(ary[2])
    comment.ValType = strings.TrimSpace(ary[1])

    if strings.HasPrefix(txt, SOAP_METHOD_FLAG) {
        comment.Type = SOAP_COMMENT_TYPE_METHOD
        group.Method = comment
    } else if strings.HasPrefix(txt, SOAP_PARAM_FLAG) {
        comment.Type = SOAP_COMMENT_TYPE_PARAM
        group.Params = append(group.Params, comment)
    } else if strings.HasPrefix(txt, SOAP_RET_FLAG) {
        comment.Type = SOAP_COMMENT_TYPE_RET
        group.Ret = append(group.Ret, comment)
    } else {
        return
    }

    return
}

func trimStrSlice(ary []string) (nary []string) {
    for i, _ := range ary {
        if len(strings.TrimSpace(ary[i])) > 0 {
            nary = append(nary, strings.TrimSpace(ary[i]))
        }
    }
    return nary
}
