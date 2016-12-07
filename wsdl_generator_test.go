package gowsdl

import (
    "fmt"
	"testing"
    "io/ioutil"
)

func TestToVar(t *testing.T) {
    soapAddr := "http://127.0.0.1:8000"
    prefix   := "service_prefix"
    namespace := "http://ws.domain.com/"
//    spec     := &InputSpec{
//        FuncSpecs: []*InputFuncSpec{
//            &InputFuncSpec{
//                Name: "oauth_token",
//                Params: []*InputParam{
//                    &InputParam{Name: "client_id", Type: "string"},
//                    &InputParam{Name: "grant_type", Type: "string"},
//                    &InputParam{Name: "username", Type: "string"},
//                    &InputParam{Name: "password", Type: "string"},
//                },
//                RetDatas: []*InputParam{
//                    &InputParam{Name: "return", Type: "string"},
//                },
//            },
//            &InputFuncSpec{
//                Name: "get_name",
//                Params: []*InputParam{
//                    &InputParam{Name: "access_token", Type: "string"},
//                    &InputParam{Name: "id", Type: "int"},
//                    &InputParam{Name: "extra", Type: "string"},
//                },
//                RetDatas: []*InputParam{
//                    &InputParam{Name: "return", Type: "string"},
//                },
//            },
//        },
//        ExceptionName: "exception",
//        ExceptionRetDatas: []*InputParam{
//            &InputParam{Name: "return", Type: "string"},
//        },
//    }

    group, err := ParseComments("./target")
    if err != nil {
        t.Fatal(err)
    }

    ename := "exception"
    edata := []*InputParam{
        &InputParam{Name: "return", Type: "string"},
    }
    spec, err := NewInputSpecViaSoapComments(group, ename, edata)
    if err != nil {
        t.Fatal(err)
    }
    wg := NewWsdlGenerator(soapAddr, prefix, namespace, spec)
    xml := wg.ToVar()
    if len(xml) < 1 {
        t.Fatal("fail")
    }
    fmt.Println("write to ws.wsdl...")
    ioutil.WriteFile("ws.wsdl", []byte(xml), 0777)
}
