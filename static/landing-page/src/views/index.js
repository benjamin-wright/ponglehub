import '../css/global.css';
import '../components/loading-page';
import '../components/nav-bar';

import {html, LitElement} from 'lit';
import {customElement, property} from 'lit/decorators.js';
import { getUserData, redirectToLogin } from '../auth';

@customElement('index-view')
export class IndexView extends LitElement {
  @property({type: Boolean})
  loading = true
  
  connectedCallback() {
    super.connectedCallback();

    getUserData()
      .then(() => {
        window.location = "/home";
      })
      .catch(err => {
        console.error(err);
        this.loading = false;
      })
  }

  async _login() {
    redirectToLogin()
  }

  content() {
    if (this.loading) {
      return html`<loading-page ?show="${this.loading}"></loading-page>`;
    } else {
      return html``;
    }
  }

  render() {
    return html`
      <nav-bar ?loading="${this.loading}" ?authorised="${false}" @login-event="${this._login}"></nav-bar>
      ${ this.content() }
    `;
  }
}
