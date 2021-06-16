const Koa = require('koa');
const Router = require('@koa/router');

const app = new Koa();
const router = new Router();

const users = [];

router.get('/users', (ctx, next) => {
  ctx.body = users;
});

app
  .use(router.routes())
  .use(router.allowedMethods());

const server = app.listen(80, '0.0.0.0', () => {
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