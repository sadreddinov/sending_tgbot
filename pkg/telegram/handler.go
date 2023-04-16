package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	logrus.Printf("сообщение от %s c текстом: %s", message.From.UserName, message.Text)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Лучшая программа для выплат ФОРС2")
	b.bot.Send(msg)
}
func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "unknown command")

	switch message.Command() {
	case "start":
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("рассылка", "send"),
			),
		)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите действие:")
		msg.ReplyMarkup = keyboard
		b.bot.Send(msg)
		return nil
	default:
		_, err := b.bot.Send(msg)
		return err
	}
}

func (b *Bot) updateCallbackQueryHandler(update *tgbotapi.CallbackQuery, updates tgbotapi.UpdatesChannel) {
	if update.Data == "send" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "введите сообщение на рассылку")
		b.bot.Send(msg)
		var msgID int
		for update := range updates {
			if update.Message != nil && update.Message.From.UserName == b.admin_name {
				msgID = update.Message.MessageID
				break
			}
			if update.Message.From.UserName != b.admin_name && update.Message.Chat.IsPrivate() {
				go b.handleMessage(update.Message)
			}
		}
		for _, v := range b.groups {
			config := tgbotapi.ChatConfigWithUser{
				ChatID:             v,
				UserID:             b.admin_chat_id,
				SuperGroupUsername: "",
			}
			memberConfig := tgbotapi.GetChatMemberConfig{config}
			if member, err := b.bot.GetChatMember(memberConfig); err == nil {
				if member.Status != "left" && member.Status != "" && member.Status != "kicked" {
					go b.sending(v, update.Message.Chat.ID, msgID)
				}
			}

		}
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "рассылка успешна")
		b.bot.Send(msg)
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("рассылка", "send"),
		),
	)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите действие:")
	msg.ReplyMarkup = keyboard
	b.bot.Send(msg)
}

func (b *Bot) sending(chatId int64, botChatID int64, msgID int) {
	msg := tgbotapi.NewCopyMessage(chatId, botChatID, msgID)
	_, err := b.bot.Send(msg)
	if err != nil {
		logrus.Errorf("Error forwarding message: %s", err.Error())
	}
}
