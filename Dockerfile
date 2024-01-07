FROM chromedp/headless-shell
COPY glados-checkin /usr/local/bin
COPY user.yaml /opt/
RUN set -ex \
    && sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list \
    && sed -i  's|security.debian.org/debian-security|mirrors.ustc.edu.cn/debian-security|g' /etc/apt/sources.list \
    && apt-get update \
    && apt-get install ca-certificates libc6 -y
