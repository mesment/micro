#使用alpine作为基础镜像
FROM alpine:latest
#在容器根目录下创建app目录
RUN mkdir /app
#将工作目录设置为/app
WORKDIR /app
#将微服务的服务端运行文件拷贝到/app下
ADD consignment-service /app/consignment-service

#运行服务端程序
CMD ["./consignment-service"]