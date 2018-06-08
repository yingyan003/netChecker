FROM img.reg.3g:15000/ubuntu-base:v3

MAINTAINER 247059543@qq.com

RUN mkdir -p /go/

#文件存在时，是否覆盖？
ADD ./ /go/

WORKDIR /go/

ENTRYPOINT ["/go/ping"]
