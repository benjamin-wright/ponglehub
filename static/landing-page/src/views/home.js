import '../css/global.css';
import '../components/loading-page';
import '../components/nav-bar';
import '../components/hello-world';

import {html, LitElement} from 'lit';
import {customElement, property} from 'lit/decorators.js';
import { logOut, getUserData } from '../auth';

@customElement('home-view')
export class HomeView extends LitElement {
  @property()
  userName = ""

  @property({type: Boolean})
  loading = true
  
  connectedCallback() {
    super.connectedCallback();

    getUserData()
      .then(data => {
        this.loading = false;
        this.userName = data.name;
      })
      .catch(err => {
        console.error(err);
        this.loading = false;
        this.userName = "";
      })
  }

  async _logOut() {
    await logOut();
    window.location = '/';
  }

  content() {
    if (this.loading) {
      return html`<loading-page ?show="${this.loading}"></loading-page>`;
    } else {
      return html`<hello-world name="${this.userName}"></hello-world>`;
    }
  }

  render() {
    return html`
      <nav-bar ?loading="${this.loading}" ?authorised="${true}" @logout-event="${this._logOut}"></nav-bar>
      ${ this.content() }
    `;
  }
}
