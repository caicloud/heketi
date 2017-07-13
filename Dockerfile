FROM ubuntu:16.04
MAINTAINER Delweng Zheng <delweng@gmail.com>

RUN apt-get update && apt-get install -y cron

RUN rm -rf /etc/cron.*

ADD ./heketi /usr/bin/heketi
ADD ./client/cli/go/heketi-cli /usr/bin/heketi-cli
ADD ./heketi-start.sh /usr/bin/heketi-start.sh
VOLUME /etc/heketi

RUN mkdir /var/lib/heketi
VOLUME /var/lib/heketi

ENTRYPOINT ["/usr/bin/heketi-start.sh"]
EXPOSE 8080
