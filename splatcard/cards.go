package main

import (
	"math/rand"
	"time"
)

type Card struct {
	Key                string
	Value              string
	ConsecutiveCorrect int
}

type CardSet struct {
	Cards []*Card
}

func NewCardSet() *CardSet {
	// Initialize with all cards
	cards := []Card{
		{Key: "Alabama", Value: "Montgomery"},
		{Key: "Alaska", Value: "Juneau"},
		{Key: "Arizona", Value: "Phoenix"},
		{Key: "Arkansas", Value: "Little Rock"},
		{Key: "California", Value: "Sacramento"},
		{Key: "Colorado", Value: "Denver"},
		{Key: "Connecticut", Value: "Hartford"},
		{Key: "Delaware", Value: "Dover"},
		{Key: "Florida", Value: "Tallahassee"},
		{Key: "Georgia", Value: "Atlanta"},
		{Key: "Hawaii", Value: "Honolulu"},
		{Key: "Idaho", Value: "Boise"},
		{Key: "Illinois", Value: "Springfield"},
		{Key: "Indiana", Value: "Indianapolis"},
		{Key: "Iowa", Value: "Des Moines"},
		{Key: "Kansas", Value: "Topeka"},
		{Key: "Kentucky", Value: "Frankfort"},
		{Key: "Louisiana", Value: "Baton Rouge"},
		{Key: "Maine", Value: "Augusta"},
		{Key: "Maryland", Value: "Annapolis"},
		{Key: "Massachusetts", Value: "Boston"},
		{Key: "Michigan", Value: "Lansing"},
		{Key: "Minnesota", Value: "Saint Paul"},
		{Key: "Mississippi", Value: "Jackson"},
		{Key: "Missouri", Value: "Jefferson City"},
		{Key: "Montana", Value: "Helena"},
		{Key: "Nebraska", Value: "Lincoln"},
		{Key: "Nevada", Value: "Carson City"},
		{Key: "New Hampshire", Value: "Concord"},
		{Key: "New Jersey", Value: "Trenton"},
		{Key: "New Mexico", Value: "Santa Fe"},
		{Key: "New York", Value: "Albany"},
		{Key: "North Carolina", Value: "Raleigh"},
		{Key: "North Dakota", Value: "Bismarck"},
		{Key: "Ohio", Value: "Columbus"},
		{Key: "Oklahoma", Value: "Oklahoma City"},
		{Key: "Oregon", Value: "Salem"},
		{Key: "Pennsylvania", Value: "Harrisburg"},
		{Key: "Rhode Island", Value: "Providence"},
		{Key: "South Carolina", Value: "Columbia"},
		{Key: "South Dakota", Value: "Pierre"},
		{Key: "Tennessee", Value: "Nashville"},
		{Key: "Texas", Value: "Austin"},
		{Key: "Utah", Value: "Salt Lake City"},
		{Key: "Vermont", Value: "Montpelier"},
		{Key: "Virginia", Value: "Richmond"},
		{Key: "Washington", Value: "Olympia"},
		{Key: "West Virginia", Value: "Charleston"},
		{Key: "Wisconsin", Value: "Madison"},
		{Key: "Wyoming", Value: "Cheyenne"},
	}

	// Shuffle the cards
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	random.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})

	// Convert to a slice of pointers
	cardPointers := make([]*Card, len(cards))
	for i := range cards {
		cardPointers[i] = &cards[i]
	}

	return &CardSet{Cards: cardPointers}
}

// GetCard returns the next card from the queue.
func (cs *CardSet) GetCard() *Card {
	if len(cs.Cards) == 0 {
		return nil
	}
	card := cs.Cards[0]
	cs.Cards = cs.Cards[1:]
	return card
}

// ReinsertCard reinserts a card back into the queue at a given position.
func (cs *CardSet) ReinsertCard(card *Card, position int) {
	if position > len(cs.Cards) {
		position = len(cs.Cards)
	}
	if position < 0 {
		position = 0
	}

	cs.Cards = append(cs.Cards[:position], append([]*Card{card}, cs.Cards[position:]...)...)
}
