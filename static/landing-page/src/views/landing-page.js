import {html, LitElement} from 'lit';
import {customElement, property} from 'lit/decorators.js';
import { logOut, getUserData, redirectToLogin } from '../auth';

@customElement('landing-page')
export class LandingPage extends LitElement {
  @property({type: Boolean})
  loading = true

  @property({type: Boolean})
  loggedIn = false

  @property()
  userName = ""
  
  connectedCallback() {
    super.connectedCallback();

    getUserData()
      .then(data => {
        this.loading = false;
        this.loggedIn = true;
        this.userName = data.name;
      })
      .catch(err => {
        console.error(err);

        this.loading = false;
        this.loggedIn = false;
        this.userName = "";
      })
  }

  async _logOut() {
    await logOut();

    this.loading = false;
    this.loggedIn = false;
    this.userName = "";
  }

  async _login() {
    redirectToLogin()
  }

  content() {
    if (this.loading || !this.loggedIn) {
      return html``;
    } else {
      return html`<hello-world name="${this.userName}"></hello-world>`;
    }
  }

  render() {
    return html`
      <nav-bar ?loading="${this.loading}" ?authorised="${this.loggedIn}" @logout-event="${this._logOut}" @login-event="${this._login}"></nav-bar>
      <loading-page ?show="${this.loading}"></loading-page>
      ${ this.content() }
    `;
  }
}
