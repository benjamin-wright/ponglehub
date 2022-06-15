import '@pongle/styles/global.css';
import '../components/nav-bar';
import '../css/style.css';

import { PongleEvents } from '@pongle/events';

class IndexPage {
    playerName: HTMLElement | null;
    gameSection: HTMLElement | null;
    navBar: HTMLElement | null;

    constructor() {
        this.playerName = document.getElementById('player-name');
        this.gameSection = document.getElementById('game-section');
        this.navBar = document.getElementById('nav-bar');
    }

    setPlayerName(name: string) {
        if (this.playerName) {
            this.playerName.innerText = name;
        }

        if (this.gameSection) {
            this.gameSection.hidden = false;
        }
    }

    onLogout(callback: () => void) {
        if (this.navBar) {
            this.navBar.addEventListener('logout-event', callback);
        }
    }
}

class IndexEvents {
    username: string
    friends: {[key: string]: string};
    games: any[]
    page: IndexPage
    events: PongleEvents

    constructor(page: IndexPage) {
        this.username = "";
        this.friends = {};
        this.games = [];
        this.page = page;
        this.events = new PongleEvents("ponglehub.co.uk");

        this.page.onLogout(() => this.logout());
    }

    async start() {
        await this.events.start(this.listen.bind(this), this.start.bind(this));
        this.events.send("auth.list-friends", {});
        this.events.send("draughts.list-games", {});
    }

    listen(type: string, data: any) {
        switch (type) {
            case "auth.whoami.response":
                this.username = data.display;
                break;
            case "auth.list-friends.response":
                this.friends = data;
                break;
            case "draughts.list-games.response":
                this.games = data.games;
                break;
            default:
                console.info(`What's this? ${type}`);
                break;
        }

        this.render();
    }

    render() {
        this.page.setPlayerName(this.username);
    }

    async logout() {
        await this.events.logout();
        this.events.login();
    }
}

const page = new IndexPage();
const indexEvents = new IndexEvents(page);

indexEvents.start();