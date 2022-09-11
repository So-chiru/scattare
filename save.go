package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var stack = make([]Message, 0)

func stack_add(msg *Message) error {
	if msg == nil {
		return nil
	}

	stack = append(stack, *msg)

	return nil
}

func validate_file() {
	if _, err := os.Stat(*OUPUT_FILE); os.IsNotExist(err) {
		file, err := os.Create(*OUPUT_FILE)
		if err != nil {
			fmt.Println(err)
			return
		}

		file.Write([]byte("[]"))
		file.Close()
	}
}

func get_file() (*os.File, error) {
	var file, err = os.OpenFile(*OUPUT_FILE, os.O_RDWR, 0755)
	if err != nil {
		return nil, err
	}

	return file, nil
}

var comma = []byte(",")

func json_save(msgs []Message) {
	validate_file()

	file, err := get_file()
	if err != nil {
		fmt.Println(err)
		return
	}

	o, err := file.Seek(-1, io.SeekEnd)
	if err != nil {
		fmt.Println(err)
		return
	}

	var b = make([]byte, 1)
	_, err = file.ReadAt(b, o)

	var is_empty = false

	if err != nil {
		fmt.Println(err)
		return
	}

	// if file is empty, mark to remove the closing bracket
	if string(b) == "]" {
		o, err = file.Seek(-1, io.SeekEnd)
		if err != nil {
			fmt.Println(err)
			return
		}

		stat, err := file.Stat()
		if err != nil {
			fmt.Println(err)
			return
		}

		var size = stat.Size()

		if size > 2 {
			file.WriteAt(comma, o)
		} else if size == 2 {
			is_empty = true
		}
	} else {
		fmt.Println("file is not a valid json array")
		return
	}

	var buf []byte = make([]byte, 0)

	// add new messages
	for index, msg := range msgs {
		m, err := json.Marshal(msg)
		if err != nil {
			fmt.Println(err)
			return
		}

		buf = append(buf, m...)

		if index < len(msgs)-1 {
			buf = append(buf, comma...)
		}
	}

	var rel_offset int64 = 0

	if is_empty {
		rel_offset = -1
	}

	// write new messages at the end of the file
	r, err := file.Seek(rel_offset, io.SeekEnd)
	if err != nil {
		fmt.Println(err)
		return
	}

	buf = append(buf, []byte("]")...)

	_, err = file.WriteAt(buf, r)
	if err != nil {
		panic(err)
	}

	buf = nil
	defer file.Close()
}

func save_worker() {
	if len(stack) == 0 {
		return
	}

	var start = time.Now()

	// copy stack into local variable
	var local = make([]Message, len(stack))
	copy(local, stack)

	// clear stack
	stack = make([]Message, 0)

	json_save(local)

	if *DEBUG_MODE {
		log.Printf("save_worker for %d messages took %s", len(local), time.Since(start))
	}

	return
}
