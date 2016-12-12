package inmem

import (
	"errors"
	"sort"
	"strings"
	"time"

	kolide_errors "github.com/kolide/kolide-ose/server/errors"
	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/patrickmn/sortutil"
)

func (orm *Datastore) NewHost(host *kolide.Host) (*kolide.Host, error) {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	for _, h := range orm.hosts {
		if host.NodeKey == h.NodeKey || host.UUID == h.UUID {
			return nil, kolide_errors.ErrExists
		}
	}

	host.ID = orm.nextID(host)
	orm.hosts[host.ID] = host

	return host, nil
}

func (orm *Datastore) SaveHost(host *kolide.Host) error {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	if _, ok := orm.hosts[host.ID]; !ok {
		return kolide_errors.ErrNotFound
	}

	for _, nic := range host.NetworkInterfaces {
		if nic.ID == 0 {
			nic.ID = orm.nextID(nic)
		}
		nic.HostID = host.ID
	}
	host.ResetPrimaryNetwork()
	orm.hosts[host.ID] = host
	return nil
}

func (orm *Datastore) DeleteHost(host *kolide.Host) error {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	if _, ok := orm.hosts[host.ID]; !ok {
		return kolide_errors.ErrNotFound
	}

	delete(orm.hosts, host.ID)
	return nil
}

func (orm *Datastore) Host(id uint) (*kolide.Host, error) {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	host, ok := orm.hosts[id]
	if !ok {
		return nil, kolide_errors.ErrNotFound
	}

	return host, nil
}

func (orm *Datastore) ListHosts(opt kolide.ListOptions) ([]*kolide.Host, error) {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	// We need to sort by keys to provide reliable ordering
	keys := []int{}
	for k, _ := range orm.hosts {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	hosts := []*kolide.Host{}
	for _, k := range keys {
		hosts = append(hosts, orm.hosts[uint(k)])
	}

	// Apply ordering
	if opt.OrderKey != "" {
		var fields = map[string]string{
			"id":                 "ID",
			"created_at":         "CreatedAt",
			"updated_at":         "UpdatedAt",
			"detail_update_time": "DetailUpdateTime",
			"hostname":           "HostName",
			"uuid":               "UUID",
			"platform":           "Platform",
			"osquery_version":    "OsqueryVersion",
			"os_version":         "OSVersion",
			"uptime":             "Uptime",
			"memory":             "PhysicalMemory",
			"mac":                "PrimaryMAC",
			"ip":                 "PrimaryIP",
		}
		if err := sortResults(hosts, opt, fields); err != nil {
			return nil, err
		}
	}

	// Apply limit/offset
	low, high := orm.getLimitOffsetSliceBounds(opt, len(hosts))
	hosts = hosts[low:high]

	return hosts, nil
}

func (orm *Datastore) EnrollHost(osQueryHostID string, nodeKeySize int) (*kolide.Host, error) {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	if osQueryHostID == "" {
		return nil, errors.New("missing host identifier from osquery for host enrollment")
	}

	nodeKey, err := kolide.RandomText(nodeKeySize)
	if err != nil {
		return nil, err
	}

	host := kolide.Host{
		OsqueryHostID:    osQueryHostID,
		NodeKey:          nodeKey,
		DetailUpdateTime: time.Unix(0, 0).Add(24 * time.Hour),
	}

	host.CreatedAt = time.Now().UTC()
	host.UpdatedAt = host.CreatedAt

	for _, h := range orm.hosts {
		if h.OsqueryHostID == osQueryHostID {
			host = *h
			break
		}
	}

	if host.ID == 0 {
		host.ID = orm.nextID(host)
	}
	orm.hosts[host.ID] = &host

	return &host, nil
}

func (orm *Datastore) AuthenticateHost(nodeKey string) (*kolide.Host, error) {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	for _, host := range orm.hosts {
		if host.NodeKey == nodeKey {
			return host, nil
		}
	}

	return nil, kolide_errors.ErrNotFound
}

func (orm *Datastore) MarkHostSeen(host *kolide.Host, t time.Time) error {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	host.UpdatedAt = t

	for _, h := range orm.hosts {
		if h.ID == host.ID {
			h.UpdatedAt = t
			break
		}
	}
	return nil
}

func (orm *Datastore) SearchHosts(query string, omit ...uint) ([]*kolide.Host, error) {
	omitLookup := map[uint]bool{}
	for _, o := range omit {
		omitLookup[o] = true
	}

	var results []*kolide.Host

	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	for _, h := range orm.hosts {
		if len(results) == 10 {
			break
		}

		if strings.Contains(h.HostName, query) && !omitLookup[h.ID] {
			results = append(results, h)
			continue
		}

		for _, nic := range h.NetworkInterfaces {
			if strings.Contains(nic.IPAddress, query) && !omitLookup[nic.HostID] {
				results = append(results, h)

				break
			}
		}
		sortutil.AscByField(h.NetworkInterfaces, "ID")
	}

	return results, nil
}

func (orm *Datastore) DistributedQueriesForHost(host *kolide.Host) (map[uint]string, error) {
	// lookup of executions for this host
	hostExecutions := map[uint]kolide.DistributedQueryExecutionStatus{}
	for _, e := range orm.distributedQueryExecutions {
		if e.HostID == host.ID {
			hostExecutions[e.DistributedQueryCampaignID] = e.Status
		}
	}

	// lookup of labels for this host (only including matching labels)
	hostLabels := map[uint]bool{}
	labels, err := orm.ListLabelsForHost(host.ID)
	if err != nil {
		return nil, err
	}
	for _, l := range labels {
		hostLabels[l.ID] = true
	}

	queries := map[uint]string{} // map campaign ID -> query string
	for _, campaign := range orm.distributedQueryCampaigns {
		if campaign.Status != kolide.QueryRunning {
			continue
		}
		for _, target := range orm.distributedQueryCampaignTargets {
			if campaign.ID == target.DistributedQueryCampaignID &&
				((target.Type == kolide.TargetHost && target.TargetID == host.ID) ||
					(target.Type == kolide.TargetLabel && hostLabels[target.TargetID])) &&
				(hostExecutions[campaign.ID] == kolide.ExecutionWaiting) {
				queries[campaign.ID] = orm.queries[campaign.QueryID].Query
			}
		}
	}

	return queries, nil
}
