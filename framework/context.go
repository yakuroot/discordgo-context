package framework

import "github.com/bwmarrin/discordgo"

type Context interface {
	Reply(options *discordgo.MessageSend) (*MessageContext, error)
	Edit(options *discordgo.MessageSend) (*MessageContext, error)
	Delete() error
}
