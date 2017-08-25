package logfile

type logLevel int

// 日志记录请求
type logRequest struct {
	buf     []byte   // log content
	stdout  bool     // 是否同时输出到标准输出
	outSync bool     // 同步输出标记
	level   logLevel // log level

	day int // 哪天, 用于隔天切换日志
}

// Cheap integer to fixed-width decimal ASCII.  Give a negative width to avoid zero-padding.
// 拷贝自标准库 log
// 实现将数字转换为固定长度的字符串输出到字节数组
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}
