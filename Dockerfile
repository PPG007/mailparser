FROM phusion/baseimage:master

CMD ["/sbin/my_init"]

EXPOSE 8080

WORKDIR /root

COPY ./dist /root/dist
COPY ./mailparser /root/mailparser
COPY ./start.sh /etc/my_init.d/start.sh
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
