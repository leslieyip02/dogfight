package game

import "server/utils"

type Player struct {
	Id       string
	Username string
}

func NewPlayer(username string) (*Player, error) {
	id, err := utils.GetShortId()
	if err != nil {
		return nil, err
	}

	player := Player{Id: id, Username: username}
	return &player, nil
}
