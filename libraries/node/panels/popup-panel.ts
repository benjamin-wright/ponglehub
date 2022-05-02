import {html, css, LitElement} from 'lit';
import {customElement} from 'lit/decorators.js';

@customElement('popup-panel')
export class PopupPanel extends LitElement {
  static styles = css`
    .background {
      position: absolute;
      top: 0;
      bottom: 0;
      left: 0;
      right: 0;
      display: flex;
      align-items: center;
      justify-content: center;
      background-color: var(--default-overlay);
    }

    .center-panel {
      background-color: var(--default-background);
      border-radius: 1em;
      padding: 1em;
    }
  `;

  render() {
    return html`
      <div class="background">
        <div class="center-panel">
          <slot></slot>
        </div>
      </div>
    `;
  }
}