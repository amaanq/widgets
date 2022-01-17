package widgets

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func GetInput(s *discordgo.Session, i *discordgo.InteractionCreate, message string) (string, error) {
	var content string
	msg, err := s.ChannelMessageSend(i.ChannelID, content)
	if err != nil {
		return "", err
	}
	defer s.ChannelMessageDelete(msg.ChannelID, msg.ID)
	for {
		select {
		case usermsg := <-nextMessageCreateC(s):
			return usermsg.Content, nil
		case <-time.After(time.Second * 20):
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
