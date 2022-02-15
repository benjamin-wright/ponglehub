import {html, css, LitElement} from 'lit';
import {customElement, property} from 'lit/decorators.js';

@customElement('nav-bar')
export class NavBar extends LitElement {
  static styles = css`
    .container {
      background-color: var(--default-foreground);
      color: var(--default-background);
      padding: 1em;
      display: flex;
      flex-direction: row;
      justify-content: space-between;
    }

    p {
      text-align: center;
      font-weight: normal;
      padding: 1em;
      margin: 0 0 2em 0;
    }

    a, a:visited {
      text-transform: uppercase;
      color: var(--default-background);
      font-weight: bold;
      text-decoration: none;
    }
  `;

  @property({type: Boolean})
  authorised = false;

  render() {
    if (this.authorised) {
      return html`
        <div class="container">
          <div><a href="/"><span>LOGO</span></a></div>
          <div><a href="/">logout</a></div>
        </div>
      `;
    } else {
      return html`
        <div class="container">
          <div><a href="/"><span>LOGO</span></a></div>
          <span></span>
        </div>
      `;
    }
  }
}