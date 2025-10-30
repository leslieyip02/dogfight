import "./Form.css";

import { useMemo, useState } from "react";
import { GoAlert } from "react-icons/go";
import { TiArrowShuffle } from "react-icons/ti";
import { generateUsername } from "unique-username-generator";

import { joinRoom } from "../api/room";
import { chooseSpriteName } from "../game/utils/sprites";

type Props = {
  setClientId: (clientId: string) => void;
  setToken: (token: string) => void;
};
const Form: React.FC<Props> = ({ setClientId, setToken }) => {
  const [username, setUsername] = useState<string>(generateUsername("-"));
  const [roomId, setRoomId] = useState<string>("");

  const [shouldShake, setShouldShake] = useState<boolean>(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const shuffleUsername = () => {
    setUsername(generateUsername("-"));
  };

  const onSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    await joinRoom(username, roomId)
      .then(response => {
        setClientId(response.clientId);
        setToken(response.token);
      })
      .catch((error: Error) => {
        setErrorMessage(error.message);

        setShouldShake(true);
        setTimeout(() => {
          setShouldShake(false);
        }, 500);
      });
  };

  const spriteSrc = useMemo(() => {
    const spriteName = chooseSpriteName(username);
    return `${spriteName}.png`;
  }, [username]);

  return (
    <form className={`form ${shouldShake ? "form--shake" : ""}`} onSubmit={onSubmit}>
      <div className="form__header">
        <img className="form__sprite-preview" src={spriteSrc} />
        <h1>dogfight</h1>
      </div>

      <div className="form__field">
        <label htmlFor="username">username:</label><br />
        <div className="form__field-wrapper">
          <input
            type="text"
            id="username"
            name="username"
            required
            aria-required="true"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
          <button type="button" onClick={shuffleUsername}>
            <TiArrowShuffle />
          </button>
        </div>
      </div>
      <div className="form__field">
        <label htmlFor="roomId">room id:</label><br />
        <div className="form__field-wrapper">
          <input
            type="text"
            id="room id"
            name="roomId"
            aria-required="true"
            value={roomId}
            onChange={(e) => setRoomId(e.target.value)}
          />
        </div>
      </div>

      <button className="form__submit" type="submit">Join</button>

      {
        errorMessage && <div className="form__error-message" role="alert">
          <GoAlert /> {errorMessage}
        </div>
      }
    </form>
  );
};

export default Form;
