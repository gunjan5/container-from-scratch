FROM iron/go:1.6.1
WORKDIR /app


# copy binary into image
COPY app /app/


ENTRYPOINT ["./app"]
