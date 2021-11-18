package framework

import (
	"github.com/bwmarrin/discordgo"
)

type InteractionContext struct {
	Session     *discordgo.Session
	Guild       *discordgo.Guild
	Channel     *discordgo.Channel
	Author      *discordgo.User
	Member      *discordgo.Member
	Interaction *discordgo.Interaction
}

func CreateNewInteractionContext(s *discordgo.Session, i *discordgo.Interaction) *InteractionContext {
	ctx := new(InteractionContext)
	ctx.Session = s
	ctx.Guild, _ = s.State.Guild(i.GuildID)
	ctx.Channel, _ = s.State.Channel(i.ChannelID)
	ctx.Author = i.User
	ctx.Member = i.Member
	ctx.Interaction = i

	return ctx
}

func (ctx InteractionContext) Reply(options *discordgo.MessageSend) (*MessageContext, error) {
	err := ctx.Session.InteractionRespond(ctx.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: options.Content,
			//Flags:      1 << 6,
			Embeds:     options.Embeds,
			Components: options.Components,
		},
	})

	if err != nil {
		return nil, err
	}

	msg, err := ctx.Session.InteractionResponse(ctx.Session.State.User.ID, ctx.Interaction)

	if err != nil {
		return nil, err
	}

	return CreateNewMessageContext(ctx.Session, msg), nil
}

func (ctx *InteractionContext) Edit(options *discordgo.MessageSend) (*MessageContext, error) {
	msg, err := ctx.Session.InteractionResponseEdit(ctx.Session.State.User.ID, ctx.Interaction, &discordgo.WebhookEdit{
		Content:    options.Content,
		Embeds:     options.Embeds,
		Components: options.Components,
	})

	if err != nil {
		return nil, err
	}

	return CreateNewMessageContext(ctx.Session, msg), nil
}

func (ctx *InteractionContext) Update(options *discordgo.MessageSend) (*MessageContext, error) {
	err := ctx.Session.InteractionRespond(ctx.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: options.Content,
			//Flags:      1 << 6,
			Embeds:     options.Embeds,
			Components: options.Components,
		},
	})

	if err != nil {
		return nil, err
	}

	return CreateNewMessageContext(ctx.Session, ctx.Interaction.Message), nil
}

func (ctx *InteractionContext) Delete() error {
	err := ctx.Session.InteractionResponseDelete(ctx.Session.State.User.ID, ctx.Interaction)

	if err != nil {
		return err
	}

	return nil
}
