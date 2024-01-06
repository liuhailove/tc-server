package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/pion/webrtc/v3"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"

	"github.com/liuhailove/tc-base-go/mediatransportutil/pkg/rtcconfig"
	"github.com/liuhailove/tc-base-go/protocol/logger"
	redisTC "github.com/liuhailove/tc-base-go/protocol/redis"
)

// CongestionControlProbeMode 拥塞控制探测模式
type CongestionControlProbeMode string

// StreamTrackerType 流轨类型
type StreamTrackerType string

const (
	generatedCLIFlagUsage = "generated"

	CongestionControlProbeModePadding CongestionControlProbeMode = "padding"
	CongestionControlProbeModeMedia   CongestionControlProbeMode = "media"

	StreamTrackerTypePacket StreamTrackerType = "packet"
	StreamTrackerTypeFrame  StreamTrackerType = "frame"

	StatsUpdateInterval = time.Second * 10
	// TelemetryStatsUpdateInterval 遥测状态更新间隔
	TelemetryStatsUpdateInterval = time.Second * 30
)

var (
	ErrKeyFileIncorrectPermission = errors.New("key file others permission must be set to 0")
	ErrKeysNotSet                 = errors.New("one of key-file or keys must be provided")
)

type Config struct {
	Port           uint32              `yaml:"port"`
	BindAddresses  []string            `yaml:"bind_addresses,omitempty"`
	PrometheusPort uint32              `yaml:"prometheus_port,omitempty"`
	Environment    string              `yaml:"environment,omitempty"`
	RTC            RTCConfig           `yaml:"rtc,omitempty"`
	Redis          redisTC.RedisConfig `yaml:"redis,omitempty"`
	Audio          AudioConfig         `yaml:"audio,omitempty"`
	Video          VideoConfig         `yaml:"video,omitempty"`
	Room           RoomConfig          `yaml:"room,omitempty"`
	TURN           TURNConfig          `yaml:"turn,omitempty"`
	Ingress        IngressConifg       `yaml:"ingress,omitempty"`
	WebHook        WebHookConfig       `yaml:"webhook,omitempty"`
	NodeSelector   NodeSelectorConfig  `yaml:"node_selector,omitempty"`
	KeyFile        string              `yaml:"key_file,omitempty"`
	Keys           map[string]string   `yaml:"keys,omitempty"`
	Region         string              `yaml:"region,omitempty"`
	SignalRelay    SignalRelayConfig   `yaml:"signal_relay,omitempty"`
	// LogLevel is deprecated
	LogLevel string        `yaml:"log_level,omitempty"`
	Logging  LoggingConfig `yaml:"logging,omitempty"`
	Limit    LimitConfig   `yaml:"limit,omitempty"`

	Development bool `yaml:"development,omitempty"`
}

type RTCConfig struct {
	rtcconfig.RTCConfig `yaml:"-"`

	TURNServers []TURNServer `yaml:"turn_servers,omitempty"`

	StrictACKS bool `yaml:"strict_acks,omitempty"`

	// NACK 缓冲的数据包数量
	PacketBufferSize int `yaml:"packet_buffer_size,omitempty"`

	// pli/fir rtcp 数据包的限制周期
	PLIThrottle PLIThrottleConfig `yaml:"pli_throttle,omitempty"`

	CongestionControl CongestionControlConfig `yaml:"congestion_control,omitempty"`

	// 允许 TCP 和 TURN/TLS 回退
	AllowTCPFallback *bool `yaml:"allow_tcp_fallback,omitempty"`

	// 出现发布错误时强制重新连接
	ReconnectOnPublicationError *bool `yaml:"reconnect_on_publication_error,omitempty"`

	// 订阅错误时强制重新连接
	ReconnectOnSubscriptionError *bool `yaml:"reconnect_on_subscription_error,omitempty"`

	// 允许时间戳调整以保持低漂移，这是实验性的
	AllowTimestampAdjustment *bool `yaml:"allow_timestamp_adjustment,omitempty"`
}

type TURNServer struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Protocol   string `yaml:"protocol"`
	Username   string `yaml:"username,omitempty"`
	Credential string `yaml:"credential,omitempty"`
}

type PLIThrottleConfig struct {
	LowQuality  time.Duration `yaml:"low_quality,omitempty"`
	MidQuality  time.Duration `yaml:"mid_quality,omitempty"`
	HighQuality time.Duration `yaml:"high_quality,omitempty"`
}

type CongestionControlProbeConfig struct {
	BaseInterval  time.Duration `yaml:"base_interval,omitempty"`
	BackoffFactor float64       `yaml:"backoff_factor,omitempty"`
	MaxInterval   time.Duration `yaml:"max_interval,omitempty"`

	SettleWait    time.Duration `yaml:"settle_wait,omitempty"`
	SettleWaitMax time.Duration `yaml:"settle_wait_max,omitempty"`

	TrendWait time.Duration `yaml:"trend_wait,omitempty"`

	OveragePct             int64         `yaml:"overage_pct,omitempty"`
	MinBps                 int64         `yaml:"min_bps,omitempty"`
	MinDuration            time.Duration `yaml:"min_duration,omitempty"`
	MaxDuration            time.Duration `yaml:"max_duration,omitempty"`
	DurationOverflowFactor float64       `yaml:"duration_overflow_factor,omitempty"`
	DurationIncreaseFactor float64       `yaml:"duration_increase_factor,omitempty"`
}

type CongestionControlChannelObserverConfig struct {
	EstimateRequiredSamples        int           `yaml:"estimate_required_samples,omitempty"`
	EstimateRequiredSamplesMin     int           `yaml:"estimate_required_samples_min,omitempty"`
	EstimateDownwardTrendThreshold float64       `yaml:"estimate_downward_trend_threshold,omitempty"`
	EstimateDownwardTrendMaxWait   time.Duration `yaml:"estimate_downward_trend_max_wait,omitempty"`
	EstimateCollapseThreshold      time.Duration `yaml:"estimate_collapse_threshold,omitempty"`
	EstimateValidityWindow         time.Duration `yaml:"estimate_validity_window,omitempty"`
	NackWindowMinDuration          time.Duration `yaml:"nack_window_min_duration,omitempty"`
	NackWindowMaxDuration          time.Duration `yaml:"nack_window_max_duration,omitempty"`
	NackRatioThreshold             float64       `yaml:"nack_ratio_threshold,omitempty"`
}

type CongestionControlConfig struct {
	Enabled                       bool                                   `yaml:"enabled,omitempty"`
	AllowPause                    bool                                   `yaml:"allow_pause,omitempty"`
	NackRatioAttenuator           float64                                `yaml:"nack_ratio_attenuator,omitempty"`
	ExpectedUsageThreshold        float64                                `yaml:"expected_usage_threshold,omitempty"`
	UseSendSideBWE                bool                                   `yaml:"send_side_bandwidth_estimation,omitempty"`
	ProbeMode                     CongestionControlProbeMode             `yaml:"padding_mode,omitempty"`
	MinChannelCapacity            int64                                  `yaml:"min_channel_capacity,omitempty"`
	ProbeConfig                   CongestionControlProbeConfig           `yaml:"probe_config,omitempty"`
	ChannelObserverProbeConfig    CongestionControlChannelObserverConfig `yaml:"channel_observer_probe_config,omitempty"`
	ChannelObserverNonProbeConfig CongestionControlChannelObserverConfig `yaml:"channel_observer_non_probe_config,omitempty"`
}

type AudioConfig struct {
	// 被视为活动的最低级别，0-127，其中 0 是最大声
	ActiveLevel uint8 `yaml:"active_level,omitempty"`
	// 要测量的百分位数，如果参与者超过 ActiveLevel 超过，则被认为是活跃的
	// MinPercentile% 的时间
	MinPercentile uint8 `yaml:"min_percentile,omitempty"`
	// 更新客户端的时间间隔，以毫秒为单位
	UpdateInterval uint32 `yaml:"update_interval,omitempty"`
	// 发送到客户端的音频级别值的平滑。
	// audioLevel 将是 `smooth_intervals` 的平均值，0 表示禁用
	SmoothIntervals uint32 `yaml:"smooth_intervals,omitempty"`
	// 为仅 opus 的音频上轨道启用红色编码下轨道
	ActiveREDEncoding bool `yaml:"active_red_encoding,omitempty"`
}

type StreamTrackerPacketConfig struct {
	SamplesRequired uint32        `yaml:"samples_required,omitempty"` // 每个周期所需的样本数
	CyclesRequired  uint32        `yaml:"cycles_required,omitempty"`  // 激活所需的周期数
	CycleDuration   time.Duration `yaml:"cycle_duration,omitempty"`
}

type StreamTrackerFrameConfig struct {
	MinFPS float64 `yaml:"min_fps"`
}

type StreamTrackerConfig struct {
	StreamTrackerType     StreamTrackerType                   `yaml:"stream_tracker_type,omitempty"`
	BitrateReportInterval map[int32]time.Duration             `yaml:"bitrate_report_interval,omitempty"`
	PacketTracker         map[int32]StreamTrackerPacketConfig `yaml:"packet_tracker,omitempty"`
	FrameTracker          map[int32]StreamTrackerFrameConfig  `yaml:"frame_tracker,omitempty"`
}

type StreamTrackersConfig struct {
	Video       StreamTrackerConfig `yaml:"video,omitempty"`
	Screenshare StreamTrackerConfig `yaml:"screenshare,omitempty"`
}

type PlayoutDelayConfig struct {
	Enabled bool `yaml:"enabled,omitempty"`
	Min     int  `yaml:"min,omitempty"`
}

type VideoConfig struct {
	DynacastPauseDelay time.Duration        `yaml:"dynacast_pause_delay,omitempty"`
	StreamTracker      StreamTrackersConfig `yaml:"stream_tracker,omitempty"`
}

type RoomConfig struct {
	// 启用自动创建房间
	AutoCreate          bool               `yaml:"auto_create,omitempty"`
	EnabledCodecs       []CodecSpec        `yaml:"enabled_codecs,omitempty"`
	MaxParticipants     uint32             `yaml:"max_participants,omitempty"`
	EmptyTimeout        uint32             `yaml:"empty_timeout,omitempty"`
	EnabledRemoteUnmute bool               `yaml:"enabled_remote_unmute,omitempty"`
	MaxMetadataSize     uint32             `yaml:"max_metadata_size,omitempty"`
	PlayoutDelay        PlayoutDelayConfig `yaml:"playout_delay,omitempty"`
}

type CodecSpec struct {
	Mime     string `yaml:"mime"`
	FmtpLine string `yaml:"fmtp_line"`
}

type LoggingConfig struct {
	logger.Config `yaml:",inline"`
	PionLevel     string `yaml:"pion_level,omitempty"`
}

type TURNConfig struct {
	Enabled             bool   `yaml:"enabled"`
	Domain              string `yaml:"domain,omitempty"`
	CertFile            string `yaml:"cert_file,omitempty"`
	KeyFile             string `yaml:"key_file,omitempty"`
	TLSPort             int    `yaml:"tls_port,omitempty"`
	UDPPort             int    `yaml:"udp_port,omitempty"`
	RelayPortRangeStart uint16 `yaml:"relay_port_range_start,omitempty"`
	RelayPortRangeEnd   uint16 `yaml:"relay_port_range_end,omitempty"`
	ExternalTLS         bool   `yaml:"external_tls,omitempty"`
}

type WebHookConfig struct {
	URLS []string `yaml:"urls"`
	// 用于webhook的key
	APIKey string `yaml:"api_key"`
}

type NodeSelectorConfig struct {
	Kind         string         `yaml:"kind"`
	SortBy       string         `yaml:"sort_by,omitempty"`
	CPULoadLimit float32        `yaml:"cpu_load_limit,omitempty"`
	SysloadLimit float32        `yaml:"sysload_limit,omitempty"`
	Regions      []RegionConfig `yaml:"regions,omitempty"`
}

type SignalRelayConfig struct {
	Enabled          bool          `yaml:"enabled"`
	RetryTimeout     time.Duration `yaml:"retry_timeout,omitempty"`
	MinRetryInterval time.Duration `yaml:"min_retry_interval,omitempty"`
	MaxRetryInterval time.Duration `yaml:"max_retry_interval,omitempty"`
	StreamBufferSize int           `yaml:"stream_buffer_size,omitempty"`
}

// RegionConfig 列出了可用区域及其纬度/经度，因此选择器会更喜欢
// 较近的区域
type RegionConfig struct {
	Name string  `yaml:"name"`
	Lat  float64 `yaml:"lat"`
	Lon  float64 `yaml:"lon"`
}

type LimitConfig struct {
	NumTracks              int32   `yaml:"num_tracks,omitempty"`
	BytesPerSec            float32 `yaml:"bytes_per_sec,omitempty"`
	SubscriptionLimitVideo int32   `yaml:"subscription_limit_video,omitempty"`
	SubscriptionLimitAudio int32   `yaml:"subscription_limit_audio,omitempty"`
}

type IngressConfig struct {
	RTMPBaseURL string `yaml:"rtmp_base_url"`
	WHIPBaseURL string `yaml:"whip_base_url"`
}

// APIConfig 不要暴露给YAML
type APIConfig struct {
	// 等待API执行的时间，默认2s
	ExecutionTimeout time.Duration

	// 检查操作完成之前等待的时间
	CheckInterval time.Duration
}

func DefaultAPIConfig() APIConfig {
	return APIConfig{
		ExecutionTimeout: 2 * time.Second,
		CheckInterval:    100 * time.Millisecond,
	}
}

var DefaultConfig = Config{
	Port: 7880,
	RTC: RTCConfig{
		RTCConfig: rtcconfig.RTCConfig{
			UseExternalIP:     false,
			TCPPort:           7881,
			UDPPort:           0,
			ICEPortRangeStart: 0,
			ICEPortRangeEnd:   0,
			STUNServers:       []string{},
		},
		PacketBufferSize: 500,
		StrictACKS:       true,
		PLIThrottle: PLIThrottleConfig{
			LowQuality:  500 * time.Millisecond,
			MidQuality:  time.Second,
			HighQuality: time.Second,
		},
		CongestionControl: CongestionControlConfig{
			Enabled:                true,
			AllowPause:             false,
			NackRatioAttenuator:    0.4,
			ExpectedUsageThreshold: 0.95,
			ProbeMode:              CongestionControlProbeModePadding,
			ProbeConfig: CongestionControlProbeConfig{
				BaseInterval:  3 * time.Second,
				BackoffFactor: 1.5,
				MaxInterval:   2 * time.Minute,

				SettleWait:    250 * time.Millisecond,
				SettleWaitMax: 10 * time.Second,

				TrendWait: 2 * time.Second,

				OveragePct:             120,
				MinBps:                 200_000,
				MinDuration:            200 * time.Millisecond,
				MaxDuration:            20 * time.Second,
				DurationOverflowFactor: 1.25,
				DurationIncreaseFactor: 1.5,
			},
			ChannelObserverProbeConfig: CongestionControlChannelObserverConfig{
				EstimateRequiredSamples:        3,
				EstimateRequiredSamplesMin:     3,
				EstimateDownwardTrendThreshold: 0.0,
				EstimateDownwardTrendMaxWait:   5 * time.Second,
				EstimateCollapseThreshold:      0,
				EstimateValidityWindow:         10 * time.Second,
				NackWindowMinDuration:          500 * time.Millisecond,
				NackWindowMaxDuration:          1 * time.Second,
				NackRatioThreshold:             0.04,
			},
			ChannelObserverNonProbeConfig: CongestionControlChannelObserverConfig{
				EstimateRequiredSamples:        12,
				EstimateRequiredSamplesMin:     8,
				EstimateDownwardTrendThreshold: -0.6,
				EstimateDownwardTrendMaxWait:   5 * time.Second,
				EstimateCollapseThreshold:      10 * time.Second,
				EstimateValidityWindow:         10 * time.Second,
				NackWindowMinDuration:          2 * time.Second,
				NackWindowMaxDuration:          3 * time.Second,
				NackRatioThreshold:             0.08,
			},
		},
	},
	Audio: AudioConfig{
		ActiveLevel:     35, // -35dBov
		MinPercentile:   40,
		UpdateInterval:  400,
		SmoothIntervals: 2,
	},
	Video: VideoConfig{
		DynacastPauseDelay: 5 * time.Second,
		StreamTracker: StreamTrackersConfig{
			Video: StreamTrackerConfig{
				StreamTrackerType: StreamTrackerTypePacket,
				BitrateReportInterval: map[int32]time.Duration{
					0: 1 * time.Second,
					1: 1 * time.Second,
					2: 1 * time.Second,
				},
				PacketTracker: map[int32]StreamTrackerPacketConfig{
					0: {
						SamplesRequired: 1,
						CyclesRequired:  4,
						CycleDuration:   500 * time.Millisecond,
					},
					1: {
						SamplesRequired: 5,
						CyclesRequired:  20,
						CycleDuration:   500 * time.Millisecond,
					},
					2: {
						SamplesRequired: 5,
						CyclesRequired:  20,
						CycleDuration:   500 * time.Millisecond,
					},
				},
				FrameTracker: map[int32]StreamTrackerFrameConfig{
					0: {
						MinFPS: 5.0,
					},
					1: {
						MinFPS: 5.0,
					},
					2: {
						MinFPS: 5.0,
					},
				},
			},
			Screenshare: StreamTrackerConfig{
				StreamTrackerType: StreamTrackerTypePacket,
				BitrateReportInterval: map[int32]time.Duration{
					0: 4 * time.Second,
					1: 4 * time.Second,
					2: 4 * time.Second,
				},
				PacketTracker: map[int32]StreamTrackerPacketConfig{
					0: {
						SamplesRequired: 1,
						CyclesRequired:  1,
						CycleDuration:   2 * time.Second,
					},
					1: {
						SamplesRequired: 1,
						CyclesRequired:  1,
						CycleDuration:   2 * time.Second,
					},
					2: {
						SamplesRequired: 1,
						CyclesRequired:  1,
						CycleDuration:   2 * time.Second,
					},
				},
				FrameTracker: map[int32]StreamTrackerFrameConfig{
					0: {
						MinFPS: 0.5,
					},
					1: {
						MinFPS: 0.5,
					},
					2: {
						MinFPS: 0.5,
					},
				},
			},
		},
	},
	Redis: redisTC.RedisConfig{},
	Room: RoomConfig{
		AutoCreate: true,
		EnabledCodecs: []CodecSpec{
			{Mime: webrtc.MimeTypeOpus},
			{Mime: "audio/red"},
			{Mime: webrtc.MimeTypeVP8},
			{Mime: webrtc.MimeTypeH264},
			// {Mime: webrtc.MimeTypeAV1},
			// {Mime: webrtc.MimeTypeVP9},
		},
		EmptyTimeout: 5 * 60,
	},
	Logging: LoggingConfig{
		PionLevel: "error",
	},
	TURN: TURNConfig{
		Enabled: false,
	},
	NodeSelector: NodeSelectorConfig{
		Kind:         "any",
		SortBy:       "random",
		SysloadLimit: 0.9,
		CPULoadLimit: 0.9,
	},
	SignalRelay: SignalRelayConfig{
		Enabled:          false,
		RetryTimeout:     7500 * time.Millisecond,
		MinRetryInterval: 500 * time.Millisecond,
		MaxRetryInterval: 4 * time.Second,
		StreamBufferSize: 1000,
	},
	Keys: map[string]string{},
}

func NewConfig(confString string, strictMode bool, c *cli.Context, baseFlags []cli.Flag) (*Config, error) {
	// 使用默认值启动
	conf := DefaultConfig
	if confString != "" {
		decoder := yaml.NewDecoder(strings.NewReader(confString))
		decoder.KnownFields(strictMode)
		if err := decoder.Decode(&conf); err != nil {
			return nil, fmt.Errorf("could not parse config: %v", err)
		}
	}

	if err := conf.RTC.Validate(conf.Development); err != nil {
		return nil, fmt.Errorf("could not validate RTC config: %v", err)
	}

	if c != nil {
		if err := conf.updateFromCLI(c, baseFlags); err != nil {
			return nil, err
		}
	}

	// 扩展文件名中的环境变量
	file, err := homedir.Expand(os.ExpandEnv(conf.KeyFile))
	if err != nil {
		return nil, err
	}
	conf.KeyFile = file

	// 如果没有设置，则设置 Turn 中继的默认值
	if conf.TURN.RelayPortRangeStart == 0 || conf.TURN.RelayPortRangeEnd == 0 {
		// 为了更容易在dev模式/docker下运行，默认有两个端口
		if conf.Development {
			conf.TURN.RelayPortRangeStart = 30000
			conf.TURN.RelayPortRangeEnd = 30002
		} else {
			conf.TURN.RelayPortRangeStart = 30000
			conf.TURN.RelayPortRangeEnd = 40000
		}
	}

	if conf.LogLevel != "" {
		conf.Logging.Level = conf.LogLevel
	}
	if conf.Logging.Level == "" && conf.Development {
		conf.Logging.Level = "debug"
	}
	if conf.Logging.PionLevel != "" {
		if conf.Logging.ComponentLevels == nil {
			conf.Logging.ComponentLevels = map[string]string{}
		}
		conf.Logging.ComponentLevels["transport.pion"] = conf.Logging.PionLevel
		conf.Logging.ComponentLevels["pion"] = conf.Logging.PionLevel
	}

	if conf.Development {
		conf.Environment = "dev"
	}

	return &conf, nil
}

func (conf *Config) IsTURNSEnabled() bool {
	if conf.TURN.Enabled && conf.TURN.TLSPort != 0 {
		return true
	}
	for _, s := range conf.RTC.TURNServers {
		if s.Protocol == "tls" {
			return true
		}
	}
	return false
}

type configNode struct {
	TypeNode  reflect.Value
	TagPrefix string
}

func (conf *Config) ToCLIFlagNames(existingFlags []cli.Flag) map[string]reflect.Value {
	existingFlagNames := map[string]bool{}
	for _, flag := range existingFlags {
		for _, flagName := range flag.Names() {
			existingFlagNames[flagName] = true
		}
	}

	flagNames := map[string]reflect.Value{}
	var currNode configNode
	nodes := []configNode{{reflect.ValueOf(conf).Elem(), ""}}
	for len(nodes) > 0 {
		currNode, nodes = nodes[0], nodes[1:]
		for i := 0; i < currNode.TypeNode.NumField(); i++ {
			// 检查 struct 字段中的 yaml 标记以获取路径
			field := currNode.TypeNode.Type().Field(i)
			yamlTagArray := strings.SplitN(field.Tag.Get("yaml"), ",", 2)
			yamlTag := yamlTagArray[0]
			isInline := false
			if len(yamlTagArray) > 1 && yamlTagArray[1] == "inline" {
				isInline = true
			}
			if (yamlTag == "" && (!isInline || currNode.TagPrefix == "")) || yamlTag == "-" {
				continue
			}
			yamlPath := yamlTag
			if currNode.TagPrefix != "" {
				if isInline {
					yamlPath = currNode.TagPrefix
				} else {
					yamlPath = fmt.Sprintf("%s.%s", currNode.TagPrefix, yamlTag)
				}
			}
			if existingFlagNames[yamlPath] {
				continue
			}

			// 将标志名称映射到值
			value := currNode.TypeNode.Field(i)
			if value.Kind() == reflect.Struct {
				nodes = append(nodes, configNode{value, yamlPath})
			} else {
				flagNames[yamlPath] = value
			}
		}
	}

	return flagNames
}

func (conf *Config) ValidateKeys() error {
	// 如果设置则首选keyfile文件
	if conf.KeyFile != "" {
		var otherFilter os.FileMode = 0007
		if st, err := os.Stat(conf.KeyFile); err != nil {
			return err
		} else if st.Mode().Perm()&otherFilter != 0000 {
			return ErrKeyFileIncorrectPermission
		}
		f, err := os.Open(conf.KeyFile)
		if err != nil {
			return err
		}
		defer func() {
			_ = f.Close()
		}()
		decoder := yaml.NewDecoder(f)
		if err = decoder.Decode(conf.Keys); err != nil {
			return err
		}
	}

	if len(conf.Keys) == 0 {
		return ErrKeysNotSet
	}

	if !conf.Development {
		for key, secret := range conf.Keys {
			if len(secret) < 32 {
				logger.Errorw("secret is too short, should be at least 32 characters for security", nil, "apiKey", key)
			}
		}
	}
	return nil
}

func GenerateCLIFlags(existingFlags []cli.Flag, hidden bool) ([]cli.Flag, error) {
	blankConfig := &Config{}
	flags := make([]cli.Flag, 0)
	for name, value := range blankConfig.ToCLIFlagNames(existingFlags) {
		kind := value.Kind()
		if kind == reflect.Ptr {
			kind = value.Type().Elem().Kind()
		}

		var flag cli.Flag
		envVar := fmt.Sprintf("TC_%s", strings.ToUpper(strings.Replace(name, ".", "-", -1)))

		switch kind {
		case reflect.Bool:
			flag = &cli.BoolFlag{
				Name:   name,
				Usage:  generatedCLIFlagUsage,
				Hidden: hidden,
			}
		case reflect.String:
			flag = &cli.StringFlag{
				Name:    name,
				EnvVars: []string{envVar},
				Usage:   generatedCLIFlagUsage,
				Hidden:  hidden,
			}
		case reflect.Int, reflect.Int32:
			flag = &cli.IntFlag{
				Name:    name,
				EnvVars: []string{envVar},
				Usage:   generatedCLIFlagUsage,
				Hidden:  hidden,
			}
		case reflect.Int64:
			flag = &cli.Int64Flag{
				Name:    name,
				EnvVars: []string{envVar},
				Usage:   generatedCLIFlagUsage,
				Hidden:  hidden,
			}
		case reflect.Uint8, reflect.Uint16, reflect.Uint32:
			flag = &cli.UintFlag{
				Name:    name,
				EnvVars: []string{envVar},
				Usage:   generatedCLIFlagUsage,
				Hidden:  hidden,
			}
		case reflect.Float32:
			flag = &cli.Float64Flag{
				Name:    name,
				EnvVars: []string{envVar},
				Usage:   generatedCLIFlagUsage,
				Hidden:  hidden,
			}
		case reflect.Float64:
			flag = &cli.Float64Flag{
				Name:    name,
				EnvVars: []string{envVar},
				Usage:   generatedCLIFlagUsage,
				Hidden:  hidden,
			}
		case reflect.Slice:
			// TODO
			continue
		case reflect.Map:
			// TODO
			continue
		default:
			return flags, fmt.Errorf("cli flag generation unsupported for config type: %s is a %s", name, kind.String())
		}

		flags = append(flags, flag)
	}

	return flags, nil
}

func (conf *Config) updateFromCLI(c *cli.Context, baseFlags []cli.Flag) error {
	generatedFlagNames := conf.ToCLIFlagNames(baseFlags)
	for _, flag := range c.App.Flags {
		flagName := flag.Names()[0]

		// the `c.App.Name != "test"` check is needed because `c.IsSet(...)` is always false in unit tests
		if !c.IsSet(flagName) && c.App.Name != "test" {
			continue
		}

		configValue, ok := generatedFlagNames[flagName]
		if !ok {
			continue
		}

		kind := configValue.Kind()
		if kind == reflect.Ptr {
			// 实例化要设置的值
			configValue.Set(reflect.New(configValue.Type().Elem()))

			kind = configValue.Type().Elem().Kind()
			configValue = configValue.Elem()
		}

		switch kind {
		case reflect.Bool:
			configValue.SetBool(c.Bool(flagName))
		case reflect.String:
			configValue.SetString(c.String(flagName))
		case reflect.Int, reflect.Int32, reflect.Int64:
			configValue.SetInt(c.Int64(flagName))
		case reflect.Uint8, reflect.Uint16, reflect.Uint32:
			configValue.SetUint(c.Uint64(flagName))
		case reflect.Float32:
			configValue.SetFloat(c.Float64(flagName))
		case reflect.Float64:
			configValue.SetFloat(c.Float64(flagName))
		// case reflect.Slice:
		// 	// TODO
		// case reflect.Map:
		// 	// TODO
		default:
			return fmt.Errorf("unsupported generated cli flag type for config: %s is a %s", flagName, kind.String())
		}
	}

	if c.IsSet("dev") {
		conf.Development = c.Bool("dev")
	}
	if c.IsSet("key-file") {
		conf.KeyFile = c.String("key-file")
	}
	if c.IsSet("keys") {
		if err := conf.unmarshalKeys(c.String("keys")); err != nil {
			return errors.New("Could not parse keys, it needs to be exactly, \"key: secret\", including the space")
		}
	}
	if c.IsSet("region") {
		conf.Region = c.String("region")
	}
	if c.IsSet("redis-host") {
		conf.Redis.Address = c.String("redis-host")
	}
	if c.IsSet("redis-password") {
		conf.Redis.Password = c.String("redis-password")
	}
	if c.IsSet("turn-cert") {
		conf.TURN.CertFile = c.String("turn-cert")
	}
	if c.IsSet("turn-key") {
		conf.TURN.KeyFile = c.String("turn-key")
	}
	if c.IsSet("node-ip") {
		conf.RTC.NodeIP = c.String("node-ip")
	}
	if c.IsSet("udp-port") {
		conf.RTC.UDPPort = uint32(c.Int("udp-port"))
	}
	if c.IsSet("bind") {
		conf.BindAddresses = c.StringSlice("bind")
	}
	return nil
}

func (conf *Config) unmarshalKeys(keys string) error {
	temp := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(keys), temp); err != nil {
		return err
	}

	conf.Keys = make(map[string]string, len(temp))

	for key, val := range temp {
		if secret, ok := val.(string); ok {
			conf.Keys[key] = secret
		}
	}
	return nil
}

// SetLogger Note: only pass in logr.Logger with default depth
func SetLogger(l logger.Logger) {
	logger.SetLogger(l, "gglive")
}

func InitLoggerFromConfig(config *LoggingConfig) {
	logger.InitFromConfig(&config.Config, "tc")
}
