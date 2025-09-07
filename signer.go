package main

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
)

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})
	wg := new(sync.WaitGroup)

	wg.Add(len(jobs))

	for _, someJob := range jobs {
		out := make(chan interface{})

		go func(in, out chan interface{}, job job) {
			defer wg.Done()
			defer close(out)
			job(in, out)
		}(in, out, someJob)

		in = out
	}

	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	wg := new(sync.WaitGroup)
	mut := new(sync.Mutex)

	for v := range in {
		wg.Add(1)
		data := strconv.Itoa(v.(int))
		go func(data string) {
			defer wg.Done()

			mut.Lock()
			signerMd5 := DataSignerMd5(data)
			mut.Unlock()

			hash := DataSignerCrc32(data) + "~" + DataSignerCrc32(signerMd5)

			fmt.Println(hash)
			out <- hash
		}(data)

	}

	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	wg := new(sync.WaitGroup)

	for v := range in {
		wg.Add(1)
		go func(data string) {
			defer wg.Done()

			var hash [6]string
			var innerWg sync.WaitGroup

			// Вычисляем каждый хэш в отдельной горутине
			for i := 0; i < 6; i++ {
				innerWg.Add(1)
				go func(i int) {
					defer innerWg.Done()
					hash[i] = DataSignerCrc32(strconv.Itoa(i) + data)
				}(i)
			}

			// Ждем завершения всех горутин для текущего значения
			innerWg.Wait()

			// Собираем результат
			result := ""
			for _, h := range hash {
				result += h
			}

			out <- result
		}(v.(string))
	}

	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	var result string
	var results []string

	for v := range in {
		results = append(results, v.(string))
	}

	sort.Strings(results)

	for i, res := range results {
		if i == 0 {
			result += res
		} else {
			result += "_" + res
		}
	}

	out <- result
}
