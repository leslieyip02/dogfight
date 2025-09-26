import type { GameJoinEventData, GameUpdatePowerupEventData } from "../game/GameEvent";

const API_URL = import.meta.env.VITE_API_URL;

export type GameStateData = {
  players: GameJoinEventData[],
  powerups: GameUpdatePowerupEventData[],
};

export async function fetchGameState(): Promise<GameStateData | null> {
  const token = localStorage.getItem("jwt");
  if (!token) {
    return Promise.reject(null);
  }
  return await fetch(`${API_URL}/room/state?token=${token}`)
    .then(response => response.json());
}
