package main

import (
	"ffCommon/encrypt"
	"ffCommon/log/log"
	"ffCommon/util"
	"ffCommon/version"
	"flag"
	"time"

	"github.com/lexical005/toml"

	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

const (
	desKey = "sG(, :Wr"
)

var config struct {
	ChannelName string

	ChannelAllVersion string

	ChannelExcelClient string
	ChannelExcelServer string

	ChannelUploadGameServer string
	ChannelUploadResServer  string

	// 0->chinese
	// 1->english
	// 2->russia
	// 3->korea
	// 4->中东
	Language string
}

var nowUploadFileName = make([]string, 0, 1)
var nowUploadFileMD5 = make([]string, 0, 1)
var nowUploadFileSize = make([]int, 0, 1)

var latestVersionInfo []interface{}

func copyAndEncryptFile(dstDir, srcName string) (dstName string, written int, err error) {
	datas, err := util.ReadFile(srcName)
	if err != nil {
		return dstName, 0, err
	}

	datas, err = encrypt.DesEncryptPaddingZero(datas, []byte(desKey))
	if err != nil {
		return dstName, 0, err
	}

	datasBase64 := encrypt.EncodeToBase64(datas, 0)

	md5 := encrypt.MD5([]byte(datasBase64), 1)
	dstName = md5 + ".TXT"

	dst, err := os.OpenFile(path.Join(dstDir, dstName), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	written, err = dst.Write([]byte(datasBase64))

	return dstName, written, err
}

// 保存到结构目录
func saveResult(result []interface{}) error {
	resNamesResult, _ := result[0].([]interface{})
	md5NamesResult, _ := result[1].([]interface{})
	resSizesResult, _ := result[2].([]interface{})
	for i := 0; i < len(resNamesResult); i++ {
		resName, _ := resNamesResult[i].(string)
		dstName, dstSize, err := copyAndEncryptFile(config.ChannelUploadResServer, path.Join(config.ChannelExcelClient, resName+".json"))
		if err != nil {
			return err
		}
		md5NamesResult[i] = dstName
		resSizesResult[i] = dstSize
	}

	dst, err := os.OpenFile(path.Join(config.ChannelUploadResServer, config.ChannelName+".TXT"), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer dst.Close()

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return err
	}

	resultBytes, err = encrypt.DesEncryptPaddingZero(resultBytes, []byte(desKey))
	if err != nil {
		return err
	}

	datasBase64 := encrypt.EncodeToBase64(resultBytes, 0)

	_, err = dst.Write([]byte(datasBase64))
	if err != nil {
		return err
	}

	return nil
}

// 获取最新发布
func getLatestPublishVersionFile() (err error) {
	versions := make([]*version.Version, 0, 1)

	rootPath := path.Join(config.ChannelAllVersion, config.ChannelName)
	rootPath, err = filepath.Abs(rootPath)
	if err != nil {
		return err
	}

	err = filepath.Walk(rootPath, func(p string, f os.FileInfo, err error) error {
		if f == nil {
			log.RunLogger.Println(err)
			return err
		}

		if !f.IsDir() {
			return nil
		}

		p, err = filepath.Abs(p)
		if err != nil {
			return err
		}

		if p == rootPath {
			return nil
		}

		v, err := version.New(f.Name())
		if err != nil {
			return err
		}

		versions = append(versions, v)
		return nil
	})

	if err != nil {
		return err
	}

	if len(versions) == 0 {
		return fmt.Errorf("getLatestPublishVersionFile: no publish version exist in channel[%v]", config.ChannelName)
	}

	// 计算最新发布的版本号
	latestPublishVersion := versions[0]
	for index := 1; index < len(versions); index++ {
		if latestPublishVersion.Compare(versions[index]) < 0 {
			latestPublishVersion = versions[index]
		}
	}

	// 读取并解析文件内容
	p := path.Join(config.ChannelAllVersion, config.ChannelName, latestPublishVersion.String(), config.ChannelName+".TXT")
	contentOri, err := util.ReadFile(p)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(contentOri, &latestVersionInfo); err != nil {
		return err
	}

	return nil
}

// 根据当前配置生成相应的数据结构
func genNowUploadInfo() (count int, err error) {
	err = filepath.Walk(config.ChannelExcelClient, func(p string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		fileNameInfos := strings.Split(f.Name(), ".")
		fileName, fileSuffix := fileNameInfos[0], fileNameInfos[len(fileNameInfos)-1]
		if fileSuffix != "json" {
			return nil
		}

		content, err := util.ReadFile(p)
		if err != nil {
			return err
		}

		nowUploadFileName = append(nowUploadFileName, fileName)
		nowUploadFileMD5 = append(nowUploadFileMD5, encrypt.MD5(content, 1)+".TXT")
		nowUploadFileSize = append(nowUploadFileSize, len(content))

		return nil
	})

	if err != nil {
		return 0, err
	}

	log.RunLogger.Println("genNowUploadInfo:")
	log.RunLogger.Println("nowUploadFileName", nowUploadFileName)
	log.RunLogger.Println("nowUploadFileMD5", nowUploadFileMD5)
	log.RunLogger.Println("nowUploadFileSize", nowUploadFileSize)
	log.RunLogger.Println("")

	return len(nowUploadFileName), nil
}

// 生成结果
func genResult() (needUpload bool, result []interface{}, err error) {
	resNamesOri, _ := latestVersionInfo[0].([]interface{})
	md5NamesOri, _ := latestVersionInfo[1].([]interface{})

	resNamesResult := make([]interface{}, 0, 2)
	md5NamesResult := make([]interface{}, 0, 2)
	resSizesResult := make([]interface{}, 0, 2)

	for i := 0; i < len(nowUploadFileName); i++ {
		found := false
		for j := 0; j < len(resNamesOri); j++ {
			if nowUploadFileName[i] == resNamesOri[j] {
				if nowUploadFileMD5[i] != md5NamesOri[j] {
					log.RunLogger.Println("md5Change", nowUploadFileName[i], nowUploadFileMD5[i], md5NamesOri[j])
					resNamesResult = append(resNamesResult, nowUploadFileName[i])
					md5NamesResult = append(md5NamesResult, nowUploadFileMD5[i])
					resSizesResult = append(resSizesResult, nowUploadFileSize[i])
				}
				found = true
				break
			}
		}
		if !found {
			log.RunLogger.Println("newFile", nowUploadFileName[i], nowUploadFileMD5[i])
			resNamesResult = append(resNamesResult, nowUploadFileName[i])
			md5NamesResult = append(md5NamesResult, nowUploadFileMD5[i])
			resSizesResult = append(resSizesResult, nowUploadFileSize[i])
		}
	}

	if len(resNamesResult) == 0 {
		return false, nil, nil
	}

	log.RunLogger.Println("")

	results := make([]interface{}, 3, 3)
	results[0] = resNamesResult
	results[1] = md5NamesResult
	results[2] = resSizesResult
	log.RunLogger.Println("\ngenResult:")
	log.RunLogger.Println(results)

	// json.Unmarshal
	return true, results, nil
}

func main() {
	defer util.PanicProtect(func(isPanic bool) {
		if isPanic {
			log.RunLogger.Println("异常退出, 以上是错误堆栈")
			<-time.After(time.Hour)
		}
	}, "ffHotResGen")

	defer func() {
		// 清理生成目录
		err := util.ClearPath(config.ChannelExcelClient, true, nil)
		if err != nil {
			log.RunLogger.Println(err)
			return
		}

		err = util.ClearPath(config.ChannelExcelServer)
		if err != nil {
			log.RunLogger.Println(err)
			return
		}
	}()

	// 命令行参数解析
	genJSONFromExcel := flag.Bool("gen", false, "gen json from excel")
	flag.Parse()

	log.RunLogger.Println("gen flag", *genJSONFromExcel)

	// 读取配置
	tomlContent, err := util.ReadFile("config.toml")
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	err = toml.Unmarshal(tomlContent, &config)
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	log.RunLogger.Printf("%+v\n", config)
	log.RunLogger.Println()

	// 清理上传目录
	err = util.ClearPath(config.ChannelUploadGameServer)
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	err = util.ClearPath(config.ChannelUploadResServer)
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	if *genJSONFromExcel {
		// 清理生成目录
		err = util.ClearPath(config.ChannelExcelClient)
		if err != nil {
			log.RunLogger.Println(err)
			return
		}

		err = util.ClearPath(config.ChannelExcelServer)
		if err != nil {
			log.RunLogger.Println(err)
			return
		}

		// 将excel生成为程序使用的json
		curDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.RunLogger.Println(err)
			return
		}

		err = os.Chdir(path.Join(curDir, "Excel"))
		if err != nil {
			log.RunLogger.Println(err)
			return
		}

		f, err := exec.Command("make_all.bat", config.Language, "dynamic").Output()
		if err != nil {
			log.RunLogger.Println(err)
			return
		}
		log.RunLogger.Println(string(f))

		err = os.Chdir(curDir)
		if err != nil {
			log.RunLogger.Println(err)
			return
		}
	}

	// 获取最新发布的版本文件
	err = getLatestPublishVersionFile()
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	// 当前配置生成相应的数据结构
	uploadFileCount, err := genNowUploadInfo()
	if err != nil {
		log.RunLogger.Println(err)
		return
	} else if uploadFileCount == 0 {
		log.RunLogger.Println("no file need to upload")
		return
	}

	// 生成更新配置
	needUpload, result, err := genResult()
	if err != nil {
		log.RunLogger.Println(err)
		return
	} else if !needUpload {
		log.RunLogger.Println("no file need to upload")
		return
	}

	// 保存到资源服务器更新上传目录
	err = saveResult(result)
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	// 将需要上传到游戏服务器的数据，直接拷贝到游戏服务器上传目录
	err = util.CopyPath(config.ChannelUploadGameServer, config.ChannelExcelServer)
	if err != nil {
		log.RunLogger.Println(err)
		return
	}
}
