package claude

const (
	Claude3Dot7SonnetLatest = "claude-3-7-sonnet-latest"
	Claude3Dot5SonnetLatest = "claude-3-5-sonnet-latest"
	Claude3Dot5HaikuLatest  = "claude-3-5-haiku-latest"
	Claude3OpusLatest       = "claude-3-opus-latest"

	Claude3Dot7Sonnet2502 = "claude-3-7-sonnet-20250219"
	Claude3Dot5Sonnet2410 = "claude-3-5-sonnet-20241022"
	Claude3Dot5Sonnet     = "claude-3-5-sonnet-20240620"
	Claude3Dot5Haiku      = "claude-3-5-haiku-20241022"
	Claude3Opus           = "claude-3-opus-20240229"
	Claude3Sonnet         = "claude-3-sonnet-20240229"
	Claude3Haiku          = "claude-3-haiku-20240307"
	Claude2Dot1           = "claude-2.1"
	Clause2Dot0           = "claude-2.0"
	Claude1Dot2Instant    = "claude-instant-1.2"
)

func Models() []string {
	return []string{
		Claude3Dot7Sonnet2502,
		Claude3Dot5Sonnet2410,
		Claude3Opus,
		Claude3Dot5Haiku,
		Claude3Dot5Sonnet,
		Claude3Sonnet,
		Claude3Haiku,
		Claude2Dot1,
		Clause2Dot0,
		Claude1Dot2Instant,
	}
}

func CurrentModels() []string {
	return []string{
		Claude3Dot7SonnetLatest,
		Claude3Dot5Haiku,
		Claude3OpusLatest,
	}
}
