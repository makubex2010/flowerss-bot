package bot

import (
	"time"

	"github.com/indes/flowerss-bot/bot/fsm"
	"github.com/indes/flowerss-bot/config"
	"github.com/indes/flowerss-bot/log"
	"github.com/indes/flowerss-bot/util"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	// UserState 用户状態，用於標示當前用户操作所在状態
	UserState map[int64]fsm.UserStatus = make(map[int64]fsm.UserStatus)

	// B telebot
	B *tb.Bot
)

func init() {
	if config.RunMode == config.TestMode {
		return
	}
	poller := &tb.LongPoller{Timeout: 10 * time.Second}
	spamProtected := tb.NewMiddlewarePoller(poller, func(upd *tb.Update) bool {
		if !isUserAllowed(upd) {
			// 檢查用户是否可以使用bot
			return false
		}

		if !CheckAdmin(upd) {
			return false
		}
		return true
	})
	log.Infow("init telegram bot",
		"token", config.BotToken,
		"endpoint", config.TelegramEndpoint,
	)

	// create bot
	var err error

	B, err = tb.NewBot(tb.Settings{
		URL:    config.TelegramEndpoint,
		Token:  config.BotToken,
		Poller: spamProtected,
		Client: util.HttpClient,
	})

	if err != nil {
		log.Fatal(err)
		return
	}
}

//Start bot
func Start() {
	if config.RunMode != config.TestMode {
		setCommands()
		setHandle()
		B.Start()
	}
}

func setCommands() {
	// 設置bot命令提示信息
	commands := []tb.Command{
		{"start", "開始使用"},
		{"sub", "訂閱rss源"},
		{"list", "當前訂閱的rss源"},
		{"unsub", "退定rss源"},
		{"unsuball", "退定所有rss源"},

		{"set", "設置rss訂閱"},
		{"setfeedtag", "設置rss訂閱標籤"},
		{"setinterval", "設置rss訂閱抓取間隔"},

		{"export", "導出訂閱為opml文件"},
		{"import", "從opml文件導入訂閱"},

		{"check", "检查我的rss訂閱状態"},
		{"pauseall", "停止抓取訂閱更新"},
		{"activeall", "開始抓取訂閱更新"},

		{"help", "使用幫助"},
		{"version", "bot版本"},
	}

	if err := B.SetCommands(commands); err != nil {
		log.Errorw("set bot commands failed", "error", err.Error())
	}
}

func setHandle() {
	B.Handle(&tb.InlineButton{Unique: "set_feed_item_btn"}, setFeedItemBtnCtr)

	B.Handle(&tb.InlineButton{Unique: "set_toggle_notice_btn"}, setToggleNoticeBtnCtr)

	B.Handle(&tb.InlineButton{Unique: "set_toggle_telegraph_btn"}, setToggleTelegraphBtnCtr)

	B.Handle(&tb.InlineButton{Unique: "set_toggle_update_btn"}, setToggleUpdateBtnCtr)

	B.Handle(&tb.InlineButton{Unique: "set_set_sub_tag_btn"}, setSubTagBtnCtr)

	B.Handle(&tb.InlineButton{Unique: "unsub_all_confirm_btn"}, unsubAllConfirmBtnCtr)

	B.Handle(&tb.InlineButton{Unique: "unsub_all_cancel_btn"}, unsubAllCancelBtnCtr)

	B.Handle(&tb.InlineButton{Unique: "unsub_feed_item_btn"}, unsubFeedItemBtnCtr)

	B.Handle("/start", startCmdCtr)

	B.Handle("/export", exportCmdCtr)

	B.Handle("/sub", subCmdCtr)

	B.Handle("/list", listCmdCtr)

	B.Handle("/set", setCmdCtr)

	B.Handle("/unsub", unsubCmdCtr)

	B.Handle("/unsuball", unsubAllCmdCtr)

	B.Handle("/ping", pingCmdCtr)

	B.Handle("/help", helpCmdCtr)

	B.Handle("/import", importCmdCtr)

	B.Handle("/setfeedtag", setFeedTagCmdCtr)

	B.Handle("/setinterval", setIntervalCmdCtr)

	B.Handle("/check", checkCmdCtr)

	B.Handle("/activeall", activeAllCmdCtr)

	B.Handle("/pauseall", pauseAllCmdCtr)

	B.Handle("/version", versionCmdCtr)

	B.Handle(tb.OnText, textCtr)

	B.Handle(tb.OnDocument, docCtr)
}
