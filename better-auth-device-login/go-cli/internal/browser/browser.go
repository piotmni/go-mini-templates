package browser

import (
	"github.com/pkg/browser"
)

// Open opens the specified URL in the default browser
func Open(url string) error {
	return browser.OpenURL(url)
}
