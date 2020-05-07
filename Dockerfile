FROM php:7.4-cli-alpine

RUN \
    apk update && apk upgrade && \
    apk add --no-cache \
    openssh-client \
    ca-certificates \
    bash \
    git

RUN curl -s https://getcomposer.org/composer-stable.phar > /usr/local/bin/composer \
    && chmod a+x /usr/local/bin/composer

RUN rm -rf /var/cache/apk/*

WORKDIR /usr/src/remote-manager

ENTRYPOINT ["./docker-entrypoint.sh"]
