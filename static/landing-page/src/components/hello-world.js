import {html, css, LitElement} from 'lit';
import {customElement, property} from 'lit/decorators.js';

@customElement('hello-world')
export class HelloWorld extends LitElement {
  static styles = css`
    p {
      text-align: center;
      font-weight: normal;
      padding: 1em;
      margin: 0 0 2em 0;
      background-color: var(--default-background);
      color: var(--default-foreground);
      border: 1px solid #666;
      text-transform: capitalize;
    }
  `;

  @property()
  name = 'no-name';

  render() {
    return html`<p>Hello ${this.name}!</p>`;
  }
}