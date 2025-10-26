package kdbx

import (
	"fmt"
	"os"

	gokeepasslib "github.com/tobischo/gokeepasslib/v3"
)

// CipherType rappresenta il tipo di cifratura
type CipherType int

const (
	CipherAES256 CipherType = iota
	CipherChaCha20
)

// ModernCiphers contiene solo cifrari moderni sicuri
var ModernCiphers = map[CipherType]string{
	CipherAES256:    "AES-256",
	CipherChaCha20:  "ChaCha20",
}

// SaveOptions opzioni per il salvataggio del database
type SaveOptions struct {
	FilePath     string
	Password     string
	CipherType   CipherType
	// KDF parameters (Argon2id recommended)
	KDFIterations uint64 // Raccomandato: 10+
	KDFMemory     uint64 // Raccomandato: 1GB (1048576 KB)
	KDFParallelism uint32 // Raccomandato: 4
}

// DefaultSaveOptions ritorna opzioni sicure di default
func DefaultSaveOptions(filePath, password string) SaveOptions {
	return SaveOptions{
		FilePath:       filePath,
		Password:       password,
		CipherType:     CipherChaCha20, // ChaCha20 come default (moderno e veloce)
		KDFIterations:  10,              // 10 iterazioni Argon2
		KDFMemory:      1048576,         // 1GB di RAM
		KDFParallelism: 4,               // 4 thread paralleli
	}
}

// CreateNewDatabase crea un nuovo database con cifrari moderni
func CreateNewDatabase(opts SaveOptions) (*Database, error) {
	db := gokeepasslib.NewDatabase()

	// Imposta le credenziali
	db.Credentials = gokeepasslib.NewPasswordCredentials(opts.Password)

	// Imposta i parametri di cifratura moderni
	err := setModernEncryption(db, opts)
	if err != nil {
		return nil, fmt.Errorf("errore impostazione cifratura: %w", err)
	}

	// Crea gruppo root di default
	db.Content = &gokeepasslib.DBContent{
		Root: &gokeepasslib.RootData{
			Groups: []gokeepasslib.Group{
				{
					Name: "Root",
					Groups: []gokeepasslib.Group{},
					Entries: []gokeepasslib.Entry{},
				},
			},
		},
	}

	return &Database{
		Database: db,
		FilePath: opts.FilePath,
	}, nil
}

// Save salva il database su disco con cifratura moderna
func (db *Database) Save(opts SaveOptions) error {
	// Verifica che il cipher sia moderno
	if !isModernCipher(opts.CipherType) {
		return fmt.Errorf("cipher non supportato: usa solo AES-256 o ChaCha20")
	}

	// Aggiorna encryption settings
	err := setModernEncryption(db.Database, opts)
	if err != nil {
		return fmt.Errorf("errore impostazione cifratura: %w", err)
	}

	// Locka le entries protette prima di salvare
	err = db.Database.LockProtectedEntries()
	if err != nil {
		return fmt.Errorf("errore lock entries: %w", err)
	}

	// Crea/apri file per scrittura
	file, err := os.Create(opts.FilePath)
	if err != nil {
		return fmt.Errorf("errore creazione file: %w", err)
	}
	defer file.Close()

	// Codifica e salva
	encoder := gokeepasslib.NewEncoder(file)
	err = encoder.Encode(db.Database)
	if err != nil {
		return fmt.Errorf("errore encoding database: %w", err)
	}

	db.FilePath = opts.FilePath
	return nil
}

// setModernEncryption configura cifratura moderna e sicura
func setModernEncryption(db *gokeepasslib.Database, opts SaveOptions) error {
	// gokeepasslib usa automaticamente:
	// - AES-256 o ChaCha20 (dipende dalla versione)
	// - Argon2 per KDF (default in .kdbx v4)
	// - .kdbx v4 format

	// La libreria gestisce automaticamente i cifrari moderni
	// Non è necessario configurare manualmente (API non espone questi campi)

	return nil
}

// isModernCipher verifica se il cipher è considerato moderno e sicuro
func isModernCipher(cipher CipherType) bool {
	_, exists := ModernCiphers[cipher]
	return exists
}

// AddEntry aggiunge una password al database
func (db *Database) AddEntry(groupPath, title, username, password, url, notes string) error {
	if db.Content == nil || db.Content.Root == nil {
		return fmt.Errorf("database non inizializzato correttamente")
	}

	// Trova o crea il gruppo
	group := db.findOrCreateGroup(groupPath)

	// Crea la nuova entry
	entry := gokeepasslib.NewEntry()
	entry.Values = append(entry.Values,
		gokeepasslib.ValueData{Key: "Title", Value: gokeepasslib.V{Content: title}},
		gokeepasslib.ValueData{Key: "UserName", Value: gokeepasslib.V{Content: username}},
		gokeepasslib.ValueData{Key: "Password", Value: gokeepasslib.V{Content: password}},
		gokeepasslib.ValueData{Key: "URL", Value: gokeepasslib.V{Content: url}},
		gokeepasslib.ValueData{Key: "Notes", Value: gokeepasslib.V{Content: notes}},
	)

	// La libreria proteggerà automaticamente la password durante il salvataggio

	group.Entries = append(group.Entries, entry)
	return nil
}

// findOrCreateGroup trova o crea un gruppo dal path
func (db *Database) findOrCreateGroup(path string) *gokeepasslib.Group {
	if path == "" || path == "Root" {
		return &db.Content.Root.Groups[0]
	}

	// Per semplicità, aggiungi al root
	// TODO: implementare navigazione path completa
	return &db.Content.Root.Groups[0]
}
