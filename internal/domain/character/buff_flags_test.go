package character

import (
	"context"
	"testing"
)

func TestBuffFlagsAffectActorState(t *testing.T) {
	actor := NewActor(4, EntityTypePlayer, 1, NewVector3(0, 0, 0), NewVector3(1, 0, 0), "flagged", 1)
	if err := actor.Start(context.Background()); err != nil {
		t.Fatalf("actor start failed: %v", err)
	}

	if actor.GetFlagState().HasFlag(FlagStateStun) {
		t.Fatalf("should not be stunned initially")
	}

	// Add a stun buff
	b := NewBuff(3001, actor, actor, 1.0)
	b.SetFlagAdd(FlagStateStun)
	actor.GetBuffManager().AddBuff(b)

	if !actor.GetFlagState().HasFlag(FlagStateStun) {
		t.Fatalf("actor should be stunned after adding stun buff")
	}

	// Remove the buff
	actor.GetBuffManager().RemoveBuff(b)

	if actor.GetFlagState().HasFlag(FlagStateStun) {
		t.Fatalf("actor should not be stunned after removing stun buff")
	}
}
