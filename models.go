package claude

const (
	Claude3Dot5Sonnet  = "claude-3-5-sonnet-20240620"
	Claude3Opus        = "claude-3-opus-20240229"
	Claude3Sonnet      = "claude-3-sonnet-20240229"
	Claude3Haiku       = "claude-3-haiku-20240307"
	Claude2Dot1        = "claude-2.1"
	Clause2Dot0        = "claude-2.0"
	Claude1Dot2Instant = "claude-instant-1.2"
)

func Models() []string {
	return []string{
		Claude3Dot5Sonnet,
		Claude3Opus,
		Claude3Sonnet,
		Claude3Haiku,
		Claude2Dot1,
		Clause2Dot0,
		Claude1Dot2Instant,
	}
}
