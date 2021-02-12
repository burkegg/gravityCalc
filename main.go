package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math"
	"net/http"
	"time"
)

type Ball struct {
	Mass       float64
	Color      string
	Pos        Vector
	Vel        Vector
	F          Vector
	Num        int `json:",omitempty"`
}

type GravityEngine struct {
	G             float64
	DeltaT        float64
	NumIterations int
	InitData      []Ball
	NextData      []Ball
	K             float64
	Trajectories  [][]Vector
}

// K is the fudge factor to avoid zero distances

// Class
func NewGravityEngine(G float64, DeltaT float64, NumIterations int) *GravityEngine {
	ge := GravityEngine{
		G,
		DeltaT,
		NumIterations,
		make([]Ball, 0),
		make([]Ball, 0),
		0.1,
		make([][]Vector, 0),
	}
	return &ge
}

// Calculate the force from ball one on ball 2
// Return a vector
func (ge *GravityEngine) CalcOneForce(b1 Ball, b2 Ball) (f Vector) {
	// f = g * m1 * m2 / s2 + ks
	s2 := DistSq(b2.Pos, b1.Pos)
	fMag := ge.G * b2.Mass * b1.Mass / (s2 + ge.K) /// woops

	direction := Subtract(b2.Pos, b1.Pos)
	unitV := Normalize(direction)
	f = ScalarMult(unitV, fMag)
	return f
	// get unit vector * multiply, then return the new vector
}

//Calculates the net force on a  ball by repeatedly calling CalcOneForce
func (ge *GravityEngine) CalcNetForceOnBall(ballNum int) (f Vector) {

	// add the forces of every other ball
	f = Vector{X: 0, Y: 0}

	for i := 0; i < len(ge.InitData); i++ {
		if i != ballNum {
			f1 := ge.CalcOneForce(ge.InitData[i], ge.InitData[ballNum])
			f = Add(f, f1)
		}
	}
	return f
}

// Moves balls one timestep
// Adds that position to history
func (ge *GravityEngine) MoveBalls() {

	for i := 0; i < ge.NumIterations; i++ {
		for ballNum := 0; ballNum < len(ge.InitData); ballNum++ {
			ge.GetNewData(ballNum)
		}
		ge.UpdateInitData()
	}
	for ballNum := 0; ballNum < len(ge.InitData); ballNum++ {
		// here we're always adding - but we want to replace and not grow the arrays
		ge.Trajectories[ballNum] = append(ge.Trajectories[ballNum], ge.NextData[ballNum].Pos)
	}
}

func (ge *GravityEngine) MakeHistory(numSteps int) {
	ge.Trajectories = make([][]Vector, len(ge.InitData))
	for stepNum := 0; stepNum < numSteps; stepNum++ {
		ge.MoveBalls() // this one moves everything one step
	}
}

// at some interval send data to the client.
func (ge *GravityEngine) PumpData(ws *websocket.Conn) {
	fmt.Println("BEGIN PUMPING")
	err := ws.WriteMessage(1, []byte("Pumpin'!"))
	if err != nil {
		log.Println(err)
	}
	for {
		time.Sleep(2000 * time.Millisecond)
		ge.MakeHistory(5)
		data, err := json.Marshal(ge.Trajectories)
		if err != nil {
			log.Println(err)
		}
		err = ws.WriteMessage(1, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func (ge *GravityEngine) GetNewData(ballNum int) {
	// get net force for ball
	// apply kinematics
	// update nextData
	// To not overwrite the init array apparently I have to make a copy
	init := make([]Ball, len(ge.InitData))
	//fmt.Printf("Length of init: %+v\n", len(ge.InitData))
	copy(init, ge.InitData)
	f := ge.CalcNetForceOnBall(ballNum)
	a := ScalarDivide(f, init[ballNum].Mass)
	a = ScalarMult(a, ge.DeltaT)
	v := Add(init[ballNum].Vel, a)
	v = ScalarMult(v, ge.DeltaT)
	p := Vector{X: init[ballNum].Pos.X, Y: init[ballNum].Pos.Y}
	p = Add(p, v)
	vect := Vector{X: p.X, Y: p.Y} // here

	ge.NextData[ballNum].Pos = vect
	ge.NextData[ballNum].Vel = Vector{X: v.X, Y: v.Y}
	//fmt.Printf("NextData[0]", ge.NextData[0])
}

func (ge *GravityEngine) UpdateInitData() {
	//fmt.Printf("Inside UpdateInitData, init: %+v\n", ge.InitData)
	//fmt.Printf("Inside UpdateInitData, Next: %+v\n", ge.NextData)

	// update one at a time until I figure out how to do it right?

	init := make([]Ball, len(ge.InitData))
	copy(init[0:], ge.NextData)
	ge.InitData = init

	//fmt.Printf("Inside UpdateInitData, init AFTER: %+v\n", ge.InitData)
	//fmt.Printf("Inside UpdateInitData, Next AFTER: %+v\n", ge.NextData)
}

func (ge *GravityEngine) InitializeNextData() {
	init := make([]Ball, len(ge.InitData))
	copy(init[0:], ge.InitData)
	ge.NextData = init
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func SetUp() (ge *GravityEngine) {
	initData := make([]Ball, 0)
	nextData := make([]Ball, 0)
	b0 := Ball{
		100,
		"green",
		Vector{X: 200, Y: 200},
		Vector{X: 0, Y: 0},
		Vector{X: 0, Y: 0},
		2,
	}
	b1 := Ball{
		100,
		"green",
		Vector{X: 200, Y: 250},
		Vector{X: 0, Y: 0},
		Vector{X: 0, Y: 0},
		3,
	}
	// Set both init and next data as init
	initData = append(initData, b0)
	initData = append(initData, b1)
	nextData = append(nextData, b0)
	nextData = append(nextData, b1)

	ge = &GravityEngine{
		G:             100,
		DeltaT:        .01,
		NumIterations: 500,
		InitData:      initData,
		NextData:      nextData,
		K:             .5,
		Trajectories: 	make([][]Vector, len(initData)),
	}
	return ge
}

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}


func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}

	// NOw I have a new connection - start pumping out data
	// Create a gravity engine, make a method for taking in the websocket and
	// let it push data out.
	ge := SetUp()
	go reader(ws)
	go ge.PumpData(ws)
}

func fakeFunc() {
	for {
		time.Sleep(2000 * time.Millisecond)
		fmt.Println("Tick")
	}
}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func SetupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	fmt.Println("Hello World")
	SetupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
	//http.HandleFunc("/hello", hello)
	//http.HandleFunc("/headers", headers)
	//http.ListenAndServe(":8090", nil)
	//engine := SetUp()
	//engine.moveBalls()
	// Do some number of moveBalls
	// each time moveBalls is called, do something with the data. Send it?
}






type Vector struct {
	X float64
	Y float64
}

// Calculate the unit vector
func Normalize(v Vector) (u Vector) {
	// if length is zero, set x = 1, y = 0
	vLen := Length(v)
	if vLen == 0 {
		u = Vector{
			1,
			0,
		}
		return u
	}
	x := v.X / vLen
	y := v.Y / vLen
	u = Vector{
		x,
		y,
	}
	u.X = v.X / vLen
	u.Y = v.Y / vLen
	return u
}

func Subtract(v1 Vector, v2 Vector) (res Vector) {
	x := v2.X - v1.X
	y := v2.Y - v1.Y
	res = Vector{
		X: x,
		Y: y,
	}
	return res
}

func Length(v Vector) (s float64) {
	s = math.Sqrt(LengthSq(v))
	return s
}

func LengthSq(v Vector) (s2 float64) {
	s2 = v.X*v.X + v.Y*v.Y
	return s2
}

func DistSq(v1 Vector, v2 Vector) (s2 float64) {
	s2 = math.Pow(v1.X-v2.X, 2) + math.Pow(v1.Y-v2.Y, 2)
	return s2
}

func VectorDist(v1 Vector, v2 Vector) (s float64) {
	s = math.Sqrt(DistSq(v1, v2))
	return s
}

// Multiply a vector by a scalar and return result
func ScalarMult(v Vector, s float64) (res Vector) {
	res.X = v.X * s
	res.Y = v.Y * s
	return res
}

func ScalarDivide(v Vector, s float64) (res Vector) {
	res.X = v.X / s
	res.Y = v.Y / s
	return res
}

func Add(v1 Vector, v2 Vector) (f Vector) {
	f.X = v1.X + v2.X
	f.Y = v1.Y + v2.Y
	return f
}
