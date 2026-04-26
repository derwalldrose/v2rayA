package v2ray

import (
	"testing"
	"time"

	"github.com/v2rayA/v2rayA/core/v2ray/where"
	"github.com/v2rayA/v2rayA/db/configure"
)

func TestAddLeastPingObservatoryForXrayUsesSingleObservatory(t *testing.T) {
	tmpl := &Template{Variant: where.Xray}
	tmpl.addLeastPingObservatory("group-a", []string{"node-a"}, "https://gstatic.com/generate_204", 10*time.Second)
	tmpl.addLeastPingObservatory("group-b", []string{"node-b"}, "https://gstatic.com/generate_204", 10*time.Second)

	if tmpl.Observatory == nil {
		t.Fatal("expected xray observatory to be created")
	}
	if tmpl.MultiObservatory != nil {
		t.Fatal("did not expect multiObservatory for xray")
	}
	selectors := tmpl.Observatory.Settings.SubjectSelector
	if len(selectors) != 2 || selectors[0] != "node-a" || selectors[1] != "node-b" {
		t.Fatalf("unexpected xray selectors: %#v", selectors)
	}

	services := tmpl.attachObservatoryService([]string{"LoggerService"}, 0)
	if len(services) != 1 || services[0] != "LoggerService" {
		t.Fatalf("xray should not attach ObservatoryService, got %#v", services)
	}
}

func TestAddLeastPingObservatoryForV2rayUsesMultiObservatory(t *testing.T) {
	tmpl := &Template{Variant: where.V2ray}
	tmpl.addLeastPingObservatory("group-a", []string{"node-a"}, "https://gstatic.com/generate_204", 10*time.Second)

	if tmpl.MultiObservatory == nil || len(tmpl.MultiObservatory.Observers) != 1 {
		t.Fatalf("expected one v2ray observer, got %#v", tmpl.MultiObservatory)
	}
	if tmpl.Observatory != nil {
		t.Fatal("did not expect single observatory for v2ray")
	}
}

func TestAddLeastLoadBurstObservatoryForXray(t *testing.T) {
	tmpl := &Template{Variant: where.Xray}
	tmpl.addLeastLoadBurstObservatory([]string{"node-a"}, "https://gstatic.com/generate_204", 15*time.Second)
	tmpl.addLeastLoadBurstObservatory([]string{"node-b"}, "https://gstatic.com/generate_204", 15*time.Second)

	if tmpl.BurstObservatory == nil {
		t.Fatal("expected xray burstObservatory to be created")
	}
	if tmpl.Observatory != nil {
		t.Fatal("did not expect regular observatory for xray leastload helper")
	}
	selectors := tmpl.BurstObservatory.SubjectSelector
	if len(selectors) != 2 || selectors[0] != "node-a" || selectors[1] != "node-b" {
		t.Fatalf("unexpected xray burst selectors: %#v", selectors)
	}
	if tmpl.BurstObservatory.PingConfig.Destination != "https://gstatic.com/generate_204" {
		t.Fatalf("unexpected burst destination: %#v", tmpl.BurstObservatory.PingConfig)
	}
	if tmpl.BurstObservatory.PingConfig.Interval != (15 * time.Second).String() {
		t.Fatalf("unexpected burst interval: %#v", tmpl.BurstObservatory.PingConfig)
	}
}

func TestRandomAndRoundRobinForXrayReuseObservatory(t *testing.T) {
	tmpl := &Template{Variant: where.Xray}
	tmpl.addLeastPingObservatory("group-a", []string{"node-a"}, "https://gstatic.com/generate_204", 10*time.Second)
	tmpl.addLeastPingObservatory("group-b", []string{"node-b"}, "https://gstatic.com/generate_204", 10*time.Second)
	if tmpl.Observatory == nil {
		t.Fatal("expected observatory for xray random/roundrobin health filtering")
	}
	if tmpl.BurstObservatory != nil {
		t.Fatal("did not expect burstObservatory for random/roundrobin")
	}
}

func TestCanonicalBalancerStrategyForXray(t *testing.T) {
	cases := []struct {
		strategy configure.ObservatoryType
		want     string
	}{
		{configure.Random, "random"},
		{configure.RoundRobin, "roundRobin"},
		{configure.LeastPing, "leastPing"},
		{configure.LeastLoad, "leastLoad"},
	}
	for _, tc := range cases {
		got, useObserverTag, ok := canonicalBalancerStrategy(where.Xray, tc.strategy)
		if !ok {
			t.Fatalf("expected strategy %q to be supported for xray", tc.strategy)
		}
		if useObserverTag {
			t.Fatalf("xray strategy %q should not use observerTag", tc.strategy)
		}
		if got != tc.want {
			t.Fatalf("strategy %q mapped to %q, want %q", tc.strategy, got, tc.want)
		}
	}
}

func TestCanonicalBalancerStrategyForV2ray(t *testing.T) {
	got, useObserverTag, ok := canonicalBalancerStrategy(where.V2ray, configure.LeastPing)
	if !ok || got != "leastPing" || !useObserverTag {
		t.Fatalf("unexpected v2ray leastPing mapping: got=%q observerTag=%v ok=%v", got, useObserverTag, ok)
	}
	got, useObserverTag, ok = canonicalBalancerStrategy(where.V2ray, configure.Random)
	if !ok || got != "random" || useObserverTag {
		t.Fatalf("unexpected v2ray random mapping: got=%q observerTag=%v ok=%v", got, useObserverTag, ok)
	}
}
