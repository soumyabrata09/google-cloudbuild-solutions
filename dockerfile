# FROM alpine
# COPY main.go /
# RUN chmod +x main.go
# CMD ["./main.go"]
FROM alpine
# VOLUME /tmp
COPY main.go /
EXPOSE 8001/tcp
ENTRYPOINT ["./main.go"]