import '@pongle/styles/global.css';
import './game-mark';

import {html, css, LitElement, TemplateResult} from 'lit';
import {customElement, property} from 'lit/decorators.js';

@customElement('game-board')
export class GameBoard extends LitElement {
  static styles = css`
    .panel {
      width: 80vw;
      height: 80vw;
      margin-right: auto;
      margin-left: auto;
    }
  
    @media (min-aspect-ratio: 1/1) {
      .panel {
        width: 80vh;
        height: 80vh; 
      }
    }
    
    .panel {
      position: relative;
      display: grid;
      grid-template-columns: 33% 34% 33%;
      grid-template-rows: 33% 34% 33%;
    }

    .board {
      position: absolute;
      top: 0%;
      bottom: 100%;
      left: 0%;
      right: 100%;
      pointer-events: none;
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

    svg {
      height: 100%;
      width: 100%;
    }
  `;

  @property({type: String})
  marks: string

  @property({type: Number})
  turn: number

  @property({type: Number})
  player: number

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
        if (this.player !== this.turn) {
          return null;
        }

        return html`<game-mark
          class="row${row}, col${column}"
          .player="${this.player}"
          @click="${() => this.select(column, row)}"
        />`;
      case "0":
        return html`<game-mark class="row${row} col${column}" player="0" selected />`;
      case "1":
        return html`<game-mark class="row${row} col${column}" player="1" selected />`;
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
        <svg class="board">
          <line x1="33%" y1="0%" x2="33%" y2="100%" stroke="black" stroke-width="2%" />
          <line x1="67%" y1="0%" x2="67%" y2="100%" stroke="black" stroke-width="2%" />
          <line x1="0%" y1="33%" x2="100%" y2="33%" stroke="black" stroke-width="2%" />
          <line x1="0%" y1="67%" x2="100%" y2="67%" stroke="black" stroke-width="2%" />
        </svg>
        ${[0,1,2].map(row => [0,1,2].map(col => this.getMark(col, row)))}
      </div>
    `;
  }
}
