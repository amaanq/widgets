package widgets

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

func (p *Paginator) AddButtonHandler(ses *discordgo.Session /*msg *discordgo.MessageSend,*/, button discordgo.Button, callback func(s *discordgo.Session, i *discordgo.InteractionCreate)) error {
	if ses == nil /* || msg == nil */ {
		return errors.New("session and message cannot be nil")
	}
	f := func() error {
		var s, last int = 0, 0
		// for idx, component := range msg.Components {
		// 	actionRow := component.(discordgo.ActionsRow)
		// 	s += len(actionRow.Components)
		// 	last = idx
		// }
		// if s == 25 {
		// 	return errors.New("max number of buttons allowed")
		// }
		// if s == 0 {
		// 	actionRow := discordgo.ActionsRow{
		// 		Components: []discordgo.MessageComponent{button},
		// 	}
		// 	msg.Components = []discordgo.MessageComponent{actionRow}
		// } else if s%5 == 0 {
		// 	actionRow := msg.Components[last].(discordgo.ActionsRow)
		// 	actionRow.Components = append(actionRow.Components, button)
		// 	msg.Components[last] = actionRow
		// } else {
		// 	actionRow := msg.Components[last].(discordgo.ActionsRow)
		// 	actionRow.Components = append(actionRow.Components, button)
		// 	msg.Components = append(msg.Components[:last], actionRow)
		// }
		p.Lock()
		for _, page := range p.Pages {
			for idx, component := range page.Components {
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
				page.Components = []discordgo.MessageComponent{actionRow}
			} else if s%5 == 0 {
				actionRow := page.Components[last].(discordgo.ActionsRow)
				actionRow.Components = append(actionRow.Components, button)
				page.Components[last] = actionRow
			} else {
				actionRow := page.Components[last].(discordgo.ActionsRow)
				actionRow.Components = append(actionRow.Components, button)
				page.Components = append(page.Components[:last], actionRow)
			}
		}
		ses.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if i.MessageComponentData().CustomID == button.CustomID {
				callback(s, i)
			}
		})
		return nil
	}
	p.Lock()
	p.customWidgetButtons = append(p.customWidgetButtons, f)
	p.Unlock()
	return nil
}
