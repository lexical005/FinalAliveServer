package main

import (
	"ffCommon/log/log"
	"ffCommon/util"

	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var globalRemoteResCache = make(map[string]*tHotResInfo)
var muGlobalRemoteResCache sync.Mutex

// 生成远端资源列表
func genRemoteResMap() {
	muGlobalRemoteResCache.Lock()
	defer muGlobalRemoteResCache.Unlock()

	remoteResChanged := false
	filepath.Walk("hotres/", func(path string, f os.FileInfo, err error) error {
		if f == nil {
			log.RunLogger.Println(err.Error())
			return err
		}
		if f.IsDir() {
			return nil
		}

		fileName := f.Name()
		fileName = strings.ToUpper(fileName)

		// 该文件已经读取过, 且文件的修改时间一致
		if hotResInfo, ok := globalRemoteResCache[fileName]; ok && hotResInfo.modTime == f.ModTime() {
			return nil
		}

		// 读取文件
		fileBuffer, e := util.ReadFile(path)
		if e != nil {
			log.RunLogger.Println("read hot res failed, file_name:" + fileName + " err:" + err.Error())
			return nil
		}

		// 计算文件内容的 md5, 判定有效性
		fileContentMd5 := md5.Sum(fileBuffer)
		strFileContentMd5 := strings.ToUpper(hex.EncodeToString(fileContentMd5[:]))
		if !strings.HasPrefix(fileName, strFileContentMd5) {
			log.RunLogger.Println("file content md5: " + strFileContentMd5 + " not match file name:" + fileName)
			return nil
		}

		// 缓存
		globalRemoteResCache[fileName] = &tHotResInfo{
			buf:     fileBuffer,
			modTime: f.ModTime(),
		}

		remoteResChanged = true
		return nil
	})

	// 热更新资源文件的 md5
	if remoteResChanged {
		log.RunLogger.Println("remote res update")
		// for channel_str, channel_info := range global_channel_info {
		// 	tmp_channel_str := strings.ToUpper(channel_str) + ".TXT"
		// 	if hot_res_info, ok := global_hot_res_cache[tmp_channel_str]; ok {
		// 		md5h := md5.New()
		// 		md5h.Write(hot_res_info.buf)
		// 		channel_info.hot_res_file_md5 = strings.ToUpper(hex.EncodeToString(md5h.Sum(nil)))
		// 	} else {
		// 		channel_info.hot_res_file_md5 = ""
		// 	}

		// 	log.RunLogger.Println(channel_str, channel_info.hot_res_file_md5)
		// }
		log.RunLogger.Println("\n")
	}
}
