FROM scratch

EXPORT 8080

ADD user-service_linux /srv/user-service

RUN ["/srv/user-service"]