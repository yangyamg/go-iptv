#!/bin/sh
set -e

log() {
    echo "$(date '+%Y/%m/%d %H:%M:%S') $*"
}

WATCH_DIR="/tmp/updata"

# 启动 license
log "启动授权服务..."
/app/license > /config/license.log 2>&1 &
LICENSE_PID=$!
log "授权服务已启动，PID=$LICENSE_PID"

# 等待授权服务可用
wait_license() {
    TIMEOUT=60
    START_TIME=$(date +%s)
    log "等待授权服务启动..."
    while true; do
        nc -z 127.0.0.1 81 >/dev/null 2>&1 && break
        [ "$TIMEOUT" -gt 0 ] && [ $(( $(date +%s) - START_TIME )) -ge $TIMEOUT ] && \
            { log "等待授权服务超时 ${TIMEOUT}s"; return 1; }
        sleep 2
    done
    log "授权服务可用"
    return 0
}

wait_license

# 启动 IPTV
start_iptv() {
    log "启动 IPTV..."
    /app/iptv &
    IPTV_PID=$!
    log "IPTV 已启动，PID=$IPTV_PID"
}

start_iptv

# 更新处理函数（按 updata.sh -> license -> iptv 顺序）
update_handler() {
    [ -d "$WATCH_DIR" ] || return

    # 1. updata.sh
    if [ -f "$WATCH_DIR/updata.sh" ]; then
        log "检测到 updata.sh，执行更新脚本..."
        chmod +x "$WATCH_DIR/updata.sh"
        sh "$WATCH_DIR/updata.sh"
        rm -f "$WATCH_DIR/updata.sh"
    fi

    # 2. license
    if [ -f "$WATCH_DIR/license" ]; then
        log "检测到 license 更新，替换并重启授权服务..."
        cp "$WATCH_DIR/license" /app/license
        kill $LICENSE_PID 2>/dev/null || true
        /app/license > /config/license.log 2>&1 &
        LICENSE_PID=$!
        rm -f "$WATCH_DIR/license"

        # 等待 license 可用
        wait_license || { log "license 启动失败，跳过 IPTV 更新"; return; }
    fi

    # 3. iptv
    if [ -f "$WATCH_DIR/iptv" ]; then
        log "检测到 IPTV 更新，替换并重启 IPTV..."
        cp "$WATCH_DIR/iptv" /app/iptv
        chmod +x /app/iptv
        kill $IPTV_PID 2>/dev/null || true
        start_iptv
        rm -f "$WATCH_DIR/iptv"
    fi
}

# 捕获单一信号 SIGUSR1 触发更新
trap 'update_handler' SIGUSR1

log "等待更新信号..."
while true; do
    sleep 60 & wait $!
done
