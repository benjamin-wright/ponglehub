import '@pongle/styles/global.css';
import '@pongle/components/nav-bar';
import '@pongle/panels/popup-panel';

import {html, css, LitElement} from 'lit';
import {customElement, state} from 'lit/decorators.js';
import {Auth} from '@pongle/auth';
import {Game} from '../services/game';

@customElement('index-view')
export class IndexView extends LitElement {
  static styles = css``;

  private auth: Auth;
  private game: Game;

  @state()
  private userName: string;
  
  @state()
  private newGame: boolean;

  @state()
  private games: string[];
  
  @state()
  private players: string[];

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
          <li>Id: ${game}</li>
        `)}
      </ul>
    `;
  }

  private newGamePopup() {
    if (!this.newGame) {
      return null;
    }

    return html`
      <popup-panel>
        <div>
          <p>Who would you like to challenge?</p>
          <button @click="${() => this.newGame = false}">OK</button>
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
      <p>Players: ${JSON.stringify(this.players)}</p>
    `;
  }
}
