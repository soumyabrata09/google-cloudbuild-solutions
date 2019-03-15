FROM alpine
# RUN chmod  *.go
COPY main.go /
CMD ["/main.go"]