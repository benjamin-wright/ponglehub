import { send, receive } from './events';

export type GameData = {
  games: string[]
}

function sleep(time: number): Promise<null> {
  return new Promise((resolve, reject) => setTimeout(resolve, time));
}

export class Game {
  private data: GameData
  private running: boolean

  constructor() {
    this.data = {
      games: []
    };
    this.running = false;
  }

  async refresh() {
    await send("naughts-and-crosses.list-games");
  }

  async start() {
    this.running = true;
    
    while (this.running) {
      let messages = await receive();
    
      messages.forEach(m => {
        switch(m.type) {
          case "naughts-and-crosses.list-games.response":
            console.info(`List response: ${JSON.stringify(m.data)}`);
            break;
          default:
            console.warn(`Unrecognised response type from server: ${m.type}`);
            break;
        }
      });

      await sleep(1);
    }
  }

  stop() {
    this.running = false;
  }
}