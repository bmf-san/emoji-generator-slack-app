package main

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/bmf-san/go-dotenv"
	"github.com/chromedp/chromedp"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// Draw is data set for drawing.
type Draw struct {
	Color   string
	BgColor string
	Line1   string
	Line2   string // optional
}

func handlerGenerator(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("tpl/multiline.tpl"))
	if r.URL.Query().Get("line2") == "" {
		tpl = template.Must(template.ParseFiles("tpl/oneline.tpl"))
	}

	tpl.Execute(w, Draw{
		Color:   r.URL.Query().Get("color"),
		BgColor: r.URL.Query().Get("bgColor"),
		Line1:   r.URL.Query().Get("line1"),
		Line2:   r.URL.Query().Get("line2"),
	})
}

func handlerEvents(w http.ResponseWriter, r *http.Request) {
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch eventsAPIEvent.Type {
	case slackevents.URLVerification:
		var res *slackevents.ChallengeResponse
		if err := json.Unmarshal(body, &res); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		if _, err := w.Write([]byte(res.Challenge)); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case slackevents.CallbackEvent:
		innerEvent := eventsAPIEvent.InnerEvent
		switch event := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			message := strings.Split(event.Text, " ")
			if len(message) < 2 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			color := message[1]

			bgColor := message[2]
			line1 := message[3]

			// line2 is optional.
			line2 := ""
			if len(message) > 4 {
				line2 = message[4]
			}

			query := url.Values{
				"color":   []string{color},
				"bgColor": []string{bgColor},
				"line1":   []string{line1},
				"line2":   []string{line2},
			}

			ctx, cancel := chromedp.NewContext(context.Background())
			defer cancel()

			// FIXME: Images are posted multiple times by running chromedp.
			var buf []byte
			if err := chromedp.Run(ctx, chromedp.Tasks{
				chromedp.Navigate(`http://localhost:9999/generator?` + query.Encode()),
				chromedp.Sleep(2 * time.Second),
				chromedp.WaitVisible(`#target`, chromedp.ByID),
				chromedp.Screenshot(`#target`, &buf, chromedp.NodeVisible, chromedp.ByID),
			}); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Failed to take a screen shot."))
				return
			}

			if err := ioutil.WriteFile("result.png", buf, 0644); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Failed to take a screen shot."))
				return
			}

			r := bytes.NewReader(buf)
			_, err = api.UploadFile(
				slack.FileUploadParameters{
					Reader:   r,
					Filename: "upload file name",
					Channels: []string{event.Channel},
				})
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Failed to post a image."))
				return
			}
		}
	}
}

func middlewareVerification(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		verifier, err := slack.NewSecretsVerifier(r.Header, os.Getenv("SLACK_SIGNING_SECRET"))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		bodyReader := io.TeeReader(r.Body, &verifier)
		body, err := ioutil.ReadAll(bodyReader)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := verifier.Ensure(); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		next.ServeHTTP(w, r)
	}
}

func main() {
	if err := dotenv.LoadEnv(); err != nil {
		log.Println(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/generator", handlerGenerator)
	mux.HandleFunc("/slack/events", middlewareVerification(handlerEvents))
	srv := &http.Server{
		Addr:    ":9999",
		Handler: mux,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
