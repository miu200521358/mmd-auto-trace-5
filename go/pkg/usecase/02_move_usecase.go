package usecase

import (
	"math"
	"strings"

	"github.com/miu200521358/mmd-auto-trace-5/pkg/config/mlog"
	"github.com/miu200521358/mmd-auto-trace-5/pkg/domain/mjson"
	"github.com/miu200521358/mmd-auto-trace-5/pkg/domain/vmd"

	"github.com/miu200521358/mmd-auto-trace-5/pkg/utils"
)

const RATIO = 1

func Move(frames *mjson.Frames, motionNum, allNum int) *vmd.VmdMotion {
	mlog.I("[%d/%d] Convert Move ...", motionNum, allNum)

	minFno := getMinFrame(frames.Frames)
	minFrame := frames.Frames[minFno]
	rootPos := mjson.Position{X: minFrame.Camera.X, Y: minFrame.Camera.Y, Z: minFrame.Camera.Z}

	bar := utils.NewProgressBar(len(frames.Frames))

	movMotion := vmd.NewVmdMotion(strings.Replace(frames.Path, "_smooth.json", "_move.vmd", -1))

	for fno, frame := range frames.Frames {
		bar.Increment()

		if frame.Confidential < 0.8 {
			continue
		}

		for jointName, pos := range frame.Joint3D {
			// ボーン名がある場合、ボーン移動モーションにも出力
			if boneName, ok := joint2bones[string(jointName)]; ok {
				bf := vmd.NewBoneFrame(float32(fno))
				bf.Position.X = pos.X
				bf.Position.Y = pos.Y * RATIO
				bf.Position.Z = pos.Z * RATIO
				movMotion.AppendBoneFrame(boneName, bf)
			}
		}

		{
			bf := vmd.NewBoneFrame(float32(fno))
			bf.Position.X = frame.Camera.X * RATIO
			bf.Position.Y = frame.Camera.Y * RATIO
			bf.Position.Z = (frame.Camera.Z - rootPos.Z) * 0.5
			movMotion.AppendBoneFrame("Camera", bf)
		}
	}

	bar.Finish()

	return movMotion
}

func getMinFrame(m map[int]mjson.Frame) int {
	minFrame := math.MaxInt

	for frame := range m {
		if frame < minFrame {
			minFrame = frame
		}
	}

	return minFrame
}

var joint2bones = map[string]string{
	"nose":            "鼻",
	"neck":            "首",
	"right_collar":    "右肩",
	"right_shoulder":  "右腕",
	"right_elbow":     "右ひじ",
	"right_wrist":     "右手首",
	"left_collar":     "左肩",
	"left_shoulder":   "左腕",
	"left_elbow":      "左ひじ",
	"left_wrist":      "左手首",
	"spine1":          "下半身",
	"right_hip":       "右足",
	"right_knee":      "右ひざ",
	"right_ankle":     "右足首",
	"left_hip":        "左足",
	"left_knee":       "左ひざ",
	"left_ankle":      "左足首",
	"right_eye":       "右目",
	"left_eye":        "左目",
	"right_ear":       "右耳",
	"left_ear":        "左耳",
	"left_big_toe":    "左つま先親",
	"left_small_toe":  "左つま先子",
	"left_heel":       "左かかと",
	"right_big_toe":   "右つま先親",
	"right_small_toe": "右つま先子",
	"right_heel":      "右かかと",
	"spine2":          "上半身",
	"spine3":          "上半身2",
	"head":            "頭",
	"pelvis":          "下半身先",
}
