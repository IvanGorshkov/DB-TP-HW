# FROM tarantool/tarantool:latest
# COPY configs/lua/app.lua /opt/tarantool
# CMD ["tarantool", "/opt/tarantool/app.lua"]

FROM postgres:13.2
RUN localedef -i ru_RU -c -f UTF-8 -A /usr/share/locale/locale.alias ru_RU.UTF-8
ENV LANG ru_RU.UTF-8
ENV LC_ALL ru_RU.UTF-8
