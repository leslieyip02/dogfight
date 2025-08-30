import { useCallback, useEffect, useLayoutEffect, useRef, useState } from "react";

import "./Game.css";
import GameEngine from "../game/GameEngine";
import p5 from "p5";
import type { GameEvent, GameInputEventData } from "../game/GameEvent";

const WS_URL = import.meta.env.VITE_WS_URL;

type Props = {
  clientId: string,
  roomId: string,
}

const Game: React.FC<Props> = ({ clientId, roomId }) => {
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

    const ws = new WebSocket(`${WS_URL}?clientId=${clientId}&roomId=${roomId}`);
    ws.onopen = async () => {
      await gameEngineRef.current?.init();
    };
    ws.onmessage = (event: MessageEvent) => {
      gameEngineRef.current?.receive(event);
    };
    setSocket(ws);
  }, [clientId, roomId, socket]);

  useLayoutEffect(() => {
    const sketch = (instance: p5) => {
      gameEngineRef.current = new GameEngine(instance, clientId, roomId, sendInput);
    };

    const instance = new p5(sketch, containerRef.current!);
    return () => instance.remove();
  }, [clientId, roomId, sendInput]);

  return <div className="game__container" ref={containerRef} />;
};

export default Game;