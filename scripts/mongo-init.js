// MongoDB 初始化脚本
// 创建游戏数据库和基础集合

// 切换到游戏数据库
db = db.getSiblingDB('mmo_game');

// 创建用户集合
db.createCollection('players');
db.createCollection('battles');
db.createCollection('items');
db.createCollection('buildings');
db.createCollection('pets');
db.createCollection('npcs');
db.createCollection('quests');
db.createCollection('rankings');

// 创建索引
db.players.createIndex({ "player_id": 1 }, { unique: true });
db.players.createIndex({ "username": 1 }, { unique: true });
db.players.createIndex({ "email": 1 }, { unique: true });
db.players.createIndex({ "level": -1 });
db.players.createIndex({ "experience": -1 });
db.players.createIndex({ "created_at": 1 });

db.battles.createIndex({ "battle_id": 1 }, { unique: true });
db.battles.createIndex({ "participants": 1 });
db.battles.createIndex({ "status": 1 });
db.battles.createIndex({ "created_at": 1 });

db.items.createIndex({ "item_id": 1 }, { unique: true });
db.items.createIndex({ "type": 1 });
db.items.createIndex({ "rarity": 1 });

db.buildings.createIndex({ "building_id": 1 }, { unique: true });
db.buildings.createIndex({ "owner_id": 1 });
db.buildings.createIndex({ "type": 1 });

db.pets.createIndex({ "pet_id": 1 }, { unique: true });
db.pets.createIndex({ "owner_id": 1 });
db.pets.createIndex({ "type": 1 });

db.npcs.createIndex({ "npc_id": 1 }, { unique: true });
db.npcs.createIndex({ "scene_id": 1 });
db.npcs.createIndex({ "type": 1 });

db.quests.createIndex({ "quest_id": 1 }, { unique: true });
db.quests.createIndex({ "player_id": 1 });
db.quests.createIndex({ "status": 1 });

db.rankings.createIndex({ "type": 1, "rank": 1 });
db.rankings.createIndex({ "player_id": 1 });

// 插入一些基础数据
db.players.insertOne({
    player_id: "admin",
    username: "admin",
    email: "admin@example.com",
    level: 1,
    experience: 0,
    gold: 1000,
    created_at: new Date(),
    updated_at: new Date()
});

print("MongoDB 初始化完成");


