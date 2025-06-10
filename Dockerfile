FROM harbor.newegg.org/base/alpine:3

LABEL maintainer="skyler.w.yang"

RUN mkdir -p /workspace/config && mkdir -p /workspace/data  && mkdir -p /workspace/certs

COPY ../backend/config/*.yaml /workspace/config

COPY ../front/projects/web/dist /workspace/web

COPY ../backend/humpback /workspace/

WORKDIR /workspace

CMD ["./humpback"]
