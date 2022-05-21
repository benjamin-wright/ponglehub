import '@pongle/styles/global.css';
import '../controls/game-summary';
import '../controls/new-game';

import {html, css, LitElement} from 'lit';
import {customElement, property} from 'lit/decorators.js';
import {GameData} from '../services/game';

@customElement('list-games')
export class ListGames extends LitElement {
  static styles = css`
    ul {
      list-style: none;
      display: flex;
      flex-wrap: wrap;
    }
  `;

  @property({type: Array})
  games: GameData[];
  
  @property()
  players: {[key: string]: string};

  private navigate(id: string) {
    window.location.href = `/game?id=${id}`;
  }

  private newGame() {
    let event = new CustomEvent('new-game', {});
    this.dispatchEvent(event);
  }

  render() {
    if (!this.games) {
      return html`<p>loading...</p>`;
    }

    return html`
      <ul>
        <li><new-game @click="${() => this.newGame()}"/></li>
        ${this.games.map(game => html`
          <li>
            <game-summary @click="${() => this.navigate(game.id)}" .game="${game}" .players="${this.players}"></game-summary>
          </li>
        `)}
      </ul>
    `;
  }
}
