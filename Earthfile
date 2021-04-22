FROM scratch

repos:
  LOCALLY
  RUN ./infra/repos.sh

init:
  LOCALLY
  RUN ./infra/start.sh
  RUN helm dep update helm/tests

clean:
  LOCALLY
  RUN ./infra/stop.sh
