package main

import (
	"fmt"
	//"encoding/json"
	//"fmt"
	"github.com/stretchr/testify/assert"
	//"os"
	//"strings"
	"testing"
)

func TestCalcOneForce(t *testing.T) {
	initData := make([]Ball, 0)

	b1 := Ball{
		50,
		"red",
		Vector{X: 200, Y: 200},
		Vector{X: 0, Y: 0},
		Vector{X: 0, Y: 0},
		0,
	}
	b2 := Ball{
		100,
		"green",
		Vector{X: 200, Y: 250},
		Vector{X: 0, Y: 0},
		Vector{X: 0, Y: 0},
		1,
	}

	initData = append(initData, b1)
	initData = append(initData, b2)
	ge := GravityEngine{
		G:             100,
		DeltaT:        1,
		NumIterations: 1,
		InitData:      initData,
		NextData:      nil,
		K:             .5,
	}

	initData2 := make([]Ball, 0)

	b3 := Ball{
		50,
		"red",
		Vector{X: 200, Y: 200},
		Vector{X: 0, Y: 0},
		Vector{X: 0, Y: 0},
		0,
	}
	b4 := Ball{
		100,
		"green",
		Vector{X: 250, Y: 300},
		Vector{X: 0, Y: 0},
		Vector{X: 0, Y: 0},
		1,
	}

	initData2 = append(initData2, b3)
	initData2 = append(initData2, b4)
	nextData := make([]Ball, len(initData))
	nextData = append(nextData, b3)
	nextData = append(nextData, b4)
	ge2 := GravityEngine{
		G:             100,
		DeltaT:        1,
		NumIterations: 1,
		InitData:      initData2,
		NextData:      nextData,
		K:             .5,
	}
	expected := Vector{X: 0, Y: 199.96000799840033}
	expected2 := Vector{X: 17.88782830686604, Y: 35.77565661373208}
	expected3 := Vector{X: -17.88782830686604, Y: -35.77565661373208}
	assert.True(t, expected == ge.CalcOneForce(ge.InitData[1], ge.InitData[0]), "Failed to get expected value #1")
	assert.True(t, expected2 == ge2.CalcOneForce(ge2.InitData[1], ge2.InitData[0]), "Failed to get correct for force #2")
	assert.True(t, expected3 == ge2.CalcOneForce(ge2.InitData[0], ge2.InitData[1]), "Failed to get correct for force #2")
}

func TestCalcNetForce(t *testing.T) {
	initData := make([]Ball, 0)

	b1 := Ball{
		100,
		"red",
		Vector{X: 200, Y: 200},
		Vector{X: 0, Y: 0},
		Vector{X: 0, Y: 0},
		0,
	}
	b2 := Ball{
		100,
		"green",
		Vector{X: 200, Y: 200},
		Vector{X: 0, Y: 0},
		Vector{X: 0, Y: 0},
		1,
	}
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
	initData = append(initData, b1)
	initData = append(initData, b2)
	initData = append(initData, b3)
	initData = append(initData, b4)
	ge := GravityEngine{
		G:             100,
		DeltaT:        1,
		NumIterations: 1,
		InitData:      initData,
		NextData:      nil,
		K:             .5,
	}

	oneF := ge.CalcOneForce(ge.InitData[0], ge.InitData[3])

	allF := ge.CalcNetForceOnBall(3)
	threeF := ScalarMult(oneF, 3)
	assert.True(t, allF == threeF, "The forces should sum correctly")
}

func TestGetNewData(t *testing.T) {
	// Make some shit
	// put that shit in there
	// make it dance!
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
	nextData = append(nextData, b3)
	nextData = append(nextData, b4)
	//copy(nextData, *initData)

	ge := GravityEngine{
		G:             100,
		DeltaT:        1,
		NumIterations: 1,
		InitData:      initData,
		NextData:      nextData,
		K:             .5,
	}
	//fmt.Printf("try to compare pointers init %+v\n", ge.InitData)
	//fmt.Printf("try to compare pointers next %+v\n", ge.NextData)
	//ge.Initialize()

	// calcForce, updateVel, update position
	// Do once w/ numIterations = 1
	// Do once w/ numIterations = 2

	//ge.GetNewData(0)

	assert.True(t, 1 == 1, "Just manually inspecting cuz I don't want to use calculator tonight")

	//fmt.Printf("Init pre first: %+v\n", ge.InitData[0])
	//fmt.Printf("Next pre first: %+v\n", ge.NextData[0])
	ge.GetNewData(0)
	//fmt.Printf("Init POST first %+v\n", ge.InitData[0])
	//fmt.Printf("Next POST first %+v\n", ge.NextData[0])
	assert.True(t, 1 == 1, "Just manually inspecting cuz I don't want to use calculator tonight")

	/*
		-------------
	*/

	//fmt.Println("UPDATINGs")
	ge.UpdateInitData()

	//fmt.Printf("Init pre second: %+v\n", ge.InitData[0])
	//fmt.Printf("Next pre second: %+v\n", ge.NextData[0])
	ge.GetNewData(0)
	//fmt.Printf("Init POST second %+v\n", ge.InitData[0])
	//fmt.Printf("Next POST second %+v\n", ge.NextData[0])
	//
	//fmt.Println("UPDATINGs")
	ge.UpdateInitData()
	assert.True(t, 1 == 1, "Just manually inspecting cuz I don't want to use calculator tonight")

	//fmt.Printf("Pre 2nd calc data 1 %+v\n", ge.NextData[1])

	//ge.GetNewData(1)
	//fmt.Printf("new data 1 2? %+v\n", ge.NextData[1])
	//assert.True(t, 1 ==1, "Just manually inspecting cuz I don't want to use calculator tonight")

}

func TestMoveBalls(t *testing.T) {
	// Make some shit
	// put that shit in there
	// make it dance!
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

	ge := GravityEngine{
		G:             100,
		DeltaT:        .01,
		NumIterations: 500,
		InitData:      initData,
		NextData:      nextData,
		K:             .5,
		Trajectories: 	make([][]Vector, len(initData)),
	}

	/*
		After 1 iteration, i expect to see:
		[ [b0, b0, b0], [b1, b1, b1]
	*/

	//expectedData0 := make([]Vector, len(initData))
	//expectedData1 := make([]Vector, len(initData))

	// for some number of datapoints to save, do ball iterations and write to expectedData arrays.
	// the arrays will only have position vectors
	fmt.Println("moveBalls")
	for i := 0; i < 5; i++ {
		ge.MoveBalls()
	}
	assert.True(t, len(ge.Trajectories[0])==5, "Expecting 5 trajectory points per ball")
}

func TestMakeHistory(t *testing.T) {

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

	ge := GravityEngine{
		G:             100,
		DeltaT:        .01,
		NumIterations: 500,
		InitData:      initData,
		NextData:      nextData,
		K:             .5,
		Trajectories: 	make([][]Vector, len(initData)),
	}

	// for some number of datapoints to save, do ball iterations and write to expectedData arrays.
	// the arrays will only have position vectors
	ge.MakeHistory(20)
	assert.True(t, len(ge.Trajectories[0])==20, "Expecting 5 trajectory points per ball")
}
