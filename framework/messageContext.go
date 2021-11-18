package framework

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type MessageContext struct {
	Session *discordgo.Session
	Guild   *discordgo.Guild
	Channel *discordgo.Channel
	Author  *discordgo.User
	Member  *discordgo.Member
	Message *discordgo.Message
}

func CreateNewMessageContext(s *discordgo.Session, m *discordgo.Message) *MessageContext {
	ctx := new(MessageContext)
	ctx.Session = s
	ctx.Guild, _ = s.State.Guild(m.GuildID)
	ctx.Channel, _ = s.State.Channel(m.ChannelID)
	ctx.Author = m.Author
	ctx.Member = m.Member
	ctx.Message = m

	return ctx
}

func (ctx MessageContext) Reply(options *discordgo.MessageSend) (*MessageContext, error) {
	msg, err := ctx.Session.ChannelMessageSendComplex(ctx.Channel.ID, options)

	if err != nil {
		return nil, err
	}

	return CreateNewMessageContext(ctx.Session, msg), nil
}

func (ctx *MessageContext) Edit(options *discordgo.MessageSend) (*MessageContext, error) {
	msg, err := ctx.Session.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Content:         &options.Content,
		Embeds:          options.Embeds,
		Components:      options.Components,
		AllowedMentions: options.AllowedMentions,
		ID:              ctx.Message.ID,
		Channel:         ctx.Channel.ID,
	})

	if err != nil {
		return nil, err
	}

	ctx.Message = msg
	return ctx, nil
}

func (ctx *MessageContext) Delete() error {
	err := ctx.Session.ChannelMessageDelete(ctx.Channel.ID, ctx.Message.ID)

	if err != nil {
		return err
	}

	return nil
}

type AwaitComponentFilter func(i *discordgo.Interaction) bool

type AwaitComponentOptions struct {
	Filter AwaitComponentFilter
	Time   int
}

func (ctx *MessageContext) AwaitMessageComponent(options AwaitComponentOptions) (*InteractionContext, bool) {
	interaction := make(chan *discordgo.Interaction)
	activated := false

	go func() {
		time.Sleep(time.Millisecond * time.Duration(options.Time))
		interaction <- nil
		return
	}()

	go func() {
		for {
			complete := make(chan *discordgo.Interaction)

			ctx.Session.AddHandlerOnce(
				func(_ *discordgo.Session, i *discordgo.InteractionCreate) {
					if i.Interaction.Message.ID != ctx.Message.ID ||
						i.Type != discordgo.InteractionMessageComponent {
						complete <- nil
						return
					}

					if filter := options.Filter(i.Interaction); !filter {
						complete <- nil
						return
					}

					complete <- i.Interaction
				},
			)

			if r := <-complete; r != nil {
				interaction <- r
				activated = true
				return
			}

			continue
		}
	}()

	result := <-interaction

	if activated {
		return CreateNewInteractionContext(ctx.Session, result), true
	}

	return nil, false
}

type AwaitMessagesFilter func(m *discordgo.Message) bool

type AwaitMessagesOptions struct {
	Filter AwaitMessagesFilter
	Max    int
	Time   int
}

func (ctx *MessageContext) AwaitMessages(options AwaitMessagesOptions) ([]*MessageContext, bool) {
	timeout := make(chan bool)
	collector := make([]*MessageContext, 0)

	go func() {
		time.Sleep(time.Millisecond * time.Duration(options.Time))
		timeout <- true
		return
	}()

	go func() {

		for {
			message := make(chan *discordgo.Message)

			ctx.Session.AddHandlerOnce(
				func(_ *discordgo.Session, m *discordgo.MessageCreate) {
					if filter := options.Filter(m.Message); filter {
						message <- m.Message
					} else {
						message <- nil
					}
				},
			)

			if r := <-message; r != nil {
				collector = append(collector, CreateNewMessageContext(ctx.Session, r))
			}

			if len(collector) == options.Max {
				timeout <- true
				return
			}

			continue
		}
	}()

	<-timeout

	if len(collector) == 0 {
		return nil, false
	}

	return collector, true
}
