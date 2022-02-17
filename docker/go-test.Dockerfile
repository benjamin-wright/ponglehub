FROM scratch
ARG EXECUTABLE
COPY ${EXECUTABLE} /text_exec
ENTRYPOINT [ "/text_exec" ]
CMD [ "-test.v" ]