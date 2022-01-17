package widgets

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestHandler(t *testing.T) {
	ButtonsWithDelete()
	AddButtonHandler(nil, &discordgo.MessageSend{
		Components: ButtonsFirstPage(),
	}, discordgo.Button{
		Label:    "hello!",
		Style:    discordgo.SuccessButton,
		CustomID: "hello",
	}, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	})
}
