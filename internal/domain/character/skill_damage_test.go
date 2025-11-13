package character

import (
	"context"
	"testing"
)

func TestSkillDamageReducesTargetHP(t *testing.T) {
	attacker := NewActor(10, EntityTypePlayer, 1, NewVector3(0, 0, 0), NewVector3(1, 0, 0), "att", 1)
	defender := NewActor(11, EntityTypePlayer, 1, NewVector3(1, 0, 0), NewVector3(-1, 0, 0), "def", 1)

	if err := attacker.Start(context.Background()); err != nil {
		t.Fatalf("attacker start: %v", err)
	}
	if err := defender.Start(context.Background()); err != nil {
		t.Fatalf("defender start: %v", err)
	}

	// Setup a simple physical skill
	sk := NewSkill(5001, attacker)
	sk.SetTimings(0.0, 0.01, 0.01) // instant, short active and cooldown
	sk.SetDamage(20, 1.0, 0.0, DamageTypePhysical)
	attacker.GetSkillManager().AddSkill(sk)

	// Record defender HP
	baseHP := defender.HP()
	if ok := attacker.GetSpell().Cast(sk.ID(), defender); !ok {
		t.Fatalf("failed to cast skill")
	}

	// Advance time to process Active window and apply damage
	_ = attacker.Update(context.Background(), 0.02)

	if defender.HP() >= baseHP {
		t.Fatalf("expected defender HP to reduce; before=%v after=%v", baseHP, defender.HP())
	}
}
