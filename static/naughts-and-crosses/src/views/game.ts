import '@pongle/styles/global.css';
import '@pongle/components/nav-bar';
import '../controls/game-board';

import { html, css, LitElement } from 'lit';
import { customElement, state } from 'lit/decorators.js';
import { convert, GameData } from '../services/game';
import { PongleEvents } from '@pongle/events';

const GAME_DATA_KEY = "game-data";

@customElement('game-view')
export class GameView extends LitElement {
  static styles = css`
    section {
      box-sizing: border-box;
      width: 100%;
      padding: 1em;
      display: flex;
      justify-content: center;
    }
  `;

  private events: PongleEvents;
  private storage: Storage;

  @state()
  private userName: string;

  @state()
  private game: GameData;

  @state()
  private gameId: string;

  @state()
  private marks: string
  
  @state()
  private players: {[key: string]: string};

  constructor() {
    super();

    this.events = new PongleEvents("ponglehub.co.uk");
    this.storage = window.localStorage;

    const params = new URLSearchParams(window.location.search);
    this.gameId = params.get("id");

    const data = this.storage.getItem(GAME_DATA_KEY);
    if (data) {
      const parsed = JSON.parse(data);

      if (this.gameId !== parsed.id) {
        this.storage.removeItem(GAME_DATA_KEY);
        return;
      }

      this.userName = parsed.username;
      this.players = parsed.players;
      this.game = parsed.game;
      this.marks = parsed.marks;
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

    if (!this.storage.getItem(GAME_DATA_KEY)) {
      this.list(this.gameId);
    }
  }

  private listen(type: string, data: any) {
    switch(type) {
      case "auth.whoami.response":
        if (this.userName && this.userName !== data) {
          this.storage.clear();
          this.list(this.gameId);
        }

        this.userName = data;
        break;
      case "auth.list-friends.response":
        this.players = data;
        break;
      case "naughts-and-crosses.load-game.response":
        this.game = convert(data.game);
        this.marks = data.marks;
        break;
      default:
        console.error(`Unrecognised response type from server: ${type}`);
        return;
    }

    this.save();
  }

  private save() {
    this.storage.setItem(GAME_DATA_KEY, JSON.stringify({
      username: this.userName,
      players: this.players,
      game: this.game,
      marks: this.marks,
      id: this.gameId,
    }));
  }

  private list(id: string) {
    this.events.send("auth.list-friends", null);
    this.events.send("naughts-and-crosses.load-game", {id});
  }

  private async logOut() {
    await this.events.logout();
  }

  render() {
    return html`
      <nav-bar .loading="${false}" .authorised="${true}" @logout-event="${this.logOut}"></nav-bar>
      <section>
        <game-board .marks="${this.marks}"></game-board>
      </section>
    `;
  }
}
