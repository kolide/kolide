package kolide

import "time"

// HostStore enrolls hosts in the datastore
type OsqueryStore interface {
	// Host methods
	EnrollHost(uuid, hostname, ip, platform string, nodeKeySize int) (*Host, error)
	AuthenticateHost(nodeKey string) (*Host, error)
	MarkHostSeen(host *Host, t time.Time) error
	GetLabelQueriesForHost(host *Host) (map[string]string, error)

	// Query methods
	NewQuery(query *Query) error

	// Label methods
	NewLabel(label *Label) error
}

type Host struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	NodeKey   string `gorm:"unique_index:idx_host_unique_nodekey"`
	HostName  string
	UUID      string `gorm:"unique_index:idx_host_unique_uuid"`
	IPAddress string
	Platform  string
	Labels    []*Label `gorm:"many2many:host_labels;"`
}

type Label struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"not null;unique_index:idx_label_unique_name"`
	QueryID   uint
	Hosts     []Host
}

type LabelQueryExecution struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Matches   bool
	LabelID   uint
	HostID    uint
}

type ScheduledQuery struct {
	ID           uint `gorm:"primary_key"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Name         string `gorm:"not null"`
	QueryID      int
	Query        Query
	Interval     uint `gorm:"not null"`
	Snapshot     bool
	Differential bool
	PackID       uint
}

type Query struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Platform  string
	Query     string   `gorm:"not null"`
	Targets   []Target `gorm:"many2many:query_targets"`
}

type TargetType int

const (
	TargetLabel TargetType = iota
	TargetHost  TargetType = iota
)

type Target struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Type      TargetType
	QueryID   uint
	TargetID  uint
}

type DistributedQueryStatus int

const (
	QueryRunning  DistributedQueryStatus = iota
	QueryComplete DistributedQueryStatus = iota
	QueryError    DistributedQueryStatus = iota
)

type DistributedQuery struct {
	ID          uint `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Query       Query
	MaxDuration time.Duration
	Status      DistributedQueryStatus
	UserID      uint
}

type DistributedQueryExecutionStatus int

const (
	ExecutionWaiting   DistributedQueryExecutionStatus = iota
	ExecutionRequested DistributedQueryExecutionStatus = iota
	ExecutionSucceeded DistributedQueryExecutionStatus = iota
	ExecutionFailed    DistributedQueryExecutionStatus = iota
)

type DistributedQueryExecution struct {
	HostID             uint
	DistributedQueryID uint
	Status             DistributedQueryExecutionStatus
	Error              string `gorm:"size:1024"`
	ExecutionDuration  time.Duration
}

type Pack struct {
	ID               uint `gorm:"primary_key"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Name             string `gorm:"not null;unique_index:idx_pack_unique_name"`
	Platform         string
	Queries          []ScheduledQuery
	DiscoveryQueries []DiscoveryQuery
}

type DiscoveryQuery struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Query     string `gorm:"size:1024" gorm:"not null"`
}

type Option struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Key       string `gorm:"not null;unique_index:idx_option_unique_key"`
	Value     string `gorm:"not null"`
	Platform  string
}

type DecoratorType int

const (
	DecoratorLoad     DecoratorType = iota
	DecoratorAlways   DecoratorType = iota
	DecoratorInterval DecoratorType = iota
)

type Decorator struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Type      DecoratorType `gorm:"not null"`
	Interval  int
	Query     string
}
