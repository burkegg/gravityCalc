package main

import (
	"fmt"
	"net/http"
)

type Ball struct {
	Mass  float64
	Color string
	Pos   Vector
	Vel   Vector
	F 	  Vector
	Num	  int`json:",omitempty"`
}

type GravityEngine struct {
	G             float64
	DeltaT        float64
	NumIterations int
	InitData      []Ball
	NextData      []Ball
	K 			  float64
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

func (ge *GravityEngine) moveBalls() {

	for i := 0; i < ge.NumIterations; i++ {
		for ballNum := 0; ballNum < len(ge.InitData); ballNum++ {
			fmt.Println("get new data")
			ge.GetNewData(ballNum)
			if ballNum == 0 {
				fmt.Printf("Got new Data ballNum %v and data %+v\n", ballNum, ge.NextData[ballNum])
			}
		}
		//fmt.Printf("Init[0] before update %+v\n", ge.InitData[0])
		//fmt.Printf("Next[0] before update %+v\n", ge.NextData[0])
		ge.UpdateInitData()
	}
}


func (ge *GravityEngine) GetNewData(ballNum int) {
	// get net force for ball
	// apply kinematics
	// update nextData
	// To not overwrite the init array apparently I have to make a copy
	init := make([]Ball, len(ge.InitData))
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

func SetUp() (engine *GravityEngine){
	initData := make([]Ball, 0)
	nextData := make([]Ball, 0)
	b3 := Ball{
		100,
		"green",
		Vector{X: 200, Y: 200},
		Vector{X: 0, Y: 0},
		Vector{X: 0, Y: 0},
		2,
	}
	b4 := Ball{
		100,
		"green",
		Vector{X: 200, Y: 250},
		Vector{X: 0, Y: 0},
		Vector{X: 0, Y: 0},
		3,
	}

	initData = append(initData, b3)
	initData = append(initData, b4)

	engine = &GravityEngine{
		G:             100,
		DeltaT:        1,
		NumIterations: 10,
		InitData:      initData,
		NextData:      nextData,
		K:             .5,
	}
	engine.InitializeNextData()
	return engine
}

func main() {

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.ListenAndServe(":8090", nil)
	engine := SetUp()
	engine.moveBalls()
	// Do some number of moveBalls
	// each time moveBalls is called, do something with the data. Send it?
}
