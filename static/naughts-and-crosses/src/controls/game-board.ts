import '@pongle/styles/global.css';

import {html, css, LitElement, TemplateResult} from 'lit';
import {customElement, property} from 'lit/decorators.js';

@customElement('game-board')
export class GameBoard extends LitElement {
  static styles = css`
    .panel {
      width: 80vw;
      height: 80vw; 
      display: grid;

      grid-template-columns: 3;
      grid-template-rows: 3;
    }

    .row0{
      grid-row: 1;
    }

    .row1{
      grid-row: 2;
    }

    .row2{
      grid-row: 3;
    }

    .col0{
      grid-column: 1;
    }

    .col1{
      grid-column: 2;
    }

    .col2{
      grid-column: 3;
    }

    .mark {
      display: flex;
      align-items: center;
      justify-content: center;
    }

    button {
      cursor: pointer;
      background: none;
      border: none;
    }
  `;

  @property({type: String})
  marks: string

  @property({type: Number})
  turn: number

  private xMark(): TemplateResult<1> {
    return html`
      <svg>
        <circle cx="50%" cy="50%" r="45%"/>
      </svg>
    `;
  }

  private select(column: number, row: number) {
    const index = row * 3 + column;
    const event = new CustomEvent<number>("select", {detail: index});

    this.dispatchEvent(event);
  }

  private getMark(column: number, row: number): TemplateResult<1> {
    if (column < 0 || column > 2) {
      throw new Error(`bad column input: ${column}`);
    }

    if (row < 0 || row > 2) {
      throw new Error(`bad row input: ${row}`);
    }

    const index = row * 3 + column;

    switch (this.marks[index]) {
      case "-":
        return html`<button
            @click="${() => this.select(column, row)}"
            class="row${row}, col${column}"
          >
            ${this.turn === 0 ? "x" : "o" }
          </button>
        `;
      case "0":
        return html`<div class="mark row${row}, col${column}">${this.xMark()}</div>`;
      case "1":
        return html`<div class="mark row${row}, col${column}"><p>o</p></div>`;
      default:
        throw new Error(`bad character: ${this.marks[index]}`);
    }
  }

  render() {
    if (!this.marks) {
      return null;
    }

    return html`
      <div class="panel">
        ${[0,1,2].map(row => [0,1,2].map(col => this.getMark(col, row)))}
      </div>
    `;
  }
}
