#!/bin/sh
set -eu

# This script is used only in containers.

repo_root=/go/src/${NAME}
dest_dir=${repo_root}/dest
release_dir=${repo_root}/release

[ ! -e ${release_dir} ] && mkdir ${release_dir} && chmod 777 ${release_dir}

apk add --no-cache openssl curl zip git

echo ""
for target in `ls ${dest_dir}`
do
  {
    cd ${dest_dir}
    archive_cmd="zip -r ${release_dir}/${target}.zip ${target}"
    echo ${archive_cmd}
    eval ${archive_cmd}
  }
done

echo ""
if [ `uname -m` = "x86_64" ]; then
  arch=amd64
else
  arch=386
fi
wget -O /tmp/ghr.zip https://github.com/tcnksm/ghr/releases/download/v0.5.3/ghr_v0.5.3_linux_${arch}.zip
unzip -d /go/bin /tmp/ghr.zip
ghr=/go/bin/ghr
ghr_cmd="${ghr} -u ${USERNAME} -r ${REPO} ${VERSION} ${release_dir}"
echo ${ghr_cmd}
eval ${ghr_cmd}

echo ""
echo release finished!
echo ""