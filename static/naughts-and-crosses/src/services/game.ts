import { PongleEvents } from "@pongle/events";

export interface GameData {
  created: string,
  id: string,
  player1: string,
  player2: string,
  turn: number
}

export function convert(g: any): GameData {
  return {
    created: g.Created,
    id: g.ID,
    player1: g.Player1,
    player2: g.Player2,
    turn: g.Turn
  }
}

export class Game {
  private events: PongleEvents;
  private storage: Storage;
  private listeners: ((property: string) => null)[];

  public live: boolean;
  public userName: string;
  public games: GameData[];
  public currentGame: GameData;
  public players: {[key: string]: string};

  constructor(host: string, storage: Storage) {
    this.live = false;
    this.events = new PongleEvents(host);
    this.storage = storage;
    this.listeners = [];
  }

  addListener(callback: (property: string) => null) {
    this.listeners.push(callback);
  }

  removeListener(callback: (property: string) => null) {
    const index = this.listeners.indexOf(callback);
    if (index >= 0) {
      this.listeners.splice(index, 1);
    }
  }

  private restore(): boolean {
    const username = this.storage.getItem("userName");
    const players = this.storage.getItem("players");
    const games = this.storage.getItem("games");

    if (username === null || players === null || games === null) {
      return false;
    }

    this.userName = username;
    this.players = JSON.parse(players);
    this.games = JSON.parse(games);

    this.inform("userName");
    this.inform("players");
    this.inform("games");

    return true;
  }

  private initialise() {
    this.events.send("auth.list-friends", null);
    this.events.send("naughts-and-crosses.list-games", null);
  }

  async start() {
    this.live = false;
    await this.events.start(
      this.event.bind(this),
      this.start.bind(this),
    );
    this.live = true;

    if (!this.restore()) {
      this.initialise();
    }
  }

  stop() {
    this.events.stop();
  }

  async logout() {
    try {
      this.storage.clear();
      await this.events.logout();
    } finally {
      this.events.login();
    }
  }

  event(type: string, data: any) {
    switch(type) {
      case "auth.whoami.response":
        if (this.userName && this.userName !== data) {
          this.storage.clear();
          this.initialise();
        }

        this.userName = data;
        this.storage.setItem("userName", this.userName);
        this.inform("userName");
        break;
      case "auth.list-friends.response":
        this.players = data;
        this.storage.setItem("players", JSON.stringify(this.players));
        this.inform("players");
        break;
      case "naughts-and-crosses.list-games.response":
        this.games = data.games.map(convert);
        this.games = this.games.sort((a, b) => Date.parse(b.created) - Date.parse(a.created));
        this.storage.setItem("games", JSON.stringify(this.games));
        this.inform("games");
        break;
      case "naughts-and-crosses.new-game.response":
        this.games = this.games.slice();
        this.games.push(convert(data.game));
        this.games = this.games.sort((a, b) => Date.parse(b.created) - Date.parse(a.created));
        this.storage.setItem("games", JSON.stringify(this.games));
        this.inform("games");
        break;
      default:
        console.error(`Unrecognised response type from server: ${type}`);
        break;
    }
  }

  private inform(property: string) {
    this.listeners.forEach(callback => callback(property));
  }

  newGame(opponent: string) {
    this.events.send("naughts-and-crosses.new-game", {opponent});
  }

  loadGame(id: string) {
    this.events.send("naughts-and-crosses.load-game", {id});
  }
}