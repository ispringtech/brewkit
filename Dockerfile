FROM docker:20.10

COPY bin/brewkit /usr/local/bin/

ENTRYPOINT ["brewkit"]