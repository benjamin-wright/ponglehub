import { PongleEvents } from "@pongle/events";

export interface GameData {
  created: string,
  id: string,
  player1: string,
  player2: string,
  turn: number,
  finished: boolean,
}

export function convert(g: any): GameData {
  return {
    created: g.Created,
    id: g.ID,
    player1: g.Player1,
    player2: g.Player2,
    turn: g.Turn,
    finished: g.Finished,
  }
}
