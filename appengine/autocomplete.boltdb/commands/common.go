package commands

import (
	"fmt"
	"time"

	"github.com/logrusorgru/aurora"

	"github.com/rickcrawford/ta/version"
)

const banner = `
TYPEAHEAD
 `

var titleMsg = fmt.Sprintf("Typeahead - Engine %s", version.Version)
var copyrightMsg = fmt.Sprintf("(C) Copyright %d Two Dot Solutions, LLC. All Rights Reserved.", time.Now().Year())
var urlMsg = fmt.Sprintf("Visit our website at %s", aurora.Cyan("https://typeahead.com"))
var longMsg = fmt.Sprintf("\n%s\n%s\n%s", titleMsg, copyrightMsg, urlMsg)
var bannerMsg = fmt.Sprintf("%s\n%s", aurora.Magenta(banner), longMsg)
