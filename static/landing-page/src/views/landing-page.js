import {html, LitElement} from 'lit';
import {customElement, property} from 'lit/decorators.js';

@customElement('landing-page')
export class LandingPage extends LitElement {
  @property({type: Boolean})
  loading = true
  
  connectedCallback() {
    super.connectedCallback()

    setTimeout(() => { this.loading = false }, 1000)
  }

  render() {
    return html`
      <nav-bar ?authorised="${!this.loading}"></nav-bar>
      <loading-page ?show="${this.loading}"></loading-page>
      <hello-world name="liam"></hello-world>
    `;
  }
}