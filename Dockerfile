#syntax=docker/dockerfile:1

FROM tarantool/tarantool:latest
COPY dhands.lua /opt/tarantool/
WORKDIR /opt/tarantool
CMD ["tarantool"] ["/opt/tarantool/dhands.lua"]

