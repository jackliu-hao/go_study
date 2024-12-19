# 基础镜像
FROM ubuntu:20.04
# 把编译后的打包进来这个镜像,放到工作目录 /app
COPY webook /app/webhook
WORKDIR /app

RUN ["chmod", "+x", "webook"]
ENTRYPOINT ["./webook"]
EXPOSE 8081