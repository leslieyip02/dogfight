import { useCallback, useEffect, useLayoutEffect, useRef, useState } from "react";

import "./Game.css";
import Engine from "../game/Engine";
import p5 from "p5";
import type { Event, InputEventData } from "../game/types/event";

const WS_URL = import.meta.env.VITE_WS_URL;

type Props = {
  clientId: string,
  token: string,
}

const Game: React.FC<Props> = ({ clientId }) => {
  const gameEngineRef = useRef<Engine | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const [socket, setSocket] = useState<WebSocket | null>(null);

  const sendInput = useCallback((data: InputEventData) => {
    if (socket?.readyState !== WebSocket.OPEN) {
      return;
    }

    const event: Event = {
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
    ws.onopen = async () => {
      await gameEngineRef.current?.init();
    };
    ws.onmessage = (event: MessageEvent) => {
      const gameEvent: Event = JSON.parse(event.data);
      gameEngineRef.current?.receive(gameEvent);
    };
    setSocket(ws);
  }, [socket]);

  useLayoutEffect(() => {
    const sketch = (instance: p5) => {
      gameEngineRef.current = new Engine(instance, clientId, sendInput);
    };

    const instance = new p5(sketch, containerRef.current!);
    return () => instance.remove();
  }, [clientId, sendInput]);

  return <div className="game__container" ref={containerRef} />;
};

export default Game;
