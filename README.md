# Parse golang comment and generate wsdl file

[![Build Status](https://drone.io/github.com/c77cc/gowsdl/status.png)](https://drone.io/github.com/c77cc/gowsdl/latest)

`gowsdl` Generate wsdl file via parse your project golang comment.


## Example
-----------------
```golang
package main

import (
    "log"
    "github.com/c77cc/gowsdl"
    "io/ioutil"
)

func main() {
    // wsdl address
    addr      := "http://127.0.0.1:8000"
    // wsdl service prefix
    prefix    := "service_prefix"
    // wsdl namespace
    namespace := "http://ws.domain.com/"
    // you project code path
    destpath := "./target"

    // also, you can implement your own parser
    group, err := gowsdl.ParseComments(destpath)
    if err != nil {
        log.Fatal(err)
    }

    ename := "exception"
    edata := []*gowsdl.InputParam{
        &gowsdl.InputParam{Name: "return", Type: "string"},
    }
    spec, err := gowsdl.NewInputSpecViaSoapComments(group, ename, edata)
    if err != nil {
        log.Fatal(err)
    }
    wg := gowsdl.NewWsdlGenerator(addr, prefix, namespace, spec)
    xml := wg.ToVar()
    if len(xml) < 1 {
        log.Fatal("fail")
    }
    log.Println("write to ws.wsdl...")
    ioutil.WriteFile("ws.wsdl", []byte(xml), 0777)
}
```
