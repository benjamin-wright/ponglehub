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

  @property({type: Boolean})
  loading = true;
  
  connectedCallback() {
    super.connectedCallback();

    this.auth
      .restore()
      .then(data => {
        this.loading = false;
        this.userName = data.name;
      })
      .catch(err => {
        console.error(err);
        window.location.href = '/';
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
    if (this.loading) {
      return html`
        <center-panel height="calc(100% - 3.1em)">
          <p>loading...</p>
        </center-panel>
      `;
    } else {
      return html`
        <center-panel height="calc(100% - 3.1em)">
          <a href="/naughts-and-crosses">
            <img src="/assets/naughts-and-crosses.png" width="128" height="128" />
          </a>
          <a href="/draughts">
            <img src="/assets/draughts.png" width="128" height="128" />
          </a>
        </center-panel>
      `;
    }
  }

  render() {
    return html`
      <nav-bar .loading="${this.loading}" .authorised="${true}" @logout-event="${this.logOut}"></nav-bar>
      ${ this.content() }
    `;
  }
}
