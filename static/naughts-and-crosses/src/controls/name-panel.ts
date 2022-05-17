import '@pongle/styles/global.css';

import {html, css, LitElement} from 'lit';
import {customElement, property} from 'lit/decorators.js';

@customElement('name-panel')
export class NamePanel extends LitElement {
  static styles = css`
    div {
      display: flex;
      height: 100%;
      justify-content: center;
      align-items: center;
      flex-direction: column;
    }

    label {
      text-transform: capitalize;
    }

    .hidden {
      visibility: hidden;
    }
  `;

  @property({type: String})
  player: string

  @property({type: Boolean})
  active: boolean


  render() {
    return html`
      <div>
        <label>${this.player}</label>
        <label class="${this.active || "hidden"}">^</label>
      </div>
    `;
  }
}
