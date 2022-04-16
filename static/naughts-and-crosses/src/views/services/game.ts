import { listen, SendEvents } from './events';

export type GameData = {
  games: string[]
}

function sleep(time: number): Promise<null> {
  return new Promise((resolve, reject) => setTimeout(resolve, time));
}

export class Game {
  private data: GameData
  private events: SendEvents

  constructor() {
    this.data = {
      games: []
    };
  }

  async refresh() {
    this.events.send("naughts-and-crosses.list-games", null);
  }

  onmessage(type: string, message: string) {
    switch(type) {
      case "naughts-and-crosses.list-games.response":
        console.info(`List response: ${message}`);
        break;
      default:
        console.warn(`Unrecognised response type from server: ${type}`);
        break;
    }
  }

  async start() {
    try {
      this.events = await listen(this.onmessage.bind(this), this.start.bind(this));
    } catch(error: any) {
      await sleep(10);
      this.start();
      console.error(error);
      return
    }

    this.events.send("naughts-and-crosses.list-games", null);
  }

  stop() {
    this.events.stopper();
  }
}