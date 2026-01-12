package usecase

import "github.com/miu200521358/mmd-auto-trace-5/pkg/domain/mjson"

func CalcMinYZ(allFrames []*mjson.Frames) (float64, float64) {
	// 最も早いフレームの最も手前の"pelvis"のZ座標を取得
	minFrame := 0

	for _, frames := range allFrames {
		if len(frames.Frames) == 0 {
			continue
		}

		for fno := range frames.Frames {
			if minFrame == 0 || fno < minFrame {
				minFrame = fno
			}
		}
	}

	minY := 0.0
	maxZ := 0.0

	for _, frames := range allFrames {
		if len(frames.Frames) == 0 {
			continue
		}

		minFrameData, ok := frames.Frames[minFrame]
		if !ok {
			continue
		}

		if pos, ok := minFrameData.GlobalJoint3D["pelvis"]; ok {
			if pos.Y < minY {
				minY = pos.Y
			}

			if pos.Z > maxZ {
				maxZ = pos.Z
			}
		}
	}

	return minY, maxZ
}
