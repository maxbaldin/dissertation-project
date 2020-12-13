package usecase_test

import (
	"testing"

	"github.com/maxbaldin/dissertation-project/src/implementation/ui/entity"
	"github.com/maxbaldin/dissertation-project/src/implementation/ui/usecase"
	"github.com/stretchr/testify/assert"
)

func TestNodeTransformer_Transform(t *testing.T) {
	data := []entity.GraphDataElement{
		{
			SourceHost:    "c1f3c0fe5e3c",
			SourceService: "service_b",
			TargetHost:    "2463cf10eb0f",
			TargetService: "service_a",
			PacketsCnt:    82551,
			SizeBytes:     (1024 * 1024 / 8) * 3600,
		},
		{
			SourceHost:    "2463cf10eb0f",
			SourceService: "service_a",
			TargetHost:    "c1f3c0fe5e3c",
			TargetService: "service_b",
			PacketsCnt:    100833,
			SizeBytes:     (1024 * 1024 / 8) * 3600,
		},
	}
	transformer := usecase.NewNodeTransformer()
	out, err := transformer.Transform(data)
	assert.NoError(t, err)
	jsonResponse, err := out.MarshalJSON()
	assert.NoError(t, err)
	jsonResponseString := string(jsonResponse)
	expected := `{"nodes":[{"data":{"id":"c1f3c0fe5e3c"}},{"data":{"id":"2463cf10eb0f"}},{"data":{"id":"c1f3c0fe5e3c-service_b","parent":"c1f3c0fe5e3c"}},{"data":{"id":"2463cf10eb0f-service_a","parent":"2463cf10eb0f"}}],"edges":[{"data":{"id":"c1f3c0fe5e3c-service_b-2463cf10eb0f-service_a","source":"c1f3c0fe5e3c-service_b","target":"2463cf10eb0f-service_a","label":"5.196MiB","width":"1px"}},{"data":{"id":"2463cf10eb0f-service_a-c1f3c0fe5e3c-service_b","source":"2463cf10eb0f-service_a","target":"c1f3c0fe5e3c-service_b","label":"30.96MiB","width":"6px"}}]}`
	assert.Equal(t, expected, jsonResponseString)
}
