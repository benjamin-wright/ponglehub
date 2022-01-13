/* global window */

export default class {
  constructor() {
    this.menuHandler = null;
    this.listenHandler = null;
  }

  menuOptions(options, handler) {
    if (this.menuHandler != null) {
      window.removeEventListener("message", this.menuHandler);
    }

    window.parent.postMessage({
      message: "setMenuData",
      value: options,
    }, "*");

    let receiveMessage = (event) => {
      const message = event.data.message;

      switch (message) {
        case "menuChoice":
          handler(event.data.value);
          break;
      }
    };

    window.addEventListener("message", receiveMessage, false);
    this.menuHandler = handler;
  }

  listenMenuOptions(handler) {
    if (this.listenHandler != null) {
      window.removeEventListener("message", this.listenHandler);
    }

    let receiveMessage = (event) => {
      const message = event.data.message;

      switch (message) {
        case "setMenuData":
          handler(event.data.value);
          break;
      }
    };

    window.addEventListener("message", receiveMessage, false);
    this.listenHandler = handler;
  }

  selectMenuOption(frame, option) {
    frame.contentWindow.postMessage({
      message: "menuChoice",
      value: option,  
    }, '*');
  }

  stop() {
    if (this.listenHandler != null) {
      window.removeEventListener("message", this.listenHandler);
    }
    this.listenHandler = null;

    if (this.menuHandler != null) {
      window.removeEventListener("message", this.menuHandler);
    }
    this.menuHandler = null;
  }
}
