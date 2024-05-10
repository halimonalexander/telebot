package telegram

const (
	MsgPageAlreadyExists = "Page already exists"

	MsgSaved = "Saved!"

	MsgStoreException = "Unable to save link, please try later"

	// ---------

	MsgNxPage = "You do not have any stored page"

	MsgFetchException = "Unable to get list, please try later"

	MsgRandomPage = "I have found something for you: %s"

	// ---------

	MsgUnknownCommand = "Command is unknown. Use /help to get list of available commands"

	MsgHello = "Hello! \n\n " + MsgHelp

	MsgHelp = `I can save and keep your pages. Also I can offer you them to read.

In order to save the page, just send me a link to it.

In order to get a random page from your list, send me command /rnd.
Caution! After that that page will be removed from your list!
`
)
