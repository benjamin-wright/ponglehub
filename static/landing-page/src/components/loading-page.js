import {html, css, LitElement} from 'lit';
import {customElement, property} from 'lit/decorators.js';

@customElement('loading-page')
export class LoadingPage extends LitElement {
  static styles = css`
    div {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 100%;
      height: calc(100% - 3em);
    }

    p {
      color: var(--default-foreground);
    }
  `;

  @property({type: Boolean})
  show = true;

  render() {
    if (this.show) {
      return html`
        <div>
          <p>loading...</p>
        </div>
      `;
    } else {
      return '';
    }
  }
}