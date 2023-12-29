package forms

type LeaguePlayerForm struct {
	// True = now, False = Later
	GameMode     string `validate:"required,oneof=NormalQP NormalDraft Ranked ARAM TFT Featured Any"`
	SeekingRole  string `validate:"oneof=Top Mid Jungle Bottom Support Any"`
	OfferingRole string `validate:"oneof=Top Mid Jungle Bottom Support Any"`
	// SeekingRank  string
	// OfferingRank string
}

// game mode:
// NormalQP, NormalDraft, Ranked, ARAM, TFT, Featured

// seeking role (only applicable in NormalQP, NormalDraft, RankedDraft):
// Top, Mid, Jungle, Bottom, Support

//
