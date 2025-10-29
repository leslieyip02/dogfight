import "./App.css";

import { useState } from "react";

import Form from "./components/Form";
import Game from "./components/Game";

const App = () => {
  const [clientId, setClientId] = useState<string | null>(null);
  const [token, setToken] = useState<string | null>(null);

  if (!clientId || !token) {
    return <Form setClientId={setClientId} setToken={setToken} />;
  }

  return <Game clientId={clientId} token={token} />;
};

export default App;
