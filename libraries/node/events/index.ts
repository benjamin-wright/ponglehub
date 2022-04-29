import { PongleStorage } from './src/storage';
import { Auth } from './src/auth';
import { Events } from './src/events';

async function sleep(milliseconds: number) {
    return new Promise(resolve => setTimeout(resolve, milliseconds));
}

export class PongleEvents {
    private storage: PongleStorage;
    private events: Events;
    private auth: Auth;

    constructor(host: string, storage: Storage) {
        this.auth = new Auth(host);
        this.events = new Events(host);
        this.storage = new PongleStorage(storage);
    }

    async start(callback: (type: string, data: any)=>void, closed: ()=>void) {
        if (!this.storage.isLoggedIn()) {
            this.auth.logIn();
        }

        while(true) {
            let attempts = 0;
            while (attempts < 3) {
                try {
                    await this.events.start(callback, closed);
                    return;
                } catch (error) {
                    attempts++;
                }
            }

            try {
                if (await this.auth.load() == null) {
                    try {
                        await this.auth.logOut();
                    } finally {
                        this.auth.logIn();
                        return;
                    }
                }
            } catch { }

            await sleep(5000);
        }
    }
    send(type: string, data: any) {
        this.events.send(type, data);
    }

    stop() {
        this.events.stop();
    }

    login() {
        this.storage.clear();
        this.auth.logIn();
    }

    async logout() {
        this.storage.clear();
        await this.auth.logOut();
    }
}
