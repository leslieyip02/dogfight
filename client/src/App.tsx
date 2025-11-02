import "./App.css";

import { useState } from "react";

import Form from "./components/Form";
import Game from "./components/Game";

const App = () => {
  const [clientId, setClientId] = useState<string | null>(null);

  if (!clientId) {
    return <Form setClientId={setClientId} />;
  }

  return <Game clientId={clientId} />;
};

export default App;
