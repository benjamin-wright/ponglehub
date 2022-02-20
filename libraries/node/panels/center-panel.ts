import {html, css, LitElement} from 'lit';
import {customElement, property} from 'lit/decorators.js';

@customElement('center-panel')
export class CenterPanel extends LitElement {
  static styles = css`
    div {
      display: flex;
      align-items: center;
      justify-content: center;
    }
  `;

  @property()
  height: string

  render() {
    return html`
      <div style="height: ${this.height}">
        <slot></slot>
      </div>
    `;
  }
}