package widgets

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func GetInput(s *discordgo.Session, i *discordgo.InteractionCreate, message string, timeout time.Duration) (string, error) {
	var content string
	msg, err := s.ChannelMessageSend(i.ChannelID, content)
	if err != nil {
		return "", err
	}
	defer s.ChannelMessageDelete(msg.ChannelID, msg.ID)
	for {
		select {
		case usermsg := <-nextMessageCreateC(s):
			if usermsg.Author.ID != i.Member.User.ID {
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

	defer s.ChannelMessageDelete(msg.ChannelID, msg.ID)

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
			return i.MessageComponentData().CustomID, nil
		case <-timeoutChan:
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
