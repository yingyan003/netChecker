FROM img.reg.3g:15000/ubuntu-base:v3

MAINTAINER liyao.miao@yeepay.com

RUN mkdir -p /go/

#文件存在时，是否覆盖？
ADD ./ /go/

WORKDIR /go/

ENTRYPOINT ["/go/ping"]
