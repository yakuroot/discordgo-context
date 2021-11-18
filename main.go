package main

import (
	"log"
	"os"
	"os/signal"
	"strings"

	command "github.com/Neoration/discordgo-context/commands"
	"github.com/Neoration/discordgo-context/framework"
	"github.com/bwmarrin/discordgo"
)

var (
	Token   = "BOT-TOKEN"
	Session *discordgo.Session
	Command *framework.CommandHandler
	Prefix  = "!"
)

func init() {
	var err error
	Session, err = discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	err = Session.Open()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	Command = framework.CreateCommandHandler()
	registerCommands()

	log.Printf("Logged in %s (%s)", Session.State.User.String(), Session.State.User.ID)
}

func main() {
	Session.AddHandler(messageCreate)

	defer Session.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("봇 종료됨")
}

func registerCommands() {
	Command.RegisterCommand("messages", []string{"awaitmessages"}, command.AwaitMessagesExample)
	Command.RegisterCommand("interaction", []string{"interactions", "awaitinteractions"}, command.AwaitInteractionExample)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot ||
		len(m.Content) <= len(Prefix) ||
		m.Content[:len(Prefix)] != Prefix {
		return
	}

	args := strings.Fields(m.Content[len(Prefix):])
	commandName := strings.ToLower(args[0])

	if cmd, exist := Command.GetCommands(commandName); exist {
		cmd(framework.CreateNewMessageContext(s, m.Message))
	}
}

func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand ||
		i.User.Bot {
		return
	}
}
