import '@pongle/styles/global.css';
import '@pongle/components/nav-bar';
import '@pongle/panels/center-panel';

import {css, html, LitElement} from 'lit';
import {customElement, property} from 'lit/decorators.js';
import {Auth} from '@pongle/auth';

@customElement('home-view')
export class HomeView extends LitElement {
  static styles = css`
    a {
      padding: 2em;
    }
  `;

  constructor() {
    super();

    this.auth = new Auth(window.localStorage);
  }

  private auth: Auth

  @property()
  userName = "";
  
  connectedCallback() {
    super.connectedCallback();

    this.auth.init()
      .then(data => {
        this.userName = data.name;
      })
      .catch(err => {
        console.error(err);
        this.auth.logIn();
      })
  }

  private async logOut() {
    try {
      await this.auth.logOut();
    } catch (err) {
      console.error('error logging out:', err)
    }

    window.location.href = '/';
  }

  private content() {
    return html`
      <center-panel height="calc(100% - 3.1em)">
        <a href="http://nac.ponglehub.co.uk">
          <img src="/assets/naughts-and-crosses.png" width="128" height="128" />
        </a>
        <a href="http://draughts.ponglehub.co.uk">
          <img src="/assets/draughts.png" width="128" height="128" />
        </a>
      </center-panel>
    `;
  }

  render() {
    return html`
      <nav-bar .loading="${false}" .authorised="${true}" @logout-event="${this.logOut}"></nav-bar>
      ${ this.content() }
    `;
  }
}
