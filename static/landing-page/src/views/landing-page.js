import {html, LitElement} from 'lit';
import {customElement, property} from 'lit/decorators.js';
import { logOut, getUserData, redirectToLogin } from '../auth';

@customElement('landing-page')
export class LandingPage extends LitElement {
  @property({type: Boolean})
  loading = true

  @property()
  userName = "no-name"
  
  connectedCallback() {
    super.connectedCallback();

    getUserData()
      .then(data => {
        this.loading = false;
        this.userName = data.name;
      })
      .catch(err => {
        console.error(err);
        redirectToLogin();
      })
  }

  async _logOut() {
    await logOut()
    redirectToLogin()
  }

  render() {
    return html`
      <nav-bar ?authorised="${!this.loading}" @logout-event="${this._logOut}"></nav-bar>
      <loading-page ?show="${this.loading}"></loading-page>
      <hello-world name="${this.userName}"></hello-world>
    `;
  }
}