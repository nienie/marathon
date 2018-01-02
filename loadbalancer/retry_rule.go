package loadbalancer

import (
	"time"

	"github.com/nienie/marathon/server"
)

const (
	//DefaultMaxRetryTimeout ...
	DefaultMaxRetryTimeout = 500 * time.Millisecond
)

//RetryRule can be cascaded, this allows adding a retry logic to an existing Rule.
type RetryRule struct {
	BaseRule
	SubRule         Rule
	maxRetryTimeout time.Duration
}

//NewRetryRule ...
func NewRetryRule(subRule Rule, maxRetryTimeout time.Duration) Rule {
	rule := &RetryRule{}
	if subRule != nil {
		rule.SubRule = subRule
	} else {
		rule.SubRule = NewRoundRobinRule()
	}
	rule.SetMaxRetryMills(maxRetryTimeout)
	return rule
}

//SetLoadBalancer ...
func (o *RetryRule) SetLoadBalancer(lb LoadBalancer) {
	o.BaseRule.SetLoadBalancer(lb)
	o.SubRule.SetLoadBalancer(lb)
}

//SetMaxRetryMills ...
func (o *RetryRule) SetMaxRetryMills(maxRetryTime time.Duration) {
	if maxRetryTime > 0 {
		o.maxRetryTimeout = maxRetryTime
	} else {
		o.maxRetryTimeout = DefaultMaxRetryTimeout
	}
}

//GetMaxRetryTimeout ...
func (o *RetryRule) GetMaxRetryTimeout() time.Duration {
	return o.maxRetryTimeout
}

//Choose ...
func (o *RetryRule) Choose(key interface{}) *server.Server {
	return o.ChooseFromLoadBalancer(o.GetLoadBalancer(), key)
}

//ChooseFromLoadBalancer ...
func (o *RetryRule) ChooseFromLoadBalancer(lb LoadBalancer, key interface{}) *server.Server {
	deadline := time.Duration(time.Now().UnixNano()) + o.maxRetryTimeout

	answer := o.SubRule.Choose(key)

	remainTime := deadline - time.Duration(time.Now().UnixNano())
	if (answer == nil || !answer.IsAlive()) && remainTime > 0 {
		timer := time.NewTimer(remainTime)
		for {
			select {
			case <-timer.C:
				break
			default:
				answer = o.SubRule.Choose(key)
				if answer != nil && answer.IsAlive() {
					break
				}
			}
		}
	}

	if answer == nil || !answer.IsAlive() {
		return nil
	}

	return answer
}
