package main

import (
	"fmt"
	"strconv"
)

// сюда писать код

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})

	for _, someJob := range jobs {

		out := make(chan interface{})

		go someJob(in, out)

		in = out
	}

}

func SingleHash(in, out chan interface{}) {
	data := (<-in).(string)
	hash := DataSignerCrc32(data) + "~" + DataSignerCrc32(DataSignerMd5(data))
	fmt.Println(hash)
	return
}

func MultiHash(in, out chan interface{}) {
	var hash [6]string
	data := (<-in).(string)
	for i := 0; i < 6; i++ {
		hash[i] = DataSignerCrc32(strconv.Itoa(i) + data)
	}

	var result string
	for i := 0; i < 6; i++ {
		result += hash[i]
	}
	return
}

func CombineResults(in, out chan interface{}) {

	return
}
