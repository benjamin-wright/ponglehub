export class NavBar extends HTMLElement {
  constructor() {
    super();

    this.authorised = false;
  }

  static get observedAttributes() {
    return ['authorised'];
  }

  attributeChangedCallback(property, oldValue, newValue) {
    switch (property) {
      case "authorised": {
        if (oldValue === newValue) return;
        this[property] = newValue;
        break;
      }
      default: {
        console.error(`Unrecognised property: ${property}`);
      }
    }
  }

  // connect component
  connectedCallback() {
    const shadow = this.attachShadow({ mode: 'closed' });
    shadow.innerHTML = `
      <style>
        .container {
        }
        p {
          text-align: center;
          font-weight: normal;
          padding: 1em;
          margin: 0 0 2em 0;
          background-color: #eee;
          border: 1px solid #666;
        }
      </style>

      <div class="container">
        <div>
          <slot>
            <span>LOGO</span>
          </slot>
        </div>
      </div>`;
  }
}