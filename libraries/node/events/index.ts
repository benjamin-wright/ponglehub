import { Auth } from './src/auth';
import { Events } from './src/events';

async function sleep(milliseconds: number) {
    return new Promise(resolve => setTimeout(resolve, milliseconds));
}

export class PongleEvents {
    private events: Events;
    private auth: Auth;

    constructor(host: string) {
        this.auth = new Auth(host);
        this.events = new Events(host);
    }

    async start(callback: (type: string, data: any)=>void, closed: ()=>void) {
        let attempt = 0;
        while(attempt < 3) {
            try {
                await this.events.start(callback, closed);
                return;
            } catch (err) {
                console.warn("Failed to connect to websocket server", err);
            }

            await sleep(1000 * attempt);
            attempt++;
        }

        try {
            await this.auth.logOut();
        } catch (err) {
            console.warn("failed to log out", err);
        }
         
        this.auth.logIn();
    }

    send(type: string, data: any) {
        this.events.send(type, data);
    }

    stop() {
        this.events.stop();
    }

    login() {
        this.auth.logIn();
    }

    async logout() {
        await this.auth.logOut();
    }
}
