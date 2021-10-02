package bot

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
)

const (
	SETUP_ID = "0"
)

func (h *BotHandler) HandleInlineQuery(q *tb.Query) {
	const (
		OPTION_NUMBER = 1
	)
	if q.Text == "" {
		return
	}
	options := make([]SetupOption, 0)
	options = append(options, SetupOption{
		Option:      "Setup Bot",
		Description: "Press here to setup homebot in your group",
	})
	results := make(tb.Results, OPTION_NUMBER)
	for i, _ := range options {
		result := &tb.ArticleResult{
			Title:       options[i].Option,
			Text:        options[i].Option,
			Description: options[i].Description,
			URL:         "",
		}

		results[i] = result
		// needed to set a unique string ID for each result
		results[i].SetResultID(strconv.Itoa(i))
	}
	err := h.Bot.Answer(q, &tb.QueryResponse{
		Results:   results,
		CacheTime: 60, // a minute
	})
	if err != nil {
		fmt.Printf("Error in inline query: %s\n", err)
	}
}
func (h *BotHandler) HandleSetup(chosen *tb.ChosenInlineResult) {
	if chosen.ResultID != SETUP_ID {
		return
	}
	_, err := h.Bot.Send(&chosen.From, "Setting up Homebot!")
	if err != nil {
		return
	}
}
