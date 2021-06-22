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

router.get('/users/:id', (ctx, next) => {
    const id = ctx.params.id;
    const user = users.find(user => user.id === id);

    if (user) {
        ctx.status = 200;
        ctx.body = {
            name: user.name,
            email: user.email,
            verified: user.verified
        };
    } else {
        ctx.status = 404;
    }
});

router.put('/users/:id', (ctx, next) => {
    const id = ctx.params.id;
    const user = users.find(user => user.id === id);
    const body = ctx.request.body;

    if (user) {
        if (body.name) {
            console.info(`${id}: name ${user.name} => ${body.name}`);
            user.name = body.name;
        }

        if (body.email) {
            console.info(`${id}: email ${user.email} => ${body.email}`);
            user.email = body.email;
        }

        if (body.password) {
            console.info(`${id}: password ${user.password} => ${body.password}`);
            user.password = body.password;
        }

        if (body.verified) {
            console.info(`${id}: verified ${user.verified} => ${body.verified}`);
            user.verified = body.verified;
        }

        ctx.status = 202;
        ctx.body = {
            name: user.name,
            email: user.email,
            verified: user.verified
        };
    } else {
        ctx.status = 404;
    }
});

router.post('/users', (ctx, next) => {
    const id = uuid.v4();
    const { name, email, password, verified } = ctx.request.body;

    users.push({
        id,
        name,
        email,
        password,
        verified: !!verified
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
