#!/bin/bash
FROM alpine
EXPOSE 8001
COPY main.go /
CMD ["/main.go"]