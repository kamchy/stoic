package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
  "github.com/kamchy/stoic/model"
  "github.com/kamchy/stoic/stoicdb"
  "github.com/kamchy/stoic"
  
  "log"
  "os"
)

func save(db stoic.Repository, qq *model.Quote) (int64, error) {
  idx, err := db.SaveQuote(*qq)
  return idx, err
}


func main() {
	a := app.New()
	w := a.NewWindow("Your stoic meditation")

	dbpath := os.Args[1]
  db, err := stoicdb.New(dbpath)
	if err != nil {
		log.Fatal(err)
	}

	hello := widget.NewLabel("Hello, practicing stoic!")
  quoteInput := widget.NewMultiLineEntry()
	quoteInput.SetPlaceHolder("Enter text...")
  authorInput := widget.NewEntry()
	authorInput.SetPlaceHolder("Enter author..")
  var qq = &model.Quote{}
	content := container.NewVBox(quoteInput, authorInput, 
    widget.NewButton("Save", func() {
      qq.Text = quoteInput.Text
      qq.Author = authorInput.Text
      quoteInput.Text = ""
      authorInput.Text = ""

		  log.Printf("Content was: %v", qq )
      save(db, qq)
	  }),
  )

  quitButton := widget.NewButton("Quit", func() {
    a.Quit()
  })

	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Meditate", func() {
			hello.SetText("AMOR FATI")
		}),
    content,    
    quitButton,
	  ),
  )

	w.ShowAndRun()
}
