package main

import (
	"ffCommon/log/log"
	"ffCommon/version"

	"strconv"
	"strings"
	"time"
)

const subpackageSize = 1024 * 1024 * 10

// 整包类型
const (
	channelPackageTypeInvalid = iota
	channelPackageTypeApk
	channelPackageTypeIpa
)

var serverConfig *tServerConfig // 服务器配置

//---------------------------------------------------------------------------
// 版本号(有效的版本号 0.0.2 主版本号.次版本号.修订版本号)
type tVersionName string

func (t tVersionName) String() string {
	return string(t)
}
func (t tVersionName) isValid() bool {
	r := strings.Split(string(t), ".")
	if len(r) != 3 {
		return false
	}

	for _, one := range r {
		if _, e := strconv.Atoi(one); e != nil {
			return false
		}
	}

	return true
}
func (t tVersionName) buildVersion() int {
	i := strings.LastIndex(string(t), ".")
	s := string(t)[i+1:]
	v, _ := strconv.Atoi(s)
	return v
}
func (t tVersionName) Less(other tVersionName) bool {
	r1 := strings.Split(string(t), ".")
	r2 := strings.Split(string(other), ".")
	for i := 0; i < len(r1); i++ {
		v1, _ := strconv.Atoi(r1[i])
		v2, _ := strconv.Atoi(r2[i])
		if v1 < v2 {
			return true
		}
	}
	return false
}

//---------------------------------------------------------------------------
// 资源服配置
type tServerConfig struct {
	listenIPPort    string // 本地监听的 ip:port
	outerIPPort     string // 客户端连接请求的 ip:port
	connectionLimit int    // 最大客户端连接数限制

	devChannel string // 测试渠道名称

	channelConfig map[string]map[string]string // 渠道配置 渠道1:{k1:v1, k2:v2}
}

//---------------------------------------------------------------------------
// 一段字节流的信息
type tBytesPackage struct {
	fileBuffer []byte // 本段字节流
	fileSize   int    // 本段大小
	md5Str     string // 本段md5
	urlStr     string // 本段的下载地址. 格式: http://192.168.0.23:8765/apk?apk=uc_1.0.105_X.apk|SUBPACKAGE_SIZE|md5_str
}

//---------------------------------------------------------------------------
// 热更新资源的描述信息
type tHotResInfo struct {
	buf     []byte
	modTime time.Time
}

//---------------------------------------------------------------------------
// 一个渠道的信息
type tChannelInfo struct {
	channelName        string           // 渠道的名称
	fullPackageVersion *version.Version // 渠道的完整包的版本号
	newestVersion      *version.Version // 渠道的最新版本号

	hasExamineVersion bool             // 是否有送审版本
	examineVersion    *version.Version // 渠道的送审版本号
	examineResponse   []byte           // 针对送审版本，渠道的返回值

	fullPackageDownResponse []byte            // 需要下载完整包时的回复
	newestResponse1         []byte            // 版本最新时的回复1-服务端没有热更新资源配置
	newestResponse2         []byte            // 版本最新时的回复2-服务端有热更新资源配置, 且与客户端的配置一致
	newestResponse3         []byte            // 版本最新时的回复3-服务端有热更新资源配置, 且与客户端的配置不一致
	newestResponseOri       map[string]string // 版本最新时的数据

	hotResFileMd5     string                  // 热更新配置文件的md5
	hotResFileModTime time.Time               // 热更新配置文件的修改时间
	hotResMap         map[string]*tHotResInfo // 热更新资源文件

	apkSubpackages []*tBytesPackage // 完整包拆分处的分包

	// hotfix                map[t_version_name]*t_bytes_package // 更新包的字节流. key: 起始版本号
	// hotfix_version_from   []t_version_name                    // 支持从哪些版本号进行更新包更新
	// hotfix_version_to     []t_version_name                    // 起始版本号对应的更新到的版本号
	// hotfix_cache_response map[t_version_name][]byte           // 缓存从该版本号更新到最新版本的回复
}

func (t *tChannelInfo) returnNewestResponse(md5 string) []byte {
	if t.hotResFileMd5 != "" {
		if t.hotResFileMd5 == md5 {
			return t.newestResponse2
		}
		return t.newestResponse3
	}
	return t.newestResponse1
}

func (t *tChannelInfo) returnExamineResponse(md5 string) []byte {
	return t.examineResponse
}

// 获得用户上传的版本号在指定渠道内的更新检查结果
func (t *tChannelInfo) checkVersion(userVersion *version.Version, md5 string) []byte {
	// 开发渠道总是返回最新
	if serverConfig.devChannel == t.channelName {
		return t.returnNewestResponse(md5)
	}

	// 服务端最新版本号与用户的版本号一致
	if t.newestVersion == userVersion {
		return t.returnNewestResponse(md5)
	}

	// 服务端送审版本号与用户的版本号一致，则知晓其为送审客户端
	if t.hasExamineVersion && userVersion.Compare(t.examineVersion) >= 0 {
		return t.returnExamineResponse(md5)
	}

	// 用户的版本号, 低于完整包的版本号
	if userVersion.Compare(t.fullPackageVersion) < 0 {
		return t.fullPackageDownResponse
	}

	// // 返回如何更新的缓存
	// if t.hotfix_cache_response != nil {
	// 	if cache, ok := t.hotfix_cache_response[userVersion]; ok {
	// 		return cache
	// 	}
	// }

	// 服务端最新版本号低于用户的版本号
	if userVersion.Compare(t.newestVersion) > 0 {
		log.RunLogger.Println("server newest verion lower than user verion newest_version: " + t.newestVersion.String() + " userVersion: " + userVersion.String())
	}

	// 出错了
	return t.returnNewestResponse(md5)
}
