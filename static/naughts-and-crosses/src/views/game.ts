import '@pongle/styles/global.css';
import '@pongle/components/nav-bar';
import '@pongle/panels/popup-panel';
import '../controls/game-summary';
import '../controls/new-game';

import {html, css, LitElement} from 'lit';
import {customElement} from 'lit/decorators.js';
import {Game} from '../services/game';

@customElement('game-view')
export class GameView extends LitElement {
  static styles = css``;

  private game: Game;

  private gameId: string;

  constructor() {
    super();
    
    this.game = new Game("ponglehub.co.uk", window.localStorage);
    this.game.addListener(this.listen.bind(this));

    const params = new URLSearchParams(window.location.search);
    this.gameId = params.get("id");
  }

  connectedCallback() {
    super.connectedCallback();
    this.game.start().then(() => {
        this.game.loadGame(this.gameId);
    });
  }

  disconnectedCallback() {
    super.disconnectedCallback();
    this.game.stop();
  }

  private listen(property: string) {
    switch (property) {
      default:
        console.info(`Ignoring unknown property: ${property}`);
        break;
    }
  }

  private async logOut() {
    await this.game.logout();
  }

  render() {
    return html`
      <nav-bar .loading="${false}" .authorised="${true}" @logout-event="${this.logOut}"></nav-bar>
      <section>
        <h1>This is your game!</h1>
      </section>
    `;
  }
}
