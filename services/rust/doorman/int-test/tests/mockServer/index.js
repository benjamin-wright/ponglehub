const Koa = require('koa');
const cors = require('koa-cors');

class MockServer {
    constructor(origin) {
        this.app = new Koa();
        this.app.use(cors({
            origin,
            credentials: true
        }));

        this.app.use(async ctx => {
            let notFound = true;
            this.routes.forEach(route => {
                if (route.path === ctx.path && route.host === ctx.host) {
                    ctx.status = route.status || 200;
                    ctx.body = route.body;
                    notFound = false;
                }

                this.calls.push({
                    method: ctx.method,
                    host: ctx.host,
                    path: ctx.path
                });
            });

            if (notFound) {
                ctx.status = 404;
                ctx.body = 'No expectation set';
            }
        });

        this.routes = [];
        this.calls = [];
    }

    start(port) {
        return new Promise((resolve, reject) => {
            this.server = this.app.listen(port | 80, err => {
                if (err) {
                    return reject(err);
                }

                resolve();
            });
        });
    }

    addRoute({ path, host, body, status }) {
        this.routes.push({
            path,
            host,
            body,
            status
        });
    }

    stop() {
        return new Promise((resolve, reject) => {
            this.server.close(err => {
                if (err) {
                    return reject(err);
                }

                resolve();
            });
        });
    }
}

module.exports = MockServer;
