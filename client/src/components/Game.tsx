import "./Game.css";

import p5 from "p5";
import { useEffect, useLayoutEffect, useRef, useState } from "react";

import Engine from "../game/Engine";
import { Event } from "../pb/event";

const WS_URL = import.meta.env.VITE_WS_URL;

type Props = {
  clientId: string,
  token: string,
}

const Game: React.FC<Props> = ({ clientId }) => {
  const gameEngineRef = useRef<Engine | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const [socket, setSocket] = useState<WebSocket | null>(null);

  useEffect(() => {
    if (socket !== null) {
      return;
    }

    const token = localStorage.getItem("jwt");
    if (!token) {
      return;
    }

    const ws = new WebSocket(`${WS_URL}?token=${token}`);
    ws.binaryType = "arraybuffer";
    ws.onopen = async () => {
      await gameEngineRef.current?.init();
    };
    ws.onmessage = (event: MessageEvent) => {
      console.log(event, event.data);
      const message = new Uint8Array(event.data);
      gameEngineRef.current?.receive(Event.decode(message));
    };
    setSocket(ws);
  }, [socket]);

  useLayoutEffect(() => {
    if (!socket) {
      return;
    }

    const sketch = (instance: p5) => {
      gameEngineRef.current = new Engine(instance, clientId, socket);
    };

    const instance = new p5(sketch, containerRef.current!);
    return () => instance.remove();
  }, [clientId, socket]);

  return <div className="game__container" ref={containerRef} />;
};

export default Game;
