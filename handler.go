package widgets

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type ButtonHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

func AddButtonHandler(ses *discordgo.Session, msg *discordgo.MessageSend, button discordgo.Button, callback ButtonHandler) error {
	if ses == nil || msg == nil {
		return errors.New("session and message cannot be nil")
	}

	var s, last int = 0, 0
	for idx, component := range msg.Components {
		actionRow := component.(discordgo.ActionsRow)
		s += len(actionRow.Components)
		last = idx
	}
	if s == 25 {
		return errors.New("max number of buttons allowed")
	}
	if s == 0 {
		actionRow := discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{button},
		}
		msg.Components = []discordgo.MessageComponent{actionRow}
	} else if s%5 == 0 {
		actionRow := msg.Components[last].(discordgo.ActionsRow)
		actionRow.Components = append(actionRow.Components, button)
		msg.Components[last] = actionRow
	} else {
		actionRow := msg.Components[last].(discordgo.ActionsRow)
		actionRow.Components = append(actionRow.Components, button)
		msg.Components = append(msg.Components[:last], actionRow)
	}
	ses.AddHandler(callback)
	return nil
}
