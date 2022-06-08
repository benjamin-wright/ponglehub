import '@pongle/styles/global.css';

import {html, css, LitElement, TemplateResult} from 'lit';
import {customElement, property} from 'lit/decorators.js';
import { timeSince } from '../services/utils';
import { GameData } from '../services/game';

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

    .panel.tied {
      background-color: var(--default-neutral);
    }

    .panel.won {
      background-color: var(--default-success);
    }

    .panel.lost {
      background-color: var(--default-failure);
    }
  `;

  @property({type: Object})
  private game: GameData;

  @property({type: Object})
  private players: {[key: string]: string};

  render() {
    const elapsed = timeSince(this.game.created);
    const player = this.players[this.game.player1] ? 1 : 0;
    const player1 = this.players[this.game.player1] || "You";
    const player2 = this.players[this.game.player2] || "You";

    if (this.game.finished) {
      let outcomeClass;
      let outcomeMessage;
      switch (true) {
        case player === this.game.turn:
          outcomeClass = "won";
          outcomeMessage = "You won";
          break;
        case this.game.turn === -1:
          outcomeClass = "tied";
          outcomeMessage = "Tied";
          break;
        default:
          outcomeClass = "lost";
          outcomeMessage = "You lost";
      }

      return html`
        <div class="panel ${outcomeClass}">
          <div class="center">
            <p><em>${player1}</em> vs <em>${player2}</em></p>
          </div>
          <div class="split">
            <p>${outcomeMessage}</p>
            <p>started ${elapsed} ago</p>
          </div>
        </div>
      `;
    }


    return html`
      <div class="panel">
        <div class="center">
            <p><em>${player1}</em> vs <em>${player2}</em></p>
        </div>
        <div class="split">
          <p>${player === this.game.turn ? "Your" : "Their"} turn</p>
          <p>started ${elapsed} ago</p>
        </div>
      </div>
    `;
  }
}
