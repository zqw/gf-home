FROM alpine:3.8

LABEL maintainer="john@johng.cn"

###############################################################################
#                                INSTALLATION
###############################################################################

# 使用国内alpine源
RUN echo http://mirrors.ustc.edu.cn/alpine/v3.8/main/ > /etc/apk/repositories

# 添加HTTPS根证书，设置系统时区 - +8时区
RUN apk update && apk add tzdata ca-certificates bash
RUN rm -rf /etc/localtime && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" > /etc/timezone

# 添加应用二进制文件
ADD ./main /bin/main
RUN chmod +x /bin/main

###############################################################################
#                                   START
###############################################################################

CMD main