import '@pongle/styles/global.css';

import {html, css, LitElement} from 'lit';
import {customElement, property} from 'lit/decorators.js';
import { timeSince } from '../services/utils';

@customElement('game-summary')
export class GameSummary extends LitElement {
  static styles = css`
    .panel {
      user-select: none;
      cursor: pointer;

      margin: 1em;
      padding: 1em;
      background: var(--default-foreground);
      color: var(--default-background);
      border: 2px solid var(--default-foreground);
      border-radius: 1em;
    }

    .panel:hover, .panel:focus {
      border: 2px solid var(--default-highlight);
    }

    em {
      color: var(--default-highlight);
      font-style: normal;
      font-weight: bold;
      text-transform: capitalize;
    }

    p {
      margin: 0;
      padding: 0;
    }

    .center {
      display: flex;
      justify-content: center;
    }

    .split {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: space-between;
    }
  `;

  @property({type: Object})
  private game: {id: string, player1: string, player2: string, turn: number, created: string};

  @property({type: Object})
  private players: {[key: string]: string};

  render() {
    let player1 = this.players[this.game.player1] || "You";
    let player2 = this.players[this.game.player2] || "You";
    let theirTurn = this.game.turn === 0 && this.players[this.game.player1];
    let elapsed = timeSince(this.game.created);

    return html`
      <div class="panel">
        <div class="center">
          <p><em>${player1}</em> vs <em>${player2}</em></p>
        </div>
        <div class="split">
          <p>${theirTurn ? "Their" : "Your"} turn</p>
          <p>started ${elapsed} ago</p>
        </div>
      </div>
    `;
  }
}
