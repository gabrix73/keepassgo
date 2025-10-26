package kdbx

import (
	"fmt"
	"os"

	gokeepasslib "github.com/tobischo/gokeepasslib/v3"
)

// Database rappresenta un database KeePass
type Database struct {
	*gokeepasslib.Database
	FilePath string
}

// Entry rappresenta una singola password/entry
type Entry struct {
	Title    string
	Username string
	Password string
	URL      string
	Notes    string
	GroupPath string // Path completo del gruppo (categoria)
}

// Group rappresenta un gruppo/categoria
type Group struct {
	Name     string
	Path     string
	Entries  []Entry
	SubGroups []Group
}

// OpenDatabase apre un database .kdbx esistente
func OpenDatabase(filePath string, password string) (*Database, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("errore apertura file: %w", err)
	}
	defer file.Close()

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(password)

	err = gokeepasslib.NewDecoder(file).Decode(db)
	if err != nil {
		return nil, fmt.Errorf("errore decodifica database: %w", err)
	}

	// Sblocca il database
	err = db.UnlockProtectedEntries()
	if err != nil {
		return nil, fmt.Errorf("errore sblocco entries: %w", err)
	}

	return &Database{
		Database: db,
		FilePath: filePath,
	}, nil
}

// GetAllEntries ottiene tutte le password dal database
func (db *Database) GetAllEntries() []Entry {
	var entries []Entry

	if db.Content != nil && db.Content.Root != nil {
		entries = db.extractEntriesFromGroup(&db.Content.Root.Groups[0], "")
	}

	return entries
}

// extractEntriesFromGroup estrae ricorsivamente entries da un gruppo
func (db *Database) extractEntriesFromGroup(group *gokeepasslib.Group, parentPath string) []Entry {
	var entries []Entry

	currentPath := parentPath
	if group.Name != "" {
		if currentPath != "" {
			currentPath += " / " + group.Name
		} else {
			currentPath = group.Name
		}
	}

	// Estrai entries del gruppo corrente
	for _, entry := range group.Entries {
		e := Entry{
			GroupPath: currentPath,
		}

		// Estrai i valori dai campi
		for _, value := range entry.Values {
			switch value.Key {
			case "Title":
				e.Title = value.Value.Content
			case "UserName":
				e.Username = value.Value.Content
			case "Password":
				e.Password = value.Value.Content
			case "URL":
				e.URL = value.Value.Content
			case "Notes":
				e.Notes = value.Value.Content
			}
		}

		entries = append(entries, e)
	}

	// Elabora ricorsivamente i sottogruppi
	for i := range group.Groups {
		subEntries := db.extractEntriesFromGroup(&group.Groups[i], currentPath)
		entries = append(entries, subEntries...)
	}

	return entries
}

// GetAllGroups ottiene tutti i gruppi/categorie
func (db *Database) GetAllGroups() []Group {
	var groups []Group

	if db.Content != nil && db.Content.Root != nil && len(db.Content.Root.Groups) > 0 {
		groups = db.extractGroups(&db.Content.Root.Groups[0], "")
	}

	return groups
}

// extractGroups estrae ricorsivamente la struttura dei gruppi
func (db *Database) extractGroups(group *gokeepasslib.Group, parentPath string) []Group {
	var groups []Group

	currentPath := parentPath
	if group.Name != "" {
		if currentPath != "" {
			currentPath += " / " + group.Name
		} else {
			currentPath = group.Name
		}
	}

	g := Group{
		Name: group.Name,
		Path: currentPath,
	}

	// Estrai entries di questo gruppo
	for _, entry := range group.Entries {
		e := Entry{
			GroupPath: currentPath,
		}

		for _, value := range entry.Values {
			switch value.Key {
			case "Title":
				e.Title = value.Value.Content
			case "UserName":
				e.Username = value.Value.Content
			case "Password":
				e.Password = value.Value.Content
			case "URL":
				e.URL = value.Value.Content
			case "Notes":
				e.Notes = value.Value.Content
			}
		}

		g.Entries = append(g.Entries, e)
	}

	// Elabora sottogruppi
	for i := range group.Groups {
		subGroups := db.extractGroups(&group.Groups[i], currentPath)
		g.SubGroups = append(g.SubGroups, subGroups...)
	}

	groups = append(groups, g)
	return groups
}

// Close chiude il database (lockando le entries protette)
func (db *Database) Close() error {
	if db.Database != nil {
		return db.Database.LockProtectedEntries()
	}
	return nil
}
