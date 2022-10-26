import '@pongle/styles/global.css';
import '../components/nav-bar';
import '../components/games-list';
import '../css/style.css';

import { PongleEvents } from '@pongle/events';

function loadElement(id: string): HTMLElement {
    const element = document.getElementById(id)
    if (element == null) {
        throw new Error(`Failed to find element with id: ${id}`);
    }

    return element
}

function loadTemplate(id: string): HTMLTemplateElement {
    const element = document.getElementById(id)
    if (element == null) {
        throw new Error(`Failed to find element with id: ${id}`);
    }

    return element as HTMLTemplateElement;
}

class IndexPage {
    playerName: HTMLElement;
    gameSection: HTMLElement;
    navBar: HTMLElement;
    newGameSection: HTMLElement;
    challengers: HTMLElement;

    newGameCallback: (id: string)=>void;

    constructor() {
        this.playerName = loadElement('player-name');
        this.gameSection = loadElement('game-section');
        this.navBar = loadElement('nav-bar');
        this.newGameSection = loadElement('new-game-section');
        this.challengers = loadElement('challengers');

        this.newGameCallback = (id: string) => {};
    }

    setPlayerName(name: string) {
        this.playerName.innerText = name;
        this.gameSection.hidden = false;
    }

    setGamesList(games: string[]) {
        // const newGame = this.templateNew.content.cloneNode(true) as DocumentFragment;
        // const newGameButton = newGame.querySelector('input');
        // if (newGameButton == null) {
        //     throw new Error('failed to find new game button in template');
        // }

        // newGameButton.onclick = () => this.showNewGamePopup(true);
        
        // this.gamesList.appendChild(newGame);
    }

    setChallengers(friends: {[key: string]: string}) {
        // while (this.challengers.firstChild) {
        //     this.challengers.removeChild(this.challengers.firstChild);
        // }

        // Object.keys(friends).forEach(id => {
        //     const node = this.templateChallenger.content.cloneNode(true) as DocumentFragment;
        //     const button = node.querySelector('input');
        //     if (button == null) {
        //         throw new Error('failed to find challenger button in template');
        //     }

        //     button.value = friends[id];
        //     button.onclick = () => this.newGameCallback(id);
            
        //     this.challengers.appendChild(node);
        // });
    }

    showNewGamePopup(show: boolean) {
        this.newGameSection.hidden = !show;
    }

    onLogout(callback: () => void) {
        this.navBar.addEventListener('logout-event', callback);
    }

    onNewGame(callback: (id: string) => void) {
        this.newGameCallback = callback;
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
        this.page.onNewGame((id: string) => this.events.send("draughts.new-game", { opponent: id }));
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
        this.page.setGamesList([]);
        this.page.setChallengers(this.friends);
    }

    async logout() {
        await this.events.logout();
        this.events.login();
    }
}

const page = new IndexPage();
const indexEvents = new IndexEvents(page);

indexEvents.start();