package ext

import "github.com/JREAMLU/j-kit/useragent"

// UA ua
type UA struct {
	IsBot          bool
	Localization   string
	IsMobile       bool
	Mozilla        string
	Platform       string
	OS             string
	Engine         string
	EngineVersion  string
	Browser        string
	BrowserVersion string
}

// ParseUserAgent parse user agent
func ParseUserAgent(ual string) UA {
	ua := useragent.New(ual)

	engine, engineVersion := ua.Engine()
	browser, browserVersion := ua.Browser()

	return UA{
		IsBot:          ua.Bot(),
		Localization:   ua.Localization(),
		IsMobile:       ua.Mobile(),
		Mozilla:        ua.Mozilla(),
		Platform:       ua.Platform(),
		OS:             ua.OS(),
		Engine:         engine,
		EngineVersion:  engineVersion,
		Browser:        browser,
		BrowserVersion: browserVersion,
	}
}
