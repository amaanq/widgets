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
	Pages               []*discordgo.MessageSend
	Index               int
	DeleteWhenDone      bool
	Loop                bool
	AuthorizedToUse     []string // user IDs
	ChannelID           string
	MessageID           string
	Running             bool
	Timeout             time.Duration
	CurrentInteraction  *discordgo.InteractionCreate
	Close               chan bool
	Session             *discordgo.Session
	cancel              func()
	customWidgetButtons []func() error
	errHandler          func(error)
}

// Add a variadic amount of user IDs to allow access to the paginator
func (p *Paginator) AllowUsers(userIDs ...string) {
	p.Lock()
	p.AuthorizedToUse = append(p.AuthorizedToUse, userIDs...)
	p.Unlock()
}

// Pass in your discord bot session and at least one discord embed.
func NewPaginator(s *discordgo.Session, messages ...*discordgo.MessageSend) *Paginator {
	return &Paginator{
		Pages:               messages,
		Index:               0,
		DeleteWhenDone:      false,
		Loop:                false,
		AuthorizedToUse:     nil,
		Running:             false,
		Session:             s,
		customWidgetButtons: []func() error{},
		errHandler:          nil,
	}
}

func (p *Paginator) AddPages(messages ...*discordgo.MessageSend) {
	p.Lock()
	p.Pages = append(p.Pages, messages...)
	p.Unlock()
}

func (p *Paginator) SetErrHandler(errfunc func(e error)) {
	p.Lock()
	p.errHandler = errfunc
	p.Unlock()
}

func (p *Paginator) Spawn(channelID string) error {
	if p.Running {
		return fmt.Errorf("already running")
	}
	if len(p.Pages) < 1 {
		return fmt.Errorf("a minimum of one page is required")
	}
	p.addHandler()
	p.Pages[0].Components = ButtonsFirstPage()
	if len(p.Pages) == 1 {
		p.Pages[0].Components = ButtonsDisabled()
	}

	for _, fnc := range p.customWidgetButtons {
		err := fnc()
		if err != nil {
			if p.errHandler != nil {
				p.errHandler(err)
			}
		}
	}

	msg, err := p.Session.ChannelMessageSendComplex(channelID, p.Pages[p.Index])

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
	if !p.isAuthorized(i.Member.User.ID) {
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
		nextMessage := p.Pages[p.Index]
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				TTS:             nextMessage.TTS,
				Content:         nextMessage.Content,
				Components:      components,
				Embeds:          nextMessage.Embeds,
				AllowedMentions: nextMessage.AllowedMentions,
				Files:           nextMessage.Files,
			},
		})
		return
	case ">>":
		p.Index = len(p.Pages) - 1
		components = ButtonsLastPage()
		nextMessage := p.Pages[p.Index]
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				TTS:             nextMessage.TTS,
				Content:         nextMessage.Content,
				Components:      components,
				Embeds:          nextMessage.Embeds,
				AllowedMentions: nextMessage.AllowedMentions,
				Files:           nextMessage.Files,
			},
		})
		return
	case "<<":
		p.Index = 0
		components = ButtonsFirstPage()
		nextMessage := p.Pages[p.Index]
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				TTS:             nextMessage.TTS,
				Content:         nextMessage.Content,
				Components:      components,
				Embeds:          nextMessage.Embeds,
				AllowedMentions: nextMessage.AllowedMentions,
				Files:           nextMessage.Files,
			},
		})
		return
	case "<":
		p.Index--
		if p.first() {
			components = ButtonsFirstPage()
		}
		nextMessage := p.Pages[p.Index]
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				TTS:             nextMessage.TTS,
				Content:         nextMessage.Content,
				Components:      components,
				Embeds:          nextMessage.Embeds,
				AllowedMentions: nextMessage.AllowedMentions,
				Files:           nextMessage.Files,
			},
		})
		return
	case "1234":
		if !p.isAuthorized(i.Member.User.ID) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Oops! This command wasn't from you!",
					Flags:   64,
				},
			})
		}
		response, err := GetInput(s, i, fmt.Sprintf("%s, enter the page you'd like to go to (from %d to %d inclusive)", i.Member.Mention(), 1, len(p.Pages)), time.Second*30)
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
		n-- // 0 indexed
		if p.outOfBounds(n) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "That's too low or high of a page number, try something smaller!",
					Flags:   64,
				},
			})
			return
		}
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
		nextMessage := p.Pages[p.Index]
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				TTS:             nextMessage.TTS,
				Content:         nextMessage.Content,
				Components:      components,
				Embeds:          nextMessage.Embeds,
				AllowedMentions: nextMessage.AllowedMentions,
				Files:           nextMessage.Files,
			},
		})
		return
	case "delete":
		s.ChannelMessageDelete(i.Message.ChannelID, i.Message.ID)
	}
}

func (p *Paginator) first() bool {
	return p.Index == 0
}
func (p *Paginator) last() bool {
	return p.Index == len(p.Pages)-1
}

func (p *Paginator) isAuthorized(userID string) bool {
	if len(p.AuthorizedToUse) == 0 {
		return true
	}
	for _, uID := range p.AuthorizedToUse {
		if uID == userID {
			return true
		}
	}
	return false
}

func (p *Paginator) outOfBounds(i int) bool {
	return i >= len(p.Pages) || i < 0
}

func (p *Paginator) SetPageFooters() {
	for index, msg := range p.Pages {
		txt := fmt.Sprintf("#[%d / %d]", index+1, len(p.Pages))
		if msg.Embeds[0].Footer != nil {
			txt += msg.Embeds[0].Footer.Text
		}
		msg.Embeds[0].Footer = &discordgo.MessageEmbedFooter{
			Text: txt,
		}
	}
}
