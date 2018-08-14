FROM scratch

EXPOSE 8080

ADD user-service_linux /srv/user-service

CMD ["/srv/user-service"]