package widgets

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestHandler(t *testing.T) {
	ButtonsWithDelete()
	err := AddButtonHandler(nil, &discordgo.MessageSend{}, discordgo.Button{
		Label:    "hello!",
		Style:    discordgo.SuccessButton,
		CustomID: "hello",
	}, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	})
	if err != nil {
		panic(err)
	}
}
