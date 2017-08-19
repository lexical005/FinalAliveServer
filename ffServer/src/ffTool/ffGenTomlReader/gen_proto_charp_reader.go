package main

import (
	"ffCommon/util"
	"fmt"
	"path/filepath"
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

        public static void Read(byte[] stream)
        {
            {FileName} = NConfig.{FileName}.Parser.ParseFrom(stream);
            Trans();
        }
    }
`

var fmtCSharpMemberReaderStruct = `
        public static {FileName}.Types.{ProtoVarType} {ProtoVarName}
        {
            get
            {
                return {FileName}.{ProtoVarName};
            }
        }
`

var fmtCSharpMemberReaderList = `
        public static Google.Protobuf.Collections.RepeatedField<{FileName}.Types.{ProtoVarType}> {ProtoVarName}
        {
            get
            {
                return {FileName}.{ProtoVarName};
            }
        }
`

var fmtCSharpMemberReaderMap = `
        public static Dictionary<{MapKeyType}, {FileName}.Types.{ProtoVarType}> {ProtoVarName}
        {
            get;
            private set;
        }
`

var fmtCSharpMapTrans = `
            {ProtoVarName} = new Dictionary<{MapKeyType}, {FileName}.Types.{ProtoVarType}>({FileName}.{ProtoVarName}Value.Count);
            for (int i = 0; i < {FileName}.{ProtoVarName}Value.Count; ++i)
            {
                {ProtoVarName}[{FileName}.{ProtoVarName}Key[i]] = {FileName}.{ProtoVarName}Value[i];
            }
`

var fmtReaderManager = `using System.Collections.Generic;

namespace NConfig
{
    public static class ConfigReaderManager
    {
        public delegate void ConfigReaderStream(System.IO.Stream stream);
        public delegate void ConfigReaderBuffer(byte[] stream);
        public class Reader
        {
            public string config;
            public ConfigReaderStream readerStream;
            public ConfigReaderBuffer readerBuffer;
        }
        public static readonly List<Reader> AllReader;
        public static Reader LanguageReader
        {
            get;
            private set;
        }

        static ConfigReaderManager()
        {
            AllReader = new List<Reader>()
            {{AllReader}
            };

            LanguageReader = new Reader()
            {
                config = "Language",
                readerStream = NConfig.LanguageReader.Read,
                readerBuffer = NConfig.LanguageReader.Read,
            };
        }
    }
}
`

var fmtReaderManagerOneReader = `
                new Reader()
                {
                    config = "{FileName}",
                    readerStream = NConfig.{FileName}Reader.Read,
                    readerBuffer = NConfig.{FileName}Reader.Read,
                },`

var regexpMainClass = regexp.MustCompile("\n  public sealed partial class ([\\w]+) : pb::IMessage")
var regexpSubClass = regexp.MustCompile("\n      public sealed partial class ([\\w]+) : pb::IMessage")

// excel文件对应的主类里的成员变量, 只有2种可能性:数组,实例.
// 数组情况下: 如果数组类型是基础类型, 则一定时字典的key数组, 因为工作簿数据数组的格式, 不是这样的
var regexpMainClassRepeatedFieldMapKey = regexp.MustCompile("\n    public pbc::RepeatedField<([\\w]+)> ([\\w]+)")

// excel文件对应的主类里的成员变量, 只有2种可能性:数组,实例.
// 工作簿数据数组的情况下, 要结合主类的前一字段的类型, 才能决定自身是独立的工作簿数据数组, 还是字典value数组
var regexpMainClassRepeatedFieldValue = regexp.MustCompile("\n    public pbc::RepeatedField<global::NConfig.([\\w]+).Types.([\\w]+)> ([\\w]+)")

// excel文件对应的主类里的成员变量, 只有2种可能性:数组,实例.
// 主类里的工作簿数据实例
var regexpMainClassStructValue = regexp.MustCompile("\n    public global::NConfig.([\\w]+).Types.([\\w]+) ([\\w]+)")

// 子类里的map成员
var regexpSubClassMapField = regexp.MustCompile("\n        public pbc::MapField<([\\w]+), ([\\w]+)> ([\\w]+)")

// 配置文件主类
type mainClassInfo struct {
	start int    // 本类型的开始
	end   int    // 下一类型的开始
	name  string // 类名

	subClass           []string                            // 子类型
	subClassMapMembers map[string][]*subClassMapMemberInfo // 子类的map成员

	member []*mainClassMemberInfo // 成员
}

// 主类成员信息
type mainClassMemberInfo struct {
	start        int    // 开始
	protoVarName string // 变量名称

	filedType string // 作为主类的成员时, 是什么类型. mapKey/mapValue/list/struct

	protoVarType string // 自身类型
}

// 子类map成员信息
type subClassMapMemberInfo struct {
	start        int    // 开始
	protoVarName string // 变量名称

	protoVarKeyType   string // 自身key类型
	protoVarValueType string // 自身value类型
}

func organizeClass(content string) (allMainClassInfo []*mainClassInfo) {
	var result [][]int

	// 类型
	{
		subClassMapFields := regexpSubClassMapField.FindAllStringSubmatchIndex(content, -1)

		// 获得excel文件对应的主类
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

				subClass:           make([]string, 0, 4),
				subClassMapMembers: make(map[string][]*subClassMapMemberInfo, 4),

				member: make([]*mainClassMemberInfo, 0, 4),
			}
			allMainClassInfo = append(allMainClassInfo, data)
		}

		// 主类里定义的子类
		result = regexpSubClass.FindAllStringSubmatchIndex(content, -1)
		result = result[1:]
		for i, one := range result {
			subClassStartIndex, name := one[0], content[one[2]:one[3]]

			for _, mainClassInfo := range allMainClassInfo {
				if subClassStartIndex > mainClassInfo.start && subClassStartIndex < mainClassInfo.end {
					// 子类属于该主类
					mainClassInfo.subClass = append(mainClassInfo.subClass, name)

					subClassEndIndex := 0
					if i+1 < len(result) {
						subClassEndIndex = result[i+1][0] - 1
					} else {
						subClassEndIndex = mainClassInfo.end - 1
					}

					// 子类里的map成员
					tmp := make([]*subClassMapMemberInfo, 0, 1)
					for _, one := range subClassMapFields {
						if one[0] > subClassStartIndex {
							if one[1] < subClassEndIndex {
								t := &subClassMapMemberInfo{
									start:        one[0],
									protoVarName: content[one[6]:one[7]],

									protoVarKeyType:   content[one[2]:one[3]],
									protoVarValueType: content[one[4]:one[5]],
								}
								tmp = append(tmp, t)
							} else {
								break
							}
						}
					}

					mainClassInfo.subClassMapMembers[name[len("St"):]] = tmp

					break
				}
			}
		}
	}

	// 从文本内解析出主类的成员信息, 并按出现顺序排序
	var allMembers map[int]*mainClassMemberInfo
	var allMembersKey []int
	{
		result1 := regexpMainClassRepeatedFieldMapKey.FindAllStringSubmatchIndex(content, -1)
		result2 := regexpMainClassRepeatedFieldValue.FindAllStringSubmatchIndex(content, -1)
		result3 := regexpMainClassStructValue.FindAllStringSubmatchIndex(content, -1)
		allMembers = make(map[int]*mainClassMemberInfo, len(result1)+len(result2)+len(result3))
		allMembersKey = make([]int, 0, len(result1)+len(result2)+len(result3))

		for _, one := range result1 {
			allMembers[one[0]] = &mainClassMemberInfo{
				start:        one[0],
				protoVarName: content[one[4]:one[5]],

				filedType: "mapKey",

				protoVarType: content[one[2]:one[3]],
			}
		}

		for _, one := range result2 {
			allMembers[one[0]] = &mainClassMemberInfo{
				start:        one[0],
				protoVarName: content[one[6]:one[7]],

				filedType: "list", // 暂时认为是list

				protoVarType: content[one[4]:one[5]],
			}
		}

		for _, one := range result3 {
			allMembers[one[0]] = &mainClassMemberInfo{
				start:        one[0],
				protoVarName: content[one[6]:one[7]],

				filedType: "struct",

				protoVarType: content[one[4]:one[5]],
			}
		}

		for key := range allMembers {
			allMembersKey = append(allMembersKey, key)
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

func genOutput(readerDir string, allMainClassInfo []*mainClassInfo) {
	// 输出
	for _, mainClassInfo := range allMainClassInfo {
		fmt.Printf("%v:%v:%v\n", mainClassInfo.name, mainClassInfo.start, mainClassInfo.end)
		fmt.Printf("%#v\n", mainClassInfo.subClass)
		for _, member := range mainClassInfo.member {
			fmt.Printf("%#v\n", member)
		}
	}

	resultManger := ""
	resultMangerAllReader := ""

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
				strMember = strings.Replace(strMember, "{ProtoVarName}", member.protoVarName, -1)
				AllMember += strMember

			} else if member.filedType == "mapKey" {
				i++
				memberValue := mainClassInfo.member[i]

				memberValue.filedType = "mapValue"

				member.protoVarName = member.protoVarName[:len(member.protoVarName)-len("Key")]
				memberValue.protoVarName = memberValue.protoVarName[:len(memberValue.protoVarName)-len("Value")]

				strMember := strings.Replace(fmtCSharpMemberReaderMap, "{FileName}", mainClassInfo.name, -1)
				strMember = strings.Replace(strMember, "{MapKeyType}", member.protoVarType, -1)
				strMember = strings.Replace(strMember, "{ProtoVarType}", memberValue.protoVarType, -1)
				strMember = strings.Replace(strMember, "{ProtoVarName}", memberValue.protoVarName, -1)
				AllMember += strMember

				mainClassMemberTrans := strings.Replace(fmtCSharpMapTrans, "{FileName}", mainClassInfo.name, -1)
				mainClassMemberTrans = strings.Replace(mainClassMemberTrans, "{MapKeyType}", member.protoVarType, -1)
				mainClassMemberTrans = strings.Replace(mainClassMemberTrans, "{ProtoVarType}", memberValue.protoVarType, -1)
				mainClassMemberTrans = strings.Replace(mainClassMemberTrans, "{ProtoVarName}", memberValue.protoVarName, -1)
				MapTrans += mainClassMemberTrans

			} else if member.filedType == "list" {
				strMember := strings.Replace(fmtCSharpMemberReaderList, "{FileName}", mainClassInfo.name, -1)
				strMember = strings.Replace(strMember, "{ProtoVarType}", member.protoVarType, -1)
				strMember = strings.Replace(strMember, "{ProtoVarName}", member.protoVarName, -1)
				AllMember += strMember
			}
		}

		// 子类map成员转换
		var fmtCSharpSubClassMapMemberTrans = `
		    for (int i = 0; i < {FileName}.{SheetName}{SheetType}.Count; ++i)
		    {
		        for (int j = 0; j < {FileName}.{SheetName}{SheetType}[i].{MemberName}Value.Count; ++j)
		        {
		            {FileName}.{SheetName}{SheetType}[i].{MemberName}.Add(
						{FileName}.{SheetName}{SheetType}[i].{MemberName}Key[j],
						{FileName}.{SheetName}{SheetType}[i].{MemberName}Value[j]
					);
		        }
		    }`

		for subClassName, tmp1 := range mainClassInfo.subClassMapMembers {
			for _, tmp2 := range tmp1 {
				for _, member := range mainClassInfo.member {
					if member.filedType != "mapKey" && member.protoVarName == subClassName {
						s := fmtCSharpSubClassMapMemberTrans
						s = strings.Replace(s, "{FileName}", mainClassInfo.name, -1)
						s = strings.Replace(s, "{SheetName}", subClassName, -1)
						s = strings.Replace(s, "{MemberName}", tmp2.protoVarName, -1)

						if member.filedType == "mapValue" {
							s = strings.Replace(s, "{SheetType}", "Value", -1)
						} else {
							s = strings.Replace(s, "{SheetType}", "", -1)
						}

						MapTrans += s
					}
				}
			}
		}

		oneReader := strings.Replace(fmtCSharpMainClassReader, "{FileName}", mainClassInfo.name, -1)
		oneReader = strings.Replace(oneReader, "{AllMember}", AllMember, -1)
		oneReader = strings.Replace(oneReader, "{MapTrans}", MapTrans, -1)
		AllReader += oneReader

		resultMangerAllReader += strings.Replace(fmtReaderManagerOneReader, "{FileName}", mainClassInfo.name, -1)
	}
	result = strings.Replace(fmtCSharpReader, "{AllReader}", AllReader, -1)

	util.WriteFile(filepath.Join(readerDir, "ConfigReader.cs"), []byte(result))

	resultManger = strings.Replace(fmtReaderManager, "{AllReader}", resultMangerAllReader, -1)
	util.WriteFile(filepath.Join(readerDir, "ConfigReaderManager.cs"), []byte(resultManger))
}

// 在proto-csharp代码的基础上, 封装读取字节流转换为proto结构体实例的代码
func genProtoCSharpReaderCode(readerDir string, protoCSharpCodePath string) {
	// 读取文件内容
	_, content := fileContent(protoCSharpCodePath)

	genOutput(readerDir, organizeClass(content))
}
