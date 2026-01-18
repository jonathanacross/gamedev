package player

import (
	"fmt"
	"image/color"
	"roguerpg/assets"
	"roguerpg/character"
	"roguerpg/core"
	"roguerpg/weapons"

	"github.com/hajimehoshi/ebiten/v2"
)

type PlayerState int

const (
	Idle PlayerState = iota
	Walking
	Hurt
	Dying
	Dead
)

const (
	PlayerSpeed      = 2.0
	PlayerBaseHealth = 6
)

type PlayerAnimationKey struct {
	State     PlayerState
	Direction core.Direction
}

type Player struct {
	character.BaseCharacter
	images      map[PlayerState]*ebiten.Image
	spriteSheet *core.SpriteSheet
	animations  map[PlayerState]map[core.Direction]*core.Animation
	state       PlayerState
	direction   core.Direction
	Velocity    core.Vector

	weapons             map[core.WeaponType]core.Weapon
	primaryWeapon       core.WeaponType
	secondaryWeapon     core.WeaponType
	activeWeapon        core.Weapon // The weapon instance currently being updated/used
	shieldHeld          bool
	hurtBoxesWithShield map[core.Direction]core.Rect

	pendingActions       []core.Action
	pendingDamageSources []*core.DamageSource

	currentStats   core.PlayerUpgrades
	futureUpgrades core.PlayerUpgrades
}

func NewPlayer() *Player {
	// Player starts with just a sword and 3 hearts
	currentStats := core.PlayerUpgrades{
		core.UpgradeTypeHeart:     3,
		core.UpgradeTypeSword:     1,
		core.UpgradeTypeBomb:      0,
		core.UpgradeTypeBoomerang: 0,
		core.UpgradeTypeShield:    0,
		core.UpgradeTypeBow:       0,
		core.UpgradeTypeWand:      0,
	}
	// Player has no upgrades yet
	futureUpgrades := core.PlayerUpgrades{
		core.UpgradeTypeHeart:     0,
		core.UpgradeTypeSword:     0,
		core.UpgradeTypeBomb:      0,
		core.UpgradeTypeBoomerang: 0,
		core.UpgradeTypeShield:    0,
		core.UpgradeTypeBow:       0,
		core.UpgradeTypeWand:      0,
	}

	ssColumns := 8
	ssRows := 6
	directionOffsets := map[core.Direction]int{
		core.Left:  8,
		core.Right: 40,
		core.Up:    24,
		core.Down:  0,
	}

	animations := map[PlayerState]map[core.Direction]*core.Animation{
		Idle:    core.NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5, 6, 7}, 0, directionOffsets, 10, true),
		Walking: core.NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5, 6, 7}, 0, directionOffsets, 10, true),
		Hurt:    core.NewDirectionAnimationMap([]int{0, 1, 2, 3}, 0, directionOffsets, 6, false),
		Dying:   core.NewDirectionAnimationMap([]int{0, 1, 2, 3, 4, 5, 6, 7}, 0, directionOffsets, 8, false),
		Dead:    core.NewDirectionAnimationMap([]int{7}, 0, directionOffsets, 8, false),
	}

	charImages := map[PlayerState]*ebiten.Image{
		Idle:    assets.PlayerIdleSpritesImage,
		Walking: assets.PlayerWalkSpritesImage,
		Hurt:    assets.PlayerHurtSpritesImage,
		Dying:   assets.PlayerDeathSpritesImage,
		Dead:    assets.PlayerDeathSpritesImage,
	}

	spriteSheet := core.NewSpriteSheet(48, 64, ssColumns, ssRows)
	pushBox := core.Rect{Left: -6, Top: -6, Right: 6, Bottom: 6}
	hurtBoxesWithShield := map[core.Direction]core.Rect{
		core.Left:  {Left: 1, Top: -6, Right: 6, Bottom: 6},
		core.Right: {Left: -6, Top: -6, Right: -1, Bottom: 6},
		core.Up:    {Left: -6, Top: 1, Right: 6, Bottom: 6},
		core.Down:  {Left: -6, Top: -6, Right: 6, Bottom: -1},
	}

	// Initialize Weapons
	weaponMap := make(map[core.WeaponType]core.Weapon)
	weaponMap[core.WeaponSword] = weapons.NewSword()
	weaponMap[core.WeaponBow] = weapons.NewBow()
	weaponMap[core.WeaponWand] = weapons.NewWand()
	weaponMap[core.WeaponShield] = weapons.NewShield()
	weaponMap[core.WeaponBomb] = weapons.NewBomb()
	weaponMap[core.WeaponBoomerang] = weapons.NewBoomerang()
	weaponMap[core.WeaponNone] = nil

	return &Player{
		BaseCharacter: character.BaseCharacter{
			BasePhysical: core.BasePhysical{
				BaseSprite: core.BaseSprite{
					Loc:        core.Location{X: 0, Y: 0},
					Image:      charImages[Idle],
					DrawOffset: core.Location{X: 25, Y: 38},
					SrcRect:    spriteSheet.Rect(0),
				},
				PushBoxOffset: pushBox,
			},
			Health:          PlayerBaseHealth,
			MaxHealth:       PlayerBaseHealth,
			Experience:      0,
			KnockbackFrames: 0,
		},
		images:              charImages,
		spriteSheet:         spriteSheet,
		animations:          animations,
		hurtBoxesWithShield: hurtBoxesWithShield,
		state:               Idle,
		direction:           core.Down,
		primaryWeapon:       core.WeaponSword,
		secondaryWeapon:     core.WeaponNone,
		currentStats:        currentStats,
		futureUpgrades:      futureUpgrades,
		weapons:             weaponMap,
		activeWeapon:        weaponMap[core.WeaponSword], // Default
	}
}

// -- PlayerContext Implementation --

func (p *Player) GetDirection() core.Direction {
	return p.direction
}

var _ core.PlayerContext = (*Player)(nil) // Verify interface compliance

func (p *Player) SetBlocking(active bool) {
	p.shieldHeld = active
}

func (p *Player) AddAction(action core.Action) {
	p.pendingActions = append(p.pendingActions, action)
}

func (p *Player) CreateDamageSource(ds *core.DamageSource) {
	p.pendingDamageSources = append(p.pendingDamageSources, ds)
}

// -- End PlayerContext --

func (p *Player) GetCurrentAnimation() *core.Animation {
	animationSet, exists := p.animations[p.state]
	if !exists {
		return nil
	}
	animation, exists := animationSet[p.direction]
	if !exists {
		return nil
	}
	return animation
}

func (p *Player) Move(moveVector core.Vector) {
	// Determine facing direction from the move vector.
	if moveVector.Y < 0 {
		p.direction = core.Up
	} else if moveVector.Y > 0 {
		p.direction = core.Down
	} else if moveVector.X < 0 {
		p.direction = core.Left
	} else if moveVector.X > 0 {
		p.direction = core.Right
	}

	p.Velocity = moveVector.Normalize().Scale(PlayerSpeed)
	p.TransitionState(Walking)
}

func (p *Player) StopMoving() {
	if p.state == Walking {
		p.TransitionState(Idle)
	}
}

func (p *Player) attackWith(weaponType core.WeaponType) {
	if p.weapons[weaponType] == nil {
		return
	}

	if p.activeWeapon != nil && p.activeWeapon != p.weapons[weaponType] && p.activeWeapon.IsAttacking() {
		return
	}

	if weapon, ok := p.weapons[weaponType]; ok {
		p.activeWeapon = weapon
		weapon.OnAttack(p)
		return
	}
}

func (p *Player) PrimaryAttack() *core.Action {
	p.attackWith(p.primaryWeapon)
	return nil
}

func (p *Player) SecondaryAttack() *core.Action {
	p.attackWith(p.secondaryWeapon)
	return nil
}

func (p *Player) SwitchWeapon(weapon core.WeaponType) {
	p.secondaryWeapon = weapon
}

func (p *Player) GetPrimaryWeapon() core.WeaponType {
	return p.primaryWeapon
}

func (p *Player) GetSecondaryWeapon() core.WeaponType {
	return p.secondaryWeapon
}

func (p *Player) GetWeaponProgress(weaponType core.WeaponType) float64 {
	if weapon, ok := p.weapons[weaponType]; ok {
		if weapon == nil {
			return 0.0
		}
		return weapon.GetCooldownProgress()
	}
	return 0.0
}

func (p *Player) AddUpgrade(upgradeType core.UpgradeType) []core.Action {
	if upgradeType == core.UpgradeTypeNone {
		return []core.Action{}
	}

	p.futureUpgrades[upgradeType]++

	if p.currentStats[upgradeType] == 0 {
		// If the player doesn't have the weapon yet, give it to them immediately
		return p.DoUpgrade(upgradeType, 0)
	}

	return []core.Action{}
}

func (p *Player) DoUpgrade(upgradeType core.UpgradeType, cost int) []core.Action {
	if upgradeType == core.UpgradeTypeNone {
		return []core.Action{}
	}
	if p.futureUpgrades[upgradeType] <= 0 {
		return []core.Action{}
	}
	if p.Experience < cost {
		return []core.Action{}
	}

	justObtained := p.currentStats[upgradeType] == 0

	p.futureUpgrades[upgradeType]--
	p.currentStats[upgradeType]++
	p.Experience -= cost

	upgradeName := ""
	switch upgradeType {
	case core.UpgradeTypeHeart:
		p.MaxHealth += 2
		p.Health += 2
		upgradeName = "Life"
	case core.UpgradeTypeSword:
		p.weapons[core.WeaponSword].SetLevel(p.currentStats[core.UpgradeTypeSword])
		upgradeName = "Sword"
	case core.UpgradeTypeBomb:
		p.weapons[core.WeaponBomb].SetLevel(p.currentStats[core.UpgradeTypeBomb])
		upgradeName = "Bomb"
	case core.UpgradeTypeBoomerang:
		p.weapons[core.WeaponBoomerang].SetLevel(p.currentStats[core.UpgradeTypeBoomerang])
		upgradeName = "Boomerang"
	case core.UpgradeTypeShield:
		p.weapons[core.WeaponShield].SetLevel(p.currentStats[core.UpgradeTypeShield])
		upgradeName = "Shield"
	case core.UpgradeTypeBow:
		p.weapons[core.WeaponBow].SetLevel(p.currentStats[core.UpgradeTypeBow])
		upgradeName = "Bow"
	case core.UpgradeTypeWand:
		p.weapons[core.WeaponWand].SetLevel(p.currentStats[core.UpgradeTypeWand])
		upgradeName = "Wand"
	}

	upgradeMessage := fmt.Sprintf("Upgraded %s to level %d!", upgradeName, p.currentStats[upgradeType])
	if justObtained {
		upgradeMessage = fmt.Sprintf("Obtained %s!", upgradeName)
	}
	return []core.Action{
		{
			Type:    core.ActionShowMessage,
			Message: upgradeMessage,
		},
	}
}

func (p *Player) IsActive() bool {
	return p.state != Hurt && p.state != Dying && p.state != Dead
}

func (p *Player) GetHurtBox() core.Rect {
	if !p.shieldHeld {
		return p.GetPushBox()
	}
	return p.hurtBoxesWithShield[p.direction].Offset(p.Loc.X, p.Loc.Y)
}

func (p *Player) handleState() []core.Action {
	animation := p.GetCurrentAnimation()
	if animation == nil {
		p.state = Idle
		return nil
	}

	if p.Health <= 0 && p.state != Dying && p.state != Dead {
		p.TransitionState(Dying)
		return nil
	}

	if p.state == Dying && animation.IsFinished() {
		p.TransitionState(Dead)
		return nil
	}

	if p.state == Dead {
		p.Velocity = core.Vector{}
		return nil
	}

	if p.state == Hurt && !p.IsKnockedBack() && animation.IsFinished() {
		p.TransitionState(Idle)
	}

	switch p.state {
	case Idle:
		p.Velocity = core.Vector{}
	case Walking:
	case Hurt, Dying, Dead:
		p.Velocity = core.Vector{}
	}
	return nil
}

func (p *Player) TransitionState(newState PlayerState) {
	if p.state == newState {
		return
	}
	p.state = newState
	if anim := p.GetCurrentAnimation(); anim != nil {
		anim.Reset()
	}
	p.BaseCharacter.Dead = (newState == Dead)
}

func (p *Player) HandleHit(ds *core.DamageSource) {
	if ds.Type == core.DamageTypeStun {
		p.ApplyStun(ds.Duration)
	} else {
		p.TakeDamage(ds.Damage)
		force := core.CalculateKnockbackForce(ds.HitBox.Center(), p.Location(), core.KnockbackForce)
		p.ApplyKnockback(force, core.KnockbackDuration)
	}
}

func (p *Player) TakeDamage(damage int) {
	if p.state == Dead || p.state == Dying || p.state == Hurt {
		return
	}
	p.TransitionState(Hurt)
	p.Health -= damage
}

func (p *Player) ApplyKnockback(force core.Vector, duration int) {
	if p.shieldHeld {
		return
	}
	p.BaseCharacter.ApplyKnockback(force, duration)
	if p.state == Dying || p.state == Dead {
		return
	}
	if p.IsKnockedBack() {
		p.TransitionState(Hurt)
	}
}

func (p *Player) ReturnBoomerang() {
	// Find the boomerang weapon in the map and call its specific method.
	// We need type assertion since the generic Weapon interface doesn't have Catch.
	if w, ok := p.weapons[core.WeaponBoomerang]; ok {
		if bw, ok := w.(*weapons.Boomerang); ok {
			bw.Catch()
		}
	}
}

func (p *Player) AddExperience(amount int) {
	p.Experience += amount
}

func (p *Player) AddHealth(amount int) {
	p.Health += amount
	if p.Health > p.MaxHealth {
		p.Health = p.MaxHealth
	}
}

func (p *Player) GetExperience() int {
	return p.Experience
}

func (p *Player) GetCurrentStats() core.PlayerUpgrades {
	return p.currentStats
}

func (p *Player) GetFutureUpgrades() core.PlayerUpgrades {
	return p.futureUpgrades
}

func (p *Player) handleObjectPickup(level core.Level) {
	playerPushBox := p.GetPushBox()
	for _, object := range level.GetObjects() {
		if interactable, ok := object.(core.Interactable); ok {
			if playerPushBox.Intersects(interactable.GetPushBox()) {
				interactable.Touch(level, p)
				break
			}
		}
	}
}

func (p *Player) Update(level core.Level, _ core.Player) core.UpdateResult {
	p.UpdateKnockback(level)
	p.handleObjectPickup(level)
	stateActions := p.handleState()

	p.Velocity.X = p.HandleTileCollisions(level, core.AxisX, p.Velocity.X)
	p.Velocity.Y = p.HandleTileCollisions(level, core.AxisY, p.Velocity.Y)

	animation := p.GetCurrentAnimation()
	animation.Update()
	p.Image = p.images[p.state]
	p.SrcRect = p.spriteSheet.Rect(animation.Frame())

	for _, w := range p.weapons {
		if w == nil {
			continue
		}
		w.Update(p, level)
	}

	actions := []core.Action{}
	if len(stateActions) > 0 {
		actions = append(actions, stateActions...)
	}
	if len(p.pendingActions) > 0 {
		actions = append(actions, p.pendingActions...)
		p.pendingActions = nil // Reset after consuming
	}

	// Convert pending damage sources to Actions
	for _, ds := range p.pendingDamageSources {
		actions = append(actions, core.Action{
			Type:         core.ActionCreateDamageSource,
			Location:     p.Location(),
			Direction:    core.DirectionToVector(p.direction),
			DamageSource: ds,
		})
	}
	// Reset damage sources after consuming
	p.pendingDamageSources = nil

	return core.UpdateResult{Actions: actions}
}

func (p *Player) CanRemove() bool {
	return false
}

func (p *Player) Draw(screen *ebiten.Image, cameraMatrix ebiten.GeoM) {
	// Draw Weapons designated as "Behind"
	if p.activeWeapon != nil && p.activeWeapon.GetRenderOrder(p.direction) == core.RenderOrderBehind {
		p.activeWeapon.Draw(screen, cameraMatrix, p.Location(), p.direction, p.GetCurrentAnimation().Frame())
	}

	// Draw Player Body
	p.BaseCharacter.Draw(screen, cameraMatrix)

	// Draw Weapons designated as "Front"
	if p.activeWeapon != nil && p.activeWeapon.GetRenderOrder(p.direction) == core.RenderOrderFront {
		p.activeWeapon.Draw(screen, cameraMatrix, p.Location(), p.direction, p.GetCurrentAnimation().Frame())
	}
}

func (p *Player) DrawDebugInfo(screen *ebiten.Image, cameraMatrix ebiten.GeoM) {
	// Note this is the same code as in BaseCharacter, but we need to re-implmement
	// it because we need to call Player.GetHurtBox() instead of Character.GetHurtBox().
	boundsColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	pushBoxColor := color.RGBA{R: 0, G: 0, B: 255, A: 255}
	hurtBoxColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	core.DrawDebugRect(screen, p.GetBounds(), boundsColor, cameraMatrix)
	core.DrawDebugRect(screen, p.GetPushBox(), pushBoxColor, cameraMatrix)
	core.DrawDebugRect(screen, p.GetHurtBox(), hurtBoxColor, cameraMatrix)
	core.DrawLocationDot(screen, p.Location(), cameraMatrix)
}
