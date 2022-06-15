export {}

class NavBar extends HTMLElement {
  constructor() {
    super();

    this.innerHTML = /*html*/`
      <style>
        nav {
          background-color: var(--default-foreground);
          color: var(--default-background);
          padding: 0;
          display: flex;
          flex-direction: row;
          justify-content: space-between;
        }

        nav a {
          text-align: center;
          font-weight: normal;
          padding: 1em;
          margin: 0;
        }

        nav a,
        nav a:visited {
          text-transform: uppercase;
          color: var(--default-background);
          font-weight: bold;
          text-decoration: none;
        }
      </style>
      <header>
        <nav>
          <a href="http://games.ponglehub.co.uk">PONGLEHUB</a>
          <a href="#">logout</a>
        </nav>
      </header>
    `;

    const links = this.getElementsByTagName('a');
    links[1].onclick = (event) => {
      event.preventDefault();
      this.dispatchEvent(new CustomEvent('logout-event', {}));
    }
  }
}

customElements.define('nav-bar', NavBar);