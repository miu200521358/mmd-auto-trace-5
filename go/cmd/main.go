package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/miu200521358/mmd-auto-trace-5/pkg/config/mlog"
	"github.com/miu200521358/mmd-auto-trace-5/pkg/usecase"
	"github.com/miu200521358/mmd-auto-trace-5/pkg/utils"
)

var logLevel string
var modelPath string
var dirPath string

func init() {
	flag.StringVar(&logLevel, "logLevel", "INFO", "set log level")
	flag.StringVar(&modelPath, "modelPath", "", "set model path")
	flag.StringVar(&dirPath, "dirPath", "", "set directory path")
	flag.Parse()

	switch logLevel {
	case "INFO":
		mlog.SetLevel(mlog.INFO)
	default:
		mlog.SetLevel(mlog.DEBUG)
	}
}

func main() {
	if modelPath == "" || dirPath == "" {
		err := fmt.Errorf("modelPath and dirPath must be provided")
		mlog.E("%v", err)
		os.Exit(1)
	}

	mlog.I("Unpack json ================")
	allFrames, err := usecase.Unpack(dirPath)
	if err != nil {
		mlog.E("Failed to unpack: %v", err)
		return
	}

	allNum := len(allFrames)

	mlog.I("[%d] Calculation Center Z ===========================", allNum)

	minY, maxZ := usecase.CalcMinYZ(allFrames)

	for i, frames := range allFrames {
		motionNum := i + 1

		// if _, err := os.Stat(filepath.Join(filepath.Dir(frames.Path), utils.GetCompleteName(frames.Path))); err == nil {
		// 	mlog.I("[%d/%d] Finished Convert Motion ===========================", motionNum, allNum)
		// 	continue
		// }

		mlog.I("[%d/%d] Convert Motion ===========================", motionNum, allNum)

		moveMotion := usecase.Move(frames, motionNum, allNum, minY, maxZ)

		if mlog.IsDebug() {
			utils.WriteVmdMotions(frames, moveMotion, dirPath, "_1move", "Move", motionNum, allNum)
		}

		rotateMotion := usecase.Rotate(moveMotion, modelPath, motionNum, allNum)

		if mlog.IsDebug() {
			utils.WriteVmdMotions(frames, rotateMotion, dirPath, "_2rotate", "Rotate", motionNum, allNum)
		}

		// legIkMotion := usecase.ConvertLegIk(rotateMotion, modelPath, motionNum, allNum)

		// if mlog.IsDebug() {
		// 	utils.WriteVmdMotions(frames, legIkMotion, dirPath, "3_legIk", "LegIK", motionNum, allNum)
		// }

		// groundMotion := usecase.FixGround(legIkMotion, modelPath, motionNum, allNum)

		// if mlog.IsDebug() {
		// 	utils.WriteVmdMotions(frames, groundMotion, dirPath, "4_ground", "Ground", motionNum, allNum)
		// }

		// heelMotion := usecase.FixHeel(frames, groundMotion, modelPath, motionNum, allNum)

		// if mlog.IsDebug() {
		// 	utils.WriteVmdMotions(frames, heelMotion, dirPath, "5_heel", "Heel", motionNum, allNum)
		// }

		// armIkMotion := usecase.ConvertArmIk(heelMotion, modelPath, motionNum, allNum)

		// utils.WriteVmdMotions(frames, armIkMotion, dirPath, "full", "Full", motionNum, allNum)

		// narrowReduceMotion := usecase.Reduce(armIkMotion, modelPath, 0.05, 0.00001, 0, "narrow", motionNum, allNum)

		// utils.WriteVmdMotions(frames, narrowReduceMotion, dirPath, "reduce_narrow", "Narrow Reduce", motionNum, allNum)

		// wideReduceMotions := usecase.Reduce(armIkMotion, modelPath, 0.07, 0.00005, 2, "wide", motionNum, allNum)

		// utils.WriteVmdMotions(frames, wideReduceMotions, dirPath, "reduce_wide", "Wide Reduce", motionNum, allNum)

		// utils.WriteComplete(dirPath, frames.Path)
	}

	mlog.I("Done!")
}
