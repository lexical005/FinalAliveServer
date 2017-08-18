package ffExcel

import (
	"ffCommon/log/log"
	"ffCommon/util"
	"fmt"
	"path"
	"strings"

	"github.com/lexical005/toml"
)

type errReasonToml struct {
	Error []struct {
		ErrCode string
		Desc    string
	}
}

var fmtErrorCSharpWhole = `namespace NConfig
{
    public static class Error
    {
        public static readonly string[] Keys;

        public static string Desc(int errCode)
        {
            if (errCode >= 0 && errCode < Keys.Length)
            {
                return LanguageReader.Error[Keys[errCode]].Value;
            }
            return "ErrCode" + errCode.ToString();
        }

        static Error()
        {
            Keys = new string[]
            {{ErrorKey}
            };
        }
    }
}
`

var fmtErrorCSharpKey = "\n                \"%v\","

// 错误码, 客户端
func genErrorCSharp(dataFilePath string, errReasonToml *errReasonToml) {
	keys := ""

	for _, one := range errReasonToml.Error {
		keys += fmt.Sprintf(fmtErrorCSharpKey, one.ErrCode)
	}

	result := strings.Replace(fmtErrorCSharpWhole, "{ErrorKey}", keys, -1)

	util.WriteFile(dataFilePath, []byte(result))
	log.RunLogger.Println(dataFilePath)
}

var fmtErrorGoPackage = `package ffError

`

var fmtErrorGoReasonHeader = `
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

var fmtErrorGoReasonLoopComment = `// %v %v
`
var fmtErrorGoReasonLoop = `var %v Error = &errReason{code: %v, desc: "%v"}
`

var fmtErrorGoCodeStart = `
var errByCode = []Error{
`

var fmtErrorGoCodeLoop = `	%v,
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

// 错误码, 服务端
func genErrorGo(dataFilePath string, errReasonToml *errReasonToml) {
	result := ""

	result += fmtErrorGoPackage

	result += fmtErrorGoReasonHeader

	for index, one := range errReasonToml.Error {
		result += fmt.Sprintf(fmtErrorGoReasonLoopComment, one.ErrCode, one.Desc)
		result += fmt.Sprintf(fmtErrorGoReasonLoop, one.ErrCode, index, one.Desc)
	}

	result += fmtErrorGoCodeStart
	for _, one := range errReasonToml.Error {
		result += fmt.Sprintf(fmtErrorGoCodeLoop, one.ErrCode)
	}
	result += fmtErrorCodeEnd

	result += strings.Replace(fmtErrorByCodeFunc, "ffError", exportConfig.serverPackageName, -1)

	util.WriteFile(dataFilePath, []byte(result))
	log.RunLogger.Println(dataFilePath)
}

func genError(excel *excel) {
	// Server
	if excel.exportServerGoCodePath != "" && exportConfig.hasGoEnv {
		tomlDataServer, err := util.ReadFile(path.Join("toml", "server", "Error.toml"))
		if err != nil {
			log.RunLogger.Println(err)
			return
		}

		// 解析导出的toml配置文件
		errReasonToml := &errReasonToml{}
		err = toml.Unmarshal(tomlDataServer, errReasonToml)
		if err != nil {
			log.RunLogger.Println(err)
			return
		}

		genErrorGo(path.Join(excel.exportServerGoCodePath, "Error.go"), errReasonToml)
	}

	// Client
	if excel.exportClientCSharpCodePath != "" {
		tomlDataClient, err := util.ReadFile(path.Join("toml", "client", "Error.toml"))
		if err != nil {
			log.RunLogger.Println(err)
			return
		}

		// 解析导出的toml配置文件
		errReasonToml := &errReasonToml{}
		err = toml.Unmarshal(tomlDataClient, errReasonToml)
		if err != nil {
			log.RunLogger.Println(err)
			return
		}

		genErrorCSharp(path.Join(excel.exportClientCSharpCodePath, "Error.cs"), errReasonToml)
	}
}
