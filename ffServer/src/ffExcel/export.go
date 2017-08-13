package ffExcel

import (
	"ffCommon/log/log"
	"ffCommon/util"

	"fmt"
	"os"
	"path"
	"path/filepath"
)

// ExcelExportLimit excel配置表导出的额外配置
type ExcelExportLimit struct {
	// Excel 针对的配置表名称
	Excel string
	// Sheet 针对配置表内哪个工作簿
	Sheet string
	// ExportLines 工作簿内哪些列导出(在配置表内已导出的基础上, 再进行此判定, 不在此列表内的, 则修正为不导出)
	ExportLines []string
	// ExportLinesRenameFrom 哪些列要重命名
	ExportLinesRenameFrom []string
	// ExportLinesRenameTo 新的名称
	ExportLinesRenameTo []string
}

// ExportConfig 导出配置
//	服务端代码语言: go, 配置文件: toml
//	客户端最终代码语言: c#, 配置文件: protobuf字节流(客户端使用临时语言go和配置文件toml, 最终生成最终所需的代码文件和配置文件)
type ExportConfig struct {
	// ServerExportGoCodePath 服务端导出的代码文件, 相对导出程序的路径, 为空时或者系统环境变量中未配置GOPATH时, 不导出
	ServerExportGoCodePath string
	// ServerExportTomlDataPath 服务端导出的toml配置文件, 相对导出程序的路径
	ServerExportTomlDataPath string
	// ServerReadTomlDataPath 服务端导出的toml配置文件, 相对读取程序的路径
	ServerReadTomlDataPath string

	// ClientExportGoCodePath 客户端导出的go代码文件, 相对导出程序的路径, 为空时, 不执不导出
	ClientExportGoCodePath string
	// ClientExportCSharpCodePath 客户端最终导出的c#代码文件, 相对导出程序的路径, 为空时, 不导出
	ClientExportCSharpCodePath string
	// ClientExportProtoBufDataPath 客户端导出的Protobuf配置文件, 相对导出程序的路径
	ClientExportProtoBufDataPath string

	// ExcelExportLimit 工作表导出时的额外配置
	ExcelExportLimit []*ExcelExportLimit

	hasGoEnv          bool   // 是否有go环境
	serverPackageName string // 根据ServerExportGoCodePath推导出来的包名
	clientPackageName string // 根据ClientExportGoCodePath推导出来的包名
}

func (ec *ExportConfig) check() error {
	ec.hasGoEnv = os.Getenv("GOPATH") != ""

	_, ec.serverPackageName = path.Split(ec.ServerExportGoCodePath)
	_, ec.clientPackageName = path.Split(ec.ClientExportGoCodePath)

	return nil
}

func (ec *ExportConfig) String() string {
	return fmt.Sprintf(`[
	ServerExportGoCodePath:%v
	ServerExportTomlDataPath:%v
	ServerReadTomlDataPath:%v
	ClientExportGoCodePath:%v
	ClientExportCSharpCodePath:%v
	ClientExportProtoBufDataPath:%v
]`,

		ec.ServerExportGoCodePath,
		ec.ServerExportTomlDataPath,
		ec.ServerReadTomlDataPath,
		ec.ClientExportGoCodePath,
		ec.ClientExportCSharpCodePath,
		ec.ClientExportProtoBufDataPath)
}

var exportConfig *ExportConfig

func clearPath(ec *ExportConfig) bool {
	result := true

	err := util.ClearPath("toml", true, nil)
	if err != nil {
		log.RunLogger.Println(err)
		result = false
	}

	if ec.ServerExportGoCodePath != "" {
		err := util.ClearPath(ec.ServerExportGoCodePath, false, []string{".go"})
		if err != nil {
			log.RunLogger.Println(err)
			result = false
		}
	}

	if ec.ServerExportTomlDataPath != "" {
		err := util.ClearPath(ec.ServerExportTomlDataPath, false, []string{".toml"})
		if err != nil {
			log.RunLogger.Println(err)
			result = false
		}
	}

	if ec.ClientExportGoCodePath != "" {
		err := util.ClearPath(ec.ClientExportGoCodePath, false, []string{".go"})
		if err != nil {
			log.RunLogger.Println(err)
			result = false
		}
	}

	if ec.ClientExportCSharpCodePath != "" {
		err := util.ClearPath(ec.ClientExportCSharpCodePath, false, []string{".cs"})
		if err != nil {
			log.RunLogger.Println(err)
			result = false
		}
	}

	if ec.ClientExportProtoBufDataPath != "" {
		err := util.ClearPath(ec.ClientExportProtoBufDataPath, false, []string{".bytes"})
		if err != nil {
			log.RunLogger.Println(err)
			result = false
		}
	}

	return result
}

// exportExcel 解析excel, 然后根据导出配置, 将解析结果保存
func exportExcel(excelFilePath string) (excel *excel, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("exportExcel excel[%v] get error:\n[%v]", excelFilePath, err)
		}
	}()

	// 解析Excel
	excel, err = parseExcel(excelFilePath)
	if err != nil {
		return nil, err
	}

	// 检查Excel是否满足导出为toml格式
	if err = checkToml(excel); err != nil {
		return nil, err
	}

	// 导出服务端
	if excel.exportToServer() {
		// 导出读取toml数据的Go代码
		tomlDataServerReadCode := genTomlDataReadCode(excel, exportConfig, "server")
		if exportConfig.hasGoEnv && exportConfig.ServerExportGoCodePath != "" {
			defFilePath := path.Join(exportConfig.ServerExportGoCodePath, excel.name+".go")
			err = util.WriteFile(defFilePath, []byte(tomlDataServerReadCode))
			if err != nil {
				return nil, err
			}
			log.RunLogger.Println(defFilePath)
		}

		// 导出toml数据
		tomlDataServer := genTomlData(excel, exportConfig, "server")
		dataFilePath := path.Join("toml", "server", excel.name+".toml")
		err = util.WriteFile(dataFilePath, []byte(tomlDataServer))
		if err != nil {
			return nil, err
		}

		if exportConfig.ServerExportTomlDataPath != "" {
			dataFilePath = path.Join(exportConfig.ServerExportTomlDataPath, excel.name+".toml")
			err = util.WriteFile(dataFilePath, []byte(tomlDataServer))
			if err != nil {
				return nil, err
			}
			log.RunLogger.Println(dataFilePath)
		}
	}

	// 导出客户端
	if excel.exportToClient() {
		// 导出读取toml数据的Go代码
		tomlDataGoReadCode := genTomlDataReadCode(excel, exportConfig, "client")
		if exportConfig.hasGoEnv && exportConfig.ClientExportGoCodePath != "" {
			defFilePath := path.Join(exportConfig.ClientExportGoCodePath, excel.name+".go")
			err = util.WriteFile(defFilePath, []byte(tomlDataGoReadCode))
			if err != nil {
				return nil, err
			}
			log.RunLogger.Println(defFilePath)
		}

		// 导出toml数据
		tomlDataServer := genTomlData(excel, exportConfig, "client")
		dataFilePath := path.Join("toml", "client", excel.name+".toml")
		err = util.WriteFile(dataFilePath, []byte(tomlDataServer))
		if err != nil {
			return nil, err
		}
	}

	return
}

// ExportExcelDir 解析指定目录内的所有excel, 然后根据导出配置, 将解析结果保存
func ExportExcelDir(excelDirPath string, _exportConfig *ExportConfig) error {
	exportConfig = _exportConfig

	// 配置检查
	err := exportConfig.check()
	if err != nil {
		return err
	}

	//  清空导出目录
	if !clearPath(exportConfig) {
		return fmt.Errorf("ExportExcelDir excelDirPath[%v] exportConfig[%v] clearPath failed", excelDirPath, exportConfig)
	}

	log.RunLogger.Println(exportConfig)

	// 遍历获得所有excel
	excelFilePaths := make([]string, 0, 1)
	filepath.Walk(excelDirPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		if filepath.Ext(f.Name()) != ".xlsx" {
			log.RunLogger.Printf("invalid excel[%v] extension: only support .xlsx\n", f.Name())
			return nil
		}

		excelFilePaths = append(excelFilePaths, path)

		return nil
	})

	// 依次导出所有excel
	allExcels := make([]*excel, 0, len(excelFilePaths))
	allValid := true
	for _, excelFilePath := range excelFilePaths {
		excel, err := exportExcel(excelFilePath)
		if err != nil {
			log.RunLogger.Println(err)
			allValid = false
		} else {
			allExcels = append(allExcels, excel)
		}
	}

	// 导出toml数据对应的Proto定义
	if exportConfig.ClientExportProtoBufDataPath != "" {
		goProto, csharpProto := genProtoDefineFromToml(allExcels, "client")
		goFilePath := path.Join("ProtoBuf", "Server", "Config.proto")
		err = util.WriteFile(goFilePath, []byte(goProto))
		if err != nil {
			return err
		}

		csharpFilePath := path.Join("ProtoBuf", "Client", "Config.proto")
		err = util.WriteFile(csharpFilePath, []byte(csharpProto))
		if err != nil {
			return err
		}
	}

	// 有错发生, 提示用户
	if !allValid {
		return fmt.Errorf("ExportExcelDir dir[%v] not all excel export success", excelDirPath)
	}

	return nil
}
