import {html, css, LitElement} from 'lit';
import {customElement, property} from 'lit/decorators.js';

@customElement('games-list')
export class GamesList extends HTMLElement {
    static styles = css`
        ul {
            list-style: none;
            display: flex;
            flex-wrap: wrap;
        }

        ul input {
            font-size: 4em;
            margin: 0.25em;
            padding: 0.1em 0.8em;
            background: var(--default-background);
            color: var(--default-foreground);
            border: 2px dashed var(--default-foreground);
            border-radius: 0.25em;
            cursor: pointer;
        }

        ul input:hover {
            background: var(--default-foreground);
            color: var(--default-background);
            border: 2px dashed var(--default-highlight);
        }
    `;

    render() {
        return html`
            <template id="new-list-item">
                <li>
                    <input type="button" value="+"/>
                </li>
            </template>
            <template id="game-list-item">
                <li>
                    <p>vs</p>
                    <p>created</p>
                </li>
            </template>
            <ul id="games-list"></ul>
        `;
    }
}