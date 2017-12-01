FROM centurylink/ca-certs:latest

EXPOSE 8080

ADD bin/app /app

CMD ["/app"]
