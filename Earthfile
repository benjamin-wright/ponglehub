FROM scratch

repos:
  LOCALLY
  RUN ./infra/repos.sh

repos-stop:
  LOCALLY
  RUN ./infra/repos-stop.sh

infra:
  LOCALLY
  RUN ./infra/start.sh
  RUN helm dep update helm/tests

infra-stop:
  LOCALLY
  RUN ./infra/stop.sh

start:
  BUILD +repos
  BUILD +infra

stop:
  BUILD +infra-stop

clean:
  BUILD +infra-stop
  BUILD +repos-stop

