import '@pongle/styles/global.css';
import '@pongle/components/nav-bar';
import '@pongle/panels/popup-panel';
import '../controls/game-summary';
import '../controls/new-game';

import {html, css, LitElement} from 'lit';
import {customElement, state} from 'lit/decorators.js';
import {PongleEvents} from '@pongle/events';
import {newGameEvent} from '../services/events';

function convert(g: any): any {
  return {
    created: g.Created,
    id: g.ID,
    player1: g.Player1,
    player2: g.Player2,
    turn: g.Turn
  }
}

@customElement('index-view')
export class IndexView extends LitElement {
  static styles = css`
    h1 {
      width: 100%;
      text-align: center;
    }

    ul {
      list-style: none;
      display: flex;
      flex-wrap: wrap;
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

  @state()
  private userName: string;
  
  @state()
  private newGame: boolean;

  @state()
  private games: {
    created: string,
    id: string,
    player1: string,
    player2: string,
    turn: number
  }[];
  
  @state()
  private players: {[key: string]: string};

  constructor() {
    super();
    
    this.events = new PongleEvents("ponglehub.co.uk", window.localStorage);
  }

  connectedCallback() {
    super.connectedCallback();
    this.listen();    
  }

  disconnectedCallback() {
    super.disconnectedCallback();
    this.events.stop();
  }

  private listen() {
    this.events.start(
      this.eventHandler.bind(this),
      this.listen.bind(this),
    )
    .then(() => {
      this.events.send("auth.list-friends", null);
      this.events.send("naughts-and-crosses.list-games", null);
    })
    .catch(error => {
      console.error("Error connecting to websocket", error);
    });
  }

  private eventHandler(type: string, data: any) {
    console.info(`Event: ${type}`);
    switch(type) {
      case "naughts-and-crosses.list-games.response":
        this.games = data.games.map(convert);
        this.games = this.games.sort((a, b) => Date.parse(b.created) - Date.parse(a.created))
        break;
      case "auth.list-friends.response":
        this.players = data;
        break;
      case "auth.whoami.response":
        this.userName = data;
        break;
      case "naughts-and-crosses.new-game.response":
        this.games.push(convert(data.game));
        this.games = this.games.sort((a, b) => Date.parse(b.created) - Date.parse(a.created))
        this.requestUpdate("games");
        break;
      default:
        console.error(`Unrecognised response type from server: ${type}`);
        break;
    }
  }

  private async logOut() {
    await this.events.logout();
    window.location.href="http://games.ponglehub.co.uk";
  }

  private listGames() {
    if (!this.games) {
      return html`<p>loading...</p>`;
    }

    return html`
      <ul>
        <li><new-game @click="${() => this.newGame = true}"/></li>
        ${this.games.map(game => html`
          <li>
            <game-summary .game="${game}" .players="${this.players}"></game-summary>
          </li>
        `)}
      </ul>
    `;
  }

  private requestNewGame(opponent: string) {
    this.newGame = false;
    this.events.send("naughts-and-crosses.new-game", newGameEvent(opponent));
  }

  private newGamePopup() {
    if (!this.newGame) {
      return null;
    }

    return html`
      <popup-panel>
        <p>Who would you like to challenge?</p>
        <ul>
          ${Object.keys(this.players).map(key => html`
            <li @click="${() => this.requestNewGame(key)}">${this.players[key]}</li>
          `)}
        </ul>
      </popup-panel>
    `
  }

  render() {
    return html`
      <nav-bar .loading="${false}" .authorised="${true}" @logout-event="${this.logOut}"></nav-bar>
      <section>
        <h1>Hi <em>${this.userName}</em>! Lets play Naughts and Crosses!</h1>
        ${this.listGames()}
        ${this.newGamePopup()}
      </section>
    `;
  }
}
