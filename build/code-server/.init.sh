#!/bin/bash 


function graceful_exit () {
    echo "receive SIGTERM, exiting..."
    pid=$(ps -ef | grep code-server | awk '{print $2}')
    kill -SIGTERM "$pid"

    exit 0
}

trap graceful_exit SIGTERM

if [ -z "$OPEN_DIR" ]; then
	OPEN_DIR="$USER_WORKSPACE" 
fi

while :;do

    # 创建用户工作空间目录
    if [ ! -d "$USER_WORKSPACE" ]; then
        echo "create $USER_WORKSPACE"
        mkdir -p $USER_WORKSPACE
    fi

    # 第一次启动工作空间,拷贝code-server的数据
    if [ ! -f "/user_data/.local/share/.first_start" ]; then
        echo "copy code-server"
        if [ ! -d "/user_data/.local/share" ]; then
            mkdir -p /user_data/.local/share
        fi
        cp -r /root/.local/share/code-server-bak /user_data/.local/share/code-server
        touch /user_data/.local/share/.first_start
    fi

    # 启动code-server
    node_id=$(ps aux | grep -E "/.workspace/code-server/lib/node /.workspace/code-server --port 9999 --host 0.0.0.0 --auth none" | grep -v grep)
    if [ -z "$node_id" ]; then
        nohup code-server --port 9999 --host 0.0.0.0 \
        --auth none --disable-update-check  --locale zh-cn \
        --user-data-dir /user_data/.local/share/code-server \
        --extensions-dir /user_data/.local/share/code-server/extensions \
        --open "$OPEN_DIR" &

        echo "start code-server success"
    fi

    sleep 3
done
