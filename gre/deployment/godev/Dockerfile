FROM golang:1.6-onbuild AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
#
## 移动到工作目录：/build
#WORKDIR /build
#
## 复制项目中的 go.mod 和 go.sum文件并下载依赖信息
#COPY go.mod .
#COPY go.sum .
#RUN go mod download
#
## 将代码复制到容器中
#COPY . .
#
## 将我们的代码编译成二进制可执行文件 tracer
#RUN go build -o tracer .
#
## 创建一个小镜像
#FROM scratch
#
## 静态文件
#COPY ./log /log
#COPY ./config /config
## 当出现static文件或者template文件的时候需要配置
## COPY ./static /static
## COPY ./template /template
#
## 从builder镜像中把/dist/tracer 拷贝到当前目录
#COPY --from=builder /build/tracer /
#
## 需要运行的命令
#ENTRYPOINT ["/tracer", "config/config.ini"]
