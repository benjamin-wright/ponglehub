import { listen, SendEvents } from './events';

export type GameData = {
  games: any[],
  players: {[key: string]: string}
}

function convert(g: any): any {
  return {
    created: g.Created,
    id: g.ID,
    player1: g.Player1,
    player2: g.Player2,
    turn: g.Turn
  }
}

export class Game {
  private data: GameData
  private events: SendEvents
  private callback: () => void

  constructor() {
    this.data = {
      games: [],
      players: {}
    };
  }

  onmessage(type: string, data: any) {
    console.info(`Event: ${type}`);

    switch(type) {
      case "naughts-and-crosses.list-games.response":
        this.data.games = data.games.map(convert);
        this.callback();
        break;
      case "auth.list-friends.response":
        this.data.players = data;
        this.callback();
        break;
      case "naughts-and-crosses.new-game.response":
        this.data.games.push(convert(data.game))
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
    this.events.send("auth.list-friends", null);
    this.events.send("naughts-and-crosses.list-games", null);
  }

  stop() {
    delete this.callback;
    this.events.stop();
  }

  games(): string[] {
    return this.data.games;
  }

  players(): {[key: string]: string} {
    return this.data.players;
  }

  newGame(opponent: string) {
    this.events.send("naughts-and-crosses.new-game", {opponent})
  }
}