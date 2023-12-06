package buffer

import "fmt"

const (
	InvalidLayerSpatial  = int32(-1) // 无效空间层
	InvalidLayerTemporal = int32(-1) // 无效时间层

	DefaultMaxLayerSpatial  = int32(2) // 默认最大空间层
	DefaultMaxLayerTemporal = int32(3) // 默认最大时间层
)

var (
	InvalidLayer = VideoLayer{
		Spatial:  InvalidLayerSpatial,
		Temporal: InvalidLayerTemporal,
	}

	DefaultMaxLayer = VideoLayer{
		Spatial:  DefaultMaxLayerSpatial,
		Temporal: DefaultMaxLayerTemporal,
	}
)

// VideoLayer 用于描述视频编码的层次结构
type VideoLayer struct {
	Spatial  int32 // 空间描述
	Temporal int32 // 时间描述
}

func (v VideoLayer) String() string {
	return fmt.Sprintf("VideoLayer{s: %d, t:%d}", v.Spatial, v.Temporal)
}

// GreaterThan 视频编码比较
func (v VideoLayer) GreaterThan(v2 VideoLayer) bool {
	return v.Spatial > v2.Spatial || (v.Spatial == v2.Spatial && v.Temporal > v2.Temporal)
}

// SpatialGraterThanOrEqual 空间层比较
func (v VideoLayer) SpatialGraterThanOrEqual(v2 VideoLayer) bool {
	return v.Spatial >= v2.Spatial
}

func (v VideoLayer) IsValid() bool {
	return v.Spatial != InvalidLayerSpatial && v.Temporal != InvalidLayerTemporal
}
