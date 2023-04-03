package job

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func testFunc(ctx context.Context, t int) (uuid.UUID, error) {
	if t <= 0 {
		return uuid.UUID{}, fmt.Errorf("%d <= 0", t)
	}
	var i int
	for i = 0; i < t; i++ {
		if ctx.Err() != nil {
			return uuid.UUID{}, ctx.Err()
		}
		time.Sleep(time.Second)
	}
	return uuid.New(), nil
}

func cleanUpFunc() {
	time.Sleep(time.Second)
}

func printJson(v any) {
	bytes, _ := json.Marshal(v)
	fmt.Println(string(bytes))
}

//func TestName(t *testing.T) {
//	ctx, cf := context.WithCancel(context.Background())
//	cch := ccjh.New(10)
//	defer func() {
//		cf()
//		if cch.Active() > 0 {
//			fmt.Println("waiting for active jobs ...")
//			for cch.Active() != 0 {
//				time.Sleep(50 * time.Millisecond)
//			}
//			fmt.Println("jobs canceled")
//		}
//	}()
//	jh := New(ctx, cch)
//	err := cch.RunAsync(10, 50*time.Millisecond)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	ctx1, cf1 := context.WithCancel(jh.Context())
//	j1 := NewJob(ctx1, cf1)
//	j1.SetTarget(func() {
//		defer cf1()
//		r, e := testFunc(ctx1, 2)
//		if e == nil {
//			e = ctx1.Err()
//		}
//		j1.SetResult(r.String(), e)
//
//	})
//	id1, err := jh.Add(j1)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println(id1)
//
//	ctx2, cf2 := context.WithCancel(jh.Context())
//	j2 := NewJob(ctx2, cf2)
//	j2.SetTarget(func() {
//		defer cf2()
//		r, e := testFunc(ctx2, 4)
//		if e == nil {
//			e = ctx2.Err()
//		}
//		j2.SetResult(r.String(), e)
//
//	})
//	id2, err := jh.Add(j2)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println(id2)
//
//	time.Sleep(250 * time.Millisecond)
//	j2.Cancel()
//	cch.Stop()
//	var l1 []model.Job
//	jh.Range(func(k uuid.UUID, v *Job) bool {
//		l1 = append(l1, v.Meta())
//		return true
//	})
//	printJson(l1)
//	for cch.Active() != 0 {
//		time.Sleep(50 * time.Millisecond)
//	}
//	var l2 []model.Job
//	jh.Range(func(k uuid.UUID, v *Job) bool {
//		l2 = append(l2, v.Meta())
//		return true
//	})
//	printJson(l2)
//}
