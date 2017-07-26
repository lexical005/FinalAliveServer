package main

import (
	"ffCommon/log/log"
	"ffCommon/util"
	"time"

	"crypto/md5"
	"encoding/hex"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

var muGlobalHotResCache sync.Mutex

// 生成热更新资源列表
func genHotResMap(channelInfo *tChannelInfo, reset bool) {
	muGlobalHotResCache.Lock()
	defer muGlobalHotResCache.Unlock()

	hotResChanged := 0

	if reset {
		hotResChanged = -1

		var t time.Time
		channelInfo.hotResFileMd5 = ""
		channelInfo.hotResFileModTime = t
		channelInfo.hotResMap = make(map[string]*tHotResInfo)
	}

	p := path.Join("data", channelInfo.channelName, "hotres", channelInfo.newestVersion.String())
	if fi, err := os.Stat(p); err != nil && os.IsExist(err) || fi != nil && fi.IsDir() {
		err := filepath.Walk(p, func(path string, f os.FileInfo, err error) error {
			if f == nil {
				log.RunLogger.Println(err.Error())
				return err
			}
			if f.IsDir() {
				return nil
			}

			fileName := f.Name()

			// 渠道热更新资源配置
			fileNameWithoutSuffix := strings.Split(fileName, ".")[0]
			if fileNameWithoutSuffix == channelInfo.channelName {
				fileName = strings.ToUpper(fileName)

				// 文件的修改时间一致
				if channelInfo.hotResFileModTime == f.ModTime() {
					return nil
				}

				// 读取文件
				fileBuffer, e := util.ReadFile(path)
				if e != nil {
					log.RunLogger.Println("read hot res failed, file_name:" + fileName + " err:" + err.Error())
					return nil
				}

				// 计算文件内容的 md5
				fileContentMd5 := md5.Sum(fileBuffer)
				channelInfo.hotResFileMd5 = strings.ToUpper(hex.EncodeToString(fileContentMd5[:]))

				// 缓存
				channelInfo.hotResMap[fileName] = &tHotResInfo{
					buf:     fileBuffer,
					modTime: f.ModTime(),
				}

				return nil
			}

			fileName = strings.ToUpper(fileName)

			// 该文件已经读取过, 且文件的修改时间一致
			if hotResInfo, ok := channelInfo.hotResMap[fileName]; ok && hotResInfo.modTime == f.ModTime() {
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
			channelInfo.hotResMap[fileName] = &tHotResInfo{
				buf:     fileBuffer,
				modTime: f.ModTime(),
			}

			hotResChanged = 1
			return nil
		})

		if err != nil {
			log.RunLogger.Println(err)
		}
	}

	// 热更新资源发生了变动
	if hotResChanged != 0 {
		if hotResChanged > 0 {
			log.RunLogger.Println(channelInfo.channelName, "hot res update with hotres")
		} else {
			log.RunLogger.Println(channelInfo.channelName, "hot res update without hotres")
		}
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
		log.RunLogger.Println("")
		log.RunLogger.Println("")
	}
}
