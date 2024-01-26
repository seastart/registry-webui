FROM alpine:3

LABEL name="registry-webui"
LABEL description="docker registry webui in onefile"
LABEL maintainer="dev@seastart.cn"

###############################################################################
#                                INSTALLATION
###############################################################################
ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH
# RUN echo "I'm building for $TARGETPLATFORM $TARGETOS $TARGETARCH"

ARG TZ=Asia/Shanghai
# set timezone
RUN apk add --no-cache tzdata
ENV TZ ${TZ}
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
# set workdir
ENV WORKDIR /var/www/registry-webui

# add executable by target architecture
ADD --chmod=755 main_linux_${TARGETARCH} $WORKDIR/main
RUN mkdir ${WORKDIR}/config

###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR
ENTRYPOINT [ "./main" ]
CMD ["--config", "./config/default.yml"]