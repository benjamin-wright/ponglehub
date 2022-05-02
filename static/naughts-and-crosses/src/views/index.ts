import '@pongle/styles/global.css';
import '@pongle/components/nav-bar';
import '@pongle/panels/popup-panel';
import '../controls/game-summary';
import '../controls/new-game';

import {html, css, LitElement} from 'lit';
import {customElement, state} from 'lit/decorators.js';
import {Game, GameData} from '../services/game';

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

  private game: Game;

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
    
    this.game = new Game("ponglehub.co.uk", window.localStorage);
    this.game.addListener(this.listen.bind(this));
  }

  connectedCallback() {
    super.connectedCallback();
    this.game.start();
  }

  disconnectedCallback() {
    super.disconnectedCallback();
    this.game.stop();
  }

  private listen(property: string) {
    switch (property) {
      case "userName":
        this.userName = this.game.userName;
        this.requestUpdate("userName");
      case "players":
        this.players = this.game.players;
        this.requestUpdate("players");
      case "games":
        this.games = this.game.games;
        this.requestUpdate("games");
      default:
        console.info(`Ignoring unknown property: ${property}`);
    }
  }

  private async logOut() {
    await this.game.logout();
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
    this.game.newGame(opponent);
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
