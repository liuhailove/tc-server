package dependencydescriptor

import (
	"fmt"
	"math"
	"strconv"
)

// DependencyDescriptorExtension is a extension payload format in
// https://aomediacodec.github.io/av1-rtp-spec/#dependency-descriptor-rtp-header-extension

func formatBitmask(b *uint32) string {
	if b == nil {
		return "-"
	}
	return strconv.FormatInt(int64(*b), 2)
}

// ------------------------------------------------------------------------------

type DependencyDescriptorExtension struct {
	Descriptor *DependencyDescriptor
	Structure  *FrameDependencyStructure
}

func (d *DependencyDescriptorExtension) Marshal() ([]byte, error) {
	return d.MarshalWithActiveChains(^uint32(0))
}

func (d *DependencyDescriptorExtension) MarshalWithActiveChains(activeChains uint32) ([]byte, error) {
	writer, err := NewDependencyDescriptorWriter(nil, d.Structure, activeChains, d.Descriptor)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, int(math.Ceil(float64(writer.ValueSizeBits())/8)))
	writer.ResetBuf(buf)
	if err = writer.Write(); err != nil {
		return nil, err
	}
	return buf, nil
}

func (d *DependencyDescriptorExtension) Unmarshal(buf []byte) (int, error) {
	reader := NewDependencyDescriptorReader(buf, d.Structure, d.Descriptor)
	return reader.Parse()
}

// ------------------------------------------------------------------------------

const (
	MaxSpatialIds    = 4
	MaxTemporalIds   = 8
	MaxDecodeTargets = 32
	MaxTemplates     = 64

	AllChainAreActive = uint32(0)

	ExtensionURI = "https://aomediacodec.github.io/av1-rtp-spec/#dependency-descriptor-rtp-header-extension"
)

// ------------------------------------------------------------------------------

type DependencyDescriptor struct {
	FirstPacketInFrame         bool
	LastPacketInFrame          bool
	FrameNumber                uint16
	FrameDependencies          *FrameDependencyTemplate
	Resolution                 *RenderResolution
	ActiveDecodeTargetsBitmask *uint32
	AttachedStructure          *FrameDependencyStructure
}

func (d *DependencyDescriptor) MarshalSize() (int, error) {
	return d.MarshalSizeWithActiveChains(^uint32(0))
}

func (d *DependencyDescriptor) MarshalSizeWithActiveChains(activeChains uint32) (int, error) {
	writer, err := NewDependencyDescriptorWriter(nil, d.AttachedStructure, activeChains, d)
	if err != nil {
		return 0, err
	}
	return int(math.Ceil(float64(writer.ValueSizeBits()) / 8)), nil
}

func (d *DependencyDescriptor) String() string {
	resolution, dependencies := "-", "-"
	if d.Resolution != nil {
		resolution = fmt.Sprintf("%+v", *d.Resolution)
	}
	if d.FrameDependencies != nil {
		dependencies = fmt.Sprintf("%+v", *d.FrameDependencies)
	}
	return fmt.Sprintf("DependencyDescriptor{FirstPacketInFrame: %v, LastPacketInFrame: %v, FrameNumber: %v, FrameDependencies: %s, Resolution: %s, ActiveDecodeTargetsBitmask: %v, AttachedStructure: %v}",
		d.FirstPacketInFrame, d.LastPacketInFrame, d.FrameNumber, dependencies, resolution, formatBitmask(d.ActiveDecodeTargetsBitmask), d.AttachedStructure)
}

// DecodeTargetIndication 帧与解码目标的关系。
type DecodeTargetIndication int

const (
	DecodeTargetNotPresent  DecodeTargetIndication = iota // DecodeTargetInfo symbol '-'
	DecodeTargetDiscardable                               // DecodeTargetInfo symbol 'D'
	DecodeTargetSwitch                                    // DecodeTargetInfo symbol 'S'
	DecodeTargetRequired                                  // DecodeTargetInfo symbol 'R'
)

func (i DecodeTargetIndication) String() string {
	switch i {
	case DecodeTargetNotPresent:
		return "-"
	case DecodeTargetDiscardable:
		return "D"
	case DecodeTargetSwitch:
		return "S"
	case DecodeTargetRequired:
		return "R"
	default:
		return "Unknown"
	}
}

// ------------------------------------------------------------------------------

type FrameDependencyTemplate struct {
	SpatialId               int
	TemporalId              int
	DecodeTargetIndications []DecodeTargetIndication
	FrameDiffs              []int
	ChainDiffs              []int
}

func (t *FrameDependencyTemplate) Clone() *FrameDependencyTemplate {
	t2 := &FrameDependencyTemplate{
		SpatialId:  t.SpatialId,
		TemporalId: t.TemporalId,
	}

	t2.DecodeTargetIndications = make([]DecodeTargetIndication, len(t.DecodeTargetIndications))
	copy(t2.DecodeTargetIndications, t.DecodeTargetIndications)

	t2.FrameDiffs = make([]int, len(t.FrameDiffs))
	copy(t2.FrameDiffs, t.FrameDiffs)

	t2.ChainDiffs = make([]int, len(t.ChainDiffs))
	copy(t2.ChainDiffs, t.ChainDiffs)

	return t2
}

// ------------------------------------------------------------------------------

type FrameDependencyStructure struct {
	StructureId      int
	NumDecodeTargets int
	NumChains        int
	// 如果使用链 (num_chains > 0)，则将解码目标索引映射到索引
	// 保护该目标的链。
	DecodeTargetProtectedByChain []int
	Resolutions                  []RenderResolution
	Templates                    []*FrameDependencyTemplate
}

func (f *FrameDependencyStructure) String() string {
	str := fmt.Sprintf("FrameDependencyStructure{StructureId: %v, NumDecodeTargets: %v, NumChains: %v, DecodeTargetProtectedByChain: %v, Resolutions: %+v, Templates: [",
		f.StructureId, f.NumDecodeTargets, f.NumChains, f.DecodeTargetProtectedByChain, f.Resolutions)
	// templates
	for _, t := range f.Templates {
		str += fmt.Sprintf("%+v, ", t)
	}
	str += "]}"

	return str
}

// ------------------------------------------------------------------------------

// RenderResolution 渲染分辨率
type RenderResolution struct {
	Width  int
	Height int
}

// ------------------------------------------------------------------------------
