FROM phusion/baseimage

# Use baseimage's init process
CMD ["/sbin/my_init"]

ADD build/server /opt/meathooks/server
ADD scripts/meathooks-server.sh /etc/service/meathooks-server/run
EXPOSE 80