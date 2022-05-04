import '@pongle/styles/global.css';

import {html, css, LitElement} from 'lit';
import {customElement} from 'lit/decorators.js';

@customElement('new-game')
export class GameSummary extends LitElement {
  static styles = css`
    input {
      font-size: 4em;
      margin: 0.25em;
      padding: 0.1em 0.8em;
      background: var(--default-background);
      color: var(--default-foreground);
      border: 2px dashed var(--default-foreground);
      border-radius: 0.25em;
      cursor: pointer;
    }

    input:hover, .panel:focus {
      background: var(--default-foreground);
      color: var(--default-background);
      border: 2px dashed var(--default-highlight);
    }
  `;

  render() {
    return html`
      <input type="button" value="+"/>
    `;
  }
}
