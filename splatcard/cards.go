package main

import (
	"math/rand"
)

type Card struct {
	Key   string
	Value string
}

type CardSet struct {
	Cards []Card
}

func NewCardSet() *CardSet {
	return &CardSet{
		Cards: []Card{
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
		},
	}
}

func (cs *CardSet) GetCard() *Card {
	idx := rand.Intn(len(cs.Cards))
	return &cs.Cards[idx]
}
