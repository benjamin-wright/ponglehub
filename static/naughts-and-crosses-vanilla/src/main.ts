import './style.css'
import { PongleEvents } from '@pongle/events';

const app = document.querySelector<HTMLDivElement>('#app')!
const events = new PongleEvents("ponglehub.co.uk");

events.login();

app.innerHTML = `
  <h1>Hello Vite!</h1>
  <a href="https://vitejs.dev/guide/features.html" target="_blank">Documentation</a>
`
