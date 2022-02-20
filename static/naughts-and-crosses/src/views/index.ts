import '@pongle/styles/global.css';
import '@pongle/components/nav-bar';
import '@pongle/panels/center-panel';

import {html, css, LitElement, TemplateResult} from 'lit';
import {customElement, state} from 'lit/decorators.js';
import {Auth} from '@pongle/auth';

@customElement('index-view')
export class IndexView extends LitElement {
  static styles = css``;

  private auth: Auth;

  @state()
  private userName: string;

  constructor() {
    super();
    
    this.auth = new Auth(window.localStorage);
  }

  connectedCallback() {
    super.connectedCallback();

    this.auth.init()
      .then(data => this.userName = data.name)
      .catch(() => this.auth.logIn());
  }

  private async logOut() {
    await this.auth.logOut();
    window.location.href="http://localhost:7000";
  }

  render() {
    return html`
      <nav-bar .loading="${false}" .authorised="${true}" @logout-event="${this.logOut}"></nav-bar>
      <h1>Lets play Naughts and Crosses!</h1>
    `;
  }
}
