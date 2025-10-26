package character

import (
	"context"
	"testing"
)

func TestSkillFSMTransitions(t *testing.T) {
	actor := NewActor(1, EntityTypePlayer, 1, NewVector3(0, 0, 0), NewVector3(1, 0, 0), "tester", 1)
	if err := actor.Start(context.Background()); err != nil {
		t.Fatalf("actor start failed: %v", err)
	}

	skill := NewSkill(1001, actor)
	skill.SetTimings(0.05, 0.05, 0.05)

	if st := skill.State(); st != SkillStateReady {
		t.Fatalf("expected Ready at start, got %v", st)
	}

	if ok := skill.StartCast(); !ok {
		t.Fatalf("StartCast should succeed in Ready state")
	}

	// Intonate window
	_ = skill.Update(context.Background(), 0.03)
	if st := skill.State(); st != SkillStateIntonate && st != SkillStateActive {
		t.Fatalf("expected Intonate or Active, got %v", st)
	}

	// Enter Active
	_ = skill.Update(context.Background(), 0.03)
	if st := skill.State(); st != SkillStateActive {
		t.Fatalf("expected Active, got %v", st)
	}

	// Enter Cooling
	_ = skill.Update(context.Background(), 0.06)
	if st := skill.State(); st != SkillStateCooling {
		t.Fatalf("expected Cooling, got %v", st)
	}

	// Back to Ready
	_ = skill.Update(context.Background(), 0.06)
	if st := skill.State(); st != SkillStateReady {
		t.Fatalf("expected Ready, got %v", st)
	}
}

func TestAttributesInitAndRegen(t *testing.T) {
	actor := NewActor(2, EntityTypePlayer, 1, NewVector3(0, 0, 0), NewVector3(1, 0, 0), "tester2", 5)
	if err := actor.Start(context.Background()); err != nil {
		t.Fatalf("actor start failed: %v", err)
	}

	fin := actor.GetAttributeManager().Final()
	if fin.MaxHP <= 0 || fin.Speed <= 0 {
		t.Fatalf("final attributes should be initialized with defaults, got MaxHP=%v Speed=%v", fin.MaxHP, fin.Speed)
	}

	if actor.HP() != fin.MaxHP {
		t.Fatalf("actor HP should init to MaxHP, got %v vs %v", actor.HP(), fin.MaxHP)
	}

	// Reduce HP and test regen
	actor.ChangeHP(-10)
	hpAfterHit := actor.HP()
	if hpAfterHit >= fin.MaxHP {
		t.Fatalf("HP should be reduced after damage")
	}

	_ = actor.Update(context.Background(), 2.0) // 2 seconds regen
	if actor.HP() <= hpAfterHit {
		t.Fatalf("HP should regenerate over time; before=%v after=%v", hpAfterHit, actor.HP())
	}
	if actor.HP() > fin.MaxHP {
		t.Fatalf("HP should not exceed MaxHP; got %v > %v", actor.HP(), fin.MaxHP)
	}
}
