FROM golang as builder
COPY src /src
WORKDIR /src
RUN set -ex \
    && go env -w GOPROXY=https://goproxy.cn \
    && CGO_ENABLED=0 go build


FROM chromedp/headless-shell
COPY --from=builder /src/glados-checkin /usr/local/bin/checker
RUN set -ex \
    && mkdir /src \
    && chmod +x /usr/local/bin/checker \
    && sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list \
    && sed -i  's|security.debian.org/debian-security|mirrors.ustc.edu.cn/debian-security|g' /etc/apt/sources.list \
    && apt-get update \
    && apt-get install ca-certificates binutils curl libc6 -y
