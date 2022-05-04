import '@pongle/styles/global.css';

import {html, css, LitElement} from 'lit';
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
  `;

  @property({type: String})
  marks: string

  private getMark(column: number, row: number): string {
    if (column < 0 || column > 2) {
      throw new Error(`bad column input: ${column}`);
    }

    if (row < 0 || row > 2) {
      throw new Error(`bad row input: ${row}`);
    }

    const index = row * 3 + column;

    return this.marks[index];
  }

  render() {
    return html`
      <div class="panel">
        ${
          [0,1,2].map(row => [0,1,2].map(col => html`
            <button class="row${row}, col${col}">${
              this.getMark(col, row)
            }</button>
          `))
        }
      </div>
    `;
  }
}
