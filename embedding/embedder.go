package embedding

import (
	"fmt"
	"log"
	"math"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/pretrained"
	ort "github.com/yalue/onnxruntime_go"
)

type Embedder struct {
	tokenizer *tokenizer.Tokenizer
	session   *ort.DynamicAdvancedSession
}

func New(modelPath, tokenizerPath string) (*Embedder, error) {
	tk, err := pretrained.FromFile(tokenizerPath)
	if err != nil {
		log.Fatal(err)
	}

	session, err := ort.NewDynamicAdvancedSession(
		modelPath,
		[]string{
			"input_ids",
			"attention_mask",
			"token_type_ids",
		},
		[]string{
			"last_hidden_state",
		},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create ONNX session: %w", err)
	}

	return &Embedder{
		tokenizer: tk,
		session:   session,
	}, nil
}

func (e *Embedder) Embed(text string, embedType EmbedType) ([]float32, error) {
	switch embedType {
	case Query:
		text = queryInstruction + text
	case Document:
	default:
		return nil, fmt.Errorf("unknown embedding type: %q", embedType)
	}

	encoding, err := e.tokenizer.EncodeSingle(text, true)
	if err != nil {
		return nil, fmt.Errorf("tokenize text: %w", err)
	}

	ids := encoding.GetIds()
	attentionMask := encoding.GetAttentionMask()
	typeIDs := encoding.GetTypeIds()
	seqLen := len(ids)

	if seqLen == 0 {
		return nil, fmt.Errorf("tokenizer returned empty sequence")
	}

	if seqLen > MaxSequenceLength {
		return nil, fmt.Errorf(
			"text produces %d tokens, maximum is %d",
			seqLen,
			MaxSequenceLength,
		)
	}

	if len(attentionMask) != seqLen {
		return nil, fmt.Errorf(
			"attention mask length mismatch: got %d, expected %d",
			len(attentionMask),
			seqLen,
		)
	}

	if len(typeIDs) != seqLen {
		return nil, fmt.Errorf(
			"token type IDs length mismatch: got %d, expected %d",
			len(typeIDs),
			seqLen,
		)
	}

	// ONNX model expects int64 tensors
	inputIDs := make([]int64, seqLen)
	mask := make([]int64, seqLen)
	types := make([]int64, seqLen)

	for i := range seqLen {
		inputIDs[i] = int64(ids[i])
		mask[i] = int64(attentionMask[i])
		types[i] = int64(typeIDs[i])
	}

	// input_ids: [batch, sequence_length]
	inputTensor, err := ort.NewTensor(
		ort.NewShape(1, int64(seqLen)),
		inputIDs,
	)
	if err != nil {
		return nil, fmt.Errorf("create input_ids tensor: %w", err)
	}
	defer inputTensor.Destroy()

	// attention_mask: [batch, sequence_length]
	maskTensor, err := ort.NewTensor(
		ort.NewShape(1, int64(seqLen)),
		mask,
	)
	if err != nil {
		return nil, fmt.Errorf("create attention_mask tensor: %w", err)
	}
	defer maskTensor.Destroy()

	// token_type_ids: [batch, sequence_length]
	typeTensor, err := ort.NewTensor(
		ort.NewShape(1, int64(seqLen)),
		types,
	)
	if err != nil {
		return nil, fmt.Errorf("create token_type_ids tensor: %w", err)
	}
	defer typeTensor.Destroy()

	// last_hidden_state: [batch, sequence_length, hidden_size]
	output, err := ort.NewEmptyTensor[float32](
		ort.NewShape(
			1,
			int64(seqLen),
			EmbeddingDimensions,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create output tensor: %w", err)
	}
	defer output.Destroy()

	err = e.session.Run(
		[]ort.Value{
			inputTensor,
			maskTensor,
			typeTensor,
		},
		[]ort.Value{
			output,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("run inference: %w", err)
	}

	outputData := output.GetData()

	expectedSize := seqLen * EmbeddingDimensions

	if len(outputData) != expectedSize {
		return nil, fmt.Errorf(
			"unexpected output size: got %d, expected %d",
			len(outputData),
			expectedSize,
		)
	}

	// BGE uses the [CLS] token representation.
	//
	// Output shape:
	//
	// [1, seqLen, 384]
	//
	// Since batch size is 1, the first 384 values
	// correspond to the [CLS] token.
	embedding := make([]float32, EmbeddingDimensions)
	copy(embedding, outputData[:EmbeddingDimensions])

	// L2 normalize the embedding.
	var sumSquares float64

	for _, value := range embedding {
		sumSquares += float64(value) * float64(value)
	}

	norm := math.Sqrt(sumSquares)

	if norm == 0 {
		return nil, fmt.Errorf("embedding has zero norm")
	}

	for i := range embedding {
		embedding[i] = float32(float64(embedding[i]) / norm)
	}

	return embedding, nil
}

func (e *Embedder) Close() error {
	if e.session != nil {
		return e.session.Destroy()
	}

	return nil
}
