package telegram

import (
	"log"
	"os"
	"strconv"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	bot           *tgbotapi.BotAPI
	groups        map[string]int64
	admin_name    string
	admin_chat_id int64
}

func NewBot(bot *tgbotapi.BotAPI) *Bot {
	var groups = make(map[string]int64, 0)
	admin_name := os.Getenv("ADMIN_NAME")
	admin_chat_id, err := strconv.ParseInt(os.Getenv("ADMIN_CHAT_ID"), int(10), int(64))
	if err != nil {
		logrus.Errorf("error while loading env variable ADMIN_CHAT_ID:", err.Error())
	}
	return &Bot{bot: bot, groups: groups, admin_name: admin_name, admin_chat_id: admin_chat_id}
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates := b.initUpdatesChannel()

	b.handleUpdates(updates)

	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.CallbackQuery != nil && update.CallbackQuery.From.UserName == b.admin_name {
			b.updateCallbackQueryHandler(update.CallbackQuery, updates)
			continue
		}

		if update.Message == nil { // If we got a message
			continue
		}

		chat := update.Message.Chat
		if chat.IsSuperGroup() || chat.IsGroup() {
			config := tgbotapi.ChatConfigWithUser{
				ChatID:             chat.ID,
				UserID:             b.admin_chat_id,
				SuperGroupUsername: "",
			}
			memberConfig := tgbotapi.GetChatMemberConfig{config}
			if member, err := b.bot.GetChatMember(memberConfig); err == nil {
				if member.Status != "left" && member.Status != "" && member.Status != "kicked" {
					if _, ok := b.groups[chat.UserName]; !ok {
						b.groups[chat.UserName] = chat.ID
						logrus.Printf("добавлен id группы: %d", chat.ID)
					} else {
						logrus.Printf("группа c id %d уже добавлена", chat.ID)
					}
				} else {
					logrus.Printf("группа c id %d не подходит", chat.ID)
				}
			}
		}

		if update.Message.IsCommand() && update.Message.From.UserName == b.admin_name {
			b.handleCommand(update.Message)
			continue
		}

		if update.Message.From.UserName != b.admin_name && update.Message.Chat.IsPrivate() {
			go b.handleMessage(update.Message)
		}
	}
}

func (b *Bot) initUpdatesChannel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}
