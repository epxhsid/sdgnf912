package embedding

/*
func embed(tk *tokenizer.Tokenizer, session *ort.DynamicAdvancedSession, text string) ([]float32, error) {
	encoding, err := tk.EncodeSingle(text, true)
	if err != nil {
		return nil, err
	}

	ids := encoding.GetIds()
	attentionMask := encoding.GetAttentionMask()
	typeIds := encoding.GetTypeIds()
	seqLen := int64(len(ids))

	// we have to specifically convert []int → []int64 because the ONNX model
	// explicitly requires int64 inputs
	inputIDs := make([]int64, len(ids))
	attention := make([]int64, len(attentionMask))
	tokenTypes := make([]int64, len(typeIds))
	for i := range ids {
		inputIDs[i] = int64(ids[i])
		attention[i] = int64(attentionMask[i])
		tokenTypes[i] = int64(typeIds[i])
	}

	// dimensions: [batch_size, sequence_length]
	inputIDsTensor, err := ort.NewTensor(
		ort.NewShape(1, seqLen),
		inputIDs,
	)
	if err != nil {
		return nil, err
	}
	defer inputIDsTensor.Destroy()

	attentionTensor, err := ort.NewTensor(
		ort.NewShape(1, seqLen),
		attention,
	)
	if err != nil {
		return nil, err
	}
	defer attentionTensor.Destroy()

	tokenTypesTensor, err := ort.NewTensor(
		ort.NewShape(1, seqLen),
		tokenTypes,
	)
	if err != nil {
		return nil, err
	}
	defer tokenTypesTensor.Destroy()

	// output tensor
	output, err := ort.NewEmptyTensor[float32](
		ort.NewShape(1, seqLen, 384),
	)
	if err != nil {
		return nil, err
	}
	defer output.Destroy()

	inputs := []ort.Value{
		inputIDsTensor,
		attentionTensor,
		tokenTypesTensor,
	}

	outputs := []ort.Value{
		output,
	}

	err = session.Run(inputs, outputs)
	if err != nil {
		return nil, err
	}

	// last_hidden_state shape:
	// [1, sequence_length, 384]
	//
	// BGE uses the first token ([CLS]) as the sentence embedding.
	data := output.GetData()

	embedding := make([]float32, 384)
	copy(embedding, data[:384])

	// L2 normalization
	var sum float32
	for _, v := range embedding {
		sum += v * v
	}

	norm := float32(math.Sqrt(float64(sum)))

	if norm == 0 {
		return nil, fmt.Errorf("zero-length embedding")
	}

	for i := range embedding {
		embedding[i] /= norm
	}

	return embedding, nil
}
*/
