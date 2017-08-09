package main

import (
	"fmt"
)

var fmtPackage = `package ffError

`

var fmtErrReasonHeader = `
import (
	"ffCommon/log/log"

	"fmt"
)

// Error Error
type Error interface {
	Code() int32
	Error() string
	String() string
}

type errReason struct {
	code int32
	desc string
}

func (ec *errReason) Code() int32 {
	return ec.code
}

func (ec *errReason) Error() string {
	return fmt.Sprintf("ffError[%d-%s]", ec.code, ec.desc)
}

func (ec *errReason) String() string {
	return fmt.Sprintf("ffError[%d-%s]", ec.code, ec.desc)
}
`

var fmtErrorReasonLoopComment = `// %v %v
`
var fmtErrorReasonLoop = `var %v Error = &errReason{code: %v, desc: "%v"}
`

var fmtErrorCodeStart = `
var errByCode = []Error{
`

var fmtErrorCodeLoop = `	%v,
`

var fmtErrorCodeEnd = `
}
`

var fmtErrorByCodeFunc = `

// ErrByCode 根据错误码获取Error
func ErrByCode(errCode int32) Error {
	if errCode >= 0 && int(errCode) < len(errByCode) {
		return errByCode[errCode]
	}

	log.FatalLogger.Printf("ffError.ErrByCode: invalid errCode[%d]", errCode)

	return ErrUnknown
}
`

func tomlToGolang(errReasonToml *errReasonToml) string {
	result := ""

	result += fmtPackage

	result += fmtErrReasonHeader

	for index, one := range errReasonToml.Error {
		result += fmt.Sprintf(fmtErrorReasonLoopComment, one.ErrCode, one.ErrDescCN)
		result += fmt.Sprintf(fmtErrorReasonLoop, one.ErrCode, index, one.ErrDescCN)
	}

	result += fmtErrorCodeStart
	for _, one := range errReasonToml.Error {
		result += fmt.Sprintf(fmtErrorCodeLoop, one.ErrCode)
	}
	result += fmtErrorCodeEnd

	result += fmtErrorByCodeFunc

	return result
}
