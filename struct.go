package gowsdl

const (
    SOAP_COMMENT_TYPE_METHOD = iota
    SOAP_COMMENT_TYPE_PARAM
    SOAP_COMMENT_TYPE_RET
)

type SoapCommentGroupInterface interface {
    GetFuncName() string
    GetMethod()   SoapCommentInterface
    GetParams()   []SoapCommentInterface
    GetRet()      []SoapCommentInterface
}

type SoapCommentInterface interface {
    GetType() int
    GetValType() string
    GetValName() string
}
