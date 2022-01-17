package widgets

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Paginator struct {
	sync.Mutex
	Pages              []*discordgo.MessageEmbed
	Index              int
	DeleteWhenDone     bool
	Loop               bool
	AuthorizedToUse    []string // user IDs
	ChannelID          string
	MessageID          string
	Running            bool
	Timeout            time.Duration
	CurrentInteraction *discordgo.InteractionCreate
	Close              chan bool
	Session            *discordgo.Session
	cancel             func()
}

// Add a variadic amount of user IDs to allow access to the paginator
func (p *Paginator) AllowUsers(userIDs ...string) {
	p.Lock()
	p.AuthorizedToUse = append(p.AuthorizedToUse, userIDs...)
	p.Unlock()
}

// Pass in your discord bot session and at least one discord embed.
func NewPaginator(s *discordgo.Session, embeds ...*discordgo.MessageEmbed) *Paginator {
	if len(embeds) == 0 {
		return nil
	}
	return &Paginator{
		Pages:           embeds,
		Index:           0,
		DeleteWhenDone:  false,
		Loop:            false,
		AuthorizedToUse: nil,
		Running:         false,
		Session:         s,
	}
}

func (p *Paginator) Spawn(channelID string, content string) error {
	if p.Running {
		return fmt.Errorf("already running")
	}

	p.addHandler()
	msg, err := p.Session.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content:    content,
		Embed:      p.Pages[p.Index],
		Components: ButtonsFirstPage(),
	})

	if err != nil {
		p.cancel()
		p.close()
		return err
	}

	p.ChannelID = channelID
	p.MessageID = msg.ID
	p.Running = true

	go func() {
		time.Sleep(p.Timeout)
		p.cancel()
		p.close()
	}()

	return nil
}

func (p *Paginator) close() {
	if p.DeleteWhenDone {
		p.Session.ChannelMessageDelete(p.ChannelID, p.MessageID)
	}
	p = nil
}

func (p *Paginator) addHandler() {
	p.Lock()
	p.cancel = p.Session.AddHandler(p.defaultPaginatorHandler)
	p.Unlock()
}

// This will only handle button component interactions for a paginator
func (p *Paginator) defaultPaginatorHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.MessageComponentData().ComponentType != discordgo.ButtonComponent || i.Message.ID != p.MessageID {
		return
	}
	p.Lock()
	defer p.Unlock()
	components := ButtonsMiddlePage()
	switch i.MessageComponentData().CustomID {
	case ">":
		p.Index++
		if p.last() {
			components = ButtonsLastPage()
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content:    i.Message.Content,
				Components: components,
				Embeds:     []*discordgo.MessageEmbed{p.Pages[p.Index]},
			},
		})
		return
	case ">>":
		p.Index = len(p.Pages) - 1
		components = ButtonsLastPage()
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content:    i.Message.Content,
				Components: components,
				Embeds:     []*discordgo.MessageEmbed{p.Pages[p.Index]},
			},
		})
		return
	case "<<":
		p.Index = 0
		components = ButtonsFirstPage()
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content:    i.Message.Content,
				Components: components,
				Embeds:     []*discordgo.MessageEmbed{p.Pages[p.Index]},
			},
		})
		return
	case "<":
		p.Index--
		if p.first() {
			components = ButtonsFirstPage()
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content:    i.Message.Content,
				Components: components,
				Embeds:     []*discordgo.MessageEmbed{p.Pages[p.Index]},
			},
		})
		return
	case "1234":
		response, err := GetInput(s, i)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Timed out!",
					Flags:   64,
				},
			})
			return
		}
		n, err := strconv.Atoi(response)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Oops! That's not a number",
					Flags:   64,
				},
			})
			return
		}
		p.Index = n
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content:    i.Message.Content,
				Components: components,
				Embeds:     []*discordgo.MessageEmbed{p.Pages[p.Index]},
			},
		})
		return
	case "delete":
		s.ChannelMessageDelete(i.Message.ChannelID, i.Message.ID) // if it errors :shrug:
	}
}

func (p *Paginator) first() bool {
	return p.Index == 0
}
func (p *Paginator) last() bool {
	return p.Index == len(p.Pages)-1
}

func (p *Paginator) isAuthorized(userID string) bool {
	for _, uID := range p.AuthorizedToUse {
		if uID == userID {
			return true
		}
	}
	return false
}

func nextMessageInteractionCreateC(s *discordgo.Session) chan *discordgo.InteractionCreate {
	out := make(chan *discordgo.InteractionCreate)
	s.AddHandlerOnce(func(_ *discordgo.Session, e *discordgo.InteractionCreate) {
		out <- e
	})
	return out
}
