const Koa = require('koa');
const logger = require('koa-logger');
const bodyParser = require('koa-bodyparser');
const Router = require('@koa/router');
const uuid = require('uuid');

const app = new Koa();
const router = new Router();

const users = [];

router.get('/users', (ctx, next) => {
    ctx.body = users;
});

router.post('/users', (ctx, next) => {
    const id = uuid.v4();
    const { name, email, password } = ctx.request.body;

    users.push({
        id,
        name,
        email,
        password
    });

    ctx.body = { id };
    ctx.status = 202;
});

router.post('/reset', (ctx, next) => {
    while (users.length > 0) {
        users.pop();
    }

    ctx.status = 200;
});

app
    .use(logger())
    .use(bodyParser())
    .use(router.routes())
    .use(router.allowedMethods());

const server = app.listen(80, () => {
    console.info('Server running on 0.0.0.0:80');
});

function shutdown() {
    logger.warn('Received kill signal (SIGTERM / SIGINT), shutting down...');

    server.close(() => {
        logger.info('Closed out remaining connections');
    });
}

process.on('SIGTERM', shutdown);
process.on('SIGINT', shutdown);
