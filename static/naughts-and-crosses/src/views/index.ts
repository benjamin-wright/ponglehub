import '@pongle/styles/global.css';
import '@pongle/components/nav-bar';
import '@pongle/panels/center-panel';

import {html, css, LitElement, TemplateResult} from 'lit';
import {customElement, state} from 'lit/decorators.js';
import {Auth} from '@pongle/auth';

@customElement('index-view')
export class IndexView extends LitElement {
  static styles = css``;

  @state()
  private loading: boolean;

  private auth: Auth;

  constructor() {
    super();
    
    this.auth = new Auth(window.localStorage);
    this.loading = !this.auth.loggedIn();
  }
  
  connectedCallback() {
    super.connectedCallback();

    if (this.loading) {
      this.auth.load()
        .then(() => this.loading = false)
        .catch(() => this.auth.logIn());
    }
  }

  private content(): TemplateResult<1> {
    if (this.loading) {
      return html`
        <center-panel height="calc(100% - 3.1em)">
          <p>loading...</p>
        </center-panel>`;
    } else {
      return html`<h1>Lets play Naughts and Crosses!</h1>`;
    }
  }

  private async logOut() {
    await this.auth.logOut();
    window.location.href="http://localhost:7000";
  }

  render() {
    return html`
      <nav-bar .loading="${this.loading}" .authorised="${true}" @logout-event="${this.logOut}"></nav-bar>
      ${ this.content() }
    `;
  }
}
