#!/bin/bash

RUN_NAME="xxx_bin"
export PATH=$PATH:/bin:/usr/local/bin:/usr/local/sbin:/usr/bin:/usr/sbin
current=$(cd $(dirname $0) && pwd )
cd $current && go mod tidy && mkdir -p output/bin output/conf output/log
find conf/ -type f ! -name "*_local.*" | xargs -I{} cp {} output/conf/
go build -o output/bin/${RUN_NAME}

echo "bin/$RUN_NAME -env=dev -run_dir=./ > ./out.log 2>&1  &" >> output/bootstrap_boe.sh
chmod +x output/bootstrap_boe.sh
echo "bin/$RUN_NAME -env=pro -run_dir=./ > ./out.log 2>&1  &" >> output/bootstrap_pro.sh
chmod +x output/bootstrap_pro.sh