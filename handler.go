package widgets

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

type ButtonHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

func AddButtonHandler(msg *discordgo.MessageSend, button *discordgo.Button, callback ButtonHandler) error {
	var s, last int = 0, 0
	actionRow := msg.Components[0].(discordgo.ActionsRow)
	for _, b := range actionRow.Components {
		button := b.(discordgo.Button)
	}
	for idx, component := range msg.Components {
		actionRow := component.(discordgo.ActionsRow)
		s += len(actionRow.Components)
		last = idx 
	}
	if s == 25 {
		return errors.New("Max number of buttons allowed")
	}
	if last != 4 {
		
	}
	return nil 
}
