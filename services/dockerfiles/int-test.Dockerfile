# syntax = docker/dockerfile:1.0-experimental
FROM node:10.15.3-alpine

WORKDIR /usr/tests/src

COPY package.json package-lock.json ./

RUN --mount=type=secret,id=npmrc,dst=/root/.npmrc \
    --mount=type=secret,id=cert,dst=/root/ca \
    npm config set -g cafile /root/ca \
    && npm ci

COPY . .

RUN npm run lint

ENTRYPOINT [ "npm" ]
CMD [ "run", "int-test" ]

