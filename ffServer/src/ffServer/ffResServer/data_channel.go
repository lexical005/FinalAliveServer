package main

import (
	"ffCommon/log/log"
	"ffCommon/util"
	"ffCommon/version"

	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type channelPackage struct {
	valid bool

	channelName string

	fileBuffer []byte
	fileSize   int

	packageVersion *version.Version
	packageType    int
}

var globalChannelInfo map[string]*tChannelInfo
var globalSubPackageCache map[string]*tBytesPackage

//---------------------------------------------------------------------------
// 生成渠道数据
func genChannelInfo(packageOuter, packageExamine *channelPackage) *tChannelInfo {
	channelInfo := &tChannelInfo{}
	channelInfo.channelName = packageOuter.channelName

	channelInfo.fullPackageVersion = packageOuter.packageVersion // 外发版本的完整包版本号
	channelInfo.newestVersion = packageOuter.packageVersion      // 外发版本的最新版本号(目前不支持更新包)
	channelInfo.examineVersion = packageOuter.packageVersion     // 送审版本与外发版本一致

	// 外发版本最新版本回复
	resultNewest := make(map[string]string)
	resultNewest["status"] = "newest"
	resultNewest["selectserver_ip"] = serverConfig.channelConfig[packageOuter.channelName]["SELECT_SERVER_IP"]
	resultNewest["res_file"] = strings.ToUpper(packageOuter.channelName)

	channelInfo.newestResponseOri = resultNewest

	resultNewest["res_status"] = "1"
	resultNewestBytes, _ := json.Marshal(resultNewest)
	channelInfo.newestResponse1 = resultNewestBytes

	resultNewest["res_status"] = "2"
	resultNewestBytes, _ = json.Marshal(resultNewest)
	channelInfo.newestResponse2 = resultNewestBytes

	resultNewest["res_status"] = "3"
	resultNewestBytes, _ = json.Marshal(resultNewest)
	channelInfo.newestResponse3 = resultNewestBytes

	// 送审版本回复
	if packageExamine.valid {
		channelInfo.hasExamineVersion = true
		channelInfo.examineVersion = packageExamine.packageVersion

		resultNewest["res_status"] = "1"
		resultNewestBytes, _ := json.Marshal(resultNewest)
		channelInfo.examineResponse = resultNewestBytes
	}

	// 外发版本整包更新回复
	responseFullDown := make(map[string]string)

	if packageOuter.packageType == channelPackageTypeApk {
		responseFullDown["status"] = "apk"
	} else if packageOuter.packageType == channelPackageTypeIpa {
		responseFullDown["status"] = "ipa"
	}
	responseFullDown["selectserver_ip"] = serverConfig.channelConfig[packageOuter.channelName]["SELECT_SERVER_IP"]

	if serverConfig.channelConfig[packageOuter.channelName]["FULL_DOWN_URL"] != "" {
		// 通过给定下载地址，由客户端自行前往下载
		responseFullDown["down_url"] = serverConfig.channelConfig[packageOuter.channelName]["FULL_DOWN_URL"]

	} else {
		// 缓存完整包更新时的回复
		var apkSubpackageURL []string
		if packageOuter.fileSize%subpackageSize == 0 {
			channelInfo.apkSubpackages = make([]*tBytesPackage, 0, (int)(packageOuter.fileSize/subpackageSize))
			apkSubpackageURL = make([]string, 0, (int)(packageOuter.fileSize/subpackageSize))
		} else {
			channelInfo.apkSubpackages = make([]*tBytesPackage, 0, (int)(packageOuter.fileSize/subpackageSize+1))
			apkSubpackageURL = make([]string, 0, (int)(packageOuter.fileSize/subpackageSize+1))
		}
		for i := 0; i < packageOuter.fileSize; i += subpackageSize {
			apkSubpackage := &tBytesPackage{}

			// 分包的字节流和大小
			if i+subpackageSize < packageOuter.fileSize {
				apkSubpackage.fileBuffer = packageOuter.fileBuffer[i : i+subpackageSize]
				apkSubpackage.fileSize = subpackageSize
			} else {
				apkSubpackage.fileBuffer = packageOuter.fileBuffer[i:packageOuter.fileSize]
				apkSubpackage.fileSize = packageOuter.fileSize - i
			}

			// 分包的md5
			md5h := md5.New()
			md5h.Write(apkSubpackage.fileBuffer)
			apkSubpackage.md5Str = strings.ToLower(hex.EncodeToString(md5h.Sum(nil)))

			// 分包的下载地址
			var subpackageName string
			subpackageName = packageOuter.channelName + "-" + packageOuter.packageVersion.String() + "-" + strconv.Itoa(i/subpackageSize) + ".apk"
			urlStr := "http://" + serverConfig.outerIPPort + "/apk?apk=" + subpackageName
			fileSizeStr := strconv.Itoa(apkSubpackage.fileSize)
			apkSubpackage.urlStr = urlStr + "|" + fileSizeStr + "|" + apkSubpackage.md5Str

			apkSubpackageURL = append(apkSubpackageURL, apkSubpackage.urlStr)

			channelInfo.apkSubpackages = append(channelInfo.apkSubpackages, apkSubpackage)

			// 存到全局分包管理器内
			globalSubPackageCache[subpackageName] = apkSubpackage
		}

		apkSubpackageURLBytes, _ := json.Marshal(apkSubpackageURL)
		responseFullDown["apk"] = string(apkSubpackageURLBytes)
	}
	responseFullDownBytes, _ := json.Marshal(responseFullDown)
	channelInfo.fullPackageDownResponse = responseFullDownBytes

	return channelInfo
}

// 解析渠道包名
func parsePackageName(channelName string, packageFullName string) (packageVersion *version.Version, packageType int, err error) {
	// 包后缀必须以 .apk 或 .ipa 结尾
	packageType = channelPackageTypeInvalid
	if packageFullName[len(packageFullName)-4:] == ".apk" {
		packageType = channelPackageTypeApk
	} else if packageFullName[len(packageFullName)-4:] == ".ipa" {
		packageType = channelPackageTypeIpa
	} else {
		err = fmt.Errorf("invalid package suffix. package_full_name[%s]\n", packageFullName)
		return
	}

	// 完整包的名称, 必须符合特定的规则: N个字符的渠道名-版本号.xxx
	nameWithoutSuffix := packageFullName[0 : len(packageFullName)-4]
	nameSplitResult := strings.Split(nameWithoutSuffix, "-")
	if len(nameSplitResult) != 2 || channelName != nameSplitResult[0] {
		err = fmt.Errorf("invalid package name format. package_full_name[%s] channel_str[%s]\n", packageFullName, channelName)
		return
	}

	// 版本号必须有效
	packageVersion, err = version.New(nameSplitResult[1])
	if err != nil {
		err = fmt.Errorf("invalid package version. package_full_name[%s] package_version[%s]\n", packageFullName, packageVersion)
		return
	}
	return
}

//---------------------------------------------------------------------------
// 生成所有渠道信息
func genAllChannels() (ok bool) {
	globalChannelInfo = make(map[string]*tChannelInfo)
	globalSubPackageCache = make(map[string]*tBytesPackage)

	ok = true
	for channelName := range serverConfig.channelConfig {
		// 遍历渠道整包
		if fi, err := os.Stat("data/" + channelName + "/package/"); err != nil && os.IsExist(err) || fi != nil && fi.IsDir() {

			packageOuter := &channelPackage{valid: false, channelName: channelName}
			packageExamine := &channelPackage{valid: false, channelName: channelName}

			err := filepath.Walk("data/"+channelName+"/package/", func(path string, f os.FileInfo, err error) error {
				if f.IsDir() {
					return err
				}

				packageVersion, packageType, err := parsePackageName(channelName, f.Name())
				if err != nil {
					return err
				}
				if !packageOuter.valid {
					packageOuter.valid = true
					packageOuter.packageVersion, packageOuter.packageType = packageVersion, packageType

					packageOuter.fileSize = int(f.Size())
					packageOuter.fileBuffer, err = util.ReadFile(path)
					if err != nil {
						err = fmt.Errorf("channel[%s] read package failed. packageName[%s]\n", channelName, f.Name())
					}
				} else {
					if packageOuter.packageVersion.Compare(packageVersion) < 0 {
						packageExamine.valid = true
						packageExamine.packageVersion, packageExamine.packageType = packageVersion, packageType

						packageExamine.fileSize = int(f.Size())
						packageExamine.fileBuffer, err = util.ReadFile(path)
						if err != nil {
							err = fmt.Errorf("channel[%s] read package failed. packageName[%s]\n", channelName, f.Name())
						}

					} else if packageVersion.Compare(packageOuter.packageVersion) < 0 {
						packageOuter, packageExamine = packageExamine, packageOuter

						packageOuter.valid = true
						packageOuter.packageVersion, packageOuter.packageType = packageVersion, packageType

						packageOuter.fileSize = int(f.Size())
						packageOuter.fileBuffer, err = util.ReadFile(path)
						if err != nil {
							err = fmt.Errorf("channel[%s] read package failed. packageName[%s]\n", channelName, f.Name())
						}
					} else {
						err = fmt.Errorf("channel[%s] outer and examine package has same version. package_version[%s]\n", channelName, packageVersion)
					}
				}
				return err
			})

			if err == nil && !packageOuter.valid {
				err = fmt.Errorf("channel[%s] has no package\n", channelName)
			}

			if err != nil {
				ok = false
				log.RunLogger.Println(err)
				continue
			}

			// 根据渠道整包生成渠道最终数据
			globalChannelInfo[packageOuter.channelName] = genChannelInfo(packageOuter, packageExamine)
		}
	}

	for _, channelInfo := range globalChannelInfo {
		channelInfo.hotResMap = make(map[string]*tHotResInfo)
		genHotResMap(channelInfo, true)
	}

	return
}
