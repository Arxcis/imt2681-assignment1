FROM alpine:latest

MAINTAINER Jonas Solsvik <jonas.solsvik@gmail.com>

WORKDIR "/opt"

ADD .docker_build/imt2681-assignment1 /opt/bin/imt2681-assignment1

CMD ["/opt/bin/imt2681-assignment1"]