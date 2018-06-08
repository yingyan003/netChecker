#基础镜像
FROM img.reg.3g:15000/ubuntu-base:v3

MAINTAINER 247059543@qq.com

RUN mkdir -p /go/

ADD ./ /go/

WORKDIR /go/

#程序运行入口。项目编译后会生成名为netChecker的二进制文件
ENTRYPOINT ["/go/netChecker"]
