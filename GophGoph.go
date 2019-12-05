package main

import (
	"fmt"         //FUNCIONES DE GOLANG
	_ "image/png" //PARA CARGAR LAS IMAGENES
	"log"
	"math/rand" //LIBRERIA MATEMÁTICA PARA NUMEROS ALEATORIOS

	"github.com/hajimehoshi/ebiten"            //PARA QUE FUNCIONE EL EBITEN
	"github.com/hajimehoshi/ebiten/ebitenutil" //PARA QUE FUNCIONE EL EBITEN
)

var img *ebiten.Image //VARIABLES DE LAS IMAGENES
var fondo *ebiten.Image
var barras *ebiten.Image
var plataforma *ebiten.Image
var hueco *ebiten.Image
var enemy *ebiten.Image
var escalera *ebiten.Image
var fondoMenu *ebiten.Image
var perdisteImg *ebiten.Image
var readySign *ebiten.Image
var goSign *ebiten.Image //VARIABLES DE LAS IMAGENES
var x float64            //COORDENADAS EN X PARA EL JUGADOR
var x2 float64           //COORDENADAS EN X PARA EL ENEMIGO
var y float64            //COORDENADAS EN Y PARA EL JUGADOR
var y2 float64           //COORDENADAS EN Y PARA EL ENEMIGO
var i float64
var num1 float64             //PARA EL NUMERO ALEATORIO DEL HUECO DEL JUGADOR
var num2 float64             //PARA EL NUMERO ALEATORIO DEL HUECO DEL ENEMIGO
var superficie float64       //VALOR DEL SUELO PARA EL JUGADOR
var superficie2 float64      //VALOR DEL SUELO PARA EL ENEMIGO
var sn1 string               //NUMERO Y POSICION DEL "HUECO" EN LAS PLATAFORMAS
var sn2 string               //VALOR TRUE/FALSE QUE APARECE PARA SABER SI LLEGO HASTA ABAJO
var sn3 string               //PARA QUE APAREZCA LA PUNTUACION
var sn4 string               //PARA EL NIVEL
var upup bool                //CUANDO LLEGA HASTA ABAJO Y SALTA HASTA ARRIBA
var puntuacion int64         //PUNTUACION
var velocidadY float64       //PARA LA VELOCIDAD QUE INCREMENTA VERTICALMENTE
var velocidadX float64       //PARA LA VELOCIDAD QUE INCREMENTA HORIZONTALMENTE
var velocidadEnemigo float64 //VELOCIDAD QUE TIENE EL ENEMIGO
var nivel int
var enemyRight bool //HACE QUE SE MUEVA A LA DERECHA AUTOMATICAMENTE
var perder bool     //CONDICION SI ES QUE SE PIERDE
var opMenu bool     //CONDICION PARA EL MENU PRINCIPAL
var time int        //TIEMPO QUE SE TOMA CON EL UPDATE

func init() { //SON LOS VALORES INICIALES CUANDO CORRE EL JUEGO
	var err error
	img, _, err = ebitenutil.NewImageFromFile("gopher.png", ebiten.FilterDefault)            //CARGA DE IMAGEN DEL JUGADOR
	fondo, _, err = ebitenutil.NewImageFromFile("fondo.png", ebiten.FilterDefault)           //CARGA DE IMAGEN DEL FONDO
	barras, _, err = ebitenutil.NewImageFromFile("barras.png", ebiten.FilterDefault)         //CARGA DE LA IMAGEN DE LOS HUECOS DEL JUGADOR
	plataforma, _, err = ebitenutil.NewImageFromFile("plataforma.png", ebiten.FilterDefault) //CARGA DE LA IMAGEN DE LAS PLATAFORMAS
	hueco, _, err = ebitenutil.NewImageFromFile("hueco.png", ebiten.FilterDefault)           //CARGA DE LA IMAGEN DE LOS HUECOS DEL JUGADOR
	enemy, _, err = ebitenutil.NewImageFromFile("enemy.png", ebiten.FilterDefault)           //CARGA DE LA IMAGEN DEL ENEMIGO
	escalera, _, err = ebitenutil.NewImageFromFile("Escalera.png", ebiten.FilterDefault)     //CARGA DE LA IMAGEN DEL HUECO DEL ENEMIGO
	fondoMenu, _, err = ebitenutil.NewImageFromFile("fondomenu.png", ebiten.FilterDefault)   //CARGA DE LA IMAGEN DEL FONDO DEL MENU
	perdisteImg, _, err = ebitenutil.NewImageFromFile("perdiste.png", ebiten.FilterDefault)  //CARGA DE LA IMAGEN DE PERDER
	readySign, _, err = ebitenutil.NewImageFromFile("Ready.png", ebiten.FilterDefault)       //CARGA DEL LETRERO DE READY
	goSign, _, err = ebitenutil.NewImageFromFile("Go.png", ebiten.FilterDefault)             //CARGA DEL LETRERO DE GO!
	if err != nil {                                                                          //POR SI OCURRE ALGUN ERROR, SE NOTIFIQUE
		log.Fatal(err)
	}
	x = 200                       //POSICIÓN INICIAL EN X DEL JUGADOR
	x2 = 200                      //POSICION INICIAL EN X DEL ENEMIGO
	y = 20                        //POSICIÓN INICIAL EN Y DEL JUGADOR
	y2 = -20                      //´POSICION INICIAL EN Y DEL ENEMIGO
	superficie = 300              //POSICIÓN INICIAL DE LAS SUPERFICIES DEL JUGADOR
	superficie2 = 300             //POSICION INICIAL DE LAS SUPERFICIES DEL ENEMIGO
	num1 = rand.Float64() * 1000  //EN NUMERO ALEATORIO INICIAL PARA EL HUECO DEL JUGADOR
	num2 = rand.Float64() * 1000  //NUMERO ALEATORIO INICIAL PARA EL HUECO DEL ENEMIGO
	sn1 = fmt.Sprintf("%f", num1) //PARA QUE APAREZCA EN LA PANTALLA
	upup = false
	sn2 = fmt.Sprintf("%t", upup)
	sn3 = fmt.Sprintf("%d", puntuacion)
	puntuacion = 0
	velocidadY = 1
	velocidadX = 0
	velocidadEnemigo = 0
	nivel = 1
	sn4 = fmt.Sprintf("%d", nivel)
	enemyRight = true
	perder = false
	opMenu = true
	time = 0
}

func menu(screen *ebiten.Image) { //FUNCION PARA EL MENU DE INICIO
	screen.DrawImage(fondoMenu, nil)          //DIBUJA LA IMAGEN DEL FONDO DEL MENU EN PANTALLA
	if ebiten.IsKeyPressed(ebiten.KeyEnter) { //SI SE PRESIONA LA TECLA ENTER
		opMenu = false
	}
}

func start(screen *ebiten.Image) { //FUNCION PARA LOS LETREROS DE READY, GO
	screen.DrawImage(fondo, nil) //DIBUJA EL FONDO DEL JUEGO
	time++                       //AUMENTA EL VALOR DEL TIEMPO
	if time < 30 {
		screen.DrawImage(readySign, nil) //MUESTRA EL LETRERO DE "READY?"
	} else if time < 60 {
		screen.DrawImage(goSign, nil) //MUESTRA EL LETRERO DE "GO!"
	} else if time >= 75 {
		juego(screen) //SE EJECUTA EL JUEGO
	}
}

func juego(screen *ebiten.Image) { //FUNCION QUE EJECUTA EL JUEGO PRINCIPAL
	opts := &ebiten.DrawImageOptions{}          //OPCIONES PARA LA IMAGEN DEL JUGADOR
	plat_ops := &ebiten.DrawImageOptions{}      //OPCIONES PARA LAS PLATAFORMAS DEL JUGADOR
	plat_ops2 := &ebiten.DrawImageOptions{}     //OPCIONES 2 PARA LAS PLATAFORMAS DEL JUGADOR
	enemyPlatOps := &ebiten.DrawImageOptions{}  //OPCIONES PARA LAS PLATAFORMAS DEL ENEMIGO
	enemyPlatOps2 := &ebiten.DrawImageOptions{} //OPCIONES PARA LAS PLATAFORMAS DEL ENEMIGO
	huecoOps := &ebiten.DrawImageOptions{}      //OPCIONES PARA LOS HUECOS DEL JUGADOR
	enemyOps := &ebiten.DrawImageOptions{}      //OPCIONES PARA EL MOVIMIENTO DEL ENEMIGO
	stairOps := &ebiten.DrawImageOptions{}      //OPCIONES PARA LOS HUECOS DEL ENEMIGO
	superficie = superficie - velocidadY        //PARA QUE SUBAN LAS PLATAFORMAS DEL JUGADOR
	superficie2 = superficie2 - velocidadY      //PARA QUE SUBAN LAS PLATAFORMAS DEL ENEMIGO

	if y >= 2500 { //SI LLEGA ABAJO upup ES VERDADERO
		upup = true
	}
	if upup == true { //CUANDO upup ES VERDADERO SALTA HASTA ARRIBA
		y -= 30
		if y <= 20 {
			upup = false
			superficie = 300
			velocidadY++
			velocidadX += 1.5
			velocidadEnemigo += 0.5
			nivel++
			sn4 = fmt.Sprintf("%d", nivel)
		}
	}
	if y != 2500 && upup == false { //CONDICIÓN PARA LA GRAVEDAD
		y += 10
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) { //CUANDO SE PRESIONA LA FLECHA A LA IZQUIERDA
		x -= 10 + velocidadX
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) { //CUANDO SE PRESIONA LA FLECHA A LA DERECHA
		x += 10 + velocidadX
	}
	if x >= 1740 { //CONDICIONES PARA QUE NO SALGA LATERALMENTE
		x = 1740
	}
	if x <= -45 { //CONDICIONES PARA QUE NO SALGA LATERALMENTE
		x = -45
	}
	if y >= superficie { //CONDICIONES PARA QUE NO SALGA VERTICALMENTE
		y = superficie
	}
	if x > num1 && x < (num1+100) && y == superficie { //ESTO ES LA POSICIÓN DEL HUECO QUE SE PONE ALEATORIAMENTE
		superficie += 300
		puntuacion += 100
		sn3 = fmt.Sprintf("%d", puntuacion) //AÑADE PUNTOS
		num1 = rand.Float64() * 1000        //GENERA OTRO NUMERO ALEATORIO
		sn1 = fmt.Sprintf("%f", num1)
	}
	if y <= -250 || y2 >= 2600 { //SI EL JUGADOR SALE DE PANTALLA HACIA ARRIBA O GANA EL ENEMIGO
		perder = true
	}

	//---------------VVV---OPCIONES DEL ENEMIGO---VVV--------------------------
	if upup == true { //SI EL JUGADOR LLEGA ABAJO, SE REINICIAN LOS VALORES INICIALES
		y2 = -40
		x2 = 0
		superficie2 = 300
	}
	if enemyRight == true { //BUSCA EL HUECO POR EL LADO DERECHO
		x2 += 5 + velocidadEnemigo
	} else {
		x2 -= 5 + velocidadEnemigo //BUSCA EL HUECO POR EL LADO IZQUIERDO
	}
	if x2 >= num2 { //1700 {
		enemyRight = false
	}
	if x2 <= num2 { //-40 {
		enemyRight = true
	}
	if y2 != 2500 { //ESTO ES LA CONDICIÓN DE LA GRAVEDAD
		y2 += 10
	}
	if y2 >= superficie2 { //CONDICIONES PARA QUE NO SALGA VERTICALMENTE
		y2 = superficie2
	}
	if x2 > num2 && x2 < (num2+100) && y2 == superficie2 { //ESTO ES LA POSICIÓN DEL HUECO QUE SE PONE ALEATORIAMENTE
		superficie2 += 300
		num2 = rand.Float64() * 1000
	}
	enemyOps.GeoM.Translate(x2, y2)                 //OPCIONES DE POSICION DEL ENEMIGO
	enemyOps.GeoM.Scale(0.25, 0.25)                 //OPCIONES DEL ESCALADO DE IMAGEN DEL ENEMIGO
	enemyPlatOps.GeoM.Translate(0, superficie2+195) //OPCIONES DE POSICION DE LAS PLATAFORMAS
	enemyPlatOps.GeoM.Scale(0.25, 0.25)             //OPCIONES DE ESCALADO DE LAS PLATAFORMAS
	enemyPlatOps2.GeoM.Translate(1000, superficie2+195)
	enemyPlatOps2.GeoM.Scale(0.25, 0.25)
	//-----------------------^^^-------------^^^-------------------------

	opts.GeoM.Translate(x, y)                      //OPCIONES DE MOVIMIENTO PARA EL JUGADOR
	opts.GeoM.Scale(0.25, 0.25)                    //OPCIONES DE LA ESCALA DE IMAGEN DEL JUGADOR
	huecoOps.GeoM.Translate(num1, superficie+195)  //OPCIONES DE POSICION DEL HUECO DEL JUGADOR
	huecoOps.GeoM.Scale(0.25, 0.25)                //OPCIONES DE ESCALA PARA LA IMAGEN DEL HUECO DEL JUGADOR
	stairOps.GeoM.Translate(num2, superficie2+195) //OPCIONES DE POSICION DEL HUECO DEL ENEMIGO
	stairOps.GeoM.Scale(0.25, 0.25)                //OPCIONES DE ESCALA PARA LA IMAGEN DEL HUECO DEL ENEMIGO
	plat_ops.GeoM.Translate(0, superficie+195)     //OPCIONES DE POSICION DE LAS PLATAFORMAS
	plat_ops.GeoM.Scale(0.25, 0.25)
	plat_ops2.GeoM.Translate(1000, superficie+195) //OPCIONES DE POSICION DE LAS PLATAFORMAS
	plat_ops2.GeoM.Scale(0.25, 0.25)
	screen.DrawImage(fondo, nil) //DIBUJA LAS IMAGENES, SE PRESENTAN COMO CAPAS
	if upup == false {
		screen.DrawImage(plataforma, plat_ops)
		screen.DrawImage(plataforma, plat_ops2)
		screen.DrawImage(plataforma, enemyPlatOps)
		screen.DrawImage(plataforma, enemyPlatOps2)
		screen.DrawImage(escalera, stairOps)
		screen.DrawImage(hueco, huecoOps)
		screen.DrawImage(enemy, enemyOps)
	}
	screen.DrawImage(img, opts)

	ebitenutil.DebugPrint(screen, "                                                  "+sn3)
	ebitenutil.DebugPrint(screen, "\n                                                  Nivel: "+sn4)
}

func perdiste(screen *ebiten.Image) { //FUNCION DE PERDER
	screen.DrawImage(perdisteImg, nil)        //DIBUJA LA IMAGEN DE PERDER
	if ebiten.IsKeyPressed(ebiten.KeyEnter) { //CUANDO SE PRESIONA ENTER
		perder = false //SE CONVIERTE EN FALSO PARA VOLVER A EMPEZAR
		x = 200        //SE REINICIAN LOS VALORES A LOS QUE ESTABAN INICIALMENTE
		x2 = 200
		y = 20
		y2 = -20
		superficie = 300
		superficie2 = 300
		num1 = rand.Float64() * 1000
		num2 = rand.Float64() * 1000
		upup = false
		puntuacion = 0
		velocidadY = 1
		velocidadX = 0
		velocidadEnemigo = 0
		nivel = 1
		sn4 = fmt.Sprintf("%d", nivel)
		time = 0
		start(screen) //SE VA A LA FUNCION DE LOS LETREROS
	}
}

func update(screen *ebiten.Image) error { //REFRESCA CONSTANTEMENTE PARA QUE APAREZCA EN PANTALLA
	if ebiten.IsDrawingSkipped() { //VALOR DE EBITENE
		return nil
	}
	if opMenu == true {
		menu(screen) //VA A LA FUNCION DEL MENU PRINCIPAL
	} else {
		start(screen) //VA A LA FUNCION DE LOS LETREROS DE READY GO
	}
	if perder == true {
		perdiste(screen) //VA A LA FUNCIÓN DE PERDER
	}
	return nil
}

func main() { //FUNCION MAIN
	if err := ebiten.Run(update, 480, 700, 1, "Goph Goph RACE"); err != nil { //LOS VALORES 2 Y 3 ENTRE PARENTESIS ES EL TAMAÑO DE LA VENTANA EN X & Y
		panic(err)
	}
}
