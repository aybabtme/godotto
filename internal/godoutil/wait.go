package godoutil

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/digitalocean/godo"
)

// WaitForActions loops through each actions in godo links and wait until they finish
func WaitForActions(ctx context.Context, cloud *godo.Client, links *godo.Links) error {
	if links == nil {
		return nil
	}
	if len(links.Actions) == 0 {
		return nil
	}

	for _, actionLink := range links.Actions {
		action, _, err := actionLink.Get(ctx, cloud)
		if err != nil {
			return err
		}
		if err := WaitForAction(ctx, cloud, action); err != nil {
			return err
		}
	}
	return nil
}

// WaitForAction waits for a single action to finish.
func WaitForAction(ctx context.Context, cloud *godo.Client, action *godo.Action) error {
	if action == nil {
		return nil
	}

	base := (4 * time.Second).Seconds()
	cap := (30 * time.Second).Seconds()
	factor := 1.5
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for attempt := 0.0; ; attempt += 1.0 {

		var err error
		action, _, err = cloud.Actions.Get(ctx, action.ID)
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
