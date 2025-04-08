package db

import (
	"context"
	"log"
	"math/rand"

	"github.com/swarajroy/gophersocial/internal/store"
)

var (
	usernames = []string{
		"Skywalker42", "CyberNova", "PixelWarrior", "QuantumKnight", "NebulaHunter",
		"EchoRider", "BlazePhoenix", "GalacticNomad", "StormRider1", "LunarWolf",
		"SolarFlareX", "DarkTitan", "MysticRogue", "IronShade", "ShadowViper",
		"NightSkyZ", "PhantomHawk", "AetherKnight", "ThunderFury", "SilentWolf",
		"CyberKnightX", "AstralGuardian", "CrimsonWraith", "FireDragonZ", "VortexReaper",
		"TechSphinx", "StormBlaze", "QuantumHawk", "GhostShadowX", "PixelDragon",
		"WarpTornado", "EchoDragon", "NebulaFury", "MysticWarden", "ShadowScythe",
		"ThunderGhost", "SolarWarden", "StarRider", "HyperNovaZ", "SilentDragon",
		"IronPhantom", "PixelTitan", "NebulaShade", "CrimsonKnightZ", "AetherReaper",
		"ShadowPhoenix", "VortexDragon", "NightStryker", "GalacticPhantom", "LunarWarden",
		"TitanWraith", "CyberVortex", "PlasmaFury", "IronRider", "EclipseBlaze",
	}

	titles = []string{
		"Whispers of the Forgotten", "The Last Horizon", "Echoes of Eternity", "Beyond the Stars", "The Silent Storm",
		"Shadows in the Mist", "The Lost Kingdom", "Tides of Fate", "The Flame of Hope", "Rise of the Phoenix",
		"Darkness Descends", "A Dance with Time", "The Final Hour", "Wings of Destiny", "Echoes of the Past",
		"Bound by Fate", "Through the Abyss", "Legacy of the Fallen", "The Shattered World", "Into the Unknown",
		"Veil of Shadows", "The Forgotten Realm", "The Hidden Empire", "Song of the Siren", "Embers of War",
		"Veins of the Earth", "Whispers in the Dark", "Tales from the Edge", "The Midnight Queen", "Beneath the Moonlight",
		"Lost in the Void", "Heart of the Storm", "Whispers of the Ancients", "Chronicles of the Lost", "Fury of the Ocean",
		"Rise of the Undying", "The Crimson Sea", "Through the Eyes of the Beast", "Guardians of the Night", "The Silver Key",
		"Secrets of the Abyss", "The Cursed Blade", "Echo of the Wolf", "Shattered Horizons", "The Mystic Path",
		"Into the Fire", "The Storm's Embrace", "Whispers of the Gods", "The Dawn of Darkness", "The Heart of the Dragon",
		"Fates Collide", "Song of the Phoenix", "Beyond the Edge", "Veiled in Time", "The Eternal War",
		"Journey to the Depths", "Rise of the Shadows", "The Fallen Star", "The Forgotten Hero", "In the Shadow of Giants",
		"Echoes of Destiny", "The Hidden Truth", "Tears of the Earth", "The Mark of the Raven", "Lords of the Underworld",
		"Blood and Steel", "The Clockwork Heart", "Whispers of the Wind", "The Shattered Sword", "Eyes of the Predator",
		"Guardians of the Fallen", "The Last Beacon", "Echoes of the Moon", "Whispers of Silence", "The Crimson Path",
		"Moonlit Reflections", "Embrace of the Beast", "The Siren's Song", "The Edge of Forever", "Shadow of the Phoenix",
		"The Forgotten Queen", "The Eternal Flame", "Into the Abyss", "Silent Footsteps", "Legacy of the Dragon",
		"Song of the Fallen", "Blood of the Earth", "Voices in the Dark", "Crown of Ice", "Tales of the Forgotten King",
		"Whispers from the Past", "The Path to Immortality", "Shadows of the Warlord", "The Silent Kingdom", "The Road to Destiny",
		"The Oracle's Curse", "Betrayal of the Crown", "Cloak of Shadows", "The Lost Citadel", "Echoes of the Abyss",
		"The Darkened Sun", "Journey through the Stars", "Guardians of the Realm", "Cursed by the Sea", "The Dragon's Breath",
		"Whispers of the Stars", "The Risen Empire", "Lost in the Shadows", "Rise of the Underdog", "Echoes of the Fallen",
		"Light in the Darkness", "The Shattered Kingdom", "Into the Wild", "Wings of the Fallen", "Secrets of the Lost Realm",
		"The Fire Within", "The Raven's Flight", "Crown of Blood", "Beyond the Veil", "The Hidden Path",
	}

	contents = []string{
		"Whispers of the Lost", "Echoes in the Void", "The Last Journey", "Beyond the Horizon", "Tales from the Dark",
		"Mystic Realms", "Shadows of the Past", "Rise of the Titan", "Heart of the Dragon", "The Silent Knight",
		"Gates of Fate", "Blood Moon Rising", "Tales of the Unknown", "The Eternal Night", "Wings of Fate",
		"Crimson Sky", "Warrior's Path", "Into the Storm", "The Hidden Realm", "The Secret of the Deep",
		"Chasing the Stars", "Guardians of the Lost", "The Silent Watcher", "Legends of the Forgotten", "The Enchanted Forest",
		"Echoes of Destiny", "Through the Flames", "The Path of Shadows", "The Shattered Crown", "The Lurking Beast",
		"Fury of the Ocean", "Journey of the Heart", "The Alchemist's Secret", "The Forbidden Kingdom", "Whispers in the Wind",
		"The Warlord's Blade", "The Raven's Call", "A Dance with Time", "The Ghost of the Mountain", "The Dark Empire",
		"Dawn of the Titans", "Secrets of the Abyss", "Rise of the Phoenix", "The Flame of Rebirth", "The Ocean's Fury",
		"Beneath the Surface", "The Last Dawn", "Whispers of the Gods", "The Final Hour", "Echoes of War",
		"Fate of the Fallen", "Wings of the Phoenix", "Crimson Tide", "The Fallen Star", "Embrace of the Void",
		"The Timekeeper's Curse", "The Oracle's Vision", "A Path to Glory", "The Hidden Fortress", "A Warrior's Oath",
		"The Darkened Sky", "The Last Sorcerer", "Tales of the Haunted", "Embrace the Darkness", "Blood of the Titans",
		"The Iron Fist", "Whispers in the Darkness", "The Crystal Blade", "The Cursed Artifact", "Echoes of the Moon",
		"Journey to the Depths", "Secrets of the Forgotten", "The King's Sacrifice", "Shadows of the Warlord", "Lords of the Void",
		"The Wanderer's Tale", "The Rise of Night", "Sword of the Fallen", "Heart of the Storm", "The Silver Moon",
		"The Last Stand", "Guardians of the Abyss", "The Eternal Flame", "The Moonlit Path", "The Cursed Knight",
		"Betrayal of the Gods", "Chronicles of the Lost", "The Siren's Song", "Tales of the Wraith", "Through the Gates of Time",
		"The Shattered Shield", "A Journey Beyond", "Whispers of the Siren", "The Demon King's Curse", "The Unbroken Path",
		"Echoes of the Lost", "The King's Wrath", "The Edge of Despair", "The Forgotten Hero", "Rise of the Underworld",
		"Secrets in the Dark", "The Phantom's Curse", "Crown of Ice", "The Flame of Despair", "The Heart of the Beast",
		"Legends of the Sea", "The Dark Sorceress", "The Last Kingdom", "Whispers of the Fallen", "Tears of the Earth",
		"Echoes in the Night", "Through the Shadows", "The Eternal Journey", "The Path to Immortality", "Guardian of the Flame",
		"The Silver Crescent", "Blood and Steel", "Tales of the Wild", "The Warrior's Code", "The Blackened Heart",
		"Chronicles of the Cursed", "Into the Abyss", "The Dragon's Breath", "Echoes of the Abyss", "The Wraith's Call",
		"A Song of Shadows", "The Throne of Darkness", "Fury of the Fallen", "The Last Hope", "Beyond the Veil",
		"Guardians of the Light", "The Shattered World", "The Lost Knight", "The Ghost of the Forest", "Shadows of Eternity",
		"Song of the Lost", "Secrets of the Storm", "Whispers of the Moon", "The Warlock's Curse", "Echoes of the Unknown",
		"The Secret of the Stars", "Path of the Warrior", "The Sorcerer's Revenge", "Legacy of the Fallen", "Into the Fire",
		"Rise of the Sorcerer", "Heart of the Dragon", "The Silent Watcher", "Fury of the Phoenix", "The Spirit of Vengeance",
		"The Nightfall Chronicles", "Whispers of the Abyss", "The King’s Gambit", "Into the Depths", "The Forsaken Land",
		"The Shattered Empire", "Riders of the Storm", "The Golden Flame", "Secrets of the Dead", "Heart of the Inferno",
		"The Last Prince", "Echoes of the Fallen King", "Rise of the Shadows", "The Sorcerer's Blood", "Through the Veil of Time",
		"The Oracle's Prophecy", "Wings of the Raven", "Blood of the Fallen", "The Darkened World", "Tales from the Abyss",
		"The Curse of the Eternal", "Fury of the Beast", "The Dark Crusade", "The Last Dragon", "Legacy of the Titan",
		"The Haunted Woods", "Whispers of the Phoenix", "Echoes in the Dark", "The Silent Soldier", "Journey into the Abyss",
		"The Blood Moon", "Echoes of Eternity", "The King's Return", "The Lost Wanderer", "A Dance with Shadows",
		"The Path of Destiny", "The Lost Chronicles", "Fate's Embrace", "The Unseen World", "The Vengeful Queen",
		"The Warden's Curse", "The Secret of the Crow", "Whispers of the Serpent", "The Dawn of War", "Shadows of the Phoenix",
		"The Night's Watch", "The Last Kingdom", "Echoes of the Hunt", "Blood of the Warrior", "Rise of the Fallen",
		"The Dark King’s Legacy", "The Phoenix’s Curse", "Tales of the Ancient", "The Moonlit Journey", "Heart of the Fallen",
		"The Frozen Throne", "Shadows of the Sorcerer", "Rise of the Warlord", "Echoes of the Blood Moon", "A Path to Glory",
		"The Legend of the Beast", "Wings of the Lost", "The Forgotten Quest", "Whispers of the Earth", "The Heart of the Storm",
		"The Queen's Curse", "The Eternal Embrace", "Through the Veil", "Legacy of the Dragon", "The Wild Hunt",
		"Whispers in the Blood", "The Eternal Beast", "The Road to Eternity", "Chronicles of the Void", "The Beast's Fury",
		"Shadows of the Hunter", "The Lost Realm", "The Moon's Embrace", "Fate's Call", "The Dark Sorcerer’s Return",
		"The Warrior's Fate", "Echoes in the Void", "The Silent Assassin", "The Eye of the Storm", "Wings of Despair",
		"Legacy of the Moon", "The Unraveling Curse", "The Raven’s Flight", "The Cursed Path", "The Fire Within",
		"The Fallen Blade", "Through the Dark", "The Kingdom of Shadows", "The Last Survivor", "Echoes from the Void",
		"Guardians of the Lost Realm", "The Silent Reaper", "Whispers of the Bloodline", "Echoes of the Warlord", "The Risen King",
		"The Dragon's Fury", "A Journey Beyond", "Secrets of the Cursed Land", "Blood of the Fallen", "The End of the Line",
		"The King's Sacrifice", "Echoes of the Lost Kingdom", "Beyond the Veil of Time", "The Light in the Darkness",
		"Chronicles of the Hidden Realm", "The Path of Shadows", "Into the Depths of Fate", "Whispers of the Fallen King",
		"The Dark Knight's Return", "Legends of the Eternal", "The King's Blood", "Tales of the Forgotten Realm", "The Dragon’s Wrath",
		"The Silent Flame", "Whispers of the Warrior", "Rise of the Phantom", "The Path to Glory", "Secrets of the Cursed Blade",
		"The Shattered Moon", "Chronicles of the Wraith", "The Final Hour", "The Witch’s Curse", "Heart of the Lost",
		"Into the Abyss", "The Silver Crescent", "Echoes in the Forest", "The Silent Sorcerer", "The Forgotten City",
		"Whispers of the Eternal Night", "The Hidden Star", "A Path Through Time", "Shadows of the Fallen", "The Warrior's Return",
		"Echoes of the Eternal Flame", "The King's Curse", "The Fallen Warrior", "The Legacy of the Phoenix", "Echoes of the Wind",
		"The Shattered Fate", "The Dragon’s Curse", "Guardians of the Dawn", "Rise of the Sorceress", "The Last Reign",
		"Whispers from the Depths", "The Silent Truth", "The Dark Road", "Echoes of the Forsaken", "Rise of the Shadow King",
	}

	tags = []string{
		"Adventure", "Fantasy", "Epic", "Mystery", "SciFi", "Action", "Thriller", "Drama", "Romance", "Horror",
		"Fiction", "Suspense", "Magic", "Historical", "Supernatural", "UrbanFantasy", "YoungAdult", "Comedy", "Family", "Superhero",
		"Dystopian", "Zombie", "Steampunk", "Cyberpunk", "HistoricalFiction", "AdventureFantasy", "DarkFantasy", "Mythology", "Witches", "Aliens",
		"Vampires", "Werewolves", "War", "Medieval", "PostApocalyptic", "Apocalypse", "MurderMystery", "TimeTravel", "ParallelUniverse", "AlternateHistory",
		"ThrillerFiction", "ActionAdventure", "SliceOfLife", "EpicFantasy", "TalesOfTheUnknown", "UrbanMystery", "CosmicHorror", "PsychologicalThriller", "SpaceOpera",
		"TeenFiction", "MagicalRealism", "FairyTales", "ClassicLiterature", "Futuristic", "ComingOfAge",
	}
)

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)
	for i := 0; i < num; i++ {
		username := usernames[i]
		users[i] = &store.User{
			Username: username,
			Email:    username + "" + "@example.com",
			//Password: "123123",
		}
	}
	return users
}

func getRandInt(lower, upper int) int {
	return rand.Intn(upper-lower) + lower
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		tagsLen := getRandInt(1, 5)
		tagsSlice := make([]string, tagsLen)
		for i := 0; i < tagsLen; i++ {
			tagsSlice[i] = tags[rand.Intn(len(tags))]
		}
		posts[i] = &store.Post{
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags:    tagsSlice,
			UserID:  user.ID,
		}
	}
	return posts
}

func generateComments(num int, posts []*store.Post, users []*store.User) []*store.Comment {
	comments := make([]*store.Comment, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		post := posts[rand.Intn(len(posts))]
		comments[i] = &store.Comment{
			PostID:  post.ID,
			UserID:  user.ID,
			Content: contents[rand.Intn(len(contents))],
		}
	}
	return comments
}

func Seed(store store.Storage) {

	ctx := context.Background()
	users := generateUsers(50)

	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Fatalf("error occurred while seeding users data! %s", err)
			return
		}
	}

	posts := generatePosts(500, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Fatalf("error occurred while seeding posts data! %s", err)
			return
		}
	}

	comments := generateComments(300, posts, users)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Fatalf("error occurred while seeding posts data! %s", err)
			return
		}
	}

	log.Println("Finished seeding the database!")

}
