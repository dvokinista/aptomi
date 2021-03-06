package helm

import (
	"fmt"
	"github.com/Aptomi/aptomi/pkg/engine/progress"
	"github.com/Aptomi/aptomi/pkg/engine/resolve"
	"github.com/Aptomi/aptomi/pkg/event"
	"github.com/Aptomi/aptomi/pkg/external"
	"github.com/Aptomi/aptomi/pkg/lang"
)

// Process is a action which gets called only once. It manages all Istio rules across all clusters, making sure they
// are up to date by creating/deleting/updating rules if/as needed
// TODO: reduce cyclomatic complexity
func (plugin *Plugin) Process(policy *lang.Policy, resolution *resolve.PolicyResolution, externalData *external.Data, eventLog *event.Log) error { // nolint: gocyclo
	// Do not run Istio
	if true {
		return nil
	}

	// todo(slukjanov): do something with progress
	prog := progress.NewNoop()

	if len(resolution.GetComponentProcessingOrder()) == 0 {
		return nil
	}

	eventLog.WithFields(
		event.Fields{},
	).Info("Figuring out which Istio rules have to be added/deleted")

	existingRules := make([]*istioRouteRule, 0)

	for _, clusterObj := range policy.GetObjectsByKind(lang.ClusterObject.Kind) {
		cluster := clusterObj.(*lang.Cluster)
		cache, err := plugin.getClusterCache(cluster, eventLog)
		if err != nil {
			return err
		}
		rules, err := cache.getExistingIstioRouteRulesForCluster()
		if err != nil {
			return err
		}
		existingRules = append(existingRules, rules...)
		prog.Advance()
	}

	// Process in the right order
	desiredRules := make(map[string][]*istioRouteRule)
	for _, key := range resolution.GetComponentProcessingOrder() {
		rules, err := plugin.getDesiredIstioRouteRulesForComponent(key, policy, resolution, externalData, eventLog)
		if err != nil {
			return fmt.Errorf("error while processing Istio Ingress for component '%s': %s", key, err)
		}
		desiredRules[key] = rules
		prog.Advance()
	}

	deleteRules := make([]*istioRouteRule, 0)
	createRules := make(map[string][]*istioRouteRule)

	// populate createRules, to make sure we will get correct number of entries for progress indicator
	for _, key := range resolution.GetComponentProcessingOrder() {
		createRules[key] = make([]*istioRouteRule, 0)
	}

	for _, existingRule := range existingRules {
		found := false
		for _, desiredRulesForComponent := range desiredRules {
			for _, desiredRule := range desiredRulesForComponent {
				if existingRule.Service == desiredRule.Service && existingRule.Cluster.Name == desiredRule.Cluster.Name {
					found = true
				}
			}
		}
		if !found {
			deleteRules = append(deleteRules, existingRule)
		}
	}

	for key, desiredRulesForComponent := range desiredRules {
		for _, desiredRule := range desiredRulesForComponent {
			found := false
			for _, existingRule := range existingRules {
				if desiredRule.Service == existingRule.Service && desiredRule.Cluster.Name == existingRule.Cluster.Name {
					found = true
				}
			}
			if !found {
				createRules[key] = append(createRules[key], desiredRule)
			}
		}
	}

	// process creations by component
	changed := false
	for _, createRulesForComponent := range createRules {
		for _, rule := range createRulesForComponent {
			eventLog.WithFields(event.Fields{}).Infof("Creating Istio rule: %s (%s)", rule.Service, rule.Cluster.Name)
			err := rule.create()
			if err != nil {
				return err
			}
			changed = true
		}
		prog.Advance()
	}

	// process deletions all at once
	for _, rule := range deleteRules {
		eventLog.WithFields(event.Fields{}).Infof("Deleting Istio rule: %s (%s)", rule.Service, rule.Cluster.Name)
		err := rule.destroy()
		if err != nil {
			return err
		}
		changed = true
	}
	prog.Advance()

	if changed {
		eventLog.WithFields(event.Fields{}).Infof("Successfully processed Istio rules")
	} else {
		eventLog.WithFields(event.Fields{}).Infof("No changes in Istio rules")
	}

	return nil
}
