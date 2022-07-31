package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kamchy/stoic"
	"github.com/kamchy/stoic/model"
	"github.com/kamchy/stoic/stoicdb"

	"fyne.io/fyne/v2/data/binding"
	"log"
)

type State struct {
	Repo      stoic.Repository
	App       fyne.App
	Window    fyne.Window
	Error     binding.String
	Content   binding.String
	Author    binding.String
	Connected binding.Bool
}

var initialString = "Welcome to the Stoic app!"
var state State = State{
	Error:     binding.BindString(&initialString),
	Connected: binding.NewBool(),
	Author:    binding.NewString(),
	Content:   binding.NewString(),
}

func save(db stoic.Repository, qq *model.Quote) (int64, error) {
	log.Printf("gui.main save: qq=%v and db is %v", qq, db)
	idx, err := db.SaveQuote(*qq)
	return idx, err
}

func openFileCb(urc fyne.URIReadCloser, err error) {

	log.Printf("openFileCB: uri is %v", urc.URI())
	dbpath := urc.URI().Path()
	if err != nil {
		state.Error.Set(err.Error())
		return
	}
	db, err := stoicdb.New(dbpath)
	if err != nil {
		state.Error.Set(err.Error())
		return
	}
	state.Error.Set("")
	state.Connected.Set(true)
	state.Repo = db
}

func createMainMenu() *fyne.MainMenu {
	mm := fyne.NewMainMenu()
	openDatabaseItem := fyne.NewMenuItem("Open", func() { dialog.ShowFileOpen(openFileCb, state.Window) })
	quitItem := fyne.NewMenuItem("Quit", fyne.CurrentApp().Quit)
	mm.Items = make([]*fyne.Menu, 0)
	mm.Items = append(mm.Items, fyne.NewMenu("Database", openDatabaseItem, quitItem))
	return mm
}

func createSaveButton(quoteInput *widget.Entry, authorInput *widget.Entry) fyne.Widget {
	bu := widget.NewButton("Save", func() {
		qq := model.Quote{}
		qq.Text = quoteInput.Text
		qq.Author = authorInput.Text
		quoteInput.Text = ""
		authorInput.Text = ""

		log.Printf("Content was: %v", qq)
		save(state.Repo, &qq)
		quoteInput.Refresh()
		authorInput.Refresh()
	})

	state.Connected.AddListener(binding.NewDataListener(func() {
		conn, err := state.Connected.Get()
		if err != nil {
			state.Error.Set(err.Error())
			return
		}
		if conn {
			bu.Enable()
		} else {
			bu.Disable()
		}

	}))
	return bu
}
func main() {
	a := app.New()
	w := a.NewWindow("Your stoic meditation")
	state.App = a
	state.Window = w

	w.SetMainMenu(createMainMenu())
	message := widget.NewLabelWithData(state.Error)

	quoteInput := widget.NewMultiLineEntry()
	quoteInput.SetPlaceHolder("Enter text...")
	authorInput := widget.NewEntry()
	authorInput.SetPlaceHolder("Enter author..")

	saveButton := createSaveButton(quoteInput, authorInput)
	content := container.NewVBox(quoteInput, authorInput, saveButton)

	quitButton := widget.NewButton("Quit", func() {
		state.App.Quit()
	})

	state.Window.SetContent(container.NewVBox(
		message,
		widget.NewButton("Meditate", func() {
			message.SetText("AMOR FATI")
		}),
		content,
		quitButton,
	),
	)
	state.Window.Resize(fyne.NewSize(400, 200))
	state.Window.ShowAndRun()
}
