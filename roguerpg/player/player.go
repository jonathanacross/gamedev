package player

import (
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

	weapons         map[core.WeaponType]core.Weapon
	primaryWeapon   core.WeaponType
	secondaryWeapon core.WeaponType
	activeWeapon    core.Weapon // The weapon instance currently being updated/used
	shieldHeld      bool

	pendingActions       []core.Action
	pendingDamageSources []*core.DamageSource

	currentStats   *core.PlayerUpgrades
	futureUpgrades *core.PlayerUpgrades
}

func NewPlayer() *Player {
	// Player starts with just a sword and 3 hearts
	currentStats := &core.PlayerUpgrades{
		Health: 3,
		Weapons: map[core.WeaponType]int{
			core.WeaponSword:     1,
			core.WeaponBomb:      0,
			core.WeaponBoomerang: 0,
			core.WeaponShield:    0,
			core.WeaponBow:       0,
			core.WeaponWand:      0,
		},
	}
	// Player has no upgrades yet
	futureUpgrades := &core.PlayerUpgrades{
		Health: 0,
		Weapons: map[core.WeaponType]int{
			core.WeaponSword:     0,
			core.WeaponBomb:      0,
			core.WeaponBoomerang: 0,
			core.WeaponShield:    0,
			core.WeaponBow:       0,
			core.WeaponWand:      0,
		},
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
	pushBox := core.Rect{
		Left:   -6,
		Top:    -6,
		Right:  6,
		Bottom: 6,
	}

	// Initialize Weapons
	weaponMap := make(map[core.WeaponType]core.Weapon)
	weaponMap[core.WeaponSword] = weapons.NewSword()
	weaponMap[core.WeaponBow] = weapons.NewBow()
	weaponMap[core.WeaponWand] = weapons.NewWand()
	weaponMap[core.WeaponShield] = weapons.NewShield()
	weaponMap[core.WeaponBomb] = weapons.NewBomb()
	weaponMap[core.WeaponBoomerang] = weapons.NewBoomerang()

	return &Player{
		BaseCharacter: character.BaseCharacter{
			BasePhysical: core.BasePhysical{
				BaseSprite: core.BaseSprite{
					Loc: core.Location{
						X: 0,
						Y: 0,
					},
					DrawOffset: core.Location{
						X: 25,
						Y: 38,
					},
					SrcRect: spriteSheet.Rect(0),
				},
				PushBoxOffset: pushBox,
			},
			Health:          PlayerBaseHealth,
			MaxHealth:       PlayerBaseHealth,
			Experience:      0,
			KnockbackFrames: 0,
		},
		images:          charImages,
		spriteSheet:     spriteSheet,
		animations:      animations,
		state:           Idle,
		direction:       core.Down,
		primaryWeapon:   core.WeaponSword,
		secondaryWeapon: core.WeaponBoomerang,
		currentStats:    currentStats,
		futureUpgrades:  futureUpgrades,
		weapons:         weaponMap,
		activeWeapon:    weaponMap[core.WeaponSword], // Default
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

func (p *Player) attackWith(weaponType core.WeaponType) *core.Action {
	if weapon, ok := p.weapons[weaponType]; ok {
		p.activeWeapon = weapon
		weapon.OnAttack(p)
		return nil
	}
	return nil
}

func (p *Player) PrimaryAttack() *core.Action {
	return p.attackWith(p.primaryWeapon)
}

func (p *Player) SecondaryAttack() *core.Action {
	return p.attackWith(p.secondaryWeapon)
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
		return weapon.GetCooldownProgress()
	}
	return 0.0
}

func (p *Player) AddUpgrade(upgradeType core.UpgradeType) {
	switch upgradeType {
	case core.UpgradeTypeHeart:
		p.futureUpgrades.Health++
		p.MaxHealth += 2
		p.Health += 2
	case core.UpgradeTypeSword:
		p.futureUpgrades.Weapons[core.WeaponSword]++
		p.weapons[core.WeaponSword].SetLevel(p.futureUpgrades.Weapons[core.WeaponSword])
	case core.UpgradeTypeBomb:
		p.futureUpgrades.Weapons[core.WeaponBomb]++
		p.weapons[core.WeaponBomb].SetLevel(p.futureUpgrades.Weapons[core.WeaponBomb])
	case core.UpgradeTypeBoomerang:
		p.futureUpgrades.Weapons[core.WeaponBoomerang]++
		p.weapons[core.WeaponBoomerang].SetLevel(p.futureUpgrades.Weapons[core.WeaponBoomerang])
	case core.UpgradeTypeShield:
		p.futureUpgrades.Weapons[core.WeaponShield]++
		p.weapons[core.WeaponShield].SetLevel(p.futureUpgrades.Weapons[core.WeaponShield])
	case core.UpgradeTypeBow:
		p.futureUpgrades.Weapons[core.WeaponBow]++
		p.weapons[core.WeaponBow].SetLevel(p.futureUpgrades.Weapons[core.WeaponBow])
	case core.UpgradeTypeWand:
		p.futureUpgrades.Weapons[core.WeaponWand]++
		p.weapons[core.WeaponWand].SetLevel(p.futureUpgrades.Weapons[core.WeaponWand])
	}
}

func (p *Player) IsActive() bool {
	return p.state != Hurt && p.state != Dying && p.state != Dead
}

func (p *Player) handleState(level core.Level) []core.Action {
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
	if p.shieldHeld {
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

func (p *Player) GetUpgrades() *core.PlayerUpgrades {
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
	stateActions := p.handleState(level)

	p.Velocity.X = p.HandleTileCollisions(level, core.AxisX, p.Velocity.X)
	p.Velocity.Y = p.HandleTileCollisions(level, core.AxisY, p.Velocity.Y)

	animation := p.GetCurrentAnimation()
	animation.Update()
	p.Image = p.images[p.state]
	p.SrcRect = p.spriteSheet.Rect(animation.Frame())

	for _, w := range p.weapons {
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
	p.BaseCharacter.DrawDebugInfo(screen, cameraMatrix)
}
