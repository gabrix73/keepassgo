package kdbx

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	// Character sets per generazione password
	LowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	UppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Digits           = "0123456789"
	SpecialChars     = "!@#$%^&*()_+-=[]{}|;:,.<>?"
)

// PasswordOptions opzioni per generare password
type PasswordOptions struct {
	Length         int
	UseLowercase   bool
	UseUppercase   bool
	UseDigits      bool
	UseSpecial     bool
	ExcludeAmbiguous bool // Escludi caratteri ambigui (0, O, l, 1, ecc.)
}

// DefaultPasswordOptions ritorna opzioni sicure di default
func DefaultPasswordOptions() PasswordOptions {
	return PasswordOptions{
		Length:         20,
		UseLowercase:   true,
		UseUppercase:   true,
		UseDigits:      true,
		UseSpecial:     true,
		ExcludeAmbiguous: false,
	}
}

// GeneratePassword genera una password casuale sicura usando crypto/rand
func GeneratePassword(opts PasswordOptions) (string, error) {
	if opts.Length < 8 {
		return "", fmt.Errorf("lunghezza password deve essere almeno 8 caratteri")
	}

	// Costruisci il set di caratteri
	charset := buildCharset(opts)
	if len(charset) == 0 {
		return "", fmt.Errorf("nessun set di caratteri selezionato")
	}

	// Genera password
	password := make([]byte, opts.Length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := range password {
		// Usa crypto/rand per sicurezza crittografica
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", fmt.Errorf("errore generazione numero casuale: %w", err)
		}
		password[i] = charset[randomIndex.Int64()]
	}

	// Verifica che contenga almeno un carattere di ogni tipo richiesto
	if !validatePassword(string(password), opts) {
		// Ricorsione per rigenerare (molto raro)
		return GeneratePassword(opts)
	}

	return string(password), nil
}

// buildCharset costruisce il set di caratteri basato sulle opzioni
func buildCharset(opts PasswordOptions) string {
	var charset string

	if opts.UseLowercase {
		if opts.ExcludeAmbiguous {
			charset += removeAmbiguous(LowercaseLetters)
		} else {
			charset += LowercaseLetters
		}
	}

	if opts.UseUppercase {
		if opts.ExcludeAmbiguous {
			charset += removeAmbiguous(UppercaseLetters)
		} else {
			charset += UppercaseLetters
		}
	}

	if opts.UseDigits {
		if opts.ExcludeAmbiguous {
			charset += removeAmbiguous(Digits)
		} else {
			charset += Digits
		}
	}

	if opts.UseSpecial {
		charset += SpecialChars
	}

	return charset
}

// removeAmbiguous rimuove caratteri ambigui
func removeAmbiguous(s string) string {
	ambiguous := "il1Lo0O"
	result := ""
	for _, c := range s {
		isAmbiguous := false
		for _, amb := range ambiguous {
			if c == amb {
				isAmbiguous = true
				break
			}
		}
		if !isAmbiguous {
			result += string(c)
		}
	}
	return result
}

// validatePassword verifica che la password contenga almeno un carattere di ogni tipo richiesto
func validatePassword(password string, opts PasswordOptions) bool {
	hasLower, hasUpper, hasDigit, hasSpecial := false, false, false, false

	for _, c := range password {
		switch {
		case isInCharset(c, LowercaseLetters):
			hasLower = true
		case isInCharset(c, UppercaseLetters):
			hasUpper = true
		case isInCharset(c, Digits):
			hasDigit = true
		case isInCharset(c, SpecialChars):
			hasSpecial = true
		}
	}

	if opts.UseLowercase && !hasLower {
		return false
	}
	if opts.UseUppercase && !hasUpper {
		return false
	}
	if opts.UseDigits && !hasDigit {
		return false
	}
	if opts.UseSpecial && !hasSpecial {
		return false
	}

	return true
}

// isInCharset verifica se un carattere è presente in un set
func isInCharset(c rune, charset string) bool {
	for _, ch := range charset {
		if c == ch {
			return true
		}
	}
	return false
}

// CalculateEntropy calcola l'entropia della password in bit
func CalculateEntropy(password string, opts PasswordOptions) float64 {
	charset := buildCharset(opts)
	charsetSize := float64(len(charset))
	passwordLength := float64(len(password))

	// Entropia = log2(charsetSize^passwordLength)
	// = passwordLength * log2(charsetSize)
	if charsetSize == 0 {
		return 0
	}

	return passwordLength * log2(charsetSize)
}

// log2 calcola logaritmo base 2
func log2(x float64) float64 {
	if x <= 0 {
		return 0
	}
	// log2(x) = ln(x) / ln(2)
	return 1.44269504089 * 2.71828182846 // Approssimazione
}
