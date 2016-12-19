package main

import (
	"math/rand"
	"time"

	"github.com/Armienn/GoLanguage/phonetics"
	"github.com/gopherjs/gopherjs/js"
)

func main() {
	js.Global.Get("gobutton").Set("onclick", generateSomePhonetics)
}

/*func printAlphabet(lang *phonetics.Phonetics, representation *phonetics.Phonetics) {
	for _, sound := range lang.Sounds {
		fmt.Print(representation.GetRepresentation(sound))
	}
	fmt.Println()
}*/

func generateSomePhonetics() {
	rand.Seed(time.Now().UTC().UnixNano())
	//rand.Seed(1)
	lang := phonetics.RandomPhonetics()
	dansk := phonetics.GetDansk()
	ipa := phonetics.GetIpa()
	//printAlphabet(dansk, dansk)
	//printAlphabet(dansk, ipa)
	//printAlphabet(ipa, dansk)
	//printAlphabet(ipa, ipa)

	//for i := 0; i < 10; i++ {
	//lang = language.RandomLanguage()
	lang.Patterns = phonetics.GetMubPatterns()
	//printAlphabet(lang, dansk)
	//printAlphabet(lang, ipa)
	//}

	text := ""
	for i := 0; i < 20; i++ {
		word := lang.RandomWord(0)
		text += dansk.GetWordRepresentation(word) + "</br>"
		text += ipa.GetWordRepresentation(word) + "</br>"
	}

	js.Global.Get("mulle").Set("innerHTML", text)
}
