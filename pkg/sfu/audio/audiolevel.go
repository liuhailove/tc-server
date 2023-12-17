package audio

import (
	"go.uber.org/atomic"
	"math"
)

const (
	silentAudioLevel = 127
	negInv20         = -1.0 / 20
)

// AudioLevelParams Audio级别参数
type AudioLevelParams struct {
	ActiveLevel     uint8
	MinPercentile   uint8
	ObserveDuration uint32
	SmoothIntervals uint32
}

// AudioLevel 跟踪参与者的音频级别
type AudioLevel struct {
	params AudioLevelParams
	//观察持续时间窗口内被视为活动的最短持续时间
	minActiveDuration uint32
	smoothFactor      float64
	activeThreshold   float64

	smoothedLevel atomic.Float64

	loudestObservedLevel uint8
	activeDuration       uint32 // ms
	observedDuration     uint32 // ms
}

func NewAudioLevel(params AudioLevelParams) *AudioLevel {
	l := &AudioLevel{
		params:               params,
		minActiveDuration:    uint32(params.MinPercentile) * params.ObserveDuration / 100,
		smoothFactor:         1,
		activeThreshold:      ConvertAudioLevel(float64(params.ActiveLevel)),
		loudestObservedLevel: silentAudioLevel,
	}
	if l.params.SmoothIntervals > 0 {
		// 指数移动平均线 (EMA)，与简单移动平均线 (SMA) 具有相同的质心
		l.smoothFactor = float64(2) / (float64(l.params.SmoothIntervals + 1))
	}

	return l
}

// Observe 观察新帧，必须从同一线程调用
func (l *AudioLevel) Observe(level uint8, durationMs uint32) {
	l.observedDuration += durationMs

	if level <= l.params.ActiveLevel {
		l.activeDuration += durationMs
		if l.loudestObservedLevel > level {
			l.loudestObservedLevel = level
		}
	}

	if l.observedDuration >= l.params.ObserveDuration {
		// 计算并重置
		if l.activeDuration >= l.minActiveDuration {
			// 根据窗口的活动量来调整观察到的最大声音级别。
			// 如果在整个持续时间内都处于活动状态，则权重将为 0
			// > 0 如果活动时间长于观察持续时间，则为 0
			// < 0 如果活跃时间小于观察持续时间
			activityWeight := 20 * math.Log10(float64(l.activeDuration)/float64(l.params.ObserveDuration))
			adjustedLevel := float64(l.loudestObservedLevel) - activityWeight
			linearLevel := ConvertAudioLevel(adjustedLevel)

			// 指数平滑以抑制瞬变
			smoothedLevel := l.smoothedLevel.Load()
			smoothedLevel += (linearLevel - smoothedLevel) * l.smoothFactor
			l.smoothedLevel.Store(smoothedLevel)
		} else {
			l.smoothedLevel.Store(0)
		}
		l.loudestObservedLevel = silentAudioLevel
		l.activeDuration = 0
		l.observedDuration = 0
	}
}

// GetLevel 返回当前舒缓的音频级别
func (l *AudioLevel) GetLevel() (float64, bool) {
	smoothedLevel := l.smoothedLevel.Load()
	active := smoothedLevel >= l.activeThreshold
	return smoothedLevel, active
}

// ConvertAudioLevel 将分贝转换回线性
func ConvertAudioLevel(level float64) float64 {
	return math.Pow(10, level*negInv20)
}
