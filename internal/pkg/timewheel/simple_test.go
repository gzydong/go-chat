package timewheel

// func TestNewSimpleTimeWheel(t *testing.T) {
//
// 	obj := NewSimpleTimeWheel(1*time.Second, 10, func(wheel *SimpleTimeWheel, key string, value any) {
// 		fmt.Println(key, value)
// 	})
//
// 	go obj.Start()
//
// 	for i := 0; i < 30; i++ {
// 		for i := 0; i < 100000; i++ {
// 			m := strutil.NewMsgId()
// 			index := i
// 			obj.Add(m, index, 1*time.Second)
// 		}
//
// 		time.Sleep(1 * time.Second)
// 	}
//
// 	time.Sleep(1 * time.Hour)
// }
