package utils

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"runtime/debug"
	"unsafe"
)

// 对外使用Batch()
// 涉及事务暂时不用这个

const batchSize = 5000

type batchFunc func(batch interface{}) (interface{}, error)

type batchFuncContext func(ctx context.Context, batch interface{}) (interface{}, error)

type BatchRunner interface {
	// Run method to run batch task, total is of type slice.
	// Recommend to use RunContext
	Run(total interface{}, f batchFunc) error

	// RunContext method to run batch task, total is of type slice
	RunContext(total interface{}, f batchFuncContext) error

	/* below are methods could be called before Run to change default behaviour */

	// Size rewrite default batchSize
	Size(size int) BatchRunner

	// Dest set dest to accumulate batch result, dest must be pointer to slice or map
	Dest(dest interface{}) BatchRunner

	// Pipe turn on/off pipeline mode.
	// In pipeline mode, a goroutine is created per batch.
	// Pipe(0) turn off pipeline mode
	// Pipe(N) limit goroutine num to N (N > 0)
	// Pipe() not limit goroutine num
	// WARNING if context has db transaction, pipe will be omitted
	Pipe(...int) BatchRunner
}

const unlimited = -1

type batchRunner struct {
	ctx   context.Context
	dest  reflect.Value
	size  int
	err   error
	limit int
	save  bool
}

func Batch() BatchRunner {
	return BatchWithContext(context.Background())
}

func BatchWithContext(ctx context.Context) BatchRunner {
	runner := &batchRunner{
		ctx:   ctx,
		dest:  reflect.Value{},
		size:  batchSize,
		err:   nil,
		limit: 0,
		save:  false,
	}

	if runner.ctx == nil {
		runner.ctx = context.Background()
	}

	return runner
}

func (b *batchRunner) Size(size int) BatchRunner {
	if size > 0 {
		b.size = size
	}
	return b
}

func (b *batchRunner) Pipe(num ...int) BatchRunner {
	if len(num) == 0 {
		b.limit = unlimited
	} else if num[0] >= 0 {
		b.limit = num[0]
	} else {
		b.err = errors.New("pipe size must greater than or equal to 0")
	}
	return b
}

func (b *batchRunner) Dest(d interface{}) BatchRunner {
	dest := reflect.ValueOf(d)
	b.dest = reflect.Indirect(dest)
	if dest.Kind() == reflect.Ptr && (b.dest.Kind() == reflect.Slice || b.dest.Kind() == reflect.Map) {
		b.save = true
	} else {
		b.err = fmt.Errorf("results should ptr to slice or map, not %s", b.dest.Kind())
	}
	return b
}

func (b *batchRunner) Run(total interface{}, f batchFunc) (err error) {
	wf := func(ctx context.Context, batch interface{}) (interface{}, error) {
		return f(batch)
	}

	return b.RunContext(total, wf)
}

func (b *batchRunner) RunContext(total interface{}, f batchFuncContext) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v, stack: %v", r, string(debug.Stack()))
		}
	}()

	if b.err != nil {
		return b.err
	}

	src := reflect.ValueOf(total)
	switch k := src.Kind(); k {
	case reflect.Slice, reflect.String, reflect.Array: // supported types
	default:
		return fmt.Errorf("batch operation not support for %s", k)
	}

	if b.save && b.dest.Len() > 0 {
		b.dest.Set(reflect.Zero(b.dest.Type()))
	}

	if b.limit != 0 {
		if hasTransactionDB(b.ctx) {
			err = b.serialRun(src, f)
		} else {
			err = b.pipelineRun(src, f)
		}
	} else {
		err = b.serialRun(src, f)
	}

	if b.save && err != nil {
		b.dest.Set(reflect.Zero(b.dest.Type()))
	}

	return
}

// ============= Inner function =============

func ptrOf(value reflect.Value) reflect.Value {
	if value.CanAddr() {
		return reflect.NewAt(value.Type(), unsafe.Pointer(value.Addr().Pointer()))
	}
	p := reflect.New(value.Type())
	p.Elem().Set(value)
	return p
}

func (b *batchRunner) serialRun(total reflect.Value, f batchFuncContext) error {
	size := b.size
	Len := total.Len()

	for cur, end := 0, size; cur < Len; cur, end = cur+size, end+size {
		if end > Len {
			end = Len
		}
		ret, err := f(b.ctx, total.Slice(cur, end).Interface())
		if err != nil {
			return err
		}
		if b.save && ret != nil {
			if err := appendDest(b.dest, ret); err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *batchRunner) pipelineRun(total reflect.Value, f batchFuncContext) (err error) {
	var (
		size = b.size
		Len  = total.Len()
		num  = (Len + size - 1) / size
	)

	var (
		done   = make(chan struct{})
		ch     = make(chan error, num)
		input  = make(chan interface{}, num)
		result chan interface{}
	)

	if b.save {
		result = make(chan interface{}, num)
	}

	limit := b.limit
	if limit == unlimited || limit > num {
		limit = num
	}
	for i := 0; i < limit; i++ {
		go func() {
			defer handlePanic(ch, result)
			for in := range input {
				select {
				case <-done:
				default:
					ret, err := f(b.ctx, in)
					ch <- err
					if b.save {
						result <- ret
					}
				}
			}
		}()
	}

	for cur, end := 0, size; cur < Len; cur, end = cur+size, end+size {
		if end > Len {
			end = Len
		}
		input <- total.Slice(cur, end).Interface()
	}
	close(input)

	for i := 0; i < num; i++ {
		if err = <-ch; err != nil {
			break
		}
		if b.save {
			if r := <-result; r != nil {
				if err = appendDest(b.dest, r); err != nil {
					break
				}
			}
		}
	}
	close(done)

	return
}

func appendDest(dest reflect.Value, val interface{}) error {
	destType := dest.Type()
	valType := reflect.TypeOf(val)
	valValue := reflect.ValueOf(val)

	// case 1: dest is map, val is the same type
	if destType.Kind() == reflect.Map {
		if valType != destType {
			return fmt.Errorf("type %T not same with dest", val)
		}
		iter := valValue.MapRange()
		for iter.Next() {
			dest.SetMapIndex(iter.Key(), iter.Value())
		}
		return nil
	}
	// case 2: dest is slice
	switch {
	case destType.Elem() == valType:
		// case 1: dest = append(dest, ret)
		dest.Set(reflect.Append(dest, reflect.ValueOf(val)))
	case destType == valType:
		// case 2: dest = append(dest, ret...)
		dest.Set(reflect.AppendSlice(dest, reflect.ValueOf(val)))
	case destType.Elem() == reflect.PtrTo(valType):
		// case 3: dest = append(dest, &ret)
		dest.Set(reflect.Append(dest, ptrOf(reflect.ValueOf(val))))
	case valType.Kind() == reflect.Slice && destType.Elem() == reflect.PtrTo(valType.Elem()):
		// case 4: dest = append(dest, &ret[0], &ret[1], ...)
		v := reflect.ValueOf(val)
		for i := 0; i < v.Len(); i++ {
			dest.Set(reflect.Append(dest, ptrOf(v.Index(i))))
		}
	default:
		return fmt.Errorf("type %T not supported", val)
	}

	return nil
}

func handlePanic(ch chan<- error, result chan<- interface{}) {
	if r := recover(); r != nil {
		ch <- fmt.Errorf("goroutine panic: %v, stack: %v", r, string(debug.Stack()))
		if result != nil {
			result <- nil
		}
	}
}

func hasTransactionDB(c context.Context) bool {
	_, ok := c.Value("TransactionDbInstance").(*gorm.DB)
	return ok
}
