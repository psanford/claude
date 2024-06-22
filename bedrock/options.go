package bedrock

type Option interface {
	set(*Client)
}
