package droplets

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/digitalocean/godo"
	"golang.org/x/net/context"
)

// waitForActions loops through each actions in godo links and wait until they finish
func waitForActions(ctx context.Context, cloud *godo.Client, links *godo.Links) error {
	if links == nil {
		return nil
	}
	if len(links.Actions) == 0 {
		return nil
	}

	for _, actionLink := range links.Actions {
		action, _, err := actionLink.Get(cloud)
		if err != nil {
			return err
		}
		if err := waitForAction(ctx, cloud, action); err != nil {
			return err
		}
	}
	return nil
}

// waitForAction waits for a single action to finish.
func waitForAction(ctx context.Context, cloud *godo.Client, action *godo.Action) error {
	if action == nil {
		return nil
	}

	base := (4 * time.Second).Seconds()
	cap := (30 * time.Second).Seconds()
	factor := 1.5
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for attempt := 0.0; ; attempt += 1.0 {

		var err error
		action, _, err = cloud.Actions.Get(action.ID)
		if err != nil {
			return err
		}
		if action.Status == "errored" {
			return errors.New(action.String())
		}
		if action.CompletedAt != nil || action.Status == "done" {
			return nil
		}
		sleepSeconds := r.Float64() * math.Min(cap, base*math.Pow(factor, attempt))
		sleep := time.Duration(sleepSeconds * float64(time.Second))
		select {
		case <-ctx.Done():
			return fmt.Errorf("timedout waiting for action %d to complete", action.ID)
		case <-time.After(sleep):
		}
	}
}
