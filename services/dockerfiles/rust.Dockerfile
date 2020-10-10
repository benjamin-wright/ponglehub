FROM alpine

ARG EXE_NAME

COPY ${EXE_NAME} /rust_binary

RUN ls -lah /

ENTRYPOINT [ "/rust_binary" ]