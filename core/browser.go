package core

import (
	"time"

	"github.com/corpix/uarand"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
	"github.com/sirupsen/logrus"
)

type BrowserOpts struct {
	IsHeadless    bool          // Use browser interface
	IsLeakless    bool          // Force to kill browser
	Timeout       time.Duration // Timeout
	LanguageCode  string
	WaitRequests  bool          // Wait requests to complete after navigation
	LeavePageOpen bool          // Leave pages and browser open
	WaitLoadTime  time.Duration // Time to wait till page loads
}

// Initialize browser parameters with default values if they are not set
func (o *BrowserOpts) Check() {
	if o.Timeout == 0 {
		o.Timeout = time.Second * 30
	}

	if o.WaitLoadTime == 0 {
		o.WaitLoadTime = time.Second * 5
	}
}

type Browser struct {
	BrowserOpts
	browserAddr string
	browser     *rod.Browser
}

func NewBrowser(opts BrowserOpts) (*Browser, error) {
	opts.Check()
	logrus.Debugf("Browser options: %+v", opts)

	path, has := launcher.LookPath()
	logrus.Debug("Browser found: ", has)

	var err error
	b := Browser{BrowserOpts: opts}
	b.browserAddr, err = launcher.New().
		Bin(path).
		Leakless(opts.IsLeakless).
		Headless(opts.IsHeadless).
		NoSandbox(true).
		Launch() // NoSandbox

	return &b, err
}

func (b *Browser) Initialize() {
	b.browser = rod.New().ControlURL(b.browserAddr)
	b.browser.MustConnect()
	b.browser.SetCookies(nil)
	logrus.Debug("Browser initialized - Address: ", b.browserAddr)
}

// Check whether browser instance is already created
func (b *Browser) IsInitialized() bool {
	if b.browserAddr != "" {
		return true
	} else {
		return false
	}
}

// Open URL
func (b *Browser) Navigate(URL string) *rod.Page {
	logrus.Debug("Navigate to: ", URL)

	// b.browser = rod.New().ControlURL(b.browserAddr)
	// b.browser.MustConnect()
	// b.browser.SetCookies(nil)

	page := stealth.MustPage(b.browser)
	wait := page.MustWaitRequestIdle()
	page.MustNavigate(URL)

	// causes bugs in google
	if b.WaitRequests {
		wait()
	}

	page.MustEmulate(devices.Device{
		UserAgent:      uarand.GetRandom(),
		AcceptLanguage: b.LanguageCode,
	})

	// Wait till page loads
	time.Sleep(b.WaitLoadTime)

	return page
}

func (b *Browser) Close() error {
	return b.browser.Close()
}
