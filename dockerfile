#!/bin/bash
FROM alpine
# RUN chmod 777 *.go
COPY main.go /
CMD ["/main.go"]