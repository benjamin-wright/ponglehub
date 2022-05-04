import '@pongle/styles/global.css';
import '@pongle/panels/popup-panel';

import {html, css, LitElement} from 'lit';
import {customElement, property} from 'lit/decorators.js';

@customElement('new-game-popup')
export class NewGamePopup extends LitElement {
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
  
  @property({type: Boolean})
  display: boolean;
  
  @property()
  players: {[key: string]: string};

  private cancel() {
    let event = new CustomEvent("cancel", {}); 
    this.dispatchEvent(event);
  }

  private requestNewGame(opponent: string) {
    let event = new CustomEvent<string>("new-game", {detail: opponent}); 
    this.dispatchEvent(event);
  }

  render() {
    if (!this.display) {
      return null;
    }
    
    return html`
      <popup-panel>
        <div class="cancel">
          <button @click="${() => this.cancel()}">X</button>
        </div>
        <h1>New Game!</h1>
        <p>Choose an opponent...</p>
        <ul>
          ${Object.keys(this.players).map(key => html`
            <li class="challenger" @click="${() => this.requestNewGame(key)}">${this.players[key]}</li>
          `)}
        </ul>
      </popup-panel>
    `
  }
}
