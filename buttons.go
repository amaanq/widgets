package widgets

import (
	"github.com/bwmarrin/discordgo"
)

var (
	// []discordgo.MessageComponent{
	delete         = false
	DefaultButtons = discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Style: discordgo.PrimaryButton,
				Emoji: discordgo.ComponentEmoji{
					Name: "⏪",
				},
				CustomID: "<<",
			},
			discordgo.Button{
				Style: discordgo.PrimaryButton,
				Emoji: discordgo.ComponentEmoji{
					Name: "◀️",
				},
				CustomID: "<",
			},
			discordgo.Button{
				Style: discordgo.PrimaryButton,
				Emoji: discordgo.ComponentEmoji{
					Name: "▶️",
				},
				CustomID: ">",
			},
			discordgo.Button{
				Style: discordgo.PrimaryButton,
				Emoji: discordgo.ComponentEmoji{
					Name: "⏩",
				},
				CustomID: ">>",
			},
			discordgo.Button{
				Style: discordgo.PrimaryButton,
				Emoji: discordgo.ComponentEmoji{
					Name: "🔢",
				},
				CustomID: "1234",
			},
		},
	}
)

func ButtonsWithoutDelete() (C []discordgo.MessageComponent) {
	C = []discordgo.MessageComponent{DefaultButtons}
	delete = false
	return
}

func ButtonsWithDelete() (C []discordgo.MessageComponent) {
	C = append([]discordgo.MessageComponent{DefaultButtons}, discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Style: 4,
				Emoji: discordgo.ComponentEmoji{
					Name: "❌",
				},
				CustomID: "delete",
			},
		},
	})
	delete = true
	return
}

func ButtonsFirstPage() (C []discordgo.MessageComponent) {
	components := []discordgo.MessageComponent{}
	for _, button := range DefaultButtons.Components {
		b := button.(discordgo.Button)
		if b.CustomID == "<<" || b.CustomID == "<" {
			b.Disabled = true
		}
		components = append(components, b)
	}
	C = []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: components,
		},
	}
	if delete {
		C = append(C, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Style: 4,
					Emoji: discordgo.ComponentEmoji{
						Name: "❌",
					},
					CustomID: "delete",
				},
			},
		})
	}
	return
}

func ButtonsMiddlePage() (C []discordgo.MessageComponent) {
	C = []discordgo.MessageComponent{DefaultButtons}
	if delete {
		C = append(C, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Style: 4,
					Emoji: discordgo.ComponentEmoji{
						Name: "❌",
					},
					CustomID: "delete",
				},
			},
		})
	}
	return
}

func ButtonsLastPage() (C []discordgo.MessageComponent) {
	components := []discordgo.MessageComponent{}
	for _, button := range DefaultButtons.Components {
		b := button.(discordgo.Button)
		if b.CustomID == ">>" || b.CustomID == ">" {
			b.Disabled = true
		}
		components = append(components, b)
	}
	C = []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: components,
		},
	}
	if delete {
		C = append(C, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Style: 4,
					Emoji: discordgo.ComponentEmoji{
						Name: "❌",
					},
					CustomID: "delete",
				},
			},
		})
	}
	return
}
