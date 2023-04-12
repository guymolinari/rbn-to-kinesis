# quanta-kinesis-consumer within a Docker Container

#FROM frolvlad/alpine-glibc
FROM arm64v8/alpine

LABEL version="1.0.0" \
  maintainer="Guy Molinari" \
  maintainer.email="molinariguy@gmail.com" \
  description="Builds Generic Docker image for use in production deployment of Quanta Kinesis Consumer" \
  repository="https://github.com/disney/quanta"

# Update 
RUN apk update
RUN apk upgrade
RUN apk add --no-cache --update alpine-sdk

# clean apk cache
RUN rm -rf /var/cache/apk/*

COPY ./bin/rbn-to-kinesis /usr/bin/rbn-to-kinesis
COPY ./Docker/entrypoint.sh /usr/bin/entrypoint.sh
COPY ./callparser/cty.dat callparser/cty.dat
RUN chmod 755 /usr/bin/rbn-to-kinesis
RUN chmod 755 /usr/bin/entrypoint.sh

# Expose default port, but can take in docker run --expose flag as additive port to expose
# Port mapping of MySQL Proxy server. Default("127.0.0.1:4000")

#VOLUME /data

ENTRYPOINT ["/usr/bin/entrypoint.sh"]
