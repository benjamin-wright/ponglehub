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
  authorised: boolean;

  @property({type: Boolean})
  loading: boolean;

  private logoutEvent(event: Event) {
    event.preventDefault();

    let e = new CustomEvent('logout-event', {});
    this.dispatchEvent(e);
  }
  
  private loginEvent(event: Event) {
    event.preventDefault();

    let e = new CustomEvent('login-event', {});
    this.dispatchEvent(e);
  }

  private button() {
    if (this.loading) {
      return html`<span></span>`;
    }

    if (this.authorised) {
      return html`<div><a href="#" @click="${this.logoutEvent}">logout</a></div>`;
    }

    return html`<div><a href="#" @click="${this.loginEvent}">login</a></div>`;
  }

  render() {
    return html`
      <div class="container">
        <div><a href="http://games.ponglehub.co.uk"><span>PONGLEHUB</span></a></div>
        ${ this.button() }
      </div>
    `;
  }
}