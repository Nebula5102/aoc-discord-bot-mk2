package discord

import (
	"github.com/Nebula5102/aoc-discord-bot-mk2/internal/config"
	"github.com/Nebula5102/aoc-discord-bot-mk2/internal/leaderboard"
	"github.com/bwmarrin/discordgo"

	"log"
	"strings"
	"time"
	"regexp"
)

type BotHandler struct {
	Session *discordgo.Session
	Tracker *leaderboard.Tracker
	cfg     *config.Config
}

func NewBotHandler(session *discordgo.Session, tracker *leaderboard.Tracker, cfg *config.Config) *BotHandler {
	return &BotHandler{
		Session: session,
		Tracker: tracker,
		cfg:     cfg,
	}
}

func (bh *BotHandler) CheckForUpdates() (bool, error) {
	log.Println("Checking for updates...")

	bh.Tracker.LastUpdate = time.Now()
	bh.Tracker.UpdateLeaderboard()
	leaderboard.StoreLeaderboard(bh.Tracker.CurrentLeaderboard)

	hadUpdates := false
	newStars, err := bh.Tracker.CheckForNewStars()
	if err != nil {
		return hadUpdates, err
	}

	newMembers, err := bh.Tracker.CheckForNewMembers()
	if err != nil {
		return hadUpdates, err
	}

	if len(newStars) > 0 {
		log.Printf("new stars: %v", newStars)
		for _, member := range newStars {
			bh.SendChannelMessage(bh.cfg.ChannelID, member+" got a star! 🌟")
		}
	}

	if len(newMembers) > 0 {
		log.Printf("new members: %v", newMembers)
		bh.SendChannelMessage(bh.cfg.ChannelID, "CHALLENGER APPROACHING!")
		for _, member := range newMembers {
			bh.SendChannelMessage(bh.cfg.ChannelID, member+" has joined the leaderboard!")
		}
	}

	if len(newStars) > 0 || len(newMembers) > 0 {
		hadUpdates = true
		formattedLeaderboard := leaderboard.FormatLeaderboard(bh.Tracker.CurrentLeaderboard)
		bh.SendChannelMessageEmbed(bh.cfg.ChannelID, formattedLeaderboard)
	}

	return hadUpdates, nil
}

func (bh *BotHandler) MessageReceived(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.ChannelID != bh.cfg.ChannelID {
		return
	}

	signup := regexp.MustCompile(`!signup`)
	start := regexp.MustCompile(`!start`)
	end := regexp.MustCompile(`!end`)

	res := "" 
	for i := 0; i<1; i++ {
		res = signup.FindString(strings.ToLower(m.Content))
		if res != "" {break}
		res = start.FindString(strings.ToLower(m.Content))
		if res != "" {break}
		res = end.FindString(strings.ToLower(m.Content))
		if res != "" {break}
	}

	if strings.ToLower(m.Content) == "!update" {
		log.Println("Update command received")
		if time.Since(bh.Tracker.LastUpdate).Minutes() > 0 {
			hadUpdates, err := bh.CheckForUpdates()
			if err != nil {
				log.Printf("error checking for updates: %v", err)
			}
			if !hadUpdates {
				bh.SendChannelMessage(bh.cfg.ChannelID, "No updates")
			}
		} else {
			bh.SendChannelMessage(bh.cfg.ChannelID, "You can only update once every 15 minutes")
		}
	} else if strings.ToLower(m.Content) == "!leaderboard" {
		log.Println("Leaderboard command received")
		formattedLeaderboard := leaderboard.FormatLeaderboard(bh.Tracker.CurrentLeaderboard)
		bh.SendChannelMessageEmbed(bh.cfg.ChannelID, formattedLeaderboard)

	} else if strings.ToLower(m.Content) == "!stars" {
		log.Println("Stars command received")
		embed := leaderboard.FormatStars(bh.Tracker.CurrentLeaderboard)
		bh.SendChannelMessageEmbed(bh.cfg.ChannelID, embed)

	} else if strings.ToLower(m.Content) == "!help" {
		sb := strings.Builder{}
		sb.WriteString("```")
		sb.WriteString("Commands:\n\n")
		sb.WriteString("!leaderboard - Shows the current Advent of Code leaderboard\n\n")
		sb.WriteString("!update - Checks for updates and shows the updated leaderboard\n\n")
		sb.WriteString("!stars - Shows the current stars\n\n")
		sb.WriteString("!comp - Displays current competition leaderboard, if there is one\n\n")
		sb.WriteString("!signup<AOCid> - Signs you up to the competition, if there is one\n\n")
		sb.WriteString("!start<Day Number> - Sets start time for the AOC day challenge\n\n")
		sb.WriteString("!end<Day Number> - Sets end time for the AOC day challenge\n\n")
		sb.WriteString("!help - Shows this message\n")
		sb.WriteString("```")
		bh.SendChannelMessage(bh.cfg.ChannelID, sb.String())
	} else if strings.ToLower(m.Content) == "!comp" {
		log.Println("Competition command received")
		sb := strings.Builder{}
		sb.WriteString("```")
		sb.WriteString("Competition has not started\n\n")
		sb.WriteString("```")
		bh.SendChannelMessage(bh.cfg.ChannelID,sb.String())
	} else if res == "!signup" {
		log.Println("Signup command received")
		re := regexp.MustCompile(`<([^>]+)>`)
		AOCUser := re.FindStringSubmatch(m.Content)
		log.Println("Update table with Discord user:",m.Author,"AOCUser:",AOCUser[1])
	} else if res == "!start" {
		log.Println("Start command received")
		re := regexp.MustCompile(`<([^>]+)>`)
		day := re.FindStringSubmatch(m.Content)
		log.Println("Update table with day time:",time.Now(),"Day challeng #:",day[1])
	} else if res == "!end" {
		log.Println("End command received")
		re := regexp.MustCompile(`<([^>]+)>`)
		day := re.FindStringSubmatch(m.Content)
		log.Println("Grab time from table and perform t.Sub(starttime)")
		log.Println("Update table with completion time:",time.Now(),"Day challeng #:",day[1])
	} 
}

func (bh *BotHandler) SendChannelMessage(channelID, message string) {
	_, err := bh.Session.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Printf("error sending message: %v", err)
	}
}

func (bh *BotHandler) SendChannelMessageEmbed(channelID string, embed *discordgo.MessageEmbed) {
	_, err := bh.Session.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		log.Printf("error sending message: %v", err)
	}
}
