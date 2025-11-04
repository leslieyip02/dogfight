import "./Game.css";

import p5 from "p5";
import { useEffect, useLayoutEffect, useRef, useState } from "react";

import Engine from "../game/Engine";
import { Event } from "../pb/event";

type Props = {
  clientId: string,
  host: string,
}

const Game: React.FC<Props> = ({ clientId, host }) => {
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

    const ws = new WebSocket(`ws://${host}/api/room/ws?token=${token}`);
    ws.binaryType = "arraybuffer";
    ws.onopen = async () => {
      await gameEngineRef.current?.init();
    };
    ws.onmessage = (event: MessageEvent) => {
      const message = new Uint8Array(event.data);
      gameEngineRef.current?.receive(Event.decode(message));
    };
    setSocket(ws);
  }, [host, socket]);

  useLayoutEffect(() => {
    if (!socket) {
      return;
    }

    const sketch = (instance: p5) => {
      gameEngineRef.current = new Engine(instance, clientId, host, socket);
    };

    const instance = new p5(sketch, containerRef.current!);
    return () => instance.remove();
  }, [clientId, host, socket]);

  return <div className="game__container" ref={containerRef} />;
};

export default Game;
