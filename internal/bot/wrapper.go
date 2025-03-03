package bot

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"log"
	"log/slog"
	"time"
)

type CommandType string

const (
	CommandStart = "start"
	CommandPlay  = "play"
)

type Wrapper struct {
	logger      *slog.Logger
	pollUpdater *ext.Updater
}

func NewBotWrapper(token string, logger *slog.Logger) *Wrapper {

	b, err := gotgbot.NewBot(token, nil)
	if err != nil {
		panic("failed to create new client instance: " + err.Error())
	}

	w := &Wrapper{}

	err = w.creteUpdater(b)
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}

	return w
}

func (w *Wrapper) Start() {
	w.pollUpdater.Idle()
}

func (w *Wrapper) Stop() {
	_ = w.pollUpdater.Stop()
}

func (w *Wrapper) creteUpdater(b *gotgbot.Bot) error {
	dispatcher := w.createDispatcher()
	updater := ext.NewUpdater(dispatcher, nil)

	err := updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})

	if err != nil {
		w.logger.Error("failed to start polling: " + err.Error())
	}

	w.pollUpdater = updater

	w.logger.Info(fmt.Sprintf("%s has been started...\n", b.User.Username))

	return nil
}

func (w *Wrapper) createDispatcher() *ext.Dispatcher {
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	dispatcher.AddHandler(handlers.NewCommand(CommandStart, start))
	dispatcher.AddHandler(handlers.NewCommand(CommandPlay, play))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("play"), play))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("game=takemypixel"), play))

	return dispatcher
}

func play(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := b.SendGame(ctx.EffectiveChat.Id, "takemypixel", nil)
	if err != nil {
		return fmt.Errorf("failed to send game: %w", err)
	}

	return nil
}

func start(b *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("Start", ctx.EffectiveMessage.GetSender().Username())
	_, err := ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Hello! I'm @%s. Lets play!", b.User.Username), &gotgbot.SendMessageOpts{
		ParseMode: "html",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
				{
					Text:         "play",
					CallbackData: "play",
				},
			}},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	return nil
}
