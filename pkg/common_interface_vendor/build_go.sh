#!/bin/bash

CURRENT_PATH=$(cd `dirname $0` && pwd)

if [[ $# -lt 1 ]]; then
    echo "usage: $0 api/service/helloworld/v1"
    exit 1
fi

targetDirName='pb_interface'
# 根据 proto 文件，生成对应语言的代码
if [ -d "${CURRENT_PATH}/$1" ];then
  rm -rf "${CURRENT_PATH}/${targetDirName}/$1"
  mkdir -p "${CURRENT_PATH}/${targetDirName}/$1"
  cp -rf ${CURRENT_PATH}/$1/*.proto "${CURRENT_PATH}/${targetDirName}/$1/"
  cd "${CURRENT_PATH}/${targetDirName}/$1" || exit 1
  # shellcheck disable=SC2045
  for file in $(ls *.proto)
  do
    if [[ "${file}" = *model* ]];then
      # kratos proto client -p "${CURRENT_PATH}/proto" "${file}"
      protoc --proto_path=. \
	       --proto_path=${CURRENT_PATH}/third_party \
 	       --go_out=paths=source_relative:. \
 	       --go-http_out=paths=source_relative:. \
 	       --go-grpc_out=paths=source_relative:. \
           --go-errors_out=paths=source_relative:. \
 	       --openapi_out=naming=proto,paths=source_relative:. \
           "${file}"   # 输出路径
    fi
  done

  for file in $(ls *.proto)
  do
    if [[ "${file}" != *model* ]];then
      protoc --proto_path=. \
	       --proto_path=${CURRENT_PATH}/third_party \
 	       --go_out=paths=source_relative:. \
 	       --go-http_out=paths=source_relative:. \
 	       --go-grpc_out=paths=source_relative:. \
           --go-errors_out=paths=source_relative:. \
 	       --openapi_out=naming=proto,paths=source_relative:. \
           "${file}" # 输出路径
    fi
  done

  # 生成 svc 的 server 代码
  if [[ "$1" != *common* ]];then
    mkdir -p "service"
    rm -rf "service/"*.go
    # shellcheck disable=SC2045
    for file in $(ls *.proto)
    do
      kratos proto server "${file}" -t "service"
    done
  fi
  rm -rf *.proto
  # 修改 go 文件中的 json tag
  cd $CURRENT_PATH && protoc-go-inject-tag -input=./${targetDirName}/$1/*.pb.go

else
  echo "文件不存在"
  exit 1
fi
