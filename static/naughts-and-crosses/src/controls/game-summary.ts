import '@pongle/styles/global.css';

import {html, css, LitElement} from 'lit';
import {customElement, property} from 'lit/decorators.js';
import { timeSince } from '../services/utils';

@customElement('game-summary')
export class GameSummary extends LitElement {
  static styles = css`
    .panel {
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
      justify-content: space-between;
    }
  `;

  @property({type: Object})
  private game: {id: string, player1: string, player2: string, turn: number, created: string};

  @property({type: Object})
  private players: {[key: string]: string};

  render() {
    console.info(this.game);
    console.info(this.players);

    let player1 = this.players[this.game.player1] || "You";
    let player2 = this.players[this.game.player2] || "You";

    let turn = this.game.turn == 0 ? player1 : player2;

    if (turn === "You") {
      turn += "r";
    } else {
      turn += "'s";
    }

    let elapsed = timeSince(this.game.created);

    return html`
      <div class="panel">
        <div class="center">
          <p><em>${player1}</em> vs <em>${player2}</em></p>
        </div>
        <div class="split">
          <p>${turn} go</p>
          <p>started ${elapsed} ago</p>
        </div>
      </div>
    `;
  }
}
