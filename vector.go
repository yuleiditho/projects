package vector

import (
	"image/color"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/vector/internal/math"
)

var emptyImage *ebiten.Image

func init() {
	const w, h = 16, 16
	emptyImage, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
	emptyImage.Fill(color.White)
}


//Representa la colección de caminos
type Path struct {
	segs [][]math.Segment
	cur  math.Point
}

// MoveTo, es el movimiento en el mapa(x, y) 
	p.cur = math.Point{X: x, Y: y}
	if len(p.segs) > 0 && len(p.segs[len(p.segs)-1]) == 0 {
		return
	}
	p.segs = append(p.segs, []math.Segment{})
}

///Cambio de las estructura de la plataformas
type platform struct {
	rect  pixel.Rect
	color color.Color
}
 //Dibuja la platoforma siguiente
func (p *platform) draw(imd *imdraw.IMDraw) {
	imd.Color = p.color
	imd.Push(p.rect.Min, p.rect.Max)

}
//Configuración de la gravedad
type gopherPhys struct {
	gravity   float64
	runSpeed  float64
	jumpSpeed float64
	ground bool
}

func (gp *gopherPhys) update(dt float64, ctrl pixel.Vec, platforms []platform) {
	//Aplicando los controles
	switch {
	case ctrl.X < 0:
		gp.vel.X = -gp.runSpeed
	case ctrl.X > 0:
		gp.vel.X = +gp.runSpeed
	default:
		gp.vel.X = 0
	}

	// Aplicando la gravedad y la velocidad
	gp.vel.Y += gp.gravity * dt
	gp.rect = gp.rect.Moved(gp.vel.Scaled(dt))

	// Checa las colosiones al regresar al estado incial
	gp.ground = false
	if gp.vel.Y <= 0 {
		for _, p := range platforms {
			if gp.rect.Max.X <= p.rect.Min.X || gp.rect.Min.X >= p.rect.Max.X {
				continue
			}
			if gp.rect.Min.Y > p.rect.Max.Y || gp.rect.Min.Y < p.rect.Max.Y+gp.vel.Y*dt {
				continue
			}
			gp.vel.Y = 0
			gp.rect = gp.rect.Moved(pixel.V(0, p.rect.Max.Y-gp.rect.Min.Y))
			gp.ground = true
		}
	}

	// salto que realiza en Y al momento de estar al estado final
	if gp.ground && ctrl.Y > 0 {
		gp.vel.Y = gp.jumpSpeed
	}
}

//Los movimientos del goph
type animState int

const (
	gopher animState = iota
	running
	jumping
)

//
type gopherAnim struct {
	anims map[string][]pixel.Rect //Mapeo principal 
	point  float64				  //Puntuación

	state   animState			//Estado
	counter float64				//Contador
	dir     float64				//Dirección

	frame pixel.Rect

	sprite *pixel.Sprite
}

//Función de las transcisiones
func (ga *gopherAnim) update(dt float64, phys *gopherPhys) {
	ga.counter += dt

	//Determinar el nuevo estado del personaje
	var newState animState 
	switch {
	case !phys.ground:
		newState = jumping
	case phys.vel.Len() == 0:
		newState = idle
	case phys.vel.Len() > 0:
		newState = running
	}

	// Resetear el contador del tiempo, si el estado ha cambiado
	if ga.state != newState {
		ga.state = newState
		ga.counter = 0
	}

	//Determinar la correcta impresión de la imagen del personaje de transición de la plataforma
	switch ga.state {
	case idle:
		ga.frame = ga.anims["Front"][0]
	case running:
		i := int(math.Floor(ga.counter / ga.rate))
		ga.frame = ga.anims["Run"][i%len(ga.anims["Run"])]
	case jumping:
		speed := phys.vel.Y
		i := int((-speed/phys.jumpSpeed + 1) / 2 * float64(len(ga.anims["Jump"])))
		if i < 0 {
			i = 0
		}
		if i >= len(ga.anims["Jump"]) {
			i = len(ga.anims["Jump"]) - 1
		}
		ga.frame = ga.anims["Jump"][i]
	}

	// set the facing direction of the gopher
	if phys.vel.X != 0 {
		if phys.vel.X > 0 {
			ga.dir = +1
		} else {
			ga.dir = -1
		}
	}
}

func (ga *gopherAnim) draw(t pixel.Target, phys *gopherPhys) {
	if ga.sprite == nil {
		ga.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}
	//Dibuja el correcto frame con la posición y dirección 
	ga.sprite.Set(ga.sheet, ga.frame)
	ga.sprite.Draw(t, pixel.IM.
		ScaledXY(pixel.ZV, pixel.V(
			phys.rect.W()/ga.sprite.Frame().W(),
			phys.rect.H()/ga.sprite.Frame().H(),
		)).
		ScaledXY(pixel.ZV, pixel.V(-ga.dir, 1)).
		Moved(phys.rect.Center()),
	)
}


