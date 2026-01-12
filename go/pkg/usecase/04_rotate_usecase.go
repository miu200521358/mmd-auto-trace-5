package usecase

import (
	"strings"

	"github.com/miu200521358/mmd-auto-trace-5/pkg/config/mlog"
	"github.com/miu200521358/mmd-auto-trace-5/pkg/domain/mmath"
	"github.com/miu200521358/mmd-auto-trace-5/pkg/domain/pmx"
	"github.com/miu200521358/mmd-auto-trace-5/pkg/domain/vmd"
	"github.com/miu200521358/mmd-auto-trace-5/pkg/infrastructure/repository"
	"github.com/miu200521358/mmd-auto-trace-5/pkg/utils"
)

func Rotate(moveMotion *vmd.VmdMotion, modelPath string, motionNum, allNum int) *vmd.VmdMotion {
	mlog.I("[%d/%d] Convert Rotate ...", motionNum, allNum)

	// モデル読み込み
	pr := repository.NewPmxRepository(false)
	data, err := pr.Load(modelPath)
	if err != nil {
		mlog.E("Failed to read pmx: %v", err)
	}
	pmxModel := data.(*pmx.PmxModel)

	bar := utils.NewProgressBar(len(boneConfigs))

	rotMotion := vmd.NewVmdMotion(strings.Replace(moveMotion.Path(), "_move.vmd", "_rotate.vmd", -1))

	moveMotion.BoneFrames.Get("下半身先").ForEach(func(fno float32, bf *vmd.BoneFrame) bool {
		centerBf := vmd.NewBoneFrame(fno)
		centerBf.Position = bf.Position.Copy()
		rotMotion.AppendBoneFrame(pmx.CENTER.String(), centerBf)
		return true
	})

	for _, boneConfig := range boneConfigs {
		bar.Increment()

		if !moveMotion.BoneFrames.Contains(boneConfig.Name) || !moveMotion.BoneFrames.Contains(boneConfig.DirectionFrom) ||
			!moveMotion.BoneFrames.Contains(boneConfig.DirectionTo) || !moveMotion.BoneFrames.Contains(boneConfig.UpFrom) ||
			!moveMotion.BoneFrames.Contains(boneConfig.UpTo) {
			continue
		}

		moveMotion.BoneFrames.Get(boneConfig.Name).ForEach(func(fno float32, bf *vmd.BoneFrame) bool {
			// モデルのボーン角度
			boneDirectionFromBone, _ := pmxModel.Bones.GetByName(boneConfig.DirectionFrom)
			boneDirectionFrom := boneDirectionFromBone.Position
			boneDirectionToBone, _ := pmxModel.Bones.GetByName(boneConfig.DirectionTo)
			boneDirectionTo := boneDirectionToBone.Position
			boneUpFromBone, _ := pmxModel.Bones.GetByName(boneConfig.UpFrom)
			boneUpFrom := boneUpFromBone.Position
			boneUpToBone, _ := pmxModel.Bones.GetByName(boneConfig.UpTo)
			boneUpTo := boneUpToBone.Position

			boneDirectionVector := boneDirectionTo.Subed(boneDirectionFrom).Normalize()
			boneUpVector := boneUpTo.Subed(boneUpFrom).Normalize()
			boneCrossVector := boneUpVector.Cross(boneDirectionVector).Normalize()

			boneQuat := mmath.NewMQuaternionFromDirection(boneDirectionVector, boneCrossVector)

			// モーションのボーン角度
			motionDirectionFromPos := moveMotion.BoneFrames.Get(boneConfig.DirectionFrom).Get(fno).Position
			motionDirectionToPos := moveMotion.BoneFrames.Get(boneConfig.DirectionTo).Get(fno).Position
			motionUpFromPos := moveMotion.BoneFrames.Get(boneConfig.UpFrom).Get(fno).Position
			motionUpToPos := moveMotion.BoneFrames.Get(boneConfig.UpTo).Get(fno).Position

			motionDirectionVector := motionDirectionToPos.Subed(motionDirectionFromPos).Normalize()
			motionUpVector := motionUpToPos.Subed(motionUpFromPos).Normalize()
			motionCrossVector := motionUpVector.Cross(motionDirectionVector).Normalize()

			motionQuat := mmath.NewMQuaternionFromDirection(motionDirectionVector, motionCrossVector)

			// キャンセルボーン角度
			cancelQuat := mmath.NewMQuaternion()
			for _, cancelBoneName := range boneConfig.Cancels {
				cancelQuat.Mul(rotMotion.BoneFrames.Get(cancelBoneName).Get(fno).Rotation)
			}

			// 調整角度
			invertQuat := mmath.NewMQuaternionFromDegrees(boneConfig.Invert.X, boneConfig.Invert.Y, boneConfig.Invert.Z)

			// ボーンフレーム登録
			rotBf := vmd.NewBoneFrame(fno)
			rotBf.Rotation = invertQuat.Mul(cancelQuat.Inverse()).Mul(motionQuat).Mul(boneQuat.Inverse()).Normalize()

			rotMotion.AppendBoneFrame(boneConfig.Name, rotBf)

			return true
		})
	}

	bar.Finish()

	return rotMotion
}

type boneConfig struct {
	Name          string
	DirectionFrom string
	DirectionTo   string
	UpFrom        string
	UpTo          string
	Cancels       []string
	Invert        *mmath.MVec3
}

var boneConfigs = []*boneConfig{
	{
		Name:          "下半身",
		DirectionFrom: "下半身",
		DirectionTo:   "下半身先",
		UpFrom:        "左足",
		UpTo:          "右足",
		Cancels:       []string{},
		Invert:        &mmath.MVec3{},
	},
	{
		Name:          "上半身",
		DirectionFrom: "上半身",
		DirectionTo:   "上半身2",
		UpFrom:        "左腕",
		UpTo:          "右腕",
		Cancels:       []string{},
		Invert:        &mmath.MVec3{},
	},
	{
		Name:          "上半身2",
		DirectionFrom: "上半身2",
		DirectionTo:   "首",
		UpFrom:        "左腕",
		UpTo:          "右腕",
		Cancels:       []string{"上半身"},
		Invert:        &mmath.MVec3{},
	},
	{
		Name:          "首",
		DirectionFrom: "首",
		DirectionTo:   "頭",
		UpFrom:        "左腕",
		UpTo:          "右腕",
		Cancels:       []string{"上半身", "上半身2"},
		Invert:        &mmath.MVec3{},
	},
	{
		Name:          "頭",
		DirectionFrom: "首",
		DirectionTo:   "頭",
		UpFrom:        "左目",
		UpTo:          "右目",
		Cancels:       []string{"上半身", "上半身2", "首"},
		Invert:        &mmath.MVec3{},
	},
	{
		Name:          "左肩",
		DirectionFrom: "左肩",
		DirectionTo:   "左腕",
		UpFrom:        "上半身2",
		UpTo:          "首",
		Cancels:       []string{"上半身", "上半身2"},
		Invert:        &mmath.MVec3{},
	},
	{
		Name:          "左腕",
		DirectionFrom: "左腕",
		DirectionTo:   "左ひじ",
		UpFrom:        "左腕",
		UpTo:          "右腕",
		Cancels:       []string{"上半身", "上半身2", "左肩"},
		Invert:        &mmath.MVec3{},
	},
	{
		Name:          "左ひじ",
		DirectionFrom: "左ひじ",
		DirectionTo:   "左手首",
		UpFrom:        "左腕",
		UpTo:          "左ひじ",
		Cancels:       []string{"上半身", "上半身2", "左肩", "左腕"},
		Invert:        &mmath.MVec3{},
	},
	{
		Name:          "左手首",
		DirectionFrom: "左手首",
		DirectionTo:   "左人指先",
		UpFrom:        "左親指１",
		UpTo:          "左小指１",
		Cancels:       []string{"上半身", "上半身2", "左肩", "左腕", "左ひじ"},
		Invert:        &mmath.MVec3{},
	},
	{
		Name:          "右肩",
		DirectionFrom: "右肩",
		DirectionTo:   "右腕",
		UpFrom:        "上半身2",
		UpTo:          "首",
		Cancels:       []string{"上半身", "上半身2"},
		Invert:        &mmath.MVec3{},
	},
	{
		Name:          "右腕",
		DirectionFrom: "右腕",
		DirectionTo:   "右ひじ",
		UpFrom:        "右肩",
		UpTo:          "右腕",
		Cancels:       []string{"上半身", "上半身2", "右肩"},
		Invert:        &mmath.MVec3{},
	},
	{
		Name:          "右ひじ",
		DirectionFrom: "右ひじ",
		DirectionTo:   "右手首",
		UpFrom:        "右腕",
		UpTo:          "右ひじ",
		Cancels:       []string{"上半身", "上半身2", "右肩", "右腕"},
		Invert:        &mmath.MVec3{},
	},
	{
		Name:          "右手首",
		DirectionFrom: "右手首",
		DirectionTo:   "右人指先",
		UpFrom:        "右親指１",
		UpTo:          "右小指１",
		Cancels:       []string{"上半身", "上半身2", "右肩", "右腕", "右ひじ"},
		Invert:        &mmath.MVec3{},
	},
	{
		Name:          "左足",
		DirectionFrom: "左足",
		DirectionTo:   "左ひざ",
		UpFrom:        "左足",
		UpTo:          "右足",
		Cancels:       []string{"下半身"},
		Invert:        &mmath.MVec3{},
	},
	{
		Name:          "左ひざ",
		DirectionFrom: "左ひざ",
		DirectionTo:   "左足首",
		UpFrom:        "左足",
		UpTo:          "右足",
		Cancels:       []string{"下半身", "左足"},
		Invert:        &mmath.MVec3{},
	},
	{
		Name:          "左足首",
		DirectionFrom: "左足首",
		DirectionTo:   "左つま先",
		UpFrom:        "左つま先親",
		UpTo:          "左つま先子",
		Cancels:       []string{"下半身", "左足", "左ひざ"},
		Invert:        &mmath.MVec3{},
	},
	{
		Name:          "右足",
		DirectionFrom: "右足",
		DirectionTo:   "右ひざ",
		UpFrom:        "左足",
		UpTo:          "右足",
		Cancels:       []string{"下半身"},
		Invert:        &mmath.MVec3{},
	},
	{
		Name:          "右ひざ",
		DirectionFrom: "右ひざ",
		DirectionTo:   "右足首",
		UpFrom:        "左足",
		UpTo:          "右足",
		Cancels:       []string{"下半身", "右足"},
		Invert:        &mmath.MVec3{},
	},
	{
		Name:          "右足首",
		DirectionFrom: "右足首",
		DirectionTo:   "右つま先",
		UpFrom:        "右つま先親",
		UpTo:          "右つま先子",
		Cancels:       []string{"下半身", "右足", "右ひざ"},
		Invert:        &mmath.MVec3{},
	},
}
