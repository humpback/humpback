FROM alpine:latest

LABEL maintainer="skyler.w.yang"

RUN mkdir -p /workspace/config && mkdir -p /workspace/data

COPY ./backend/config/*.yaml /workspace/config

COPY ./front/projects/web/dist /workspace/web

COPY ./backend/humpback /workspace/

WORKDIR /workspace

CMD ["./humpback"]
