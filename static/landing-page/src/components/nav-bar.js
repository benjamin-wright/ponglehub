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

  @property({type: Boolean})
  loading = true;

  _logoutEvent(event) {
    event.preventDefault();

    let e = new CustomEvent('logout-event', {});
    this.dispatchEvent(e);
  }
  
  _loginEvent(event) {
    event.preventDefault();

    let e = new CustomEvent('login-event', {});
    this.dispatchEvent(e);
  }

  button() {
    if (this.loading) {
      return html`<span></span>`;
    }

    if (this.authorised) {
      return html`<div><a href="#" @click="${this._logoutEvent}">logout</a></div>`;
    }

    return html`<div><a href="#" @click="${this._loginEvent}">login</a></div>`;
  }

  render() {
    return html`
      <div class="container">
        <div><a href="/"><span>LOGO</span></a></div>
        ${ this.button() }
      </div>
    `;
  }
}