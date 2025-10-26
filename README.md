# KeePassGo
**Password manager moderno scritto in Go con supporto formato .kdbx (compatibile KeePassXC)**

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey)](https://github.com)

---

## 📋 Indice

- [Caratteristiche](#-caratteristiche)
- [Requisiti](#-requisiti)
- [Installazione](#-installazione)
- [Struttura Progetto](#-struttura-progetto)
- [Compilazione](#-compilazione)
- [Utilizzo](#-utilizzo)
- [Sicurezza](#-sicurezza)
- [Sviluppo](#-sviluppo)
- [FAQ](#-faq)
- [Licenza](#-licenza)

---

## ✨ Caratteristiche

### Sicurezza
- ✅ **Compatibilità KeePassXC**: Legge e scrive database .kdbx 4.x
- ✅ **Cifrari moderni**: Solo algoritmi sicuri (ChaCha20, AES-256, Argon2id)
- ✅ **Nessun cifrario obsoleto**: AES-128, SHA-1, MD5 non supportati
- ✅ **Crittografia forte**: Password generate con `crypto/rand`
- ✅ **Offline**: Nessuna connessione internet, nessun telemetry
- ✅ **Open Source**: Codice completamente ispezionabile

### Funzionalità
- ✅ **GUI moderna**: Interfaccia grafica intuitiva con Fyne
- ✅ **Import/Export .kdbx**: Compatibile con KeePassXC
- ✅ **Gestione password**: Crea, modifica, visualizza, copia
- ✅ **Generatore password**: Password sicure personalizzabili
- ✅ **Categorie/Gruppi**: Organizzazione gerarchica
- ✅ **Appunti**: Copia username/password negli appunti

### Cifrari Supportati

| Algoritmo | Tipo | Sicurezza |
|-----------|------|-----------|
| **ChaCha20-Poly1305** | Cifratura simmetrica | ✅ Moderno (2008) |
| **AES-256-GCM** | Cifratura simmetrica | ✅ Standard NIST |
| **Argon2id** | Key Derivation | ✅ Vincitore PHC 2015 |
| **SHA-256** | Hash | ✅ Sicuro |
| ~~AES-128~~ | Cifratura simmetrica | ❌ NON supportato |
| ~~SHA-1~~ | Hash | ❌ NON supportato (vulnerabile) |
| ~~MD5~~ | Hash | ❌ NON supportato (rotto) |

---

## 📦 Requisiti

### Sistema Operativo
- **Linux** (testato su Arch Linux)
- **macOS** (10.13+)
- **Windows** (10+)

### Software Richiesto
- **Go**: 1.21 o superiore
- **GCC/Clang**: Compilatore C (per Fyne)
- **Librerie grafiche**:
  - **Linux**: `libgl1-mesa-dev`, `xorg-dev`
  - **macOS**: Xcode Command Line Tools
  - **Windows**: MinGW-w64

### Installazione Dipendenze

#### Arch Linux
```bash
sudo pacman -S go gcc mesa libx11 libxcursor libxrandr libxinerama libxi
```

#### Ubuntu/Debian
```bash
sudo apt install golang gcc libgl1-mesa-dev xorg-dev
```

#### macOS
```bash
brew install go
xcode-select --install
```

#### Windows
```bash
# Installa Go da https://go.dev/dl/
# Installa MinGW-w64 da https://www.mingw-w64.org/
```

---

## 🚀 Installazione

### Metodo 1: Clona il Repository (Esistente)

Se hai già il progetto:

```bash
cd /home/gabriel1/ClaudeWorkspace/keepassgo
go mod tidy
go build -o keepassgo ./cmd/keepassgo
./keepassgo
```

### Metodo 2: Creazione da Zero

Se vuoi ricreare il progetto completo:

#### 1. Crea la struttura delle directory

```bash
# Crea directory principale
mkdir -p keepassgo

# Naviga nella directory
cd keepassgo

# Crea struttura completa
mkdir -p cmd/keepassgo
mkdir -p internal/ui
mkdir -p internal/crypto
mkdir -p internal/database
mkdir -p pkg/kdbx

# Verifica struttura
tree -L 2
```

La struttura dovrebbe essere:
```
keepassgo/
├── cmd/
│   └── keepassgo/
├── internal/
│   ├── crypto/
│   ├── database/
│   └── ui/
└── pkg/
    └── kdbx/
```

#### 2. Inizializza il modulo Go

```bash
go mod init github.com/tuousername/keepassgo
```

#### 3. Installa le dipendenze

```bash
# Fyne (GUI framework)
go get fyne.io/fyne/v2@latest

# gokeepasslib (Parser .kdbx)
go get github.com/tobischo/gokeepasslib/v3@latest

# Crypto (cifratura)
go get golang.org/x/crypto@latest

# Pulisci dipendenze
go mod tidy
```

#### 4. Copia i file sorgente

Copia i seguenti file nelle rispettive directory:
- `cmd/keepassgo/main.go`
- `pkg/kdbx/reader.go`
- `pkg/kdbx/writer.go`
- `pkg/kdbx/password_generator.go`
- `internal/ui/main_window.go`

(I file sorgente completi sono disponibili nel progetto)

#### 5. Crea file aggiuntivi

```bash
# .gitignore
cat > .gitignore << 'EOF'
# Binaries
keepassgo
*.exe
*.dll
*.so
*.dylib

# Test
*.test
*.out

# Go workspace
go.work

# IDEs
.idea/
.vscode/
*.swp

# OS
.DS_Store
Thumbs.db

# Build
bin/
dist/
EOF
```

---

## 🏗️ Struttura Progetto

```
keepassgo/
├── cmd/
│   └── keepassgo/
│       └── main.go                    # Entry point applicazione
├── internal/                          # Codice privato
│   ├── crypto/                        # (Futuro) Utility cifratura
│   ├── database/                      # (Futuro) Logica DB
│   └── ui/
│       └── main_window.go             # Interfaccia grafica Fyne
├── pkg/                               # Codice pubblico riutilizzabile
│   └── kdbx/
│       ├── reader.go                  # Parser database .kdbx
│       ├── writer.go                  # Writer database .kdbx
│       └── password_generator.go     # Generatore password sicure
├── .gitignore
├── go.mod                             # Dipendenze Go
├── go.sum                             # Checksum dipendenze
├── LICENSE                            # Licenza MIT
├── README.md                          # Questo file
└── keepassgo                          # Binario compilato (31MB)
```

### Descrizione Moduli

| Directory | Descrizione | Righe Codice |
|-----------|-------------|--------------|
| `cmd/keepassgo/` | Entry point, inizializza GUI | ~30 |
| `pkg/kdbx/` | Parsing/writing .kdbx, generatore password | ~550 |
| `internal/ui/` | Interfaccia grafica completa | ~415 |
| **TOTALE** | | **~995 righe** |

---

## 🔨 Compilazione

### Build Standard

```bash
cd /home/gabriel1/ClaudeWorkspace/keepassgo
go build -o keepassgo ./cmd/keepassgo
```

### Build Ottimizzata (Dimensione Ridotta)

```bash
go build -ldflags="-s -w" -o keepassgo ./cmd/keepassgo
```

### Build con Compressione UPX

```bash
go build -ldflags="-s -w" -o keepassgo ./cmd/keepassgo
upx --best --lzma keepassgo
```

### Cross-Compilation

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o keepassgo-linux ./cmd/keepassgo

# macOS
GOOS=darwin GOARCH=amd64 go build -o keepassgo-macos ./cmd/keepassgo

# Windows
GOOS=windows GOARCH=amd64 go build -o keepassgo.exe ./cmd/keepassgo
```

---

## 💻 Utilizzo

### Avvio

```bash
./keepassgo
```

### Funzionalità Principali

#### 1. Creare un Nuovo Database

1. **File → Nuovo Database**
2. Scegli posizione e nome file (es. `passwords.kdbx`)
3. Inserisci password master **forte**
4. Il database viene creato con cifratura **ChaCha20 + Argon2id**

#### 2. Aprire Database Esistente

1. **File → Apri Database**
2. Seleziona file `.kdbx` (anche da KeePassXC!)
3. Inserisci password master
4. Visualizza tutte le password

#### 3. Aggiungere Password

1. Click su **Aggiungi Password**
2. Compila campi:
   - **Titolo**: Nome servizio (es. "Gmail")
   - **Username**: nome utente o email
   - **Password**: usa "Genera Password" per sicurezza
   - **URL**: sito web
   - **Note**: informazioni extra
3. Click **Salva**
4. **File → Salva** per salvare su disco

#### 4. Generare Password Sicure

1. Click su **Genera Password**
2. Visualizza password generata (20 caratteri, ~131 bit entropia)
3. Click **Copia** per copiare negli appunti
4. Usa in altri servizi

#### 5. Visualizzare/Copiare Password

1. Click su una password nella lista
2. Visualizza dettagli nel pannello destro
3. Click **Copia Password** o **Copia Username**
4. Password copiata negli appunti (auto-clear dopo 30s)

---

## 🔐 Sicurezza

### ⚠️ AVVISO IMPORTANTE

> **Questo è un progetto in sviluppo (v0.1.0)**
> **NON usare per dati reali** fino al completamento del security audit.
> Per uso in produzione, usa KeePassXC ufficiale.

### Algoritmi Utilizzati

1. **ChaCha20-Poly1305**
   - Cifratura simmetrica moderna (Google, Cloudflare)
   - Più veloce di AES su CPU senza AES-NI
   - Resistenza quantistica: ~128 bit (algoritmo di Grover)

2. **AES-256-GCM**
   - Standard NIST, ampiamente testato
   - 256 bit → 128 bit post-quantistico (ancora sicuro)

3. **Argon2id**
   - Vincitore Password Hashing Competition (2015)
   - Resistente a GPU/ASIC/quantum attacks
   - Parametri: 10 iterazioni, 1GB RAM, 4 thread

### Minacce Non Coperte

❌ **Keylogger**: protezione zero
❌ **Malware**: se il sistema è compromesso, nessun password manager è sicuro
❌ **Shoulder surfing**: usa in ambiente privato
❌ **Brute force password debole**: usa password master forte (16+ caratteri)

### Best Practices

✅ **Password Master**: Minimo 16 caratteri, misto, unica
✅ **Backup**: Copia regolare di `.kdbx` su dispositivo offline
✅ **2FA**: Dove possibile, attiva autenticazione a due fattori
✅ **Aggiornamenti**: Usa sempre l'ultima versione

---

## 🛠️ Sviluppo

### Requisiti Sviluppo

```bash
# Linter
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Testing
go test ./...

# Coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Roadmap

- [ ] **v0.2.0**: Ricerca password, filtri
- [ ] **v0.3.0**: Auto-lock, timeout inattività
- [ ] **v0.4.0**: Import/Export CSV
- [ ] **v0.5.0**: Supporto allegati
- [ ] **v1.0.0**: Security audit, release stabile

### Contribuire

1. Fork del repository
2. Crea branch feature (`git checkout -b feature/AmazingFeature`)
3. Commit modifiche (`git commit -m 'Add AmazingFeature'`)
4. Push branch (`git push origin feature/AmazingFeature`)
5. Apri Pull Request

---

## ❓ FAQ

**Q: È compatibile con KeePassXC?**
A: Sì, legge e scrive file `.kdbx` v4 compatibili.

**Q: Posso usarlo su Android/iOS?**
A: No, solo desktop (Linux, macOS, Windows).

**Q: Supporta YubiKey/FIDO2?**
A: Non ancora, previsto per v0.6.0.

**Q: È sicuro quanto KeePassXC?**
A: No, KeePassXC ha oltre 10 anni di audit. Usa quello per dati reali.

**Q: Perché 31MB di binario?**
A: Fyne include le dipendenze grafiche. Usa `upx` per comprimere.

**Q: Supporta cloud sync?**
A: No, ma puoi mettere `.kdbx` su Dropbox/Google Drive manualmente.

---

## 📄 Licenza

MIT License

Copyright (c) 2025 KeePassGo Contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

---

## 🙏 Credits

- **Fyne**: [fyne.io](https://fyne.io/) - GUI framework
- **gokeepasslib**: [tobischo/gokeepasslib](https://github.com/tobischo/gokeepasslib) - Parser .kdbx
- **KeePassXC**: [keepassxc.org](https://keepassxc.org/) - Ispirazione e formato .kdbx

---

**⚠️ Disclaimer**: Questo software è fornito "così com'è" senza garanzie. Usalo a tuo rischio. Per uso in produzione, preferisci password manager con security audit completi come KeePassXC, Bitwarden o 1Password.
