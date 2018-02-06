FROM debian:latest

RUN apt-get update \
    && apt-get install -y curl mysql-client unixodbc libpq5

RUN SPHINX='sphinx-3.0.1-7fec4f6-linux-amd64' \
    && curl http://sphinxsearch.com/files/${SPHINX}.tar.gz | tar -xzvf - \
    && mv /${SPHINX} /sphinx

WORKDIR /sphinx
COPY sphinx.conf etc/sphinx.conf

VOLUME [ "/var/data" ]
VOLUME [ "/var/log" ]
EXPOSE 9312
EXPOSE 9306
CMD [ "/sphinx/bin/searchd", "--nodetach", "--config", "/sphinx/etc/sphinx.conf" ]
