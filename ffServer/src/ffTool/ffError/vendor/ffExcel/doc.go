// Package ffExcel 定义了一套机制，实现将excel配置表转换为程序使用的配置。
// excel表配置格式：
//  sheet名称为readme的，程序将忽略，供表格说明使用
//  sheet的前四行，分别为本列配置说明，字段key，字段类型，客户端/服务端配置
//  sheet名称，必须以_list(sheet内容将被导出为数组)，_map(sheet内容是字典，键的名称固定为Key，类型只能为int或string)，_struct(sheet内容只有一行，用以表述配置的值)
package ffExcel
