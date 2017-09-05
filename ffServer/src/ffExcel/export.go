package ffExcel

import (
	"cellvalue"
	"ffCommon/log/log"
	"ffCommon/util"
	"sort"
	"strings"

	"fmt"
	"os"
	"path"
	"path/filepath"
)

// ExcelExportType excel配置表导出方式配置, 默认config, 还支持error,enum
type ExcelExportType struct {
	// Excel 针对的配置表名称
	Excel string
	// Type 导出方式
	Type string

	// ServerExportGoCodePath 服务端导出的代码文件, 相对导出程序的路径, 为空时或者系统环境变量中未配置GOPATH时, 不导出
	ServerExportGoCodePath string
	// ClientExportCSharpCodePath 客户端最终导出的c#代码文件, 相对导出程序的路径, 为空时, 不导出
	ClientExportCSharpCodePath string
}

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
	// ServerReadTomlDataPath 服务端导出的toml配置文件, 相对读取程序的路径
	ServerReadTomlDataPath string

	// ClientExportGoCodePath 客户端导出的go代码文件, 相对导出程序的路径, 为空时, 不执不导出
	ClientExportGoCodePath string
	// ClientExportCSharpCodePath 客户端最终导出的c#代码文件, 相对导出程序的路径, 为空时, 不导出
	ClientExportCSharpCodePath string
	// ClientExportProtoBufDataPath 客户端导出的Protobuf配置文件, 相对导出程序的路径
	ClientExportProtoBufDataPath string

	// ExcelExportType 导出方式配置
	ExcelExportType []*ExcelExportType

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
	ServerReadTomlDataPath:%v
	ClientExportGoCodePath:%v
	ClientExportCSharpCodePath:%v
	ClientExportProtoBufDataPath:%v
]`,

		ec.ServerExportGoCodePath,
		ec.ServerReadTomlDataPath,
		ec.ClientExportGoCodePath,
		ec.ClientExportCSharpCodePath,
		ec.ClientExportProtoBufDataPath)
}

var exportConfig *ExportConfig

func clearPath(ec *ExportConfig) bool {
	result := true

	var err error
	err = util.ClearPath("toml", true, nil)
	if err != nil {
		log.RunLogger.Println(err)
		result = false
	}

	if ec.ServerExportGoCodePath != "" {
		util.Walk(ec.ServerExportGoCodePath, func(f os.FileInfo) error {
			// 忽略文件夹以及非go文件
			name := f.Name()
			if f.IsDir() || !strings.HasSuffix(name, ".go") {
				return nil
			}

			// 非配置文件
			if name[0] == strings.ToLower(name)[0] {
				return nil
			}

			// 删除生成的文件
			return os.Remove(filepath.Join(ec.ServerExportGoCodePath, f.Name()))
		})
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
func exportExcel(excel *excel) (err error) {
	// 导出服务端
	if excel.exportToServer() {
		// 导出读取toml数据的Go代码
		if excel.exportType == "config" {
			if exportConfig.hasGoEnv && excel.exportServerGoCodePath != "" {
				tomlDataServerReadCode := genTomlDataReadCode(excel, exportConfig, "server")
				defFilePath := path.Join(excel.exportServerGoCodePath, excel.name+".go")
				err = util.WriteFile(defFilePath, []byte(tomlDataServerReadCode))
				if err != nil {
					return err
				}
				log.RunLogger.Println(defFilePath)
			}
		}

		// 导出toml数据
		tomlDataServer := genTomlData(excel, exportConfig, "server")
		dataFilePath := path.Join("toml", "server", excel.name+".toml")
		err = util.WriteFile(dataFilePath, []byte(tomlDataServer))
		if err != nil {
			return err
		}
	}

	// 导出客户端
	if excel.exportToClient() {
		// 导出读取toml数据的Go代码
		if excel.exportType == "config" {
			if exportConfig.hasGoEnv && excel.exportClientGoCodePath != "" {
				tomlDataGoReadCode := genTomlDataReadCode(excel, exportConfig, "client")
				defFilePath := path.Join(excel.exportClientGoCodePath, excel.name+".go")
				err = util.WriteFile(defFilePath, []byte(tomlDataGoReadCode))
				if err != nil {
					return err
				}
				log.RunLogger.Println(defFilePath)
			}
		}

		// 导出toml数据
		tomlDataServer := genTomlData(excel, exportConfig, "client")
		dataFilePath := path.Join("toml", "client", excel.name+".toml")
		err = util.WriteFile(dataFilePath, []byte(tomlDataServer))
		if err != nil {
			return err
		}
	}

	// 错误码
	if excel.exportType == "error" {
		genError(excel)
	} else if excel.exportType == "enum" {
		genEnum(excel)
	}

	return
}

// ExportExcelDir 解析指定目录内的所有excel, 然后根据导出配置, 将解析结果保存
func ExportExcelDir(excelDirPath string, _exportConfig *ExportConfig) error {
	exportConfig = _exportConfig

	allValid := true

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

	// 遍历获得所有excel路径
	allExcelPath := make([]string, 0, 16)
	allExcels := make([]*excel, 0, 16)
	allConfigExcels := make([]*excel, 0, 16)
	err = filepath.Walk(excelDirPath, func(path string, f os.FileInfo, err error) error {
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

		log.RunLogger.Println("found excel:", f.Name())

		allExcelPath = append(allExcelPath, path)

		return nil
	})

	if err != nil {
		return err
	}
	sort.Strings(allExcelPath)

	// 优先解析 Error 和 Enum
	for i := len(allExcelPath) - 1; i >= 0; i = i - 1 {
		path := allExcelPath[i]
		if strings.HasSuffix(path, "Error.xlsx") || strings.HasSuffix(path, "Enum.xlsx") {

			// 解析Excel
			excel, err := parseExcel(path)
			if err != nil {
				return err
			}

			// 检查Excel是否满足导出为toml格式
			if err = checkToml(excel); err != nil {
				return err
			}

			// 导出类型
			excel.exportType = "config"
			excel.exportServerGoCodePath = exportConfig.ServerExportGoCodePath
			excel.exportClientGoCodePath = exportConfig.ClientExportGoCodePath
			excel.exportClientCSharpCodePath = exportConfig.ClientExportCSharpCodePath
			for _, exportTypeConfig := range exportConfig.ExcelExportType {
				if exportTypeConfig.Excel == excel.name {
					excel.exportType = exportTypeConfig.Type
					excel.exportServerGoCodePath = exportTypeConfig.ServerExportGoCodePath
					excel.exportClientCSharpCodePath = exportTypeConfig.ClientExportCSharpCodePath
					break
				}
			}

			err = exportExcel(excel)
			if err != nil {
				return err
			}

			allExcelPath = append(allExcelPath[:i], allExcelPath[i+1:]...)
		}
	}

	//
	cellvalue.InitEnum(path.Join("toml", "server", "Enum.toml"))

	// 配置表
	for _, path := range allExcelPath {
		// 解析Excel
		excel, err := parseExcel(path)
		if err != nil {
			return err
		}

		// 检查Excel是否满足导出为toml格式
		if err = checkToml(excel); err != nil {
			return err
		}

		// 导出类型
		excel.exportType = "config"
		excel.exportServerGoCodePath = exportConfig.ServerExportGoCodePath
		excel.exportClientGoCodePath = exportConfig.ClientExportGoCodePath
		excel.exportClientCSharpCodePath = exportConfig.ClientExportCSharpCodePath
		for _, exportTypeConfig := range exportConfig.ExcelExportType {
			if exportTypeConfig.Excel == excel.name {
				excel.exportType = exportTypeConfig.Type
				excel.exportServerGoCodePath = exportTypeConfig.ServerExportGoCodePath
				excel.exportClientCSharpCodePath = exportTypeConfig.ClientExportCSharpCodePath
				break
			}
		}

		err = exportExcel(excel)
		if err != nil {
			log.RunLogger.Println(err)
			allValid = false
		}

		allExcels = append(allExcels, excel)
		allConfigExcels = append(allConfigExcels, excel)
	}

	// 生成服务端toml数据读取代码
	if exportConfig.ServerExportGoCodePath != "" && exportConfig.hasGoEnv {
		genReadAllTomlCode(allExcels)
	}

	// 导出toml数据对应的Proto定义
	if exportConfig.ClientExportProtoBufDataPath != "" {
		goProto, csharpProto := genProtoDefineFromToml(allConfigExcels, "client")
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
