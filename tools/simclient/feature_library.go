package simclient

// helper returns pointer to bool literal
func boolPtr(v bool) *bool {
	return &v
}

// featureLibrary maps feature identifiers to reusable action sequences.
var featureLibrary = map[string][]ScenarioActionConfig{
	"system.heartbeat": {
		{Name: "system.heartbeat", Message: "system.heartbeat", ExpectResponse: boolPtr(true)},
	},
	"player.login": {
		{Name: "player.login", Message: "player.login", ExpectResponse: boolPtr(true)},
	},
	"player.logout": {
		{Name: "player.logout", Message: "player.logout", ExpectResponse: boolPtr(true)},
	},
	"player.move": {
		{Name: "player.move", Message: "player.move", ExpectResponse: boolPtr(true)},
	},
	"player.info": {
		{Name: "player.info", Message: "player.info", ExpectResponse: boolPtr(true)},
	},
	"player.stats": {
		{Name: "player.stats", Message: "player.stats", ExpectResponse: boolPtr(true)},
	},
	"player.basic": {
		{Name: "player.login", Message: "player.login", ExpectResponse: boolPtr(true)},
		{Name: "player.info", Message: "player.info", ExpectResponse: boolPtr(true)},
		{Name: "player.stats", Message: "player.stats", ExpectResponse: boolPtr(true)},
	},
	"battle.create": {
		{Name: "battle.create", Message: "battle.create", ExpectResponse: boolPtr(true)},
	},
	"battle.join": {
		{Name: "battle.join", Message: "battle.join", ExpectResponse: boolPtr(true)},
	},
	"battle.action": {
		{Name: "battle.action", Message: "battle.action", ExpectResponse: boolPtr(true)},
	},
	"battle.status": {
		{Name: "battle.status", Message: "battle.status", ExpectResponse: boolPtr(true)},
	},
	"battle.basic": {
		{Name: "battle.create", Message: "battle.create", ExpectResponse: boolPtr(true)},
		{Name: "battle.status", Message: "battle.status", ExpectResponse: boolPtr(true)},
	},
	"pet.summon": {
		{Name: "pet.summon", Message: "pet.summon", ExpectResponse: boolPtr(true)},
	},
	"pet.dismiss": {
		{Name: "pet.dismiss", Message: "pet.dismiss", ExpectResponse: boolPtr(true)},
	},
	"pet.info": {
		{Name: "pet.info", Message: "pet.info", ExpectResponse: boolPtr(true)},
	},
	"pet.status": {
		{Name: "pet.status", Message: "pet.status", ExpectResponse: boolPtr(true)},
	},
	"pet.basic": {
		{Name: "pet.info", Message: "pet.info", ExpectResponse: boolPtr(true)},
		{Name: "pet.status", Message: "pet.status", ExpectResponse: boolPtr(true)},
	},
	"building.create": {
		{Name: "building.create", Message: "building.create", ExpectResponse: boolPtr(true)},
	},
	"building.upgrade": {
		{Name: "building.upgrade", Message: "building.upgrade", ExpectResponse: boolPtr(true)},
	},
	"building.status": {
		{Name: "building.status", Message: "building.status", ExpectResponse: boolPtr(true)},
	},
	"building.basic": {
		{Name: "building.status", Message: "building.status", ExpectResponse: boolPtr(true)},
	},
	"social.friend_list": {
		{Name: "social.friend_list", Message: "social.friend_list", ExpectResponse: boolPtr(true)},
	},
	"social.friend_remove": {
		{Name: "social.friend_remove", Message: "social.friend_remove", ExpectResponse: boolPtr(true)},
	},
	"social.chat": {
		{Name: "social.chat", Message: "social.chat", ExpectResponse: boolPtr(true)},
	},
	"social.team_basic": {
		{Name: "social.team_create", Message: "social.team_create", ExpectResponse: boolPtr(true)},
		{Name: "social.team_info", Message: "social.team_info", ExpectResponse: boolPtr(true)},
		{Name: "social.team_join", Message: "social.team_join", ExpectResponse: boolPtr(true)},
		{Name: "social.team_leave", Message: "social.team_leave", ExpectResponse: boolPtr(true)},
	},
	"item.use": {
		{Name: "item.use", Message: "item.use", ExpectResponse: boolPtr(true)},
	},
	"quest.accept": {
		{Name: "quest.accept", Message: "quest.accept", ExpectResponse: boolPtr(true)},
	},
	"quest.progress": {
		{Name: "quest.progress", Message: "quest.progress", ExpectResponse: boolPtr(true)},
	},
	"quest.complete": {
		{Name: "quest.complete", Message: "quest.complete", ExpectResponse: boolPtr(true)},
	},
}
