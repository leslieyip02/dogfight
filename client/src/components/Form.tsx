import "./Form.css";

import { useState } from "react";
import { generateUsername } from "unique-username-generator";

import { joinRoom } from "../api/room";

type Props = {
  setClientId: (clientId: string) => void;
  setToken: (token: string) => void;
};

const Form: React.FC<Props> = ({ setClientId, setToken }) => {
  const [username, setUsername] = useState<string>(generateUsername("-"));
  const [roomId, setRoomId] = useState<string>("");

  const onSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    await joinRoom(username, roomId)
      .then(response => {
        setClientId(response.clientId);
        setToken(response.token);
      })
      .catch(error => {
        // TODO: feedback to user
        console.log(error);
      });
  };

  return (
    <>
      <form className="form" onSubmit={onSubmit}>
        <h1 className="form__header">dogfight</h1>
        <div className="form__field">
          <label htmlFor="username">username:</label><br />
          <input
            type="text"
            id="username"
            name="username"
            required
            aria-required="true"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
        </div>
        <div className="form__field">
          <label htmlFor="roomId">room id:</label><br />
          <input
            type="text"
            id="room id"
            name="roomId"
            aria-required="true"
            value={roomId}
            onChange={(e) => setRoomId(e.target.value)}
          />
        </div>
        <button className="form__submit" type="submit">Join</button>
      </form>
    </>
  );

};

export default Form;
