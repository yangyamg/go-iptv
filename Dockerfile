# ================================
# Builder - 只负责准备不同架构的可执行文件
# ================================
FROM --platform=$BUILDPLATFORM alpine:latest AS builder

ARG TARGETARCH
WORKDIR /app

# 只复制实际需要用于选择架构的文件
COPY license_all/ license_all/
COPY iptv_dist/ iptv_dist/
COPY start_all/ start_all/

RUN case "$TARGETARCH" in \
      amd64) cp ./license_all/license_amd64 ./license \
          && cp ./iptv_dist/iptv_amd64 ./iptv \
          && cp ./start_all/start_amd64 ./start ;; \
      arm64) cp ./license_all/license_arm64 ./license \
          && cp ./iptv_dist/iptv_arm64 ./iptv \
          && cp ./start_all/start_arm64 ./start ;; \
      arm) cp ./license_all/license_arm ./license \
          && cp ./iptv_dist/iptv_arm ./iptv \
          && cp ./start_all/start_arm ./start ;; \
      *) echo "未知架构: $TARGETARCH" && exit 1 ;; \
    esac

RUN chmod +x iptv license start


# ================================
# Final Image
# ================================
FROM alpine:3.22

ENV TZ=Asia/Shanghai
WORKDIR /app
VOLUME /config
EXPOSE 80 8080

# 安装依赖，减少层
RUN apk add --no-cache \
      openjdk8 \
      bash \
      curl \
      tzdata \
      sqlite \
      ffmpeg \
    && cp /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone

# 静态资源（不常改变）尽早 COPY，这样缓存会利用到
COPY client /client
COPY apktool/ /usr/bin/
COPY static/ /app/static
COPY config.yml /app/config.yml
COPY README.md  /app/README.md
COPY logo/ /app/logo
COPY dictionary.txt /app/dictionary.txt

RUN chmod -R 777 /usr/bin/apktool*

# 其他不常改变的文件
COPY database/ /app/database
COPY alias.json /app/alias.json
COPY ChangeLog.md /app/ChangeLog.md
COPY Version /app/Version
COPY license_all/Version_lic /app/Version_lic

# 最后 COPY 程序部分（经常更新）
COPY --from=builder /app/iptv .
COPY --from=builder /app/license .
COPY --from=builder /app/start .

CMD ["./start"]
