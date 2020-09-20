FROM golang:1.15.2 AS gobuilder
COPY upload.go /workspace/
WORKDIR /workspace
RUN CGO_ENABLED=0 go build -o app

FROM alpine:latest
LABEL purpose=file-upload port=8199
COPY --from=gobuilder /workspace/app /opt/app
EXPOSE 8199
ENTRYPOINT ["/opt/app"]


