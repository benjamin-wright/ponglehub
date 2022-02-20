import '@pongle/styles/global.css';
import '@pongle/components/nav-bar';
import '@pongle/panels/center-panel';

import {html, css, LitElement, TemplateResult} from 'lit';
import {customElement, state} from 'lit/decorators.js';
import {Auth} from '@pongle/auth';

@customElement('index-view')
export class IndexView extends LitElement {
  static styles = css`
    h1 {
      color: var(--default-foreground);
    }

    .highlight {
      color: var(--default-highlight);
      text-shadow: 1px 1px 3px var(--default-foreground);
    }
  `;

  @state()
  private loading: boolean;

  private auth: Auth;

  constructor() {
    super();
    
    this.auth = new Auth(window.localStorage);
    this.loading = this.auth.loading();
  }
  
  connectedCallback() {
    super.connectedCallback();

    if (this.auth.loggedIn()) {
      window.location.href = "/home";
      return
    }

    if (this.auth.loading()) {
      this.auth.load()
      .then(() => {
          window.location.href = "/home";
        }).catch(err => {
          console.warn('failed to load user data:', err);
          this.loading = false;
        });
    }
  }

  private async login() {
    this.auth.logIn();
  }

  private content(): TemplateResult<1> {
    if (this.loading) {
      return html`<p>loading...</p>`;
    } else {
      return html`<h1>Welcome to <span class="highlight">PONGLEHUB GAMES</span>!</h1>`;
    }
  }

  render() {
    return html`
      <nav-bar .loading="${this.loading}" .authorised="${false}" @login-event="${this.login}"></nav-bar>
      <center-panel height="calc(100% - 3.1em)">
        ${ this.content() }
      </center-panel>
    `;
  }
}
