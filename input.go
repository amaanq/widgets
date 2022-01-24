package widgets

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func GetInput(s *discordgo.Session, channelID, userID, message string, timeout time.Duration) (string, error) {
	msg, err := s.ChannelMessageSend(channelID, message)
	if err != nil {
		return "", err
	}
	defer s.ChannelMessageDelete(msg.ChannelID, msg.ID)
	for {
		select {
		case usermsg := <-nextMessageCreateC(s):
			if usermsg.Author.ID != userID {
				continue
			}
			s.ChannelMessageDelete(usermsg.ChannelID, usermsg.ID)
			return usermsg.Message.Content, nil
		case <-time.After(timeout):
			return "", fmt.Errorf("timed out")
		}
	}
}

// your message must already have the options, with the custom ID having the choice selected, these MUST be unique
func GetInputFromInteraction(s *discordgo.Session, channelID, userID string, message *discordgo.MessageSend, timeout time.Duration) (string, error) {
	msg, err := s.ChannelMessageSendComplex(channelID, message)
	if err != nil {
		return "", err
	}

	//components := msg.Components[0].(*discordgo.ActionsRow).Components
	newComponents := []discordgo.MessageComponent{}

	currActionRow := discordgo.ActionsRow{}
	mainComponents := []discordgo.MessageComponent{}

	timeoutChan := make(chan int)
	go func() {
		time.Sleep(timeout)
		timeoutChan <- 0
	}()

	for {
		select {
		case i := <-nextInteractionCreateC(s):
			if i.Member.User.ID != userID {
				continue
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
			})
			for _, mainComponent := range msg.Components {
				components := mainComponent.(*discordgo.ActionsRow).Components
				for _, comp := range components {
					button := comp.(*discordgo.Button)
					if i.MessageComponentData().CustomID == button.CustomID {
						button.Style = discordgo.SuccessButton
					}
					button.Disabled = true
					newComponents = append(newComponents, button)
				}
				currActionRow.Components = append(currActionRow.Components, newComponents...)
				mainComponents = append(mainComponents, currActionRow)
				newComponents = []discordgo.MessageComponent{}
				currActionRow = discordgo.ActionsRow{}
			}
			defer s.ChannelMessageEditComplex(&discordgo.MessageEdit{
				Content:    &msg.Content,
				Components: mainComponents,
				Embeds:     msg.Embeds,

				ID:      msg.ID,
				Channel: msg.ChannelID,
			})
			return i.MessageComponentData().CustomID, nil
		case <-timeoutChan:
			for _, mainComponent := range msg.Components {
				components := mainComponent.(*discordgo.ActionsRow).Components
				for _, comp := range components {
					button := comp.(*discordgo.Button)
					button.Disabled = true
					newComponents = append(newComponents, button)
				}
				currActionRow.Components = append(currActionRow.Components, newComponents...)
				mainComponents = append(mainComponents, currActionRow)
				newComponents = []discordgo.MessageComponent{}
				currActionRow = discordgo.ActionsRow{}
			}
			defer s.ChannelMessageEditComplex(&discordgo.MessageEdit{
				Content:    &msg.Content,
				Components: mainComponents,
				Embeds:     msg.Embeds,

				ID:      msg.ID,
				Channel: msg.ChannelID,
			})
			return "", fmt.Errorf("timed out")
		}
	}
}

// Credits https://github.com/Necroforger/dgwidgets/blob/master/util.go
func nextMessageCreateC(s *discordgo.Session) chan *discordgo.MessageCreate {
	out := make(chan *discordgo.MessageCreate)
	s.AddHandlerOnce(func(_ *discordgo.Session, e *discordgo.MessageCreate) {
		out <- e
	})
	return out
}

func nextInteractionCreateC(s *discordgo.Session) chan *discordgo.InteractionCreate {
	out := make(chan *discordgo.InteractionCreate)
	s.AddHandlerOnce(func(_ *discordgo.Session, e *discordgo.InteractionCreate) {
		out <- e
	})
	return out
}
