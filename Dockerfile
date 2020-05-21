FROM golang:1.14-alpine as build
RUN apk --no-cache add tzdata ca-certificates
WORKDIR /app
ADD . .
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -mod=vendor -a -ldflags "-s" -installsuffix cgo -o bin/app src/main/*.go

FROM scratch as final
COPY --from=build /app/bin/app .
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8080
CMD ["/app"]
