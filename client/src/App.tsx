import "./App.css";

import { useState } from "react";

import Form from "./components/Form";
import Game from "./components/Game";

const App = () => {
  const [clientId, setClientId] = useState<string | null>(null);
  const [host, setHost] = useState<string | null>(null);

  if (!clientId || !host) {
    return <Form setClientId={setClientId} setHost={setHost} />;
  }

  return <Game clientId={clientId} host={host} />;
};

export default App;
