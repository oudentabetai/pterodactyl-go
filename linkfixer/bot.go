package linkfixer

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var urlPattern = regexp.MustCompile(`https?://[^\s<>]+`)

var domainMapping = map[string]string{
	"twitter.com":       "fxtwitter.com",
	"www.twitter.com":   "fxtwitter.com",
	"x.com":             "fxtwitter.com",
	"www.x.com":         "fxtwitter.com",
	"instagram.com":     "ddinstagram.com",
	"www.instagram.com": "ddinstagram.com",
	"tiktok.com":        "vxtiktok.com",
	"www.tiktok.com":    "vxtiktok.com",
}

func Register(s *discordgo.Session) {
	s.AddHandler(onMessageCreate)
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m == nil || m.Author == nil || m.Author.Bot {
		return
	}

	matches := urlPattern.FindAllString(m.Content, -1)
	if len(matches) == 0 {
		return
	}

	fixedLinks := make([]string, 0, len(matches))
	for _, rawLink := range matches {
		trimmedLink := strings.TrimRight(rawLink, ".,!?)]}")
		fixedLink, ok := rewriteLink(trimmedLink)
		if !ok {
			continue
		}
		fixedLinks = append(fixedLinks, fixedLink)
	}

	if len(fixedLinks) == 0 {
		return
	}

	content := "🔗 LinkFixer:\n" + strings.Join(fixedLinks, "\n")
	_, _ = s.ChannelMessageSendReply(m.ChannelID, content, m.Reference())
}

func rewriteLink(rawLink string) (string, bool) {
	parsed, err := url.Parse(rawLink)
	if err != nil || parsed.Host == "" {
		return "", false
	}

	replacedHost, ok := domainMapping[strings.ToLower(parsed.Host)]
	if !ok {
		return "", false
	}
	parsed.Host = replacedHost
	return parsed.String(), true
}
