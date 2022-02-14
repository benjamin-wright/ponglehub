export class HelloWorld extends HTMLElement {
  constructor() {
    super();
    this.name = 'no-name';
  }

  static get observedAttributes() {
    return ['name'];
  }

  attributeChangedCallback(property, oldValue, newValue) {
    switch (property) {
      case "name": {
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
        p {
          text-align: center;
          font-weight: normal;
          padding: 1em;
          margin: 0 0 2em 0;
          background-color: #eee;
          border: 1px solid #666;
        }
      </style>

      <p>Hello ${this.name}!</p>`;
  }
}