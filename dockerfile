FROM alpine
RUN chmod +x main.go
COPY main.go /
CMD ["/main.go"]