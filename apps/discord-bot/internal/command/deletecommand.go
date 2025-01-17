package command

import (
	"context"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/senchabot-opensource/monorepo/apps/discord-bot/internal/service"
	"github.com/senchabot-opensource/monorepo/apps/discord-bot/internal/service/streamer"
	"github.com/senchabot-opensource/monorepo/config"
	"github.com/senchabot-opensource/monorepo/packages/gosenchabot"
)

func (c *commands) DeleteCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, service service.Service) {
	options := i.ApplicationCommandData().Options

	switch options[0].Name {
	case "stream-anno-default-channel":
		ok, err := service.DeleteDiscordBotConfig(ctx, i.GuildID, "stream_anno_default_channel")
		if err != nil {
			log.Printf("Error while deleting Discord bot config: %v", err)
			ephemeralRespond(s, i, config.ErrorMessage+"#0001")
			return
		}

		if !ok {
			ephemeralRespond(s, i, config.ErrorMessage+"#0002")
			return
		}
		ephemeralRespond(s, i, "Varsayılan Twitch canlı yayın duyuru kanalı ayarı kaldırıldı.")

	case "stream-anno-default-content":
		_, err := service.SetDiscordBotConfig(ctx, i.GuildID, "stream_anno_default_content", "")
		if err != nil {
			log.Printf("Error while setting Discord bot config: %v", err)
			ephemeralRespond(s, i, config.ErrorMessage+"#0001")
			return
		}

		ephemeralRespond(s, i, "Yayın duyuru mesajı içeriği varsayılan olarak ayarlandı: `{stream.user}, {stream.category} yayınına başladı! {stream.url}`")
	case "stream-anno-custom-content":
		options = options[0].Options
		twitchUsername := options[0].StringValue()
		twitchUsername = gosenchabot.ParseTwitchUsernameURLParam(twitchUsername)

		ok, err := service.UpdateTwitchStreamerAnnoContent(ctx, twitchUsername, i.GuildID, nil)
		if err != nil {
			log.Printf("Error while deleting streamer announcement custom content: %v", err)
			ephemeralRespond(s, i, config.ErrorMessage+"#0001")
			return
		}

		if !ok {
			ephemeralRespond(s, i, config.ErrorMessage+"#0001")
			return
		}

		ephemeralRespond(s, i, twitchUsername+" kullanıcı adlı Twitch yayıncısına özgü yayın duyuru mesajı silindi.")

	case "streamer":
		options = options[0].Options
		twitchUsername := options[0].StringValue()
		twitchUsername = gosenchabot.ParseTwitchUsernameURLParam(twitchUsername)

		response0, uInfo := streamer.GetTwitchUserInfo(twitchUsername, c.twitchAccessToken)
		if response0 != "" {
			ephemeralRespond(s, i, response0)
			return
		}

		ok, err := service.DeleteDiscordTwitchLiveAnno(ctx, uInfo.ID, i.GuildID)
		if err != nil {
			ephemeralRespond(s, i, config.ErrorMessage+"#XXXX")
			return
		}

		if !ok {
			ephemeralRespond(s, i, "`"+twitchUsername+"` kullanıcı adlı Twitch yayıncısı veritabanında bulunamadı.")
			return
		}

		streamers := streamer.GetStreamersData(i.GuildID)
		delete(streamers, uInfo.Login)
		ephemeralRespond(s, i, "`"+uInfo.Login+"` kullanıcı adlı Twitch streamer veritabanından silindi.")

	case "stream-event-channel":
		options = options[0].Options
		channelId := options[0].ChannelValue(s).ID
		channelName := options[0].ChannelValue(s).Name

		ok, err := service.DeleteAnnouncementChannel(ctx, channelId)
		if err != nil {
			ephemeralRespond(s, i, config.ErrorMessage+"#XXYX")
			log.Println("Error while deleting announcement channel:", err)
			return
		}
		if !ok {
			ephemeralRespond(s, i, "`"+channelName+"` isimli yazı kanalı yayın etkinlik yazı kanalları listesinde bulunamadı.")
			return
		}
		ephemeralRespond(s, i, "`"+channelName+"` isimli yazı kanalı yayın etkinlik yazı kanalları listesinden kaldırıldı.")
	}
}
