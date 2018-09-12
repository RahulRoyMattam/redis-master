package redis

import (
	"context"

	"github.com/pkg/errors"
)

//doAll : Call the Do command on all the valid redis connections
//if all the requests to redis errored out, then lastError will have the error of the last request.
//if atleast one request to redis did not error out, lastError will be nil.
func doAll(ctx context.Context, commandName string, args ...interface{}) (responses []*response, lastError error) {
	resChan := make(chan *response, len(Instances))
	requests := 0
	lastError = errors.New("doAll:default error")
	for _, inst := range Instances {
		//make instance consistently accessible in the closure
		inst := inst
		if inst.ConnActive {
			requests++
			go func() {
				reply, err := inst.doSafe(ctx, commandName, args...)
				resChan <- &response{
					instanceKey: inst.Key,
					reply:       reply,
					err:         err,
				}
			}()
		}
	}

	for {
		select {
		case res := <-resChan:
			responses = append(responses, res)

			//setting lastError
			if res.err != nil && lastError != nil {
				lastError = res.err
			}
			if res.err == nil {
				lastError = nil
			}

			if len(responses) == requests {
				return responses, lastError
			}
		case <-ctx.Done():
			return responses, errors.New("redis-master : the forward request to all redis instances timed out")
		}
	}
}

//doBest : Call the Do command on the best valid redis connections
func doBest(ctx context.Context, commandName string, args ...interface{}) (ret interface{}, err error) {
	if BestInstance != nil && BestInstance.ConnActive {
		ret, err := BestInstance.doSafe(ctx, commandName, args...)
		return ret, err
	}
	return ret, errors.New("redis-master: could not find the best redis instance")
}

type response struct {
	instanceKey string
	reply       interface{}
	err         error
}
