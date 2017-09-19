FROM ubuntu:16.04
MAINTAINER Delweng Zheng <delweng@gmail.com>

RUN apt-get update \
    && apt-get install -y cron \
    && apt-get autoremove -y \
    && rm -rf /var/cache/apt/archives \
    && rm -rf /etc/cron.* \
    && mkdir /var/lib/heketi


ADD ./heketi /usr/bin/heketi
ADD ./client/cli/go/heketi-cli /usr/bin/heketi-cli
ADD ./heketi-start.sh /usr/bin/heketi-start.sh
ADD ./heketi-init.sh /usr/bin/heketi-init.sh
VOLUME [/etc/heketi, /var/lib/heketi]

ENTRYPOINT ["/usr/bin/heketi-start.sh"]
EXPOSE 8080
