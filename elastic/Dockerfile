FROM docker.elastic.co/elasticsearch/elasticsearch:6.5.4
RUN set -x && curl -L -o /usr/local/bin/jq https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64 \
    && sha256sum -c - <<<'af986793a515d500ab2d35f8d2aecd656e764504b789b66d7e1a0b727a124c44 /usr/local/bin/jq' \
    && chmod +x /usr/local/bin/jq
COPY elasticsearch.yml /usr/share/elasticsearch/config/
COPY setup-snapshot.bash /usr/local/bin/
COPY backup.bash /usr/local/bin/
