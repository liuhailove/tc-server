package buffer

import (
	"github.com/liuhailove/tc-base-go/protocol/tc"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRidConversion(t *testing.T) {
	type RidAndLayer struct {
		rid   string
		layer int32
	}
	tests := []struct {
		name       string
		trackInfo  *tc.TrackInfo
		ridToLayer map[string]RidAndLayer
	}{
		{
			"no track info",
			nil,
			map[string]RidAndLayer{
				"":                {rid: QuarterResolution, layer: 0},
				QuarterResolution: {rid: QuarterResolution, layer: 0},
				HalfResolution:    {rid: HalfResolution, layer: 1},
				FullResolution:    {rid: FullResolution, layer: 2},
			},
		},
		{
			"no layers",
			&tc.TrackInfo{},
			map[string]RidAndLayer{
				"":                {rid: QuarterResolution, layer: 0},
				QuarterResolution: {rid: QuarterResolution, layer: 0},
				HalfResolution:    {rid: HalfResolution, layer: 1},
				FullResolution:    {rid: FullResolution, layer: 2},
			},
		},
		{
			"single layer, low",
			&tc.TrackInfo{
				Layers: []*tc.VideoLayer{
					{Quality: tc.VideoQuality_LOW},
				},
			},
			map[string]RidAndLayer{
				"":                {rid: QuarterResolution, layer: 0},
				QuarterResolution: {rid: QuarterResolution, layer: 0},
				HalfResolution:    {rid: QuarterResolution, layer: 0},
				FullResolution:    {rid: QuarterResolution, layer: 0},
			},
		},
		{
			"single layer, medium",
			&tc.TrackInfo{
				Layers: []*tc.VideoLayer{
					{Quality: tc.VideoQuality_MEDIUM},
				},
			},
			map[string]RidAndLayer{
				"":                {rid: QuarterResolution, layer: 0},
				QuarterResolution: {rid: QuarterResolution, layer: 0},
				HalfResolution:    {rid: QuarterResolution, layer: 0},
				FullResolution:    {rid: QuarterResolution, layer: 0},
			},
		},
		{
			"single layer, high",
			&tc.TrackInfo{
				Layers: []*tc.VideoLayer{
					{Quality: tc.VideoQuality_HIGH},
				},
			},
			map[string]RidAndLayer{
				"":                {rid: QuarterResolution, layer: 0},
				QuarterResolution: {rid: QuarterResolution, layer: 0},
				HalfResolution:    {rid: QuarterResolution, layer: 0},
				FullResolution:    {rid: QuarterResolution, layer: 0},
			},
		},
		{
			"two layers, low and medium",
			&tc.TrackInfo{
				Layers: []*tc.VideoLayer{
					{Quality: tc.VideoQuality_LOW},
					{Quality: tc.VideoQuality_MEDIUM},
				},
			},
			map[string]RidAndLayer{
				"":                {rid: QuarterResolution, layer: 0},
				QuarterResolution: {rid: QuarterResolution, layer: 0},
				HalfResolution:    {rid: HalfResolution, layer: 1},
				FullResolution:    {rid: HalfResolution, layer: 1},
			},
		},
		{
			"two layers, low and high",
			&tc.TrackInfo{
				Layers: []*tc.VideoLayer{
					{Quality: tc.VideoQuality_LOW},
					{Quality: tc.VideoQuality_HIGH},
				},
			},
			map[string]RidAndLayer{
				"":                {rid: QuarterResolution, layer: 0},
				QuarterResolution: {rid: QuarterResolution, layer: 0},
				HalfResolution:    {rid: HalfResolution, layer: 1},
				FullResolution:    {rid: HalfResolution, layer: 1},
			},
		},
		{
			"two layers, medium and high",
			&tc.TrackInfo{
				Layers: []*tc.VideoLayer{
					{Quality: tc.VideoQuality_MEDIUM},
					{Quality: tc.VideoQuality_HIGH},
				},
			},
			map[string]RidAndLayer{
				"":                {rid: QuarterResolution, layer: 0},
				QuarterResolution: {rid: QuarterResolution, layer: 0},
				HalfResolution:    {rid: HalfResolution, layer: 1},
				FullResolution:    {rid: HalfResolution, layer: 1},
			},
		},
		{
			"three layers",
			&tc.TrackInfo{
				Layers: []*tc.VideoLayer{
					{Quality: tc.VideoQuality_LOW},
					{Quality: tc.VideoQuality_MEDIUM},
					{Quality: tc.VideoQuality_HIGH},
				},
			},
			map[string]RidAndLayer{
				"":                {rid: QuarterResolution, layer: 0},
				QuarterResolution: {rid: QuarterResolution, layer: 0},
				HalfResolution:    {rid: HalfResolution, layer: 1},
				FullResolution:    {rid: FullResolution, layer: 2},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for testRid, expectedResult := range test.ridToLayer {
				actualLayer := RidToSpatialLayer(testRid, test.trackInfo)
				require.Equal(t, expectedResult.layer, actualLayer)

				actualRid := SpatialLayerToRid(actualLayer, test.trackInfo)
				require.Equal(t, expectedResult.rid, actualRid)
			}
		})
	}
}

func TestQualityConversion(t *testing.T) {
	type QualityAndLayer struct {
		quality tc.VideoQuality
		layer   int32
	}
	tests := []struct {
		name           string
		trackInfo      *tc.TrackInfo
		qualityToLayer map[tc.VideoQuality]QualityAndLayer
	}{
		{
			"no track info",
			nil,
			map[tc.VideoQuality]QualityAndLayer{
				tc.VideoQuality_LOW:    {quality: tc.VideoQuality_LOW, layer: 0},
				tc.VideoQuality_MEDIUM: {quality: tc.VideoQuality_MEDIUM, layer: 1},
				tc.VideoQuality_HIGH:   {quality: tc.VideoQuality_HIGH, layer: 2},
			},
		},
		{
			"no layers",
			&tc.TrackInfo{},
			map[tc.VideoQuality]QualityAndLayer{
				tc.VideoQuality_LOW:    {quality: tc.VideoQuality_LOW, layer: 0},
				tc.VideoQuality_MEDIUM: {quality: tc.VideoQuality_MEDIUM, layer: 1},
				tc.VideoQuality_HIGH:   {quality: tc.VideoQuality_HIGH, layer: 2},
			},
		},
		{
			"single layer, low",
			&tc.TrackInfo{
				Layers: []*tc.VideoLayer{
					{Quality: tc.VideoQuality_LOW},
				},
			},
			map[tc.VideoQuality]QualityAndLayer{
				tc.VideoQuality_LOW:    {quality: tc.VideoQuality_LOW, layer: 0},
				tc.VideoQuality_MEDIUM: {quality: tc.VideoQuality_LOW, layer: 0},
				tc.VideoQuality_HIGH:   {quality: tc.VideoQuality_LOW, layer: 0},
			},
		},
		{
			"single layer, medium",
			&tc.TrackInfo{
				Layers: []*tc.VideoLayer{
					{Quality: tc.VideoQuality_MEDIUM},
				},
			},
			map[tc.VideoQuality]QualityAndLayer{
				tc.VideoQuality_LOW:    {quality: tc.VideoQuality_MEDIUM, layer: 0},
				tc.VideoQuality_MEDIUM: {quality: tc.VideoQuality_MEDIUM, layer: 0},
				tc.VideoQuality_HIGH:   {quality: tc.VideoQuality_MEDIUM, layer: 0},
			},
		},
		{
			"single layer, high",
			&tc.TrackInfo{
				Layers: []*tc.VideoLayer{
					{Quality: tc.VideoQuality_HIGH},
				},
			},
			map[tc.VideoQuality]QualityAndLayer{
				tc.VideoQuality_LOW:    {quality: tc.VideoQuality_HIGH, layer: 0},
				tc.VideoQuality_MEDIUM: {quality: tc.VideoQuality_HIGH, layer: 0},
				tc.VideoQuality_HIGH:   {quality: tc.VideoQuality_HIGH, layer: 0},
			},
		},
		{
			"two layers, low and medium",
			&tc.TrackInfo{
				Layers: []*tc.VideoLayer{
					{Quality: tc.VideoQuality_LOW},
					{Quality: tc.VideoQuality_MEDIUM},
				},
			},
			map[tc.VideoQuality]QualityAndLayer{
				tc.VideoQuality_LOW:    {quality: tc.VideoQuality_LOW, layer: 0},
				tc.VideoQuality_MEDIUM: {quality: tc.VideoQuality_MEDIUM, layer: 1},
				tc.VideoQuality_HIGH:   {quality: tc.VideoQuality_MEDIUM, layer: 1},
			},
		},
		{
			"two layers, low and high",
			&tc.TrackInfo{
				Layers: []*tc.VideoLayer{
					{Quality: tc.VideoQuality_LOW},
					{Quality: tc.VideoQuality_HIGH},
				},
			},
			map[tc.VideoQuality]QualityAndLayer{
				tc.VideoQuality_LOW:    {quality: tc.VideoQuality_LOW, layer: 0},
				tc.VideoQuality_MEDIUM: {quality: tc.VideoQuality_HIGH, layer: 1},
				tc.VideoQuality_HIGH:   {quality: tc.VideoQuality_HIGH, layer: 1},
			},
		},
		{
			"two layers, medium and high",
			&tc.TrackInfo{
				Layers: []*tc.VideoLayer{
					{Quality: tc.VideoQuality_MEDIUM},
					{Quality: tc.VideoQuality_HIGH},
				},
			},
			map[tc.VideoQuality]QualityAndLayer{
				tc.VideoQuality_LOW:    {quality: tc.VideoQuality_MEDIUM, layer: 0},
				tc.VideoQuality_MEDIUM: {quality: tc.VideoQuality_MEDIUM, layer: 0},
				tc.VideoQuality_HIGH:   {quality: tc.VideoQuality_HIGH, layer: 1},
			},
		},
		{
			"three layers",
			&tc.TrackInfo{
				Layers: []*tc.VideoLayer{
					{Quality: tc.VideoQuality_LOW},
					{Quality: tc.VideoQuality_MEDIUM},
					{Quality: tc.VideoQuality_HIGH},
				},
			},
			map[tc.VideoQuality]QualityAndLayer{
				tc.VideoQuality_LOW:    {quality: tc.VideoQuality_LOW, layer: 0},
				tc.VideoQuality_MEDIUM: {quality: tc.VideoQuality_MEDIUM, layer: 1},
				tc.VideoQuality_HIGH:   {quality: tc.VideoQuality_HIGH, layer: 2},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for testQuality, expectedResult := range test.qualityToLayer {
				actualLayer := VideoQualityToSpatialLayer(testQuality, test.trackInfo)
				require.Equal(t, expectedResult.layer, actualLayer)

				actualQuality := SpatialLayerToVideoQuality(actualLayer, test.trackInfo)
				require.Equal(t, expectedResult.quality, actualQuality)
			}
		})
	}
}

func TestVideoQualityToRidConversion(t *testing.T) {
	tests := []struct {
		name         string
		trackInfo    *tc.TrackInfo
		qualityToRid map[tc.VideoQuality]string
	}{
		{
			"no track info",
			nil,
			map[tc.VideoQuality]string{
				tc.VideoQuality_LOW:    QuarterResolution,
				tc.VideoQuality_MEDIUM: HalfResolution,
				tc.VideoQuality_HIGH:   FullResolution,
			},
		},
		{
			"no layers",
			&tc.TrackInfo{},
			map[tc.VideoQuality]string{
				tc.VideoQuality_LOW:    QuarterResolution,
				tc.VideoQuality_MEDIUM: HalfResolution,
				tc.VideoQuality_HIGH:   FullResolution,
			},
		},
		{
			"single layer, low",
			&tc.TrackInfo{
				Layers: []*tc.VideoLayer{
					{Quality: tc.VideoQuality_LOW},
				},
			},
			map[tc.VideoQuality]string{
				tc.VideoQuality_LOW:    QuarterResolution,
				tc.VideoQuality_MEDIUM: QuarterResolution,
				tc.VideoQuality_HIGH:   QuarterResolution,
			},
		},
		{
			"single layer, medium",
			&tc.TrackInfo{
				Layers: []*tc.VideoLayer{
					{Quality: tc.VideoQuality_MEDIUM},
				},
			},
			map[tc.VideoQuality]string{
				tc.VideoQuality_LOW:    QuarterResolution,
				tc.VideoQuality_MEDIUM: QuarterResolution,
				tc.VideoQuality_HIGH:   QuarterResolution,
			},
		},
		{
			"single layer, high",
			&tc.TrackInfo{
				Layers: []*tc.VideoLayer{
					{Quality: tc.VideoQuality_HIGH},
				},
			},
			map[tc.VideoQuality]string{
				tc.VideoQuality_LOW:    QuarterResolution,
				tc.VideoQuality_MEDIUM: QuarterResolution,
				tc.VideoQuality_HIGH:   QuarterResolution,
			},
		},
		{
			"two layers, low and medium",
			&tc.TrackInfo{
				Layers: []*tc.VideoLayer{
					{Quality: tc.VideoQuality_LOW},
					{Quality: tc.VideoQuality_MEDIUM},
				},
			},
			map[tc.VideoQuality]string{
				tc.VideoQuality_LOW:    QuarterResolution,
				tc.VideoQuality_MEDIUM: HalfResolution,
				tc.VideoQuality_HIGH:   HalfResolution,
			},
		},
		{
			"two layers, low and high",
			&tc.TrackInfo{
				Layers: []*tc.VideoLayer{
					{Quality: tc.VideoQuality_LOW},
					{Quality: tc.VideoQuality_HIGH},
				},
			},
			map[tc.VideoQuality]string{
				tc.VideoQuality_LOW:    QuarterResolution,
				tc.VideoQuality_MEDIUM: HalfResolution,
				tc.VideoQuality_HIGH:   HalfResolution,
			},
		},
		{
			"two layers, medium and high",
			&tc.TrackInfo{
				Layers: []*tc.VideoLayer{
					{Quality: tc.VideoQuality_MEDIUM},
					{Quality: tc.VideoQuality_HIGH},
				},
			},
			map[tc.VideoQuality]string{
				tc.VideoQuality_LOW:    QuarterResolution,
				tc.VideoQuality_MEDIUM: QuarterResolution,
				tc.VideoQuality_HIGH:   HalfResolution,
			},
		},
		{
			"three layers",
			&tc.TrackInfo{
				Layers: []*tc.VideoLayer{
					{Quality: tc.VideoQuality_LOW},
					{Quality: tc.VideoQuality_MEDIUM},
					{Quality: tc.VideoQuality_HIGH},
				},
			},
			map[tc.VideoQuality]string{
				tc.VideoQuality_LOW:    QuarterResolution,
				tc.VideoQuality_MEDIUM: HalfResolution,
				tc.VideoQuality_HIGH:   FullResolution,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for testQuality, expectedRid := range test.qualityToRid {
				actualRid := VideoQualityToRid(testQuality, test.trackInfo)
				require.Equal(t, expectedRid, actualRid)
			}
		})
	}
}
