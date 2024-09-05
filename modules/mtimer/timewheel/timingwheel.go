package timewheel

type TimingWheel struct {
	interval      int64 //刻度的时间间隔,最高精度是毫秒
	scales        int64 //刻度数
	upperInterval int64 //上一级时间轮的刻度间隔
	step          int64 //当前时间指针的指向
	lv            int64 //多维时间轮索引
	buckets       []*Bucket
	overflowWheel *TimingWheel
}

func NewTimingWheel(interval int64, lv int64, scales int64) *TimingWheel {
	tw := &TimingWheel{
		interval: interval,
		lv:       lv,
		scales:   scales,
		buckets:  make([]*Bucket, scales),
	}

	for i := int64(0); i < scales; i++ {
		tw.buckets[i] = newBucket(i)
	}

	return tw
}

func (tw *TimingWheel) Add(t *Timer) error {
	return tw.regTimer(t)
}

func (tw *TimingWheel) Update(t *Timer) {

}

func (tw *TimingWheel) Remove(t *Timer) {
	b := tw.buckets[t.bIndex]
	b.Del(t.id)
}

// regOverflowWheel 注册下一级时间轮
// register overflow wheel
func (tw *TimingWheel) regOverflowWheel(overflowWheel *TimingWheel) {
	tw.overflowWheel = overflowWheel
}

func (tw *TimingWheel) regTimer(t *Timer) error {
	step, pos, circle := tw.calcPositionAndCircle(t.delay)
	if circle > 0 {
		if tw.overflowWheel == nil {
			t.round = circle
			tw.setByBucket(t, pos, step)
			return nil
		} else {
			t.delay -= tw.calcDelayAdjustment()
			return tw.overflowWheel.regTimer(t)
		}
	}

	//将任务加入到当前时间轮
	tw.setByBucket(t, pos, step)
	return nil
}

func (tw *TimingWheel) tick(cmd Command, forward bool) {
	bucket := tw.buckets[tw.step]
	var diff int //heap 偏移量

	//取出当前时间轮指针指向的刻度上的所有定时器
	for {
		timer := bucket.Peek(diff)
		if timer == nil {
			break
		}

		if timer.round > 0 {
			diff++
			timer.round--
			continue
		}
		if cmd(timer) {
			bucket.Pop(diff)
		} else {
			diff++
		}
	}

	if forward {
		//指针前进一步
		tw.step = (tw.step + 1) % tw.scales
	}
}

// setBucket 设置定时器到指定的bucket
// setBucket sets the timer to the specified bucket
func (tw *TimingWheel) setByBucket(t *Timer, pos, step int64) {
	//计算定时器的延迟时间
	t.delay -= step * tw.interval
	t.setBucketIndex(pos)
	t.setTimeWheelIndex(tw.lv)
	tw.buckets[pos].Set(t)
}

// delay 单位 ms
func (tw *TimingWheel) calcPositionAndCircle(delay int64) (step, pos, circle int64) {
	step = delay / tw.interval
	circle = step / tw.scales
	pos = (tw.step + step) % tw.scales
	return
}

func (tw *TimingWheel) movingForward() bool {
	return (tw.step+1)%tw.scales == 0
}

// calcDelayAdjustment 计算延迟调整值
// calcDelayAdjustment calculates the delay adjustment value
func (tw *TimingWheel) calcDelayAdjustment() int64 {
	return (tw.scales - tw.step - 1) * tw.interval
}
