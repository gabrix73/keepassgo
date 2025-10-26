package main

import (
	"fmt"

	"fyne.io/fyne/v2/app"

	"github.com/gabriel1/keepassgo/internal/ui"
)

const (
	AppName    = "KeePassGo"
	AppVersion = "0.1.0"
)

func main() {
	fmt.Printf("%s v%s - Modern Password Manager\n", AppName, AppVersion)
	fmt.Println("Compatible with KeePassXC .kdbx format")
	fmt.Println("Using modern cryptography only (ChaCha20, AES-256, Argon2id)")
	fmt.Println()

	// Crea applicazione
	myApp := app.NewWithID("com.gabriel1.keepassgo")

	// Crea e mostra finestra principale
	mainWindow := ui.NewMainWindow(myApp)
	mainWindow.Show()
}
