import '@pongle/styles/global.css';
import '@pongle/components/nav-bar';
import '@pongle/panels/popup-panel';
import '../controls/list-games';
import '../controls/new-game-popup';

import { html, css, LitElement } from 'lit';
import { customElement, state } from 'lit/decorators.js';
import { convert, GameData } from '../services/game';
import { PongleEvents } from '@pongle/events';

const INDEX_DATA_KEY = "index-data";

@customElement('index-view')
export class IndexView extends LitElement {
  static styles = css`
    h1 {
      width: 100%;
      text-align: center;
    }

    em {
      color: var(--default-foreground);
      font-style: normal;
      font-weight: bold;
      text-transform: capitalize;
    }

    section {
      padding: 1em;
    }
  `;

  private events: PongleEvents;
  private storage: Storage;

  @state()
  private userName: string;
  
  @state()
  private newGame: boolean;

  @state()
  private games: GameData[];
  
  @state()
  private players: {[key: string]: string};

  constructor() {
    super();

    this.events = new PongleEvents("ponglehub.co.uk");
    this.storage = window.localStorage;

    const data = this.storage.getItem(INDEX_DATA_KEY);
    if (data) {
      const parsed = JSON.parse(data);

      this.userName = parsed.username;
      this.players = parsed.players;
      this.games = parsed.games;
    }
  }

  connectedCallback() {
    super.connectedCallback();
    this.start();
  }

  disconnectedCallback() {
    super.disconnectedCallback();
    this.events.stop();
  }

  private async start() {
    await this.events.start(
      this.listen.bind(this),
      this.start.bind(this),
    );

    if (!this.storage.getItem(INDEX_DATA_KEY)) {
      this.list();
    }
  }

  private listen(type: string, data: any) {
    switch(type) {
      case "auth.whoami.response":
        if (this.userName && this.userName !== data) {
          this.storage.clear();
          this.list();
        }

        this.userName = data;
        break;
      case "auth.list-friends.response":
        this.players = data;
        break;
      case "naughts-and-crosses.list-games.response":
        this.games = data.games.map(convert);
        this.games = this.games.sort((a, b) => Date.parse(b.created) - Date.parse(a.created));
        break;
      case "naughts-and-crosses.new-game.response":
        this.games = this.games.slice();
        this.games.push(convert(data.game));
        this.games = this.games.sort((a, b) => Date.parse(b.created) - Date.parse(a.created));
        break;
      default:
        console.error(`Unrecognised response type from server: ${type}`);
        return;
    }

    this.save();
  }

  private save() {
    this.storage.setItem(INDEX_DATA_KEY, JSON.stringify({
      username: this.userName,
      players: this.players,
      games: this.games
    }));
  }

  private list() {
    this.events.send("auth.list-friends", null);
    this.events.send("naughts-and-crosses.list-games", null);
  }

  private async logOut() {
    await this.events.logout();
  }

  private requestNewGame(opponent: string) {
    this.newGame = false;
    this.events.send("naughts-and-crosses.new-game", {opponent});
  }

  render() {
    return html`
      <nav-bar .loading="${false}" .authorised="${true}" @logout-event="${this.logOut}"></nav-bar>
      <section>
        <h1>Hi <em>${this.userName}</em>! Lets play Naughts and Crosses!</h1>
        <list-games
          .games="${this.games}"
          .players="${this.players}"
          @new-game="${() => this.newGame = true}"
        ></list-games>
        <new-game-popup
          .display="${this.newGame}"
          .players="${this.players}"
          @cancel="${() => this.newGame = false}"
          @new-game="${(event: CustomEvent<string>) => this.requestNewGame(event.detail)}"
        ></new-game-popup>
      </section>
    `;
  }
}
