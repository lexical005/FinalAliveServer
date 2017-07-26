package ffExcel

import (
	"ffCommon/log/log"
	"ffCommon/util"

	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// ExportConfig 导出配置
type ExportConfig struct {
	ServerExportDefPath  string // 服务端导出的定义文件, 相对导出程序的路径, 为空时或者系统环境变量中未配置GOPATH时, 不导出
	ServerExportDataType string // 服务端导出的配置文件的类型
	ServerExportDataPath string // 服务端导出的配置文件, 相对导出程序的路径
	ServerReadDataPath   string // 服务端导出的配置文件, 相对读取程序的路径

	ClientExportDataType string // 客户端导出的配置文件的类型
	ClientExportDataPath string // 客户端导出的配置文件, 相对导出程序的路径

	hasGoEnv    bool   // 是否有go环境
	packageName string // 根据ServerExportDefPath推导出来的包名
}

func (ec *ExportConfig) check() error {
	ec.hasGoEnv = os.Getenv("GOPATH") != ""

	if ec.ServerExportDataType != "toml" {
		return fmt.Errorf("ExportConfig ServerExportDataType[%v] not support", ec.ServerExportDataType)
	}

	if ec.ClientExportDataType != "lua" {
		return fmt.Errorf("ExportConfig ClientExportDataType[%v] not support", ec.ClientExportDataType)
	}

	_, ec.packageName = path.Split(ec.ServerExportDefPath)

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
	ServerExportDefPath:%v
	ServerExportDataType:%v
	ServerExportDataPath:%v
	ClientExportDataType:%v
	ClientExportDataPath:%v
]`,
		ec.ServerExportDefPath,
		ec.ServerExportDataType,
		ec.ServerExportDataPath,
		ec.ClientExportDataType,
		ec.ClientExportDataPath)
}

// ExportExcel 解析excel, 然后根据导出配置, 将解析结果保存
func ExportExcel(excelFilePath string, exportConfig *ExportConfig) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("ExportExcel excel[%v] get error:\n[%v]\n", excelFilePath, err)
		}
	}()

	excel, err := parseExcel(excelFilePath)
	if err != nil {
		return err
	}

	if exportConfig.ServerExportDataType == "toml" {
		tomlDef, tomlData, err := genToml(excel, exportConfig)
		if err != nil {
			return err
		}

		gofilename := strings.ToLower(excel.name)
		tomlfilename := excel.name

		// 导出定义
		if exportConfig.hasGoEnv && exportConfig.ServerExportDefPath != "" {
			// 输出到临时目录
			defFilePath := path.Join(exportConfig.ServerExportDefPath, gofilename+".go")
			err = util.WriteFile(defFilePath, []byte(tomlDef))
			if err != nil {
				return err
			}
			log.RunLogger.Println(defFilePath)
		}

		// 导出配置
		dataFilePath := path.Join(exportConfig.ServerExportDataPath, tomlfilename+".toml")
		err = util.WriteFile(dataFilePath, []byte(tomlData))
		if err != nil {
			return err
		}
		log.RunLogger.Println(dataFilePath)
	}

	if exportConfig.ClientExportDataType == "lua" {
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
