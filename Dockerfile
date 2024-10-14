# 使用镜像GO语言的镜像，对项目进行打包
FROM golang:1.23.2-alpine3.20 AS builder
# 定好工作目录
WORKDIR /app
# 将哪个路径下的文件进行打包，并且被复制到哪里
COPY . .
# 进入到容器后执行的命名
RUN go build -o main main.go

# RUN stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env . 
# 暴露那些端口
EXPOSE 8080
# 执行文件
CMD [ "/app/main" ]