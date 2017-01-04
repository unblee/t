#!/bin/sh
set -eu

# This script is used only in containers.

dest_dir=/go/src/${NAME}/dest

[ ! -e ${dest_dir} ] && mkdir ${dest_dir} && chmod 777 ${dest_dir}

for os in windows linux darwin; do
  for arch in amd64 386; do
    {
      out_dir=${dest_dir}/${NAME}-${VERSION}-${os}-${arch}
      [ ! -e ${out_dir} ] && mkdir ${out_dir} && chmod 777 ${out_dir}
      [ ${os} = "windows" ] && suffix=.exe || suffix=""
      build_cmd=GOOS="${os} GOARCH=${arch} go build -a -tags netgo -installsuffix netgo -ldflags \"${LDFLAGS}\" -o ${out_dir}/${NAME}${suffix}"
      echo ${build_cmd}
      eval ${build_cmd}
    } &
  done
done
wait

echo
echo cross-build finished!
echo