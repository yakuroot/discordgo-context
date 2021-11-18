package command

import (
	"github.com/Neoration/discordgo-context/framework"
	"github.com/bwmarrin/discordgo"
)

func AwaitInteractionExample(ctx framework.Context) {
	getComponents := func(disabled bool) []discordgo.MessageComponent {
		return []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Click Here!",
						Style:    discordgo.PrimaryButton,
						CustomID: "custom-id",
						Disabled: disabled,
					},
				},
			},
		}
	}

	msg, err := ctx.Reply(&discordgo.MessageSend{
		Content:    "Press the button below within 10 seconds.",
		Components: getComponents(false),
	})

	if err != nil {
		return
	}

	interaction, active := msg.AwaitMessageComponent(framework.AwaitComponentOptions{
		Filter: func(i *discordgo.Interaction) bool {
			return i.MessageComponentData().CustomID == "custom-id"
		},
		Time: 10e3,
	})

	if !active {
		msg.Edit(&discordgo.MessageSend{
			Content:    "Timeout.",
			Components: getComponents(true),
		})

		return
	}

	interaction.Update(&discordgo.MessageSend{
		Content:    "Button!",
		Components: getComponents(true),
	})
}
