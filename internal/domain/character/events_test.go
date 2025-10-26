package character

import (
	"context"
	"sync"
	"testing"
)

type fakePublisher struct {
	mu     sync.Mutex
	events []DomainEvent
}

func (f *fakePublisher) Publish(e DomainEvent) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.events = append(f.events, e)
}
func (f *fakePublisher) CountByName(name string) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	c := 0
	for _, e := range f.events {
		if e.EventName() == name {
			c++
		}
	}
	return c
}

func TestDamageDealtEventIsPublished(t *testing.T) {
	attacker := NewActor(21, EntityTypePlayer, 1, NewVector3(0, 0, 0), NewVector3(1, 0, 0), "att", 1)
	defender := NewActor(22, EntityTypePlayer, 1, NewVector3(1, 0, 0), NewVector3(-1, 0, 0), "def", 1)
	if err := attacker.Start(context.Background()); err != nil {
		t.Fatalf("att start: %v", err)
	}
	if err := defender.Start(context.Background()); err != nil {
		t.Fatalf("def start: %v", err)
	}

	pub := &fakePublisher{}
	// Attach publisher to defender (OnHurt publishes)
	defender.SetEventPublisher(pub)

	sk := NewSkill(6001, attacker)
	sk.SetTimings(0.0, 0.01, 0.01)
	sk.SetDamage(10, 1.0, 0.0, DamageTypePhysical)
	attacker.GetSkillManager().AddSkill(sk)

	if ok := attacker.GetSpell().Cast(sk.ID(), defender); !ok {
		t.Fatalf("cast failed")
	}
	_ = attacker.Update(context.Background(), 0.02)

	if pub.CountByName("DamageDealt") == 0 {
		t.Fatalf("expected DamageDealt event to be published")
	}
}
