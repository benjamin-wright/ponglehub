import '@pongle/styles/global.css';
import '@pongle/components/nav-bar';
import '@pongle/panels/popup-panel';
import '../controls/list-games';
import '../controls/new-game-popup';

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

    .challenger {
      border: 2px solid var(--default-foreground);
      border-radius: 1em;
      padding: 1em;
      text-transform: capitalize;
      color: var(--default-foreground);
      cursor: pointer;
      user-select: none;
    }
    
    .challenger:focus, .challenger:hover {
      border: 2px solid var(--default-highlight);
      background: var(--default-foreground);
      color: var(--default-background);
    }

    .cancel {
      display: flex;
      justify-content: right;
    }

    .cancel button {
      background: none;
      border: none;
      color: red;
      cursor: pointer;
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
        break;
      case "players":
        this.players = this.game.players;
        this.requestUpdate("players");
        break;
      case "games":
        this.games = this.game.games;
        this.requestUpdate("games");
        break;
      default:
        console.info(`Ignoring unknown property: ${property}`);
        break;
    }
  }

  private async logOut() {
    await this.game.logout();
  }

  private requestNewGame(opponent: string) {
    this.newGame = false;
    this.game.newGame(opponent);
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
