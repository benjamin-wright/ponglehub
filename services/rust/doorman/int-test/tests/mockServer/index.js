const Koa = require('koa');
const cors = require('koa-cors');

class MockServer {
    constructor() {
        this.app = new Koa();
        this.app.use(cors());

        this.app.use(async ctx => {
            console.log(`Recieved ${ctx.method} request to ${ctx.path}`);

            let notFound = true;
            this.routes.forEach(route => {
                if (route.path === ctx.path) {
                    console.log(ctx.cookies.get('pongle_auth'));

                    ctx.status = route.status || 200;
                    ctx.body = route.body;
                    notFound = false;
                }
            });

            if (notFound) {
                ctx.status = 404;
                ctx.body = 'No expectation set';
            }
        });

        this.routes = [];
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

    addRoute(path, body, status) {
        this.routes.push({
            path,
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
