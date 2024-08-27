FROM alpine
COPY image-tool /usr/local/bin/image-tool

ENTRYPOINT ["image-tool"]
