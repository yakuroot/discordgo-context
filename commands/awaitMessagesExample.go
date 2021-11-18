package command

import (
	"strings"

	"github.com/Neoration/discordgo-context/framework"
	"github.com/bwmarrin/discordgo"
)

func AwaitMessagesExample(ctx framework.Context) {
	msg, err := ctx.Reply(&discordgo.MessageSend{
		Content: "Please send at least 10 messages in 30 seconds with more than 5 letters.",
	})

	if err != nil {
		return
	}

	msgs, actived := msg.AwaitMessages(framework.AwaitMessagesOptions{
		Filter: func(m *discordgo.Message) bool {
			return len(m.Content) > 5
		},
		Max:  10,
		Time: 30e3,
	})

	if !actived {
		msg.Edit(&discordgo.MessageSend{
			Content: "Timeout",
		})
		return
	}

	result := make([]string, 0)
	for _, v := range msgs {
		result = append(result, v.Message.Content)
	}

	msg.Edit(&discordgo.MessageSend{
		Content: "Sensed message: " + strings.Join(result, ", "),
	})
}
