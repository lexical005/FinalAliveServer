# -*- coding: utf-8 -*-
import os
import sys

SERVICE_LANGUAGE_GO = 'Go'
SERVICE_LANGUAGE_SUPPORT = [SERVICE_LANGUAGE_GO, ]


USAGE_PROMPT = 'make-gPRC.py service_type service_language(%s)' % (SERVICE_LANGUAGE_SUPPORT,)

# 检查是否包含指定的服务
def check_service_hold(row, service_type):
	startIndex = row.find('[')
	endIndex = row.find(']')
	if startIndex != -1 and endIndex != -1:
		all_service_type = row[startIndex + 1:endIndex]
		for one in all_service_type.split(' '):
			if one == service_type:
				return True
	return False


# 读取原始rpc定义文件
def read_rpc_define(service_type):
	f = open("rpc-define.proto", "r")

	SERVICE_START_START, SERVICE_START_END, SERVICE_END = '// ServiceStart[', ']\n', '// ServiceEnd'

	service_content = ''
	rowStatus = 'ServiceNone'	# ServiceStart
	rowHold = True
	while True:
		row = f.readline()
		if not row:
			break

		if rowStatus == 'ServiceNone':
			if row.startswith(SERVICE_START_START) and row.endswith(SERVICE_START_END):
				rowStatus = 'ServiceStart'
				rowHold = check_service_hold(row, service_type)

		elif rowStatus == 'ServiceStart':
			if row.startswith(SERVICE_END):
				rowStatus = 'ServiceNone'

		if rowHold:
			service_content += row

	f.close()

	return service_content


# 创建目录
def create_path(dir_name):
	cur_dir = os.path.split(os.path.realpath(__file__))[0]
	dir_path = os.path.join(cur_dir, dir_name)
	if not os.path.exists(dir_path):
		os.mkdir(dir_path)


# 写文件
def write_proto_file(rpc_content, package_name, file_name):
	syntax = 'syntax = "proto3";\n\n'
	package = 'package %s;\n\n\n\n' % (package_name,)
	file_content = syntax + package + rpc_content

	f = open(file_name, "w")
	f.write(file_content)
	f.close()


# 删文件
def delete_proto_file(file_name):
	if os.path.exists(file_name):
		os.remove(file_name)


# go语言rpc生成
def gen_service_language_go(service_type):
	rpc_content = read_rpc_define(service_type)
	package_name = 'ffRPC%s' % (service_type,)
	file_name = 'ffRPC_%s.proto' % (service_type,)

	create_path(package_name)
	
	cur_dir = os.getcwd()
	os.chdir(package_name)


	write_proto_file(rpc_content, package_name, file_name)
	os.system('%s --go_out=plugins=grpc:. %s' % (os.path.join(cur_dir, 'protoc-3.2.0.exe'), file_name))
	os.system('go install ffRPC/%s' % (package_name,))
	delete_proto_file(file_name)


	os.chdir(cur_dir)


# 出错时，提示错误信息，并停止执行
def go_error(error_msg, callback):
	print error_msg
	os.system('color 4a')
	os.system('pause')
	if callback != None:
		callback()
	exit(1)


# 语言的导出实现
SERVICE_LANGUAGE_GEN_FUNC = {
	SERVICE_LANGUAGE_GO:gen_service_language_go,
}


if __name__ == '__main__':
	service_type, service_language = None, None

	if len(sys.argv) > 2:
		service_type = sys.argv[1]
		service_language = sys.argv[2]

	if service_language not in SERVICE_LANGUAGE_SUPPORT:
		go_error('invalid service_language[%s]. usage:\n%s' % (service_language, USAGE_PROMPT), None)

	SERVICE_LANGUAGE_GEN_FUNC[service_language](service_type)
	print('gen %s-%s success\n' % (service_type, service_language))
