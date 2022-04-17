import { listen, SendEvents } from './events';

export type GameData = {
  games: string[],
  players: string[]
}

export class Game {
  private data: GameData
  private events: SendEvents
  private callback: () => void

  constructor() {
    this.data = {
      games: [],
      players: []
    };
  }

  async refresh() {
    this.events.send("naughts-and-crosses.list-games", null);
  }

  onmessage(type: string, data: any) {
    switch(type) {
      case "naughts-and-crosses.list-games.response":
        this.data.games = data.games;
        this.callback();
        break;
      case "auth.list-users.response":
        this.data.players = data.userids;
        this.callback();
        break;
      default:
        console.error(`Unrecognised response type from server: ${type}`);
        break;
    }
  }

  async start(callback: () => void) {
    this.callback = callback;

    this.events = listen(this.onmessage.bind(this));
    this.events.send("auth.list-users", null);
    this.events.send("naughts-and-crosses.list-games", null);
  }

  stop() {
    delete this.callback;
    this.events.stop();
  }

  games(): string[] {
    return this.data.games;
  }

  players(): string[] {
    return this.data.players;
  }
}