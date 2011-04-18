/* assumptions: we are the only players in the game
 * later this can be refined if you have an idea of what other players have or don't have.
 * 
 * 1. get inputs
 * a. your hand (required)
 * b. the board (0 cards, 3 , 4, or 5 cards)
 *         Other numbers of cards represent errors
 * 
 * 2. for all possible poker hands that can be formed by your hand:
 * calculate the odds of getting that hand (100% if you already have it)
 * 
 * 3. print all odds in a nice format
 */

package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	PARSE_STATE_EAT_VAL = iota
	PARSE_STATE_EAT_VAL_SAW_1
	PARSE_STATE_EAT_SUIT
)

const (
	CLUBS = iota
	DIAMONDS
	HEARTS
	SPADES
)

func usage() {
	fmt.Fprintf(os.Stderr,
`%s: the Texas Hold Em' poker odds calculator.

This program calculates your 'outs' for a Texas Hold Em' poker hand.
Texas Hold Em' is a popular version of poker where each player receives
exactly two secret cards. Then there are five rounds of betting.

The format used to specify cards is as follows:
[type][suit]
There are 4 suits:
C = clubs, D = diamonds, H = hearts, S = spades
There are 13 different card types:
1 = A = ace, K = king, Q = queen, J = jack, 2 = 2, ... 10 = 10

Usage:
-a [your hand as a whitespace-separated list of cards]
-b [the board as a whitespace-separated list of cards]
If no -b is given, it will be assumed that no cards are on the board.
-h this help message

Usage Example:
%s -a KS\ QS
Find the outs you have pre-flop with a king and queen of spades.
`, os.Args[0], os.Args[0])
}

type card struct {
	val int
	suit int
}

func valToStr(v int) (string) {
	switch {
	case v == 1:
		return "A"
	case v == 11:
		return "J"
	case v == 12:
		return "Q"
	case v == 13:
		return "K"
	}
	return fmt.Sprintf("%d", v)
}

func suitToStr(s int) (string) {
	switch {
	case s == CLUBS:
		return "♣C"
	case s == DIAMONDS:
		return "♦D"
	case s == HEARTS:
		return "♥H"
	case s == SPADES:
		return "♠S"
	}
	return ""
}

func (p *card) String() string {
	return fmt.Sprintf("%s%s", valToStr(p.val), suitToStr(p.suit))
}

func strToCard(str string, cnt *int) (myCard *card) {
	myCard = new(card)
	var parseState = PARSE_STATE_EAT_VAL
	for ;*cnt < len(str); {
		var c = str[*cnt]
		*cnt++
		switch {
		case parseState == PARSE_STATE_EAT_VAL:
			switch {
			case c == ' ' || c == '\t':
				continue
			case c == '1':
				parseState = PARSE_STATE_EAT_VAL_SAW_1
			case c >= '2' && c <= '9':
				myCard.val = (int)(c - '0')
				parseState = PARSE_STATE_EAT_SUIT
			case c == 'J':
				myCard.val = 11
				parseState = PARSE_STATE_EAT_SUIT
			case c == 'Q':
				myCard.val = 12
				parseState = PARSE_STATE_EAT_SUIT
			case c == 'K':
				myCard.val = 13
				parseState = PARSE_STATE_EAT_SUIT
			case c == 'A':
				myCard.val = 1
				parseState = PARSE_STATE_EAT_SUIT
			default:
				*cnt = -1
				return nil
			}
		case parseState == PARSE_STATE_EAT_VAL_SAW_1:
			switch {
			case c == '0':
				myCard.val = 10
				parseState = PARSE_STATE_EAT_SUIT
			default:
				*cnt = -1
				return nil
			}
		case parseState == PARSE_STATE_EAT_SUIT:
			switch {
			case c == 'C':
				myCard.suit = CLUBS
			case c == 'D':
				myCard.suit = DIAMONDS
			case c == 'H':
				myCard.suit = HEARTS
			case c == 'S':
				myCard.suit = SPADES
			default:
				*cnt = -1
				return nil
			}
			return myCard
		}
	}
	*cnt = -1
    return nil
}

func (p *card) Compare(rhs *card) int {
	if (p.suit < rhs.suit) {
		return -1;
	}
	if (p.suit > rhs.suit) {
		return 1;
	}
	if (p.val < rhs.val) {
		return -1;
	}
	if (p.val > rhs.val) {
		return 1;
	}
	return 0;
}

func hasDuplicates(c []*card) *card {
	for i := range(c) {
		for j := range(c) {
			if i == j {
				continue
			}
			if (c[i].Compare(c[j]) == 0) {
				return c[i]
			}
		}
	}
	return nil
}

func strToCards(str string) (cards []*card, cnt int) {
	for cnt = 0; cnt != -1; {
		var c = strToCard(str, &cnt)
		if (c != nil) {
			cards = append(cards,c)
		}
	}
	return
}

func cardsToStr(c []*card) (string) {
	ret := ""
	sep := ""
	for i := range(c) {
		ret += fmt.Sprintf("%s%s", sep, c[i].String())
		sep = ", "
	}
	return ret
}

func intsToStr(s []int) (string) {
	ret := ""
	sep := ""
	for i := range(s) {
		ret += fmt.Sprintf("%s%d", sep, s[i])
		sep = ", "
	}
	return ret
}

func checkBoardLength(l int) {
	var validLens = []int { 0, 3, 4, 5 }
	for i := range(validLens) {
		if (int(l) == validLens[i]) {
			return
		}
	}

	fmt.Printf("illegal board length. Expected a length of %s, " +
	"but your board length was %d.\n", intsToStr(validLens), l)
	os.Exit(1)
}

func generateAllVals(i *int, cards *[52]*card, suit int) {
	for val := 1; val <= 13; val++ {
		cards[*i] = new(card)
		cards[*i].suit = suit
		cards[*i].val = val
		*i++
	}
}

func generateAllCards(cards *[52]*card) {
	i := 0
	generateAllVals(&i, cards, CLUBS)
	generateAllVals(&i, cards, DIAMONDS)
	generateAllVals(&i, cards, HEARTS)
	generateAllVals(&i, cards, SPADES)
}

// find next k-combination
// assume x has form x'01^a10^b in binary
func nextCombination(x *int64) bool {
	u := *x & -*x // extract rightmost bit 1; u =  0'00^a10^b
	v := u + *x // set last non-trailing bit 0, and clear to the right; v=x'10^a00^b
	if (v==0) { // then overflow in v, or x==0
		return false; // signal that next k-combination cannot be represented
	}
	*x = v + (((v^*x)/u)>>2); // v^x = 0'11^a10^b, (v^x)/u = 0'0^b1^{a+2}, and x ← x'100^b1^a
	if (*x >= (1<<52)) {
		return false; // too big
	}
	return true; // successful completion
}

func combinationToCards(comb int64, allCards *[52]*card, holeC *[]*card,
						boardC *[]*card) ([]*card) {
	var ret []*card = make([]*card, 5)
	copy(ret[:], *holeC)
	copy(ret[len(*holeC):], *boardC)
	var n = len(*holeC) + len(*boardC)
	for i := range(allCards) {
		if (((1 << uint(i)) & comb) != 0) {
			ret[n] = &(*allCards[i])
			n++
			if (n >= 5) {
				return ret
			}
		}
	}
	fmt.Printf("combinationToHand: logic error: got to unreachable point\n")
	os.Exit(1)
	return ret
}

const {
	HIGH_CARD = iota
	PAIR
	TWO_PAIR
	THREE_OF_A_KIND
	STRAIGHT
	FLUSH
	FULL_HOUSE
	FOUR_OF_A_KIND
	STRAIGHT_FLUSH
}

type hand struct {
	cards []*card
	val [2]int
	flushSuit int
	ty int
}

// hmm. would be better to sort the cards appropriately beforehand. It would
// make straight detection easier.
func makeHand(cards []*card) *hand {
	ret := new(hand)
	ret.cards = cards
	var vals = make(map[int] int)
	var suits = make(map[int] int)
	var order []int
	for i := range(cards) {
		c := cards[i]
		vals[c.val] = vals[c.val] + 1
		suits[c.suit] = vals[c.suit] + 1
	}

	// check for flush
	for i := range(suits) {
		if (suits[i] >= 4) {
			ret.flushSuit = i
		}
	}
	// check for straight flush
	var orderedVals := make([]int, len(vals))
	j := 0
	for k,_ := range (vals) {
		orderedVals[j] = k
		j++
	}
	sort.SortInts(&orderedVals)
	prev = -1
	runLen = 0
	for i := range(orderedVals) {
		if (prev + 1 == orderedVals[i]) {
			runLen++;
		}
		else {
			runLen = 0
		}
	}
	if ((runLen >= 5) && (ret.flushSuit != 0)) {
		ret.val[0] = orderedVals[len(orderedVals) - 5]
		ret.ty = STRAIGHT_FLUSH
		return ret
	}

	var freqs = make(map[int] []int)
	for k,v := range(vals) {
		if (v > 4) {
			fmt.Printf("got %d of a kind for value %d (max is 4)\n", v, k)
			os.Exit(0)
		}
		curFreqs := freqs[v]
		m := 0
		for m = 0; m < len(curFreqs); m++ {
			if (curFreqs[m] >= k)
				break
		}
		newFreqs := curFreqs[:m]
		append(newFreqs, k)
		append(newFreqs, curFreqs[m:])
		freqs[v] = newFreqs
	}

	// four of a kind
	if (len(freqs[4]) > 0) {
		ret.ty = FOUR_OF_A_KIND
		ret.val[0] = freqs[4][0]
		return ret
	}

	// full house
	if (len(freqs[3]) > 0) {
		if (len(freqs[3]) > 1) {
			ret.val[0] = freqs[3][0]
			ret.val[1] = freqs[3][1]
			ret.ty = FULL_HOUSE
		}
		else if (len(freqs[2]) > 0) {
			ret.val[0] = freqs[3][0]
			ret.val[1] = freqs[2][0]
			ret.ty = FULL_HOUSE
		}
	}

	// flush
	if (ret.flushSuit != 0) {
		ret.ty = FLUSH
		return ret
	}

	// straight
	if (runLen >= 5) {
		// where does the straight start?
		ret.val[0] = orderedVals[len(orderedVals) - 5]
		ret.ty = STRAIGHT
		return ret
	}

	// three of a kind
	if (len(freqs[3]) > 0) {
		ret.val[0] = freqs[3][0]
		ret.ty = THREE_OF_A_KIND
		return ret
	}

	// two pairs
	if (len(freqs[2]) >= 2) {
		ret.val[0] = freqs[2][0]
		ret.val[1] = freqs[2][1]
		ret.ty = TWO_PAIR
		return ret
	}

	// a pair
	if (len(freqs[2]) >= 1) {
		ret.val[0] = freqs[2][0]
		ret.ty = PAIR
		return ret
	}

	// I guess not.
	ret.ty = HIGH_CARD
	return ret
}

func main() {
	flag.Usage = usage
	var verbose = flag.Bool("v", false, "verbose")
	var help = flag.Bool("h", false, "help")
	var hole = flag.String("a", "", "your two hole cards")
	var board = flag.String("b", "", "the board")
	flag.Parse()

	if (*help) {
		usage()
		os.Exit(0)
	}
	if (*hole == "") {
		fmt.Printf("You must give two hole cards with -a\n")
		usage()
		os.Exit(1)
	}
	holeC, errIdx := strToCards(*hole)
	if (errIdx != -1) {
		fmt.Printf("Error parsing your hole cards: parse error at character %d\n",
				   errIdx)
		os.Exit(1)
	}

	if (*verbose) {
		fmt.Printf("Your hole cards: %s\n", cardsToStr(holeC));
	}

	boardC, bErrIdx := strToCards(*board)
	if (bErrIdx != -1) {
		fmt.Printf("parse error at character %d\n", bErrIdx)
	}
	checkBoardLength(len(boardC))
	if (*verbose) {
		fmt.Printf("The board: %s\n", cardsToStr(boardC));
	}

	var c = make([]*card, len(boardC) + len(holeC))
	copy(c, boardC)
	copy(c[len(boardC):], holeC)
	dupe := hasDuplicates(c)
	if (dupe != nil) {
		fmt.Printf("The card %s appears more than once!\n", dupe)
		os.Exit(1)
	}

	// generate all cards
	var allCards [52]*card
	generateAllCards(&allCards)

	var comb int64 = 31
	switch (len(boardC)) {
	case 0:
		comb = 31
	case 3:
		comb = 3
	case 4:
		comb = 1
	case 5:
		comb = 0
	default:
		fmt.Printf("invalid board length %d\n", len(boardC))
		os.Exit(1)
	}
	for ;nextCombination(&comb); {
		handC := combinationToCards(comb, &allCards, &holeC, &boardC)
		fmt.Printf("%d: %s\n", comb, cardsToStr(handC))
		var hand = makeHand(handC)
		if (hand != nil) {
			fmt.Printf("%s", hand.String())
		}
	}

}
