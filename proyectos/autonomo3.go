package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// Usuario encapsulado
type Usuario struct {
	CI               string
	NombresCompletos string
}

// Producto encapsulado
type Producto struct {
	Nombre   string
	Precio   float64
	Cantidad int
	Imagen   string
}

// Orden encapsulada
type Orden struct {
	Codigo    string
	Usuario   Usuario
	Productos []Producto
	Total     float64
	Pagado    bool
}

// Variables globales
var usuarios []Usuario
var productos = []Producto{
	{"Cerveza IPA", 5.50, 100, "/static/img/ipa.jpg"},
	{"Cerveza Stout", 6.00, 80, "/static/img/stout.jpg"},
	{"Cerveza Lager", 4.50, 150, "/static/img/lager.jpg"},
	{"Cerveza Pilsner", 3.50, 200, "/static/img/pilsner.jpg"},
	{"Cerveza Porter", 6.50, 90, "/static/img/porter.jpg"},
	{"Cerveza Ale", 7.00, 70, "/static/img/ale.jpg"},
	{"Cerveza Blonde", 5.00, 110, "/static/img/blonde.jpg"},
	{"Cerveza Amber", 5.75, 120, "/static/img/amber.jpg"},
	{"Cerveza Wheat", 4.75, 140, "/static/img/wheat.jpg"},
	{"Cerveza Dubbel", 8.50, 60, "/static/img/dubbel.jpg"},
}
var ordenes []Orden

// Generar código aleatorio para la orden
func GenerarCodigoOrden() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("ORD-%d", rand.Intn(100000))
}

// Calcular Total de una lista de productos
func CalcularTotal(productos []Producto) float64 {
	total := 0.0
	for _, p := range productos {
		total += p.Precio * float64(p.Cantidad)
	}
	return total
}

// Página de inicio
func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Inicio</title>
		<link rel="stylesheet" href="/static/css/styles.css">
	</head>
	<body>
		<div class="container">
			<h1>Bienvenido a la Tienda de Cervezas</h1>
			<a href="/agregarUsuario" class="btn-submit">Comenzar</a>
		</div>
	</body>
	</html>
	`
	t := template.Must(template.New("home").Parse(tmpl))
	t.Execute(w, nil)
}

// Página para agregar usuario
func agregarUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		nombre := r.FormValue("nombre")
		ci := r.FormValue("ci")
		if nombre == "" || ci == "" {
			http.Error(w, "Nombre y cédula son obligatorios", http.StatusBadRequest)
			return
		}
		usuarios = append(usuarios, Usuario{CI: ci, NombresCompletos: nombre})
		http.Redirect(w, r, "/inventario?ci="+ci, http.StatusSeeOther)
		return
	}
	tmpl := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Agregar Usuario</title>
		<link rel="stylesheet" href="/static/css/styles.css">
	</head>
	<body>
		<div class="container">
			<h1>Agregar Usuario</h1>
			<form action="/agregarUsuario" method="POST">
				<label for="nombre">Nombre:</label>
				<input type="text" id="nombre" name="nombre" required>
				<label for="ci">Cédula:</label>
				<input type="text" id="ci" name="ci" required>
				<button type="submit" class="btn-submit">Agregar</button>
			</form>
		</div>
	</body>
	</html>
	`
	t := template.Must(template.New("agregarUsuario").Parse(tmpl))
	t.Execute(w, nil)
}

// Página de inventario para realizar pedido
func inventarioHandler(w http.ResponseWriter, r *http.Request) {
	ci := r.URL.Query().Get("ci")
	var usuario *Usuario
	for _, u := range usuarios {
		if u.CI == ci {
			usuario = &u
			break
		}
	}
	if usuario == nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPost {
		var pedidoProductos []Producto
		for i := range productos {
			cantidadStr := r.FormValue("cantidad_" + productos[i].Nombre)
			cantidad, err := strconv.Atoi(cantidadStr)
			if err != nil || cantidad < 0 || cantidad > productos[i].Cantidad {
				continue
			}
			if cantidad > 0 {
				productos[i].Cantidad -= cantidad
				pedidoProductos = append(pedidoProductos, Producto{
					Nombre:   productos[i].Nombre,
					Precio:   productos[i].Precio,
					Cantidad: cantidad,
				})
			}
		}
		total := CalcularTotal(pedidoProductos)
		codigo := GenerarCodigoOrden()
		orden := Orden{
			Codigo:    codigo,
			Usuario:   *usuario,
			Productos: pedidoProductos,
			Total:     total,
			Pagado:    false,
		}
		ordenes = append(ordenes, orden)
		http.Redirect(w, r, "/detallePedido?codigo="+codigo, http.StatusSeeOther)
		return
	}

	tmpl := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Inventario</title>
		<link rel="stylesheet" href="/static/css/styles.css">
	</head>
	<body>
		<div class="container">
			<h1>Realizar Pedido</h1>
			<form action="/inventario?ci={{.Usuario.CI}}" method="POST">
				<div class="product-list">
					{{range .Productos}}
					<div class="product-item">
						<img src="{{.Imagen}}" alt="{{.Nombre}}" class="product-image">
						<h3>{{.Nombre}}</h3>
						<p>Precio: ${{.Precio}}</p>
						<p>Cantidad Disponible: {{.Cantidad}}</p>
						<label for="cantidad_{{.Nombre}}">Cantidad:</label>
						<input type="number" id="cantidad_{{.Nombre}}" name="cantidad_{{.Nombre}}" min="0" max="{{.Cantidad}}" value="0">
					</div>
					{{end}}
				</div>
				<button type="submit" class="btn-submit">Realizar Pedido</button>
			</form>
		</div>
	</body>
	</html>
	`
	data := struct {
		Usuario   *Usuario
		Productos []Producto
	}{
		Usuario:   usuario,
		Productos: productos,
	}
	t := template.Must(template.New("inventario").Parse(tmpl))
	t.Execute(w, data)
}

// Página de detalle del pedido
func detallePedidoHandler(w http.ResponseWriter, r *http.Request) {
	codigo := r.URL.Query().Get("codigo")
	for _, orden := range ordenes {
		if orden.Codigo == codigo {
			tmpl := `
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Detalle del Pedido</title>
				<link rel="stylesheet" href="/static/css/styles.css">
			</head>
			<body>
				<div class="container">
					<h1>Detalle del Pedido</h1>
					<p>Número de Orden: <b>{{.Codigo}}</b></p>
					<p>Usuario: {{.Usuario.NombresCompletos}}</p>
					<h2>Productos:</h2>
					<ul>
						{{range .Productos}}
						<li>{{.Cantidad}} x {{.Nombre}} - ${{printf "%.2f" .Precio}}</li>
						{{end}}
					</ul>
					<p>Total: ${{printf "%.2f" .Total}}</p>
					<a href="/pago?codigo={{.Codigo}}" class="btn">Realizar Pago</a>
				</div>
			</body>
			</html>
			`
			t := template.Must(template.New("detallePedido").Parse(tmpl))
			t.Execute(w, orden)
			return
		}
	}
	http.Error(w, "Pedido no encontrado", http.StatusNotFound)
}

// Página para realizar pago
func pagoHandler(w http.ResponseWriter, r *http.Request) {
	codigo := r.URL.Query().Get("codigo")
	for i, orden := range ordenes {
		if orden.Codigo == codigo {
			ordenes[i].Pagado = true
			tmpl := `
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Pago Realizado</title>
				<link rel="stylesheet" href="/static/css/styles.css">
			</head>
			<body>
				<div class="container">
					<h1>Pago Realizado</h1>
					<p>Gracias por tu compra. El pedido con número <b>{{.Codigo}}</b> ha sido pagado con éxito.</p>
					<a href="/" class="btn">Volver al Inicio</a>
				</div>
			</body>
			</html>
			`
			t := template.Must(template.New("pagoRealizado").Parse(tmpl))
			t.Execute(w, orden)
			return
		}
	}
	http.Error(w, "Pedido no encontrado", http.StatusNotFound)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/agregarUsuario", agregarUsuarioHandler)
	http.HandleFunc("/inventario", inventarioHandler)
	http.HandleFunc("/detallePedido", detallePedidoHandler)
	http.HandleFunc("/pago", pagoHandler)

	fmt.Println("Servidor corriendo en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
