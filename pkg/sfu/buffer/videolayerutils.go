package buffer

import (
	"github.com/liuhailove/tc-base-go/protocol/logger"
	"github.com/liuhailove/tc-base-go/protocol/tc"
)

const (
	QuarterResolution = "q" // 四分之一分辨率
	HalfResolution    = "h" // 半分辨率
	FullResolution    = "f" // 全分辨率
)

// LayerPresenceFromTrackInfo 层显示来自轨迹信息
func LayerPresenceFromTrackInfo(trackInfo *tc.TrackInfo) *[tc.VideoQuality_HIGH + 1]bool {
	if trackInfo == nil || len(trackInfo.Layers) == 0 {
		return nil
	}

	var layerPresence [tc.VideoQuality_HIGH + 1]bool
	for _, layer := range trackInfo.Layers {
		layerPresence[layer.Quality] = true
	}

	return &layerPresence
}

// RidToSpatialLayer 网格层转空间层
func RidToSpatialLayer(rid string, trackInfo *tc.TrackInfo) int32 {
	lp := LayerPresenceFromTrackInfo(trackInfo)
	if lp == nil {
		switch rid {
		case QuarterResolution:
			return 0
		case HalfResolution:
			return 1
		case FullResolution:
			return 2
		default:
			return 0
		}
	}

	switch rid {
	case QuarterResolution:
		switch {
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_MEDIUM] && lp[tc.VideoQuality_HIGH]:
			fallthrough
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_MEDIUM]:
			fallthrough
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_HIGH]:
			fallthrough
		case lp[tc.VideoQuality_MEDIUM] && lp[tc.VideoQuality_HIGH]:
			return 0

		default:
			// 只有一个质量发布，可以是任何
			return 0
		}
	case HalfResolution:
		switch {
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_MEDIUM] && lp[tc.VideoQuality_HIGH]:
			fallthrough
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_MEDIUM]:
			fallthrough
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_HIGH]:
			fallthrough
		case lp[tc.VideoQuality_MEDIUM] && lp[tc.VideoQuality_HIGH]:
			return 1

		default:
			// 只有一个质量发布，可以是任何
			return 0
		}

	case FullResolution:
		switch {
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_MEDIUM] && lp[tc.VideoQuality_HIGH]:
			return 2

		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_MEDIUM]:
			logger.Warnw("unexpected rid f with only two qualities, low and medium", nil)
			return 1
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_HIGH]:
			logger.Warnw("unexpected rid f with only two qualities, low and high", nil)
			return 1
		case lp[tc.VideoQuality_MEDIUM] && lp[tc.VideoQuality_HIGH]:
			logger.Warnw("unexpected rid f with only two qualities, medium and high", nil)
			return 1

		default:
			// 只有一个质量发布，可以是任何
			return 0
		}

	default:
		// no rid, should be single layer
		return 0
	}
}

// SpatialLayerToRid 空间层转网格层
func SpatialLayerToRid(layer int32, trackInfo *tc.TrackInfo) string {
	lp := LayerPresenceFromTrackInfo(trackInfo)
	if lp == nil {
		switch layer {
		case 0:
			return QuarterResolution
		case 1:
			return HalfResolution
		case 2:
			return FullResolution
		default:
			return QuarterResolution
		}
	}

	switch layer {
	case 0:
		switch {
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_MEDIUM] && lp[tc.VideoQuality_HIGH]:
			fallthrough
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_MEDIUM]:
			fallthrough
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_HIGH]:
			fallthrough
		case lp[tc.VideoQuality_MEDIUM] && lp[tc.VideoQuality_HIGH]:
			return QuarterResolution

		default:
			return QuarterResolution
		}

	case 1:
		switch {
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_MEDIUM] && lp[tc.VideoQuality_HIGH]:
			fallthrough
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_MEDIUM]:
			fallthrough
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_HIGH]:
			fallthrough
		case lp[tc.VideoQuality_MEDIUM] && lp[tc.VideoQuality_HIGH]:
			return HalfResolution

		default:
			return QuarterResolution
		}

	case 2:
		switch {
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_MEDIUM] && lp[tc.VideoQuality_HIGH]:
			return FullResolution

		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_MEDIUM]:
			logger.Warnw("unexpected layer 2 with only two qualities, low and medium", nil)
			return HalfResolution
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_HIGH]:
			logger.Warnw("unexpected layer 2 with only two qualities, low and high", nil)
			return HalfResolution
		case lp[tc.VideoQuality_MEDIUM] && lp[tc.VideoQuality_HIGH]:
			logger.Warnw("unexpected layer 2 with only two qualities, medium and high", nil)
			return HalfResolution

		default:
			return QuarterResolution
		}

	default:
		return QuarterResolution
	}
}

// VideoQualityToRid 视频质量转网格
func VideoQualityToRid(quality tc.VideoQuality, trackInfo *tc.TrackInfo) string {
	return SpatialLayerToRid(VideoQualityToSpatialLayer(quality, trackInfo), trackInfo)
}

func SpatialLayerToVideoQuality(layer int32, trackInfo *tc.TrackInfo) tc.VideoQuality {
	lp := LayerPresenceFromTrackInfo(trackInfo)
	if lp == nil {
		switch layer {
		case 0:
			return tc.VideoQuality_LOW
		case 1:
			return tc.VideoQuality_MEDIUM
		case 2:
			return tc.VideoQuality_HIGH
		default:
			return tc.VideoQuality_OFF
		}
	}

	switch layer {
	case 0:
		switch {
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_MEDIUM] && lp[tc.VideoQuality_HIGH]:
			fallthrough
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_MEDIUM]:
			fallthrough
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_HIGH]:
			fallthrough
		case lp[tc.VideoQuality_LOW]:
			return tc.VideoQuality_LOW

		case lp[tc.VideoQuality_MEDIUM] && lp[tc.VideoQuality_HIGH]:
			fallthrough
		case lp[tc.VideoQuality_MEDIUM]:
			return tc.VideoQuality_MEDIUM

		default:
			return tc.VideoQuality_HIGH
		}

	case 1:
		switch {
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_MEDIUM] && lp[tc.VideoQuality_HIGH]:
			fallthrough
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_MEDIUM]:
			return tc.VideoQuality_MEDIUM
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_HIGH]:
			fallthrough
		case lp[tc.VideoQuality_MEDIUM] && lp[tc.VideoQuality_HIGH]:
			return tc.VideoQuality_HIGH

		default:
			logger.Errorw("invalid layer", nil, "layer", layer, "trackInfo", trackInfo)
			return tc.VideoQuality_HIGH
		}

	case 2:
		switch {
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_MEDIUM] && lp[tc.VideoQuality_HIGH]:
			return tc.VideoQuality_HIGH

		default:
			logger.Errorw("invalid layer", nil, "layer", layer, "trackInfo", trackInfo)
			return tc.VideoQuality_HIGH
		}
	}

	return tc.VideoQuality_OFF
}

// VideoQualityToSpatialLayer 视频质量转空间层
func VideoQualityToSpatialLayer(quality tc.VideoQuality, trackInfo *tc.TrackInfo) int32 {
	lp := LayerPresenceFromTrackInfo(trackInfo)
	if lp == nil {
		switch quality {
		case tc.VideoQuality_LOW:
			return 0
		case tc.VideoQuality_MEDIUM:
			return 1
		case tc.VideoQuality_HIGH:
			return 2
		default:
			return InvalidLayerSpatial
		}
	}

	switch quality {
	case tc.VideoQuality_LOW:
		switch {
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_MEDIUM] && lp[tc.VideoQuality_HIGH]:
			fallthrough
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_MEDIUM]:
			fallthrough
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_HIGH]:
			fallthrough
		case lp[tc.VideoQuality_MEDIUM] && lp[tc.VideoQuality_HIGH]:
			fallthrough
		default: // only one quality published, could be any
			return 0
		}

	case tc.VideoQuality_MEDIUM:
		switch {
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_MEDIUM] && lp[tc.VideoQuality_HIGH]:
			fallthrough
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_MEDIUM]:
			fallthrough
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_HIGH]:
			return 1
		case lp[tc.VideoQuality_MEDIUM] && lp[tc.VideoQuality_HIGH]:
			return 0

		default: // only one quality published, could be any
			return 0
		}

	case tc.VideoQuality_HIGH:
		switch {
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_MEDIUM] && lp[tc.VideoQuality_HIGH]:
			return 2
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_MEDIUM]:
			fallthrough
		case lp[tc.VideoQuality_LOW] && lp[tc.VideoQuality_HIGH]:
			fallthrough
		case lp[tc.VideoQuality_MEDIUM] && lp[tc.VideoQuality_HIGH]:
			return 1
		default: // only one quality published, could be any
			return 0
		}
	}

	return InvalidLayerSpatial
}
