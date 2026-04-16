// Package handlers breachList
package handlers

// breachedCredentials maps breached usernames to their leaked plaintext passwords.
var breachedCredentials = map[string]string{ //nolint:gochecknoglobals
	"Benthe1954":    "^Jt^pLkzW2",
	"Jack1969":      "_yRw7uqk4h",
	"Pernille1949":  "^e3cTMrM4p",
	"Elna1996":      "t8YYsYRu$0",
	"Charlotte2020": ")1EFaslG$O",
	"Hans1965":      "YY0Bihn%g$",
	"Hanne2000":     "EpHYH0Vie(",
	"Doris1985":     "ziTGrweO^4",
	"Naja1991":      "Ne%1GsjYgx",
	"Albert1966":    "HCN2@Nw1@#",
	"Weena1968":     "@8XXVyc#*4",
	"Emil2002":      "$3mLd3ui3V",
	"Uffe1985":      "Qv33LISxU(",
	"Stephan1949":   "EozcbEk@&5",
	"Hannah1983":    "*$ow)Kub*9",
	"Jarl2024":      "^U1qZ(0ryO",
	"Boe1991":       "4v)8Wsr$FJ",
	"Sine2008":      "zNI#3G4sL%",
	"Olivia1968":    "^_Sno42&m5",
	"Abelone1972":   "YM*43dDbR)",
}

// isBreached reports whether the username/password pair was in the production breach.
func isBreached(username, password string) bool {
	leaked, ok := breachedCredentials[username]
	return ok && leaked == password
}
