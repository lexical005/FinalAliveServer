package logfile

import (
	"ffCommon/util"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

type latestLogFileInfo struct {
	name  string
	year  string
	month string
	day   string
	count string

	size int64
}

// update 更新
//	data: regexp.MustCompile(fmt.Sprintf("%v ([\\d]+)-([\\d]+)-([\\d]+) ([\\d]+)\\.log", l.filePrefix))返回的结果
func (info *latestLogFileInfo) update(data []string, size int64) {
	if data == nil {
		return
	}

	if info.name < data[0] {
		info.name, info.year, info.month, info.day, info.count = data[0], data[1], data[2], data[3], data[4]
		info.size = size
	}
}

// latestName 最新名称
//	返回值: 最新日志文件应使用的名称, 此日志文件剩余可写入大小
func (info *latestLogFileInfo) latestName(filePrefix string, fileLenLimit int, forceSwitch bool) (string, int, int) {
	// 时间
	now := time.Now()
	year, month, day := now.Date()
	count := 1
	outLenLimit := int64(fileLenLimit)

	if info.name != "" {
		lyear, _ := strconv.Atoi(info.year)
		lmonth, _ := strconv.Atoi(info.month)
		lday, _ := strconv.Atoi(info.day)
		lcount, _ := strconv.Atoi(info.count)

		// 同一天
		if year == lyear && int(month) == lmonth && day == lday {
			if !forceSwitch {
				if info.size >= outLenLimit {
					// 最后一个日志文件已满了
					count = lcount + 1
				} else {
					// 最后一个日志文件未满, 继续写入
					count = lcount
					outLenLimit -= info.size
				}
			} else {
				// 强制换一个文件
				count = lcount + 1
			}
		}
	}

	bufFileName := make([]byte, 0, 512)
	buf := &bufFileName

	// 前缀
	*buf = append(*buf, filePrefix...)
	*buf = append(*buf, ' ')

	// 日期与时间
	itoa(buf, year, 4)
	*buf = append(*buf, '-')
	itoa(buf, int(month), 2)
	*buf = append(*buf, '-')
	itoa(buf, day, 2)
	*buf = append(*buf, ' ')

	itoa(buf, count, 0)

	// 后缀
	*buf = append(*buf, ".log"...)
	return string(*buf), int(outLenLimit), day
}

func latestName(filePath string, filePrefix string, fileLenLimit int, forceSwitch bool) (string, int, int) {
	// 遍历之前的日志文件
	r := regexp.MustCompile(fmt.Sprintf("%v ([\\d]+)-([\\d]+)-([\\d]+) ([\\d]+)\\.log", filePrefix))
	latestLogFileInfo := &latestLogFileInfo{}
	util.Walk(filePath, func(info os.FileInfo) error {
		if info.IsDir() {
			return nil
		}

		latestLogFileInfo.update(r.FindStringSubmatch(info.Name()), info.Size())
		return nil
	})

	// 实际应该使用的名称, 以及写入长度限制
	return latestLogFileInfo.latestName(filePrefix, fileLenLimit, forceSwitch)
}
