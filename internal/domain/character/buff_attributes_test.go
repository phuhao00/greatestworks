package character

import (
	"context"
	"testing"
)

func TestBuffAttributeModifiersAffectFinalsAndSpeed(t *testing.T) {
	actor := NewActor(3, EntityTypePlayer, 1, NewVector3(0, 0, 0), NewVector3(1, 0, 0), "buffed", 1)
	if err := actor.Start(context.Background()); err != nil {
		t.Fatalf("actor start failed: %v", err)
	}

	am := actor.GetAttributeManager()
	baseMaxHP := am.Final().MaxHP
	baseSpeed := am.Final().Speed

	// Create a buff with +50 MaxHP and +10% MaxHP, and +20% Speed
	b := NewBuff(2001, actor, actor, 5.0)
	b.SetModifier(AttributeModifier{MaxHPAdd: 50, MaxHPMul: 0.1, SpeedMul: 0.2})
	actor.GetBuffManager().AddBuff(b)

	fin := am.Final()
	if fin.MaxHP <= baseMaxHP {
		t.Fatalf("buff should increase MaxHP, got base=%v final=%v", baseMaxHP, fin.MaxHP)
	}

	// Speed is applied to actor each update; tick once to refresh speed from finals
	_ = actor.Update(context.Background(), 0.016)
	if actor.Speed() <= baseSpeed {
		t.Fatalf("buff should increase movement speed, got base=%v current=%v", baseSpeed, actor.Speed())
	}
}
