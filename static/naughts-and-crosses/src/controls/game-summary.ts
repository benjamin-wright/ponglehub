import '@pongle/styles/global.css';

import {html, css, LitElement} from 'lit';
import {customElement, property} from 'lit/decorators.js';

@customElement('game-summary')
export class GameSummary extends LitElement {
  static styles = css`
    div {
      margin: 1em;
      padding: 1em;
      background: var(--default-foreground);
      color: var(--default-background);
      border-radius: 1em;
    }

    em {
      color: var(--default-highlight);
      font-style: normal;
      font-weight: bold;
    }

    p {
      margin: 0;
      padding: 0;
    }
  `;

  @property({type: Object})
  private game: {id: string, player1: string, player2: string, turn: number, created: Date};

  @property({type: Object})
  private players: {[key: string]: string};

  render() {
    console.info(this.game);
    console.info(this.players);

    let player1 = this.players[this.game.player1] || "You";
    let player2 = this.players[this.game.player2] || "You";

    return html`
      <div>
        <p><em>${player1}</em> vs <em>${player2}</em></p>
        <p>Turn: ${this.game.turn}</p>
        <p>Created: ${this.game.created}</p>
      </div>
    `;
  }
}
