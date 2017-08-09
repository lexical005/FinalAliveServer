package ffExcel

import (
	"ffCommon/log/log"
	"ffCommon/util"

	"fmt"
	"os"
	"path"
	"path/filepath"
)

// ExportConfig 导出配置
//	服务端代码语言: go, 配置文件: toml
//	客户端最终代码语言: c#, 配置文件: protobuf字节流(客户端使用临时语言go和配置文件toml, 最终生成最终所需的代码文件和配置文件)
type ExportConfig struct {
	// ServerExportCodePath 服务端导出的代码文件, 相对导出程序的路径, 为空时或者系统环境变量中未配置GOPATH时, 不导出
	ServerExportCodePath string
	// ServerExportDataPath 服务端导出的配置文件, 相对导出程序的路径
	ServerExportDataPath string
	// ServerReadDataPath 服务端导出的配置文件, 相对读取程序的路径
	ServerReadDataPath string

	// ClientExportCodePath 服务端导出的代码文件, 相对导出程序的路径, 为空时, 不导出
	ClientExportCodePath string
	// ClientExportDataPath 客户端导出的配置文件, 相对导出程序的路径
	ClientExportDataPath string

	hasGoEnv    bool   // 是否有go环境
	packageName string // 根据ServerExportCodePath推导出来的包名
}

func (ec *ExportConfig) check() error {
	ec.hasGoEnv = os.Getenv("GOPATH") != ""

	_, ec.packageName = path.Split(ec.ServerExportCodePath)

	return nil
}

func (ec *ExportConfig) clearPath() bool {
	result := true
	if ec.ServerExportDataPath != "" {
		err := util.ClearPath(ec.ServerExportDataPath)
		if err != nil {
			log.RunLogger.Println(err)
			result = false
		}
	}

	if ec.ClientExportDataPath != "" {
		err := util.ClearPath(ec.ClientExportDataPath)
		if err != nil {
			log.RunLogger.Println(err)
			result = false
		}
	}

	return result
}

func (ec *ExportConfig) String() string {
	return fmt.Sprintf(`[
	ServerExportCodePath:%v
	ServerExportDataPath:%v
	ClientExportDataPath:%v
]`,
		ec.ServerExportCodePath,
		ec.ServerExportDataPath,
		ec.ClientExportDataPath)
}

// ExportExcel 解析excel, 然后根据导出配置, 将解析结果保存
func ExportExcel(excelFilePath string, exportConfig *ExportConfig) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("ExportExcel excel[%v] get error:\n[%v]", excelFilePath, err)
		}
	}()

	// 解析Excel
	excel, err := parseExcel(excelFilePath)
	if err != nil {
		return err
	}

	// 检查Excel是否满足导出为toml格式
	if err = checkToml(excel); err != nil {
		return err
	}

	// 导出服务端
	{
		// 导出服务端读取代码
		tomlDataServerReadCode := genTomlDataReadCode(excel, exportConfig, "server")
		if exportConfig.hasGoEnv && exportConfig.ServerExportCodePath != "" {
			defFilePath := path.Join(exportConfig.ServerExportCodePath, excel.name+".go")
			err = util.WriteFile(defFilePath, []byte(tomlDataServerReadCode))
			if err != nil {
				return err
			}
			log.RunLogger.Println(defFilePath)
		}

		// 导出服务端配置
		tomlDataServer := genTomlData(excel, exportConfig, "server")
		dataFilePath := path.Join(exportConfig.ServerExportDataPath, excel.name+".toml")
		err = util.WriteFile(dataFilePath, []byte(tomlDataServer))
		if err != nil {
			return err
		}
		log.RunLogger.Println(dataFilePath)

		dataFilePath = path.Join("toml", "server", excel.name+".toml")
		err = util.WriteFile(dataFilePath, []byte(tomlDataServer))
		if err != nil {
			return err
		}
	}

	// 导出客户端
	{
		// // 导出客户端读取代码
		// tomlDataServerReadCode := genTomlDataReadCode(excel, exportConfig, "client")
		// if exportConfig.hasGoEnv && exportConfig.ClientExportDataPath != "" {
		// 	defFilePath := path.Join(exportConfig.ClientExportDataPath, excel.name+".go")
		// 	err = util.WriteFile(defFilePath, []byte(tomlDataServerReadCode))
		// 	if err != nil {
		// 		return err
		// 	}
		// 	log.RunLogger.Println(defFilePath)
		// }

		// 导出客户端配置
		tomlDataServer := genTomlData(excel, exportConfig, "client")
		dataFilePath := path.Join(exportConfig.ClientExportDataPath, excel.name+".toml")
		err = util.WriteFile(dataFilePath, []byte(tomlDataServer))
		if err != nil {
			return err
		}
		log.RunLogger.Println(dataFilePath)

		dataFilePath = path.Join("toml", "client", excel.name+".toml")
		err = util.WriteFile(dataFilePath, []byte(tomlDataServer))
		if err != nil {
			return err
		}
	}

	return nil
}

// ExportExcelDir 解析指定目录内的所有excel, 然后根据导出配置, 将解析结果保存
func ExportExcelDir(excelDirPath string, exportConfig *ExportConfig) error {
	// 配置检查
	err := exportConfig.check()
	if err != nil {
		return err
	}

	//  清空导出目录
	if !exportConfig.clearPath() {
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
	allValid := true
	for _, excelFilePath := range excelFilePaths {
		err := ExportExcel(excelFilePath, exportConfig)
		if err != nil {
			log.RunLogger.Println(err)
			allValid = false
		}
	}

	// 有错发生, 提示用户
	if !allValid {
		return fmt.Errorf("ExportExcelDir dir[%v] not all excel export success", excelDirPath)
	}

	return nil
}
