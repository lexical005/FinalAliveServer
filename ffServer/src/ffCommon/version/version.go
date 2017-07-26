package version

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// VersionMain 主版本号
	VersionMain = iota

	// VersionSub 子版本号
	VersionSub

	// VersionTotal 累计版本号
	VersionTotal
)

// Version 版本号(0.0.2 主版本号.次版本号.累计版本号)
type Version struct {
	vers [3]int
}

// from 从字符串构建版本号
func (v *Version) from(ver string) (err error) {
	sVers := strings.Split(ver, ".")
	if len(sVers) != 3 {
		return fmt.Errorf("Version.from: invalid ver[%v]", ver)
	}

	i64, err := strconv.Atoi(sVers[VersionMain])
	if err != nil {
		return fmt.Errorf("Version.from: invalid ver[%v]", ver)
	}
	v.vers[VersionMain] = i64

	i64, err = strconv.Atoi(sVers[VersionSub])
	if err != nil {
		return fmt.Errorf("Version.from: invalid ver[%v]", ver)
	}
	v.vers[VersionSub] = i64

	i64, err = strconv.Atoi(sVers[VersionTotal])
	if err != nil {
		return fmt.Errorf("Version.from: invalid ver[%v]", ver)
	}
	v.vers[VersionTotal] = i64

	return nil
}

// Compare 比较2个版本号
// 0    v = other
// -1   v < other
// 1    v > other
func (v *Version) Compare(other *Version) int {
	// 比较主版本号
	if v.vers[VersionMain] < other.vers[VersionMain] {
		return -1
	} else if v.vers[VersionMain] > other.vers[VersionMain] {
		return 1
	}

	// 比较子版本号
	if v.vers[VersionSub] < other.vers[VersionSub] {
		return -1
	} else if v.vers[VersionSub] > other.vers[VersionSub] {
		return 1
	}

	// 比较累计版本号
	if v.vers[VersionTotal] < other.vers[VersionTotal] {
		return -1
	} else if v.vers[VersionTotal] > other.vers[VersionTotal] {
		return 1
	}

	return 0
}

func (v *Version) String() string {
	return fmt.Sprintf("%v.%v.%v", v.vers[VersionMain], v.vers[VersionSub], v.vers[VersionTotal])
}

// New 将字符串版本号转为程序用的版本号
func New(ver string) (v *Version, err error) {
	v = &Version{}
	if err = v.from(ver); err != nil {
		return nil, err
	}

	return v, nil
}
