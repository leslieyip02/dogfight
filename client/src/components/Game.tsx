import { useCallback, useEffect, useLayoutEffect, useRef, useState } from "react";

import "./Game.css";
import GameEngine from "../game/GameEngine";
import p5 from "p5";
import type { GameEvent, GameInputEventData } from "../game/GameEvent";

const WS_URL = import.meta.env.VITE_WS_URL;

type Props = {
  clientId: string,
  token: string,
}

const Game: React.FC<Props> = ({ clientId }) => {
  const gameEngineRef = useRef<GameEngine | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const [socket, setSocket] = useState<WebSocket | null>(null);

  const sendInput = useCallback((data: GameInputEventData) => {
    if (socket?.readyState !== WebSocket.OPEN) {
      return;
    }

    const event: GameEvent = {
      type: "input",
      data: data,
    };
    socket?.send(JSON.stringify(event));
  }, [socket]);

  useEffect(() => {
    if (socket !== null) {
      return;
    }

    const token = localStorage.getItem("jwt");
    if (!token) {
      return;
    }

    const ws = new WebSocket(`${WS_URL}?token=${token}`);
    ws.onopen = gameEngineRef.current?.init ?? null;
    ws.onmessage = (event: MessageEvent) => {
      const gameEvent: GameEvent = JSON.parse(event.data);
      gameEngineRef.current?.receive(gameEvent);
    };
    setSocket(ws);
  }, [socket]);

  useLayoutEffect(() => {
    const sketch = (instance: p5) => {
      gameEngineRef.current = new GameEngine(instance, clientId, sendInput);
    };

    const instance = new p5(sketch, containerRef.current!);
    return () => instance.remove();
  }, [clientId, sendInput]);

  return <div className="game__container" ref={containerRef} />;
};

export default Game;
