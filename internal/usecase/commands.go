package usecase

import (
	"discord_bot/internal/repository"
	"discord_bot/internal/utils"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	constandvar "discord_bot/internal/constant"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	htgotts "github.com/hegedustibor/htgo-tts"
)

func HandleMessageVoice(s *discordgo.Session, m *discordgo.MessageCreate) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           "default",
	}))
	messageJson, err := json.Marshal(m)
	if err != nil {
		fmt.Println("Error marshalling message to json: ", err)
	}
	queueURL := constandvar.QueueURL

	if !hasPrefix(m.Content, constandvar.PrefixSay) {
		fmt.Println("Message does not have prefix")
		return
	}

	// execute join function
	botJoinChannel(s, m)
	repository.SendMsg(sess, &queueURL, string(messageJson))
	if m.Author.ID == s.State.User.ID {
		fmt.Println("Bot is already in the voice channel")
		return
	}
	fmt.Printf("Message sent to queue with message:\"%s\"\n", m.Content)
	go say(s, m, sess, queueURL)
}

func say(s *discordgo.Session, m *discordgo.MessageCreate, sess *session.Session, queueURL string) {
	for {
		structInfo := &discordgo.MessageCreate{}
		messages, err := repository.GetMessages(sess, &queueURL)
		if err != nil {
			fmt.Println("Error getting messages: ", err)
			return
		}
		if len(messages.Messages) == 0 {
			fmt.Println("\nNo messages in queue")
			return
		}
		// convert message string to *discordgo.MessageCreate
		byteMessage := []byte(*messages.Messages[0].Body)
		err = json.Unmarshal(byteMessage, &structInfo)
		if err != nil {
			fmt.Println("Error while unmarshalling message to struct: ", err)
			return
		}

		structInfo.Message.Content = utils.DeleteTextInString(structInfo.Message.Content, constandvar.PrefixSay)
		// get name user and add to text
		structInfo.Message.Content = fmt.Sprintf("%s nói là: %s", structInfo.Member.Nick, structInfo.Message.Content)
		// delete message from channel
		err = s.ChannelMessageDelete(m.ChannelID, structInfo.ID)
		if err != nil {
			fmt.Println("Error deleting message from channel: ", err)
			return
		}

		message := utils.DeleteTextInString(structInfo.Message.Content, constandvar.PrefixSay)

		// rechat message with bold text
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("**%s**", message))

		// play audio and cleanup
		err = playAudioAndCleanup(s, m, structInfo)
		if err != nil {
			fmt.Println("Error playing audio and cleanup: ", err)
			return
		}
		// delete message
		err = repository.DeleteMessage(sess, &queueURL, messages.Messages[0].ReceiptHandle)
		if err != nil {
			fmt.Println("Error deleting message from queue: ", err)
		}
	}
}

func joinVoiceChannel(s *discordgo.Session, guildID, channelID string) error {
	_, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	return err
}

func hasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

func playAudioAndCleanup(s *discordgo.Session, m *discordgo.MessageCreate, structInfo *discordgo.MessageCreate) error {
	// code generate audio file from text and save it to the disk as "hello.mp3"
	speech := htgotts.Speech{Folder: "audio", Language: "vi"}
	text := strings.TrimPrefix(structInfo.Message.Content, constandvar.PrefixSay)
	filename := "output"

	// Convert the text to speech and save the result to the MP3 file
	_, err := speech.CreateSpeechFile(text, filename)
	if err != nil {
		fmt.Println("Error creating speech file: ", err)
		return err
	}
	// voice connection
	dgv, err := s.ChannelVoiceJoin(m.GuildID, m.ChannelID, false, true)
	if err != nil {
		fmt.Println("Error joining voice channel: ", err)
		return err
	}

	// play and wait for the playback to finish
	dgvoice.PlayAudioFile(dgv, fmt.Sprintf("%s/%s", "audio", "output.mp3"), make(chan bool))

	// delete audio file
	err = os.Remove(fmt.Sprintf("%s/%s", "audio", "output.mp3"))
	if err != nil {
		fmt.Println("Error deleting audio file: ", err)
	}
	return nil
}
func botJoinChannel(s *discordgo.Session, m *discordgo.MessageCreate) {
	// refactor this to a function
	authorInChannel := false
	for _, guild := range s.State.Guilds {
		for _, vs := range guild.VoiceStates {
			if vs.UserID == m.Author.ID {
				authorInChannel = true
				if err := joinVoiceChannel(s, guild.ID, vs.ChannelID); err != nil {
					fmt.Println("Error joining voice channel: ", err)
				}
				break
			}
		}
	}
	if !authorInChannel {
		s.ChannelMessageSend(m.ChannelID, "You need to be in a voice channel to use this command!")
	}
}
