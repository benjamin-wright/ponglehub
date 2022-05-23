import '@pongle/styles/global.css';
import '@pongle/panels/popup-panel';

import {html, css, LitElement, TemplateResult} from 'lit';
import {customElement, property} from 'lit/decorators.js';

@customElement('victory-popup')
export class VictoryPopup extends LitElement {
  static styles = css`
    h1 {
      width: 100%;
      text-align: center;
      text-transform: capitalize; 
    }

    .ok {
      display: flex;
      justify-content: right;
    }

    .ok button {
      background: none;
      border: none;
      color: red;
      cursor: pointer;
    }
  `;
  
  @property({type: Boolean})
  display: boolean;

  @property({type: String})
  player: string;

  private ok() {
    let event = new CustomEvent("ok", {}); 
    this.dispatchEvent(event);
  }

  private message(): TemplateResult<1> {
    if (this.player === "you") {
        return html`<p>Nicely done, you absolute legend!</p>`;
    } else {
        return html`<p>Good game, better luck next time!</p>`;
    }
  }

  render() {
    if (!this.display) {
      return null;
    }
    
    return html`
      <popup-panel>
        <h1>${this.player} Won!</h1>
        ${this.message()}
        <div class="ok">
          <button @click="${() => this.ok()}">OK</button>
        </div>
      </popup-panel>
    `
  }
}
