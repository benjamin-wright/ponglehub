import '@pongle/styles/global.css';
import '@pongle/components/nav-bar';
import '@pongle/panels/popup-panel';
import '../controls/game-summary';

import {html, css, LitElement} from 'lit';
import {customElement, state} from 'lit/decorators.js';
import {Auth} from '@pongle/auth';
import {Game} from '../services/game';

@customElement('index-view')
export class IndexView extends LitElement {
  static styles = css`
    ul {
      list-style: none;
    }
  `;

  private auth: Auth;
  private game: Game;

  @state()
  private userName: string;
  
  @state()
  private newGame: boolean;

  @state()
  private games: string[];
  
  @state()
  private players: {[key: string]: string};

  constructor() {
    super();
    
    this.auth = new Auth(window.localStorage);
    this.game = new Game();
  }

  connectedCallback() {
    super.connectedCallback();

    this.auth.init()
      .then(data => this.userName = data.name)
      .catch(() => this.auth.logIn());

    this.game.start(() => {
      this.games = this.game.games();
      this.players = this.game.players();
      this.requestUpdate();
    });
  }

  disconnectedCallback() {
    super.disconnectedCallback();

    this.game.stop();
  }

  private async logOut() {
    await this.auth.logOut();
    window.location.href="http://games.ponglehub.co.uk";
  }

  private listGames() {
    if (!this.games) {
      return html`<p>loading...</p>`;
    }

    return html`
      <p>games:</p>
      <ul>
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
    this.game.newGame(opponent)
  }

  private newGamePopup() {
    if (!this.newGame) {
      return null;
    }

    return html`
      <popup-panel>
        <div>
          <p>Who would you like to challenge?</p>
          <ul>
            ${Object.keys(this.players).map(key => html`
              <li @click="${() => this.requestNewGame(key)}">${this.players[key]}</li>
            `)}
          </ul>
        </div>
      </popup-panel>
    `
  }

  render() {
    return html`
      <nav-bar .loading="${false}" .authorised="${true}" @logout-event="${this.logOut}"></nav-bar>
      <h1>Hi ${this.userName}! Lets play Naughts and Crosses!</h1>
      <button @click="${() => this.newGame = true}">New Game</button>
      ${this.listGames()}
      ${this.newGamePopup()}
    `;
  }
}
