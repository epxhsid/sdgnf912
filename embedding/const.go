package embedding

const (
	EmbeddingDimensions = 384
	MaxSequenceLength   = 512

	queryInstruction = "Represent this sentence for searching relevant passages: "
)

type EmbedType string

const (
	Document EmbedType = "document"
	Query    EmbedType = "query"
)
