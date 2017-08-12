package main

import (
	"ffCommon/util"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var fmtCSharpReader = `using System.Collections.Generic;

namespace NConfig
{{AllReader}
}
`

var fmtCSharpMainClassReader = `
    public static class {FileName}Reader
    {
        private static {FileName} {FileName};
{AllMember}

        private static void Trans()
        {{MapTrans}
        }

        public static void Read(System.IO.Stream stream)
        {
            {FileName} = NConfig.{FileName}.Parser.ParseFrom(stream);
            Trans();
        }
	}
`

var fmtCSharpMemberReaderStruct = `
        public static {FileName}.Types.{ProtoVarType} {OriginalVarType}
        {
            get
            {
                return {FileName}.{OriginalVarType};
            }
		}
`

var fmtCSharpMemberReaderList = `
        public static Google.Protobuf.Collections.RepeatedField<{FileName}.Types.{ProtoVarType}> {OriginalVarType}
        {
            get
            {
                return {FileName}.{OriginalVarType};
            }
        }
`

var fmtCSharpMemberReaderMap = `
        public static Dictionary<int, {FileName}.Types.{ProtoVarType}> {OriginalVarType}
        {
            get;
            private set;
        }
`

var fmtCSharpMapTrans = `
            {OriginalVarType} = new Dictionary<{MapKeyType}, {FileName}.Types.{ProtoVarType}>({FileName}.{OriginalVarType}Value.Count);
            for (int i = 0; i < {FileName}.{OriginalVarType}Value.Count; ++i)
            {
                {OriginalVarType}[{FileName}.{OriginalVarType}Key[i]] = {FileName}.{OriginalVarType}Value[i];
            }
`

var regexpMainClass = regexp.MustCompile(`
  public sealed partial class ([\w]+) : pb::IMessage`)
var regexpSubClass = regexp.MustCompile(`
      public sealed partial class ([\w]+) : pb::IMessage`)
var regexpMainClassRepeatedFieldMapKey = regexp.MustCompile(`
    public pbc::RepeatedField<([\w]+)> ([\w]+)`)
var regexpMainClassRepeatedFieldValue = regexp.MustCompile(`
    public pbc::RepeatedField<global::NConfig.([\w]+).Types.([\w]+)> ([\w]+)`)
var regexpMainClassStructValue = regexp.MustCompile(`
    public global::NConfig.([\w]+).Types.([\w]+) ([\w]+)`)

// 配置文件主类
type mainClassInfo struct {
	start int    // 本类型的开始
	end   int    // 下一类型的开始
	name  string // 类名

	subClass []string // 子类型

	member []*mainClassMemberInfo // 成员
}

// 主类成员信息
type mainClassMemberInfo struct {
	start           int    // 开始
	originalVarType string // 变量名称

	filedType string // 作为主类的成员时, 是什么类型. map/list/struct

	protoVarType string // 自身类型
}

func organizeClass(content string) (allMainClassInfo []*mainClassInfo) {
	var result [][]int

	// 类型
	{
		// 主类
		result = regexpMainClass.FindAllStringSubmatchIndex(content, -1)
		result = result[1:]
		allMainClassInfo = make([]*mainClassInfo, 0, len(result)) // 忽略第一个 Grammar
		for count, one := range result {
			start := one[0]
			end := one[1]
			name := content[one[2]:one[3]]

			if count+1 < len(result) {
				end = result[count+1][0]
			} else {
				end = len(content)
			}

			data := &mainClassInfo{
				start: start,
				end:   end,
				name:  name,

				subClass: make([]string, 0, 4),
				member:   make([]*mainClassMemberInfo, 0, 4),
			}
			allMainClassInfo = append(allMainClassInfo, data)
		}

		// 主类里定义的子类
		result = regexpSubClass.FindAllStringSubmatchIndex(content, -1)
		for _, one := range result[1:] {
			start, end, name := one[0], one[1], content[one[2]:one[3]]
			for _, mainClassInfo := range allMainClassInfo {
				if start > mainClassInfo.start && end < mainClassInfo.end {
					// 子类属于该主类
					mainClassInfo.subClass = append(mainClassInfo.subClass, name)
					break
				}
			}
		}
	}

	// 从文本内解析出成员信息, 并按出现顺序排序
	var allMembers map[int]*mainClassMemberInfo
	var allMembersKey []int
	{
		result1 := regexpMainClassRepeatedFieldMapKey.FindAllStringSubmatchIndex(content, -1)
		result2 := regexpMainClassRepeatedFieldValue.FindAllStringSubmatchIndex(content, -1)
		result3 := regexpMainClassStructValue.FindAllStringSubmatchIndex(content, -1)
		allMembers = make(map[int]*mainClassMemberInfo, len(result1)+len(result2)+len(result3))
		allMembersKey = make([]int, len(result1)+len(result2)+len(result3))

		for _, one := range result1 {
			allMembers[one[0]] = &mainClassMemberInfo{
				start:           one[0],
				originalVarType: content[one[4]:one[5]],

				filedType: "map",

				protoVarType: content[one[2]:one[3]],
			}
		}

		for _, one := range result2 {
			allMembers[one[0]] = &mainClassMemberInfo{
				start:           one[0],
				originalVarType: content[one[6]:one[7]],

				filedType: "list", // 暂时认为是list

				protoVarType: content[one[4]:one[5]],
			}
		}

		for _, one := range result3 {
			allMembers[one[0]] = &mainClassMemberInfo{
				start:           one[0],
				originalVarType: content[one[6]:one[7]],

				filedType: "struct",

				protoVarType: content[one[4]:one[5]],
			}
		}

		index := 0
		for key := range allMembers {
			allMembersKey[index] = key
			index++
		}
		sort.Ints(allMembersKey)
	}

	// 将成员信息, 规整到主类内
	{
		for _, start := range allMembersKey {
			memeber := allMembers[start]
			for _, mainClassInfo := range allMainClassInfo {
				if start > mainClassInfo.start && start < mainClassInfo.end {
					mainClassInfo.member = append(mainClassInfo.member, memeber)
					break
				}
			}
		}
	}

	return
}

func genOutput(savePath string, allMainClassInfo []*mainClassInfo) {
	// 输出
	for _, mainClassInfo := range allMainClassInfo {
		fmt.Printf("%v:%v:%v\n", mainClassInfo.name, mainClassInfo.start, mainClassInfo.end)
		fmt.Printf("%#v\n", mainClassInfo.subClass)
		for _, member := range mainClassInfo.member {
			fmt.Printf("%#v\n", member)
		}
	}

	result := ""
	AllReader := ""
	for _, mainClassInfo := range allMainClassInfo {
		AllMember := ""
		MapTrans := ""

		for i := 0; i < len(mainClassInfo.member); i++ {
			member := mainClassInfo.member[i]
			if member.filedType == "struct" {

				strMember := strings.Replace(fmtCSharpMemberReaderStruct, "{FileName}", mainClassInfo.name, -1)
				strMember = strings.Replace(strMember, "{ProtoVarType}", member.protoVarType, -1)
				strMember = strings.Replace(strMember, "{OriginalVarType}", member.originalVarType, -1)
				AllMember += strMember

			} else if member.filedType == "map" {
				i++
				memberValue := mainClassInfo.member[i]

				member.originalVarType = member.originalVarType[:len(member.originalVarType)-len("Key")]
				memberValue.originalVarType = memberValue.originalVarType[:len(memberValue.originalVarType)-len("Value")]

				strMember := strings.Replace(fmtCSharpMemberReaderMap, "{FileName}", mainClassInfo.name, -1)
				strMember = strings.Replace(strMember, "{ProtoVarType}", memberValue.protoVarType, -1)
				strMember = strings.Replace(strMember, "{OriginalVarType}", memberValue.originalVarType, -1)
				AllMember += strMember

				memberTrans := strings.Replace(fmtCSharpMapTrans, "{FileName}", mainClassInfo.name, -1)
				memberTrans = strings.Replace(memberTrans, "{MapKeyType}", member.protoVarType, -1)
				memberTrans = strings.Replace(memberTrans, "{ProtoVarType}", memberValue.protoVarType, -1)
				memberTrans = strings.Replace(memberTrans, "{OriginalVarType}", memberValue.originalVarType, -1)
				MapTrans += memberTrans

			} else if member.filedType == "list" {

				strMember := strings.Replace(fmtCSharpMemberReaderList, "{FileName}", mainClassInfo.name, -1)
				strMember = strings.Replace(strMember, "{ProtoVarType}", member.protoVarType, -1)
				strMember = strings.Replace(strMember, "{OriginalVarType}", member.originalVarType, -1)
				AllMember += strMember
			}
		}

		oneReader := strings.Replace(fmtCSharpMainClassReader, "{FileName}", mainClassInfo.name, -1)
		oneReader = strings.Replace(oneReader, "{AllMember}", AllMember, -1)
		oneReader = strings.Replace(oneReader, "{MapTrans}", MapTrans, -1)
		AllReader += oneReader
	}
	result = strings.Replace(fmtCSharpReader, "{AllReader}", AllReader, -1)

	util.WriteFile(savePath, []byte(result))
}

// 在proto-csharp代码的基础上, 封装读取字节流转换为proto结构体实例的代码
func genProtoCSharpReaderCode(savePath string, protoCSharpCodePath string) {
	// 读取文件内容
	_, content := fileContent(protoCSharpCodePath)

	genOutput(savePath, organizeClass(content))
}
