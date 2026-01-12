package usecase

import (
	"strings"

	"github.com/miu200521358/mmd-auto-trace-5/pkg/config/mlog"
	"github.com/miu200521358/mmd-auto-trace-5/pkg/domain/mjson"
	"github.com/miu200521358/mmd-auto-trace-5/pkg/domain/mmath"
	"github.com/miu200521358/mmd-auto-trace-5/pkg/domain/vmd"

	"github.com/miu200521358/mmd-auto-trace-5/pkg/utils"
)

const SCALE = 0.1259496 * 100

func Move(frames *mjson.Frames, motionNum, allNum int, minY, maxZ float64) *vmd.VmdMotion {
	mlog.I("[%d/%d] Convert Move ...", motionNum, allNum)

	bar := utils.NewProgressBar(len(frames.Frames))

	movMotion := vmd.NewVmdMotion(strings.Replace(frames.Path, ".json", "_move.vmd", -1))

	for fno, frame := range frames.Frames {
		bar.Increment()

		for jointName, pos := range frame.Joint3D {
			// ボーン名がある場合、ボーン移動モーションにも出力
			if boneName, ok := joint2bones[string(jointName)]; ok {
				bf := vmd.NewBoneFrame(float32(fno))
				bf.Position = &mmath.MVec3{X: pos.X, Y: -pos.Y - minY, Z: pos.Z - maxZ}
				bf.Position.MulScalar(SCALE)
				movMotion.AppendBoneFrame(boneName, bf)
			}
		}

		// 追加で計算するボーン
		{
			bf := vmd.NewBoneFrame(float32(fno))
			bf.Position = movMotion.BoneFrames.Get("右足").Get(float32(fno)).Position.Added(movMotion.BoneFrames.Get("左足").Get(float32(fno)).Position).DivedScalar(2)
			movMotion.AppendBoneFrame("下半身先", bf)
		}
		{
			bf := vmd.NewBoneFrame(float32(fno))
			bf.Position = movMotion.BoneFrames.Get("上半身").Get(float32(fno)).Position.Copy()
			movMotion.AppendBoneFrame("下半身", bf)
		}
	}

	bar.Finish()

	return movMotion
}

var joint2bones = map[string]string{
	"pelvis":          "上半身",
	"spine2":          "上半身2",
	"spine3":          "上半身3",
	"neck":            "首",
	"head":            "頭",
	"right_collar":    "右肩",
	"right_shoulder":  "右腕",
	"right_elbow":     "右ひじ",
	"right_wrist":     "右手首",
	"left_collar":     "左肩",
	"left_shoulder":   "左腕",
	"left_elbow":      "左ひじ",
	"left_wrist":      "左手首",
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
}
