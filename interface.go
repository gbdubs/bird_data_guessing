package bird_data_guessing

import "github.com/gbdubs/attributions"

type Input struct {
	LatinName   string
	EnglishName string
}

type Output struct {
	Data        Data
	DebugData   DebugData
	Attribution attributions.Attribution
}

type Data struct {
	EnglishName string
	LatinName   string
	// Food
	WheatScore  int
	WormScore   int
	BerryScore  int
	RatScore    int
	FishScore   int
	NectarScore int
	// Habitat
	ForestScore int
	GrassScore  int
	WaterScore  int
	// Nest Type
	CupScore      int
	GroundScore   int
	PlatformScore int
	SlotScore     int
	// Misc Properties
	Wingspan      int
	ClutchSize    int
	FlockingScore int
	PredatorScore int
	FunFact       string
}

// Debugging Types below this line

type DebugData struct {
	Food        Food
	NestType    NestType
	ClutchSize  Property
	Wingspan    Property
	Habitat     Habitat
	IsFlocking  Property
	IsPredatory Property
	FunFact     Property
}

type Habitat struct {
	Forest Property
	Water  Property
	Grass  Property
}

type Food struct {
	Worm   Property
	Wheat  Property
	Berry  Property
	Fish   Property
	Rat    Property
	Nectar Property
}

type NestType struct {
	Ground   Property
	Cup      Property
	Slot     Property
	Platform Property
}

type Property struct {
	Strength    int
	Context     string
	StringValue string
	IntValue    int
}

func (p *Property) Add(o *Property) {
	p.IntValue += o.IntValue
	p.StringValue += o.StringValue
	p.Strength += o.Strength
	p.Context += o.Context
}
