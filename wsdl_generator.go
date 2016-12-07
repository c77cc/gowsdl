package gowsdl

import (
    "fmt"
    "strings"
)

type WsdlGenerator struct {
    soapAddr  string
    inputSpec *InputSpec

    elements  []string
    complexs  []string
    messages  []string
    porttypes []string
    bindings  []string
    service     string

    servicePrefix string
    namespace     string

    serviceName string
    bindName string
    portName string
}

type InputSpec struct {
    FuncSpecs         []*InputFuncSpec
    ExceptionName     string
    ExceptionRetDatas []*InputParam
}

type InputParam struct {
    Name string
    Type string
}

type InputFuncSpec struct {
    Name     string
    Params   []*InputParam
    RetDatas []*InputParam
}

func NewWsdlGenerator(addr, prefix, namespace string, spec *InputSpec) *WsdlGenerator {
    wg := &WsdlGenerator{}
    wg.inputSpec = spec
    wg.soapAddr  = addr
    wg.namespace = namespace
    wg.servicePrefix = prefix
    wg.serviceName = wg.servicePrefix + "_soap_service"
    wg.bindName = wg.servicePrefix + "_soap_binding"
    wg.portName = wg.servicePrefix + "_soap_port"
    return wg
}

func NewInputSpecViaSoapComments(commentGroups []SoapCommentGroupInterface, ename string, edata []*InputParam) (spec *InputSpec, err error) {
    spec = &InputSpec{
        ExceptionName: ename,
        ExceptionRetDatas: edata,
    }

    for _, commentGroup := range commentGroups {
        funcSpec := &InputFuncSpec{Name: commentGroup.GetMethod().GetValName()}
        for _, comment := range commentGroup.GetParams() {
            if comment.GetType() == SOAP_COMMENT_TYPE_PARAM {
                param := &InputParam{Name: comment.GetValName(), Type: comment.GetValType()}
                funcSpec.Params = append(funcSpec.Params, param)
            }
        }
        for _, comment := range commentGroup.GetRet() {
            if comment.GetType() == SOAP_COMMENT_TYPE_RET {
                data := &InputParam{Name: comment.GetValName(), Type: comment.GetValType()}
                funcSpec.RetDatas = append(funcSpec.RetDatas, data)
            }
        }
        spec.FuncSpecs = append(spec.FuncSpecs, funcSpec)
    }
    return
}

func (wg *WsdlGenerator) ToVar() (ret string) {
    for _, funcSpec := range wg.inputSpec.FuncSpecs {
        name  := funcSpec.Name
        rname := wg.getFuncRespName(name)
        wg.genElement(name)
        wg.genElement(rname)

        wg.genComplexType(name, funcSpec.Params)
        wg.genComplexType(rname, funcSpec.RetDatas)

        wg.genMessage(name)
        wg.genMessage(rname)

        wg.genPortType(name)
        wg.genBinding(name)
    }

    wg.genService()

    estr := wg.joinElements()
    cstr := wg.joinComplexTypes()
    mstr := wg.joinMessages()
    pstr := wg.joinPortTypes()
    bstr := wg.joinBindings()

    ret  = fmt.Sprintf(
        `<?xml version="1.0" encoding="utf-8"?>

<wsdl:definitions xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:tns="%s" xmlns:soap="http://schemas.xmlsoap.org/wsdl/soap/" xmlns:ns1="http://schemas.xmlsoap.org/soap/http" name="%s" targetNamespace="%s" xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:s0="http://tempuri.org/encodedTypes">

  <wsdl:types> 
    <xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema" attributeFormDefault="qualified" elementFormDefault="qualified" targetNamespace="%s">  

    %s
    %s

    </xs:schema> 
  </wsdl:types>  

  %s

  %s

  %s

  %s

</wsdl:definitions>
        `,
        wg.namespace, wg.serviceName, wg.namespace, wg.namespace, estr, cstr, mstr, pstr, bstr, wg.service,
    )
    return
}

func (wg *WsdlGenerator) joinElements() string {
    wg.genElement(wg.inputSpec.ExceptionName)
    return strings.Join(wg.elements, "\n")
}

func (wg *WsdlGenerator) joinComplexTypes() string {
    wg.genComplexType(wg.inputSpec.ExceptionName, wg.inputSpec.ExceptionRetDatas)
    return strings.Join(wg.complexs, "\n")
}

func (wg *WsdlGenerator) joinMessages() string {
    wg.genMessage(wg.inputSpec.ExceptionName, wg.inputSpec.ExceptionName)
    return strings.Join(wg.messages, "\n")
}

func (wg *WsdlGenerator) joinPortTypes() string {
    str := strings.Join(wg.porttypes, "\n")
    return fmt.Sprintf(
        `
          <wsdl:portType name="%s"> 
            %s
          </wsdl:portType>  
        `, wg.portName, str,
    )
}

func (wg *WsdlGenerator) joinBindings() string {
    str := strings.Join(wg.bindings, "\n")
    return fmt.Sprintf(
        `
          <wsdl:binding name="%s" type="tns:%s"> 
            <soap:binding style="document" transport="http://schemas.xmlsoap.org/soap/http"/>  
            %s
          </wsdl:binding>  
        `, wg.bindName, wg.portName, str,
    )
}

func (wg *WsdlGenerator) genElement(elename string) {
    wg.elements = append(wg.elements,
        fmt.Sprintf(`<xs:element name="%s" type="tns:%s"/>`, elename, elename),
    )
}

func (wg *WsdlGenerator) genComplexType(name string, params []*InputParam) {
    var seqs []string
    for _, p := range params {
        seqs = append(seqs, fmt.Sprintf(`<xs:element minOccurs="0" name="%s" type="xs:%s"/>`, p.Name, p.Type))
    }
    wg.complexs = append(wg.complexs, fmt.Sprintf(
        `
          <xs:complexType name="%s"> 
            <xs:sequence> 
                %s
            </xs:sequence> 
          </xs:complexType>  
        `,
        name, strings.Join(seqs, "\n"),
    ))
}

func (wg *WsdlGenerator) genMessage(name string, paramname ...string) {
    pname := "parameters"
    if len(paramname) > 0 {
        pname = paramname[0]
    }
    wg.messages = append(wg.messages, fmt.Sprintf(
        `
          <wsdl:message name="%s"> 
            <wsdl:part element="tns:%s" name="%s"></wsdl:part> 
          </wsdl:message>  
        `,
        name, name, pname,
    ))
}

func (wg *WsdlGenerator) genPortType(name string) {
    rname := wg.getFuncRespName(name)
    ename := wg.inputSpec.ExceptionName
    wg.porttypes = append(wg.porttypes, fmt.Sprintf(
        `
            <wsdl:operation name="%s">
              <wsdl:input message="tns:%s" name="%s"/>
              <wsdl:output message="tns:%s" name="%s"/>
              <wsdl:fault message="tns:%s" name="%s"/>
            </wsdl:operation>
        `,
        name, name, name,
        rname, rname,
        ename, ename,
    ))
}

func (wg *WsdlGenerator) genBinding(name string) {
    rname := wg.getFuncRespName(name)
    ename := wg.inputSpec.ExceptionName
    wg.bindings = append(wg.bindings, fmt.Sprintf(
        `
            <wsdl:operation name="%s"> 
              <soap:operation soapAction="%s" style="document"/>  
              <wsdl:input name="%s"> 
                <soap:body use="literal"/> 
              </wsdl:input>  
              <wsdl:output name="%s"> 
                <soap:body use="literal"/> 
              </wsdl:output>  
              <wsdl:fault name="%s"> 
                <soap:fault name="%s" use="literal"/> 
              </wsdl:fault> 
            </wsdl:operation> 
        `,
        name, name, name,
        rname,
        ename, ename,
    ))
}

func (wg *WsdlGenerator) genService() {
    wg.service = fmt.Sprintf(`
          <wsdl:service name="%s"> 
            <wsdl:port binding="tns:%s" name="%s"> 
                <soap:address location="%s"/> 
            </wsdl:port> 
          </wsdl:service> 
    `, wg.serviceName, wg.bindName, wg.bindName, wg.soapAddr)
}

func (wg *WsdlGenerator) getFuncRespName(name string) string {
    return name + "_response"
}
