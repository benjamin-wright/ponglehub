import '@pongle/styles/global.css';

import {html, css, LitElement, TemplateResult} from 'lit';
import {customElement, property} from 'lit/decorators.js';

@customElement('game-mark')
export class GameMark extends LitElement {
  static styles = css`
    .mark {
      display: flex;
      align-items: center;
      justify-content: center;
    }

    button {
      cursor: pointer;
      background: none;
      border: none;
      width: 100%;
      height: 100%;
      stroke: white;
    }

    button svg {
      display: none;
    }

    button:hover {
      background: var(--default-highlight);
      stroke: red;
    }

    button:hover svg {
      display: block;
    }
    
    svg {
      height: 100%;
      width: 100%;
    }
  `;

  @property({type: Number})
  player: number

  @property({type: Boolean})
  selected: boolean

  private xMark(color: string): TemplateResult<1> {
    return html`
      <svg>
        <line x1="15%" y1="15%" x2="85%" y2="85%" stroke="${color}" stroke-width="10%" />
        <line x1="15%" y1="85%" x2="85%" y2="15%" stroke="${color}" stroke-width="10%" />
      </svg>
    `;
  }

  private oMark(color: string): TemplateResult<1> {
    return html`
      <svg style="color: blue">
        <circle cx="50%" cy="50%" r="35%" fill="none" stroke="${color}" stroke-width="10%" />
      </svg>
    `;
  }

  render() {
    const color = this.selected ? "black" : "blue";
    const mark = this.player === 0 ? this.xMark(color) : this.oMark(color);

    if (this.selected) {
      return mark;
    }

    return html`<button>${mark}</button>`;
  }
}
