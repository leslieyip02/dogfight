import type { SnapshotEventData } from "../game/types/event";

const API_URL = import.meta.env.VITE_API_URL;

export async function fetchGameSnapshot(): Promise<SnapshotEventData | null> {
  const token = localStorage.getItem("jwt");
  if (!token) {
    return Promise.reject(null);
  }
  return await fetch(`${API_URL}/room/state?token=${token}`)
    .then(response => response.json());
}
