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

var TerrainTileset = loadImage("level/terrain.png")
var WallBlobTileset = loadImage("level/walls_blob.png")

var PlayerIdleSpritesImage = loadImage("player/idle.png")
var PlayerWalkSpritesImage = loadImage("player/walk.png")
var PlayerHurtSpritesImage = loadImage("player/hurt.png")
var PlayerDeathSpritesImage = loadImage("player/death.png")
var PlayerAttackSwordSpritesImage = loadImage("player/attack_sword.png")
var PlayerAttackShieldSpritesImage = loadImage("player/attack_shield.png")
var PlayerAttackBowSpritesImage = loadImage("player/attack_bow.png")
var PlayerAttackWandSpritesImage = loadImage("player/attack_wand.png")

var BatSpritesImage = loadImage("enemies/bat.png")
var SpikeTurleSpritesImage = loadImage("enemies/spike_turtle.png")
var BlobSpritesImage = loadImage("enemies/blob.png")
var GoblinSpritesImage = loadImage("enemies/goblin.png")
var GhostSpritesImage = loadImage("enemies/ghost.png")

var BombSpritesImage = loadImage("objects/bomb.png")
var ArrowSpritesImage = loadImage("objects/arrow.png")
var BombExplosionSpritesImage = loadImage("objects/explosion.png")
var BoomerangSpritesImage = loadImage("objects/boomerang.png")
var StarSpritesImage = loadImage("objects/star.png")
var StairsSpritesImage = loadImage("objects/stairs.png")
var ChestSpritesImage = loadImage("objects/chest.png")
var HeartSpritesImage = loadImage("objects/heart.png")

var UiHealthHeartImage = loadImage("ui/heart.png")
var UiWeaponSelectWindowImage = loadImage("ui/weapon_select_window.png")
var UiIconsImage = loadImage("ui/icons.png")
var UiSelectRectImage = loadImage("ui/select_rect.png")

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
