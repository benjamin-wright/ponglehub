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
            try {
                console.info("getting events socket...");
                await this.events.start(callback, closed);
                console.info("done!");
                return;
            } catch (error) {
                console.info("failed to get events socket", error)
            }

            try {
                console.info("checking logged in...");
                if (await this.auth.load() == null) {
                    try {
                        console.info('logging out...');
                        await this.auth.logOut();
                        console.info('done!');
                    } finally {
                        console.info('navigating to login...');
                        this.auth.logIn();
                        return;
                    }
                }
            } catch (error) {
                console.info("failed to load user info", error);
            }
            
            console.info("Next round!");
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
