package v1alpha1

import "github.com/opensergo/opensergo-control-plane/pkg/api/v1alpha1"

const (
	// Type of persistence
	MEMORY = 0 // memory storage
	LFS    = 1 // local file system storage
	DFS    = 2 // distribution file system storage, eg: remote file system, object storage ...

	// Storage full Strategy
	BLOCK = 0 // block when full
	DROP  = 1 // drop new data when full
	LAST  = 2 // drop the oldest data when full
	ERROR = 3 // return error when full

	// Back off policy
	BackOffPolicyLinear      = 0 // Linear growth strategy
	BackOffPolicyExponential = 1 // Exponential growth strategy is based on 2 by default
)

// EventChannel Event pipeline corresponding message queue subject
type EventChannel struct {
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Required
	ID string `json:"id"`

	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Required
	URL string `json:"url"`

	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Required
	MQTopicName string `json:"mqTopicName"`
}

// EventProcessorRef logic processor
type EventProcessorRef struct {
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Required
	UniqueID string `json:"uniqueID"`

	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Required
	// Processor type such as Service represents service
	Kind string `json:"kind"`

	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Required
	// The processor name, such as events-demo, represents the service named events-demo
	Name string `json:"name"`
}

// EventSource event source
type EventSource struct {
	Ref EventProcessorRef `json:"ref"`
}

// EventTrigger event trigger
type EventTrigger struct {
	Ref EventProcessorRef `json:"ref"`
}

// EventComponents components of event components
type EventComponents struct {
	Channels []EventChannel `json:"channels"`
	Sources  []EventSource  `json:"sources"`
	Triggers []EventTrigger `json:"triggers"`
}

// PersistenceAddress address of persistence
type PersistenceAddress struct {
	// The storage address is determined according to the PersistenceType
	// For example, local file directory/usr/local/data remote file system address xxx. xxx. xxx. xxx
	Address string `json:"address"`
}

// Persistence detail of persistence
type Persistence struct {
	// +kubebuilder:validation:Type=int
	// +kubebuilder:validation:Required
	PersistenceType int `json:"persistenceType"`

	// +kubebuilder:validation:Type=int64
	// +kubebuilder:validation:Required
	PersistenceSize int64 `json:"persistenceSize"`

	// +kubebuilder:validation:Type=int
	// +kubebuilder:validation:Required
	FullStrategy int `json:"fullStrategy"`
}

// RetryRule rule of retry
type RetryRule struct {
	// +kubebuilder:validation:Type=integer
	// +kubebuilder:validation:Format=int64
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Required
	RetryMax int64 `json:"retryMax"`

	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Required
	// The delay between retries uses the string of the time period specified in the ISO-8601 rule
	// for example, P6S represents the duration of 6s
	BackOffDelay string `json:"backOffDelay"`

	// +kubebuilder:validation:Type=int
	// +kubebuilder:validation:Required
	BackOffPolicyType int `json:"backOffPolicyType"`
}

// DeadLetterStrategy dead letter message strategy for EventTrigger
type DeadLetterStrategy struct {
	// +kubebuilder:validation:Type=bool
	// +kubebuilder:validation:Required
	Enable bool `json:"enable"`

	// +kubebuilder:validation:Type=integer
	// +kubebuilder:validation:Format=int64
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Required
	// Threshold value of retry triggering dead letter
	RetryTriggerThreshold int64 `json:"retryTriggerThreshold"`

	// When enable_ Effective when block=true
	// EventChannel indicates the channel storing dead letter production or consumption messages
	// Persistence indicates use local or remote storage
	// If StoreEventChannel and StorePersistence both have value, StoreEventChannel has higher priority.
	StoreEventChannel EventChannel `json:"storeEventChannel"`
	StorePersistence  Persistence  `json:"storePersistence"`
}

// EventRuntimeStrategy runtime strategy for event source or trigger
// centralized setting FaultTolerance RateLimit CircuitBreaker
type EventRuntimeStrategy struct {
	FaultToleranceRule       v1alpha1.FaultToleranceRuleSpec       `json:"faultToleranceRule"`
	RateLimitStrategy        v1alpha1.RateLimitStrategySpec        `json:"rateLimitStrategy"`
	CircuitBreakerStrategy   v1alpha1.CircuitBreakerStrategySpec   `json:"circuitBreakerStrategy"`
	ConcurrencyLimitStrategy v1alpha1.ConcurrencyLimitStrategySpec `json:"concurrencyLimitStrategy"`

	RetryRule RetryRule `json:"retryRule"`
}

// EventSourceStrategy strategy for event source producer
type EventSourceStrategy struct {
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Required
	EventSourceID string `json:"eventSourceID"`

	// +kubebuilder:validation:Type=bool
	// +kubebuilder:validation:Optional
	AsyncSend bool `json:"asyncSend"`

	FaultTolerantStorage Persistence `json:"faultTolerantStorage"`

	RuntimeStrategy EventRuntimeStrategy `json:"runtimeStrategy"`
}

// EventTriggerStrategy strategy for event trigger or consumer
type EventTriggerStrategy struct {
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Required
	EventTriggerID string `json:"eventTriggerID"`

	// +kubebuilder:validation:Type=int64
	// +kubebuilder:validation:Required
	ReceiveBufferSize int64 `json:"receiveBufferSize"`

	// +kubebuilder:validation:Type=bool
	// +kubebuilder:validation:Optional
	EnableIdempotence bool `json:"enableIdempotence"`

	RuntimeStrategy EventRuntimeStrategy `json:"runtimeStrategy"`

	DeadLetterStrategy DeadLetterStrategy `json:"deadLetterStrategy"`
}

// EventStrategies strategy for event
type EventStrategies struct {
	SourceStrategies  []EventSourceStrategy  `json:"sourceStrategies"`
	TriggerStrategies []EventTriggerStrategy `json:"triggerStrategies"`

	DefaultSourceRuntimeStrategy  EventRuntimeStrategy `json:"defaultSourceRuntimeStrategy"`
	DefaultTriggerRuntimeStrategy EventRuntimeStrategy `json:"defaultTriggerRuntimeStrategy"`
	DefaultDeadLetterStrategy     DeadLetterStrategy   `json:"defaultDeadLetterStrategy"`
}

// EventFilter filter before trigger
type EventFilter struct {
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Optional
	// CloudEvents SQL Expression is based on cloudevents sql expression
	CESQL string `json:"cesql"`
}

// EventRouter event router
type EventRouter struct {
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Required
	SourceID string `json:"sourceID"`

	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Required
	TriggerID string `json:"triggerID"`

	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Required
	ChannelID string `json:"channelID"`

	Filter EventFilter `json:"filter"`
}

// EventRouterRules event routing rules include source trigger filter
type EventRouterRules struct {
	RouterRules []EventRouter `json:"routerRules"`
}
