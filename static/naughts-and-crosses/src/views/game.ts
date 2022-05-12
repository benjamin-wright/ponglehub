import '@pongle/styles/global.css';
import '@pongle/components/nav-bar';
import '../controls/game-board';

import { html, css, LitElement } from 'lit';
import { customElement, state } from 'lit/decorators.js';
import { convert, GameData } from '../services/game';
import { PongleEvents } from '@pongle/events';

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

  @state()
  private userName: string;
  
  @state()
  private userId: string;

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

    const params = new URLSearchParams(window.location.search);
    this.gameId = params.get("id");
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

    this.list(this.gameId);
  }

  private listen(type: string, data: any) {
    switch(type) {
      case "auth.whoami.response":
        if (this.userName && this.userName !== data.display) {
          this.list(this.gameId);
        }

        this.userName = data.display;
        this.userId = data.id;
        break;
      case "auth.list-friends.response":
        this.players = data;
        break;
      case "naughts-and-crosses.load-game.response":
        this.game = convert(data.game);
        this.marks = data.marks;
        break;
      case "naughts-and-crosses.mark.response":
        this.game = convert(data.game);
        this.marks = data.marks;
        break;
      case "naughts-and-crosses.load-game.rejection.response":
        window.location.href = "../naughts-and-crosses"
        break;
      default:
        console.error(`Unrecognised response type from server: ${type}`);
        return;
    }
  }

  private list(id: string) {
    this.events.send("auth.list-friends", null);
    this.events.send("naughts-and-crosses.load-game", {id});
  }

  private select(index: number) {
    this.events.send("naughts-and-crosses.mark", {game: this.gameId, position: index})
  }

  private async logOut() {
    await this.events.logout();
    this.events.login();
  }

  private getTurn(): number {
    if(!this.game) {
      return 0;
    }

    return this.game.turn;
  }

  private getPlayer(): number {
    if(!this.game) {
      return 0;
    }

    return this.game.player1 === this.userId ? 0 : 1;
  }

  render() {
    return html`
      <nav-bar .loading="${false}" .authorised="${true}" @logout-event="${this.logOut}"></nav-bar>
      <section>
        <game-board .turn="${this.getTurn()}" .player="${this.getPlayer()}" .marks="${this.marks}" @select="${(event: CustomEvent<number>) => this.select(event.detail)}"></game-board>
      </section>
    `;
  }
}
