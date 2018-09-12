package redis

import (
	"sort"

	"github.com/pkg/errors"
)

var (
	//BestInstance is a reference to the best redis instance available in terms of load balancing.
	BestInstance   *Instance
	updateBalancer = make(chan *Instance)
)

func loadbalance() {
	candidates := make(map[string]*Instance)
	for {
		select {
		case i := <-updateBalancer:
			if i.ConnActive {
				if BestInstance == nil {
					BestInstance = i
				}
				if _, ok := candidates[i.Key]; !ok {
					candidates[i.Key] = i
				}
				if i.AvailableSpace > BestInstance.AvailableSpace {
					BestInstance = i
				}
			} else {
				delete(candidates, i.Key)
				if BestInstance.Key == i.Key {
					var err error
					if BestInstance, err = findNextBest(candidates); err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	}
}

func findNextBest(cmap map[string]*Instance) (*Instance, error) {
	var cslice []*Instance
	for _, v := range cmap {
		cslice = append(cslice, v)
	}
	sort.SliceStable(cslice, func(i, j int) bool { return cslice[i].AvailableSpace > cslice[j].AvailableSpace })
	if len(cslice) > 0 {
		return cslice[0], nil
	}
	return nil, errors.New("redis-master load balancer: could not find the best candidate")
}
