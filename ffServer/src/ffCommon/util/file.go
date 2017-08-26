package util

import (
	"ffCommon/log/log"

	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// CompressZIP 实现文件或文件夹压缩
// frm: 待压缩的文件或文件夹
// dst: 压缩结果的文件存储路径，不会追加或修改文件后缀
func CompressZIP(frm, dst string) error {
	// 创建目标文件所需的路径
	dstDir := filepath.Dir(dst)
	if fi, _ := os.Stat(dstDir); fi == nil {
		os.MkdirAll(dstDir, 0777)
	}

	buf := bytes.NewBuffer(make([]byte, 0, 10*1024*1024)) // 创建一个读写缓冲
	myzip := zip.NewWriter(buf)                           // 用压缩器包装该缓冲

	// 遍历源
	err := filepath.Walk(frm, func(path string, info os.FileInfo, err error) error {
		var file []byte
		if err != nil {
			return filepath.SkipDir
		}

		header, err := zip.FileInfoHeader(info) // 转换为zip格式的文件信息
		if err != nil {
			return filepath.SkipDir
		}

		header.Name, _ = filepath.Rel(filepath.Dir(frm), path)
		if !info.IsDir() {
			// 确定采用的压缩算法（这个是内建注册的deflate）
			header.Method = 8
			file, err = ReadFile(path) // 获取文件内容
			if err != nil {
				return filepath.SkipDir
			}
		} else {
			file = nil
		}

		// 上面的部分如果出错都返回filepath.SkipDir
		// 下面的部分如果出错都直接返回该错误
		// 目的是尽可能的压缩目录下的文件，同时保证zip文件格式正确
		w, err := myzip.CreateHeader(header) // 创建一条记录并写入文件信息
		if err != nil {
			return err
		}

		_, err = w.Write(file) // 非目录文件会写入数据，目录不会写入数据
		if err != nil {        // 因为目录的内容可能会修改
			return err // 最关键的是我不知道咋获得目录文件的内容
		}
		return nil
	})

	if err != nil {
		return err
	}

	myzip.Close() // 关闭压缩器，让压缩器缓冲中的数据写入buf

	file, err := os.Create(dst) // 建立zip文件
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = buf.WriteTo(file) // 将buf中的数据写入文件
	if err != nil {
		return err
	}

	return nil
}

// MergeFileNameSlice 合并目录内名称相似的文件为大文件，并输出到子目录output目录内
// a.b.c.d1, a.b.c.d1 ==> a.b.c
// 文件操作函数：http://www.cnblogs.com/sevenyuan/archive/2013/02/28/2937275.html
type MergeFileNameSlice []string

func (p MergeFileNameSlice) Len() int {
	return len(p)
}
func (p MergeFileNameSlice) Less(i, j int) bool {
	fpi := p[i]
	p1 := strings.LastIndex(fpi, ".")
	if p1 != -1 {
		fpi = fpi[:p1]
	}

	fpj := p[j]
	p2 := strings.LastIndex(fpj, ".")
	if p2 != -1 {
		fpj = fpj[:p2]
	}

	if fpi == fpj {
		if len(p[i]) == len(p[j]) {
			return p[i] < p[j]
		}
		return len(p[i]) < len(p[j])
	}
	return fpi < fpj
}
func (p MergeFileNameSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// MergeFiles 拆分文件合并为整文件
func MergeFiles(dir string) (err error) {
	// 输出目录
	outputDir := filepath.Join(dir, "output")

	// 创建目标文件所需的路径
	if fi, _ := os.Stat(outputDir); fi != nil {
		os.RemoveAll(outputDir)
	}
	os.MkdirAll(outputDir, 0777)

	fileNames := make([]string, 0, 8)

	// 获取目录内文件名列表
	err = filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		fileNames = append(fileNames, f.Name())
		return nil
	})

	if err != nil {
		return err
	}

	// 排序
	sort.Sort(MergeFileNameSlice(fileNames))

	// a.b.c.d1, a.b.c.d1 ==> a.b.c
	mapMergeFiles := make(map[string][]string, 4)
	for _, fileName := range fileNames {
		mergeFileName := fileName

		lastPointIndex := strings.LastIndex(fileName, ".")
		if lastPointIndex != -1 {
			mergeFileName = fileName[:lastPointIndex]
		}

		if _, ok := mapMergeFiles[mergeFileName]; !ok {
			mapMergeFiles[mergeFileName] = make([]string, 0, 1)
		}
		mapMergeFiles[mergeFileName] = append(mapMergeFiles[mergeFileName], fileName)
	}

	// 保存合并后的文件
	for mergeFileName, sourceFileNames := range mapMergeFiles {
		log.RunLogger.Println(mergeFileName)
		for _, one := range sourceFileNames {
			log.RunLogger.Println(one)
		}
		log.RunLogger.Println("")

		fw, err := os.Create(filepath.Join(outputDir, mergeFileName))
		if err != nil {
			return err
		}
		defer fw.Close()

		for _, sourceFileName := range sourceFileNames {
			contents, err := ReadFile(filepath.Join(dir, sourceFileName))
			if err != nil {
				return err
			}

			fw.Write(contents)
		}
	}

	return err
}

// IsPathExist 指定路径是否可达, 可用于判定文件或文件夹是否存在
// 第一个返回值: 是否可达
// 第二个返回值: 不可达时, 是否发生了错误
// 使用方式:
//		优先处理error, 然后再根据路径是否可达进行逻辑处理
func IsPathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// CreatePath 当指定目录路径不可达时, 创建目录路径
func CreatePath(path string) error {
	isExist, err := IsPathExist(path)
	if err != nil {
		return err
	}
	if isExist {
		return nil
	}
	return os.MkdirAll(path, 0777)
}

// RemovePath 删除指定目录路径（连带删除目录本身）（删除后立即创建同名目录且立即访问，则可能出错）
func RemovePath(path string) error {
	return os.RemoveAll(path)
}

// RemoveFile removes the named file or directory.
// If there is an error, it will be of type *PathError.
func RemoveFile(name string) error {
	return os.Remove(name)
}

// ClearPath 按条件清空指定目录路径
//	targetDir: 要清空的文件夹
//	clearSubDir: 是否清空子目录. 清空子目录时, 但又由于fileSuffixLimit限制, 导致子目录下某些文件未被删除时, 将导致函数执行失败
//	fileSuffixLimit: 要删除文件的后缀, 不限定时, 即全部删除
func ClearPath(targetDir string, clearSubDir bool, fileSuffixLimit []string) error {
	// 路径不存在时，创建
	// 路径存在且访问报错，则返回错误
	// 路径存在且正常访问，则遍历删除其中的所有元素
	if ok, err := IsPathExist(targetDir); err != nil {
		return err
	} else if !ok {
		return CreatePath(targetDir)
	}

	targetDir, err := filepath.Abs(targetDir)
	if err != nil {
		return nil
	}

	// 遍历源
	var waitDelSubDirs []string
	waitDelFiles := make([]string, 0, 1)
	if clearSubDir {
		// 清除子目录
		waitDelSubDirs = make([]string, 0, 1)
		err = filepath.Walk(targetDir, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			p, err = filepath.Abs(p)
			if err != nil {
				return err
			}

			if info.IsDir() {
				// 记录子目录，待删除
				if p != targetDir {
					waitDelSubDirs = append(waitDelSubDirs, p)
				}
				return nil
			}

			// 文件添加到待待删除文件列表内
			if fileSuffixLimit != nil {
				for _, limit := range fileSuffixLimit {
					if len(limit) > 0 && strings.HasSuffix(info.Name(), limit) {
						waitDelFiles = append(waitDelFiles, p)
					}
				}
			} else {
				waitDelFiles = append(waitDelFiles, p)
			}

			return err
		})

	} else {

		// 不遍历子目录
		err = Walk(targetDir, func(info os.FileInfo) error {
			if info.IsDir() {
				return nil
			}

			if fileSuffixLimit != nil {
				for _, limit := range fileSuffixLimit {
					if len(limit) > 0 && strings.HasSuffix(info.Name(), limit) {
						waitDelFiles = append(waitDelFiles, filepath.Join(targetDir, info.Name()))
					}
				}
			} else {
				waitDelFiles = append(waitDelFiles, filepath.Join(targetDir, info.Name()))
			}

			return nil
		})
	}

	for _, f := range waitDelFiles {
		err = os.Remove(f)
		if err != nil {
			return err
		}
	}

	if waitDelSubDirs != nil {
		for _, d := range waitDelSubDirs {
			err = os.Remove(d)
			if err != nil {
				return err
			}
		}
	}

	return err
}

// CopyFile 文件拷贝，目标目录必须已经存在
func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()

	CreatePath(path.Dir(dstName))

	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

// CopyPath 目录拷贝
func CopyPath(to, frm string) (err error) {
	// 遍历源
	err = filepath.Walk(frm, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		_, err = CopyFile(path.Join(to, info.Name()), p)
		return err
	})

	return err
}

// CreateFile 创建文件, 调用者负责Close
func CreateFile(filePath string) (*os.File, error) {
	err := CreatePath(filepath.Dir(filePath))
	if err != nil {
		return nil, err
	}

	// FileMode: rwxrwxrwx, 此处是 rw-rw-rw-, 即0666, 代表所有用户可读写, 但不能执行
	return os.Create(filePath)
}

// AppedFile 追加到文件末尾(不存在时创建), 调用者负责Close
func AppedFile(filePath string) (*os.File, error) {
	err := CreatePath(filepath.Dir(filePath))
	if err != nil {
		return nil, err
	}

	// FileMode: rwxrwxrwx, 此处是 rw-rw-rw-, 即0666, 代表所有用户可读写, 但不能执行
	return os.OpenFile(filePath, os.O_CREATE|os.O_APPEND, 0666)
}

// WriteFile writes data to a file named by filename.
// If the file does not exist, WriteFile creates it with permissions perm;
// otherwise WriteFile truncates it before writing.
func WriteFile(filename string, data []byte) error {
	err := CreatePath(filepath.Dir(filename))
	if err != nil {
		return err
	}

	// FileMode: rwxrwxrwx, 此处是 rw-rw-rw-, 即0666, 代表所有用户可读写, 但不能执行
	return ioutil.WriteFile(filename, data, 0666)
}

// ReadFile reads the file named by filename and returns the contents.
// A successful call returns err == nil, not err == EOF. Because ReadFile
// reads the whole file, it does not treat an EOF from Read as an error
// to be reported.
func ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

// Walk walks the file tree rooted at root, calling walkFn for each file or directory in the tree, including root.
// 不会递归遍历子目录(filepath.Walk会递归遍历子目录)
func Walk(root string, walkFn func(f os.FileInfo) error) error {
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}

	for _, f := range files {
		if err = walkFn(f); err != nil {
			return err
		}
	}
	return nil
}
