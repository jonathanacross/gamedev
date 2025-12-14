package assets

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed *
var assetsFS embed.FS

var TerrainTileset = loadImage("terrain.png")
var WallBlobTileset = loadImage("walls_blob.png")
var DungeonObjectsTileset = loadImage("dungeon_objects.png")
var PlayerIdleSpritesImage = loadImage("player_idle.png")
var PlayerWalkSpritesImage = loadImage("player_walk.png")
var PlayerHurtSpritesImage = loadImage("player_hurt.png")
var PlayerDeathSpritesImage = loadImage("player_death.png")
var PlayerAttackSwordSpritesImage = loadImage("player_attack_sword.png")
var PlayerAttackShieldSpritesImage = loadImage("player_attack_shield.png")
var PlayerAttackBowSpritesImage = loadImage("player_attack_bow.png")
var BatSpritesImage = loadImage("bat.png")
var BlobSpritesImage = loadImage("blob.png")
var GoblinSpritesImage = loadImage("goblin.png")
var GhostSpritesImage = loadImage("ghost.png")
var BombSpritesImage = loadImage("bomb.png")
var ArrowSpritesImage = loadImage("arrow.png")
var BombExplosionSpritesImage = loadImage("explosion.png")
var BoomerangSpritesImage = loadImage("boomerang.png")
var StarSpritesImage = loadImage("star_small.png")
var HealthHeartImage = loadImage("heart.png")
var WeaponSelectWindowImage = loadImage("weapon_select_window.png")
var UiIconsImage = loadImage("ui_icons.png")
var TextFaceSource = loadFaceSource("m5x7.ttf")

func loadImage(name string) *ebiten.Image {
	f, err := assetsFS.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

func loadFaceSource(name string) *text.GoTextFaceSource {
	f, err := assetsFS.ReadFile(name)
	if err != nil {
		panic(err)
	}

	face, err := text.NewGoTextFaceSource(bytes.NewReader(f))
	if err != nil {
		panic(err)
	}
	return face
}
