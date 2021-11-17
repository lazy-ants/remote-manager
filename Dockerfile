FROM php:8.0-cli-alpine

RUN apk add --no-cache \
      openssh-client \
        ca-certificates \
        bash \
        git

RUN docker-php-ext-install pcntl

RUN curl -s https://getcomposer.org/composer-stable.phar > /usr/local/bin/composer \
    && chmod a+x /usr/local/bin/composer

COPY docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

WORKDIR /usr/src/remote-manager

ENTRYPOINT ["docker-entrypoint.sh"]
