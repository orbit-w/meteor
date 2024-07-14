package math

/*
   @Author: orbit-w
   @File: geometry
   @2023 11月 周二 11:27
*/

func Vector(a, b [2]float64) [2]float64 {
	return [2]float64{b[0] - a[0], b[1] - a[1]}
}

func CrossProduct(v, c [2]float64) float64 {
	return v[0]*c[1] - v[1]*c[0]
}

func SegmentIntersect(a, b, c, d [2]float64) bool {
	vab := Vector(a, b)
	vac := Vector(a, c)
	vad := Vector(a, d)
	vcd := Vector(c, d)
	vca := reverse(vac)
	vcb := Vector(c, b)
	cpAbAc := CrossProduct(vab, vac)
	cpAbAd := CrossProduct(vab, vad)
	cpA := cpAbAc * cpAbAd
	cpCdCa := CrossProduct(vcd, vca)
	cpCdCb := CrossProduct(vcd, vcb)
	cpB := cpCdCa * cpCdCb
	if cpA < 0 && cpB < 0 {
		return true
	} else if cpA == 0 && cpB == 0 {
		//共线
		if cpAbAc == 0 && cpAbAd == 0 {
			var i int32
			//a, y轴平行
			if vab[0] == 0 {
				i = 1
			}
			minA, maxA := sortFloat64(a[i], b[i])
			minB, maxB := sortFloat64(c[i], d[i])
			return !(maxA < minB || minA > maxB)
		} else {
			//两线段其中一端点重合
			return true
		}
	} else if cpA == 0 {
		return cpB < 0
	} else if cpB == 0 {
		return cpA < 0
	}
	return false
}

func reverse(t [2]float64) [2]float64 {
	return [2]float64{-t[0], -t[1]}
}

func sortFloat64(a, b float64) (min, max float64) {
	if a > b {
		min = b
		max = a
	} else {
		min = a
		max = b
	}
	return
}
