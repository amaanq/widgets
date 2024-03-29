package widgets

import (
	"github.com/bwmarrin/discordgo"
)

var (
	// []discordgo.MessageComponent{
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
	return
}

func ButtonsWithDelete() (C []discordgo.MessageComponent) {
	C = append([]discordgo.MessageComponent{DefaultButtons}, discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Style: discordgo.DangerButton,
				Emoji: discordgo.ComponentEmoji{
					Name: "❌",
				},
				CustomID: "delete",
			},
		},
	})
	return
}

func ButtonsDisabled(delete bool) (C []discordgo.MessageComponent) {
	components := []discordgo.MessageComponent{}
	for _, button := range DefaultButtons.Components {
		b := button.(discordgo.Button)
		if b.CustomID == "<<" || b.CustomID == "<" || b.CustomID == ">>" || b.CustomID == ">" || b.CustomID == "1234" {
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
					Style: discordgo.DangerButton,
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

func ButtonsFirstPage(delete bool) (C []discordgo.MessageComponent) {
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
					Style: discordgo.DangerButton,
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

func ButtonsMiddlePage(delete bool) (C []discordgo.MessageComponent) {
	C = []discordgo.MessageComponent{DefaultButtons}
	if delete {
		C = append(C, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Style: discordgo.DangerButton,
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

func ButtonsLastPage(delete bool) (C []discordgo.MessageComponent) {
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
					Style: discordgo.DangerButton,
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
