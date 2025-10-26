package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/gabriel1/keepassgo/pkg/kdbx"
)

// MainWindow rappresenta la finestra principale
type MainWindow struct {
	App      fyne.App
	Window   fyne.Window
	Database *kdbx.Database

	// UI Components
	entryList    *widget.List
	detailsPanel *fyne.Container
	entries      []kdbx.Entry
}

// NewMainWindow crea una nuova finestra principale
func NewMainWindow(app fyne.App) *MainWindow {
	win := app.NewWindow("KeePassGo - Modern Password Manager")

	mw := &MainWindow{
		App:     app,
		Window:  win,
		entries: []kdbx.Entry{},
	}

	mw.setupUI()
	win.Resize(fyne.NewSize(1000, 600))
	win.CenterOnScreen()

	return mw
}

// setupUI configura l'interfaccia utente
func (mw *MainWindow) setupUI() {
	// Menu bar
	menu := mw.createMenu()
	mw.Window.SetMainMenu(menu)

	// Lista password (sinistra)
	mw.entryList = mw.createEntryList()

	// Pannello dettagli (destra)
	mw.detailsPanel = mw.createDetailsPanel()

	// Layout principale: split view
	splitView := container.NewHSplit(
		container.NewBorder(
			widget.NewLabel("Password"),
			mw.createToolbar(),
			nil,
			nil,
			mw.entryList,
		),
		mw.detailsPanel,
	)
	splitView.SetOffset(0.3)

	mw.Window.SetContent(splitView)
}

// createMenu crea il menu dell'applicazione
func (mw *MainWindow) createMenu() *fyne.MainMenu {
	// File menu
	openItem := fyne.NewMenuItem("Apri Database", mw.openDatabase)
	newItem := fyne.NewMenuItem("Nuovo Database", mw.newDatabase)
	saveItem := fyne.NewMenuItem("Salva", mw.saveDatabase)
	quitItem := fyne.NewMenuItem("Esci", func() {
		mw.App.Quit()
	})

	fileMenu := fyne.NewMenu("File", openItem, newItem, saveItem, fyne.NewMenuItemSeparator(), quitItem)

	// Help menu
	aboutItem := fyne.NewMenuItem("Info", mw.showAbout)
	helpMenu := fyne.NewMenu("Aiuto", aboutItem)

	return fyne.NewMainMenu(fileMenu, helpMenu)
}

// createToolbar crea la toolbar con i pulsanti principali
func (mw *MainWindow) createToolbar() *fyne.Container {
	addBtn := widget.NewButton("Aggiungi Password", mw.addEntry)
	delBtn := widget.NewButton("Elimina", mw.deleteEntry)
	genBtn := widget.NewButton("Genera Password", mw.generatePassword)

	return container.NewHBox(addBtn, delBtn, genBtn)
}

// createEntryList crea la lista delle password
func (mw *MainWindow) createEntryList() *widget.List {
	list := widget.NewList(
		func() int {
			return len(mw.entries)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			entry := mw.entries[id]
			label.SetText(fmt.Sprintf("%s (%s)", entry.Title, entry.Username))
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		mw.showEntryDetails(id)
	}

	return list
}

// createDetailsPanel crea il pannello dei dettagli
func (mw *MainWindow) createDetailsPanel() *fyne.Container {
	welcomeLabel := widget.NewLabel("Apri o crea un database per iniziare")
	welcomeLabel.Wrapping = fyne.TextWrapWord

	info := widget.NewRichTextFromMarkdown(`
## KeePassGo

Password manager moderno con supporto .kdbx

### Caratteristiche di sicurezza:
- ✅ Solo cifrari moderni (AES-256, ChaCha20)
- ✅ Argon2id per key derivation
- ✅ Nessun supporto per cifrari obsoleti
- ✅ Compatibile con KeePassXC

### Come iniziare:
1. **File → Nuovo Database** per creare un nuovo database
2. **File → Apri Database** per aprire un .kdbx esistente
3. Usa **Aggiungi Password** per creare nuove entry
`)

	return container.NewVBox(
		welcomeLabel,
		widget.NewSeparator(),
		info,
	)
}

// openDatabase apre un database esistente
func (mw *MainWindow) openDatabase() {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil || reader == nil {
			return
		}
		defer reader.Close()

		filePath := reader.URI().Path()

		// Chiedi password
		mw.promptPassword("Inserisci password del database", func(password string) {
			db, err := kdbx.OpenDatabase(filePath, password)
			if err != nil {
				dialog.ShowError(fmt.Errorf("Errore apertura database: %w", err), mw.Window)
				return
			}

			mw.Database = db
			mw.entries = db.GetAllEntries()
			mw.entryList.Refresh()

			dialog.ShowInformation("Successo",
				fmt.Sprintf("Database aperto: %d password trovate", len(mw.entries)),
				mw.Window)
		})
	}, mw.Window)
}

// newDatabase crea un nuovo database
func (mw *MainWindow) newDatabase() {
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil || writer == nil {
			return
		}
		defer writer.Close()

		filePath := writer.URI().Path()

		// Chiedi password per il nuovo database
		mw.promptPassword("Crea password master per il database", func(password string) {
			opts := kdbx.DefaultSaveOptions(filePath, password)
			db, err := kdbx.CreateNewDatabase(opts)
			if err != nil {
				dialog.ShowError(fmt.Errorf("Errore creazione database: %w", err), mw.Window)
				return
			}

			err = db.Save(opts)
			if err != nil {
				dialog.ShowError(fmt.Errorf("Errore salvataggio database: %w", err), mw.Window)
				return
			}

			mw.Database = db
			mw.entries = []kdbx.Entry{}
			mw.entryList.Refresh()

			dialog.ShowInformation("Successo",
				fmt.Sprintf("Nuovo database creato: %s\nCifratura: ChaCha20 + Argon2id", filePath),
				mw.Window)
		})
	}, mw.Window)
}

// saveDatabase salva il database corrente
func (mw *MainWindow) saveDatabase() {
	if mw.Database == nil {
		dialog.ShowError(fmt.Errorf("Nessun database aperto"), mw.Window)
		return
	}

	// Chiedi password per conferma
	mw.promptPassword("Conferma password per salvare", func(password string) {
		opts := kdbx.DefaultSaveOptions(mw.Database.FilePath, password)

		err := mw.Database.Save(opts)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Errore salvataggio: %w", err), mw.Window)
			return
		}

		dialog.ShowInformation("Successo", "Database salvato", mw.Window)
	})
}

// promptPassword mostra un dialog per inserire la password
func (mw *MainWindow) promptPassword(title string, callback func(string)) {
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.PlaceHolder = "Password"

	dialog.ShowForm(title, "OK", "Annulla",
		[]*widget.FormItem{
			widget.NewFormItem("Password", passwordEntry),
		},
		func(ok bool) {
			if ok && passwordEntry.Text != "" {
				callback(passwordEntry.Text)
			}
		},
		mw.Window,
	)
}

// addEntry aggiunge una nuova password
func (mw *MainWindow) addEntry() {
	if mw.Database == nil {
		dialog.ShowError(fmt.Errorf("Apri o crea un database prima"), mw.Window)
		return
	}

	titleEntry := widget.NewEntry()
	usernameEntry := widget.NewEntry()
	passwordEntry := widget.NewPasswordEntry()
	urlEntry := widget.NewEntry()
	notesEntry := widget.NewMultiLineEntry()

	dialog.ShowForm("Nuova Password", "Salva", "Annulla",
		[]*widget.FormItem{
			widget.NewFormItem("Titolo", titleEntry),
			widget.NewFormItem("Username", usernameEntry),
			widget.NewFormItem("Password", passwordEntry),
			widget.NewFormItem("URL", urlEntry),
			widget.NewFormItem("Note", notesEntry),
		},
		func(ok bool) {
			if ok {
				err := mw.Database.AddEntry("",
					titleEntry.Text,
					usernameEntry.Text,
					passwordEntry.Text,
					urlEntry.Text,
					notesEntry.Text,
				)
				if err != nil {
					dialog.ShowError(err, mw.Window)
					return
				}

				// Ricarica entries
				mw.entries = mw.Database.GetAllEntries()
				mw.entryList.Refresh()
				dialog.ShowInformation("Successo", "Password aggiunta", mw.Window)
			}
		},
		mw.Window,
	)
}

// deleteEntry elimina la password selezionata
func (mw *MainWindow) deleteEntry() {
	dialog.ShowInformation("TODO", "Funzionalità in sviluppo", mw.Window)
}

// generatePassword genera una password casuale
func (mw *MainWindow) generatePassword() {
	opts := kdbx.DefaultPasswordOptions()
	password, err := kdbx.GeneratePassword(opts)
	if err != nil {
		dialog.ShowError(err, mw.Window)
		return
	}

	passwordLabel := widget.NewEntry()
	passwordLabel.SetText(password)
	passwordLabel.Disable()

	entropy := kdbx.CalculateEntropy(password, opts)

	dialog.ShowCustom("Password Generata",
		"OK",
		container.NewVBox(
			widget.NewLabel(fmt.Sprintf("Lunghezza: %d caratteri", opts.Length)),
			widget.NewLabel(fmt.Sprintf("Entropia: ~%.0f bit", entropy)),
			widget.NewSeparator(),
			passwordLabel,
			widget.NewButton("Copia", func() {
				mw.Window.Clipboard().SetContent(password)
				dialog.ShowInformation("Copiato", "Password copiata negli appunti", mw.Window)
			}),
		),
		mw.Window,
	)
}

// showEntryDetails mostra i dettagli di una password
func (mw *MainWindow) showEntryDetails(id int) {
	if id < 0 || id >= len(mw.entries) {
		return
	}

	entry := mw.entries[id]

	titleLabel := widget.NewLabelWithStyle(entry.Title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	usernameEntry := widget.NewEntry()
	usernameEntry.SetText(entry.Username)
	usernameEntry.Disable()

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetText(entry.Password)
	passwordEntry.Disable()

	urlEntry := widget.NewEntry()
	urlEntry.SetText(entry.URL)
	urlEntry.Disable()

	notesEntry := widget.NewMultiLineEntry()
	notesEntry.SetText(entry.Notes)
	notesEntry.Disable()

	copyPasswordBtn := widget.NewButton("Copia Password", func() {
		mw.Window.Clipboard().SetContent(entry.Password)
		dialog.ShowInformation("Copiato", "Password copiata negli appunti", mw.Window)
	})

	copyUsernameBtn := widget.NewButton("Copia Username", func() {
		mw.Window.Clipboard().SetContent(entry.Username)
	})

	mw.detailsPanel.Objects = []fyne.CanvasObject{
		titleLabel,
		widget.NewForm(
			widget.NewFormItem("Username", usernameEntry),
			widget.NewFormItem("Password", passwordEntry),
			widget.NewFormItem("URL", urlEntry),
			widget.NewFormItem("Note", notesEntry),
		),
		container.NewHBox(copyPasswordBtn, copyUsernameBtn),
	}

	mw.detailsPanel.Refresh()
}

// showAbout mostra la finestra "Informazioni"
func (mw *MainWindow) showAbout() {
	dialog.ShowCustom("Informazioni su KeePassGo", "OK",
		widget.NewRichTextFromMarkdown(`
## KeePassGo v0.1.0

Password manager moderno scritto in Go

### Sicurezza
- Cifrari: AES-256, ChaCha20
- KDF: Argon2id
- Compatibile con KeePassXC

### Licenza
MIT License

### Sviluppato con
- Go 1.25+
- Fyne v2
- gokeepasslib
`),
		mw.Window,
	)
}

// Show mostra la finestra
func (mw *MainWindow) Show() {
	mw.Window.ShowAndRun()
}
