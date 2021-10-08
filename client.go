package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func get_str() string {
	li := ""
	for i := 0; i < 5; i++ {
		li = li + strconv.Itoa(rand.Intn(10))
	}
	return li
}

func get_key() string {
	b := ""
	for i := 0; i < 10; i++ {
		b = b + strconv.Itoa(rand.Intn(9) + 1)
	}
	return b
}

type Session_protector struct {
	__hash string
}

func (self Session_protector) __calc_hash(session_key string, val int) string {
	switch val {
	case 1:
		result := ""
		ret := ""
		for idx := 0; idx < 5; idx++ {
			result = result + string(session_key[idx])
		}
		i, _ := strconv.Atoi(result)
		result = "00" + strconv.Itoa(i%97)
		for idx := len(result) - 2; idx < len(result); idx++ {
			ret = ret + string(result[idx])
		}
		return ret
	case 2:
		result := ""
		for idx := 0; idx < len(session_key); idx++ {
			result = result + string(session_key[len(session_key)-idx-1])
		}
		return result
	case 3:
		return session_key[len(session_key)-5:] + session_key[0:5]
	case 4:
		num := 0
		for i := 1; i < 9; i++ {
			num = num + int(session_key[i]) + 41
		}
		return string(num)
	case 5:
		ch := ""
		result := 0
		for idx := 0; idx < len(session_key); idx++ {
			ch = string(int(int(session_key[idx]) ^ 43))
			if _, err := strconv.Atoi(ch); err != nil {
				ch = string(int(ch[0]))
			}
			num, _ := strconv.Atoi(ch)
			result = result + num
		}
		return strconv.Itoa(result)
	default:
		result, _ := strconv.Atoi(session_key)
		return strconv.Itoa(result + val)
	}
}

func (self Session_protector) next_session_key(session_key string) string {
	if self.__hash == "" {
		fmt.Println("Hash is empty")
		return get_key()
	}
	for idx := 0; idx < len(self.__hash); idx++ {
		i := string(self.__hash[idx])
		if _, err := strconv.Atoi(i); err != nil {
			fmt.Println("Letter!")
			return get_key()
		}
	}
	result := 0
	ret := ""
	for idx := 0; idx < len(self.__hash); idx++ {
		num, _ := strconv.Atoi(string(self.__hash[idx]))
		k, _ := strconv.Atoi(self.__calc_hash(session_key, num))
		result = result + k
	}
	for idx := 0; idx < 10 && idx < len(strconv.Itoa(result)); idx++ {
		ret = ret + string((strconv.Itoa(result))[idx])
	}
	x := ""
	ret = "0000000000" + ret
	for idx := len(ret) - 10; idx < len(ret); idx++ {
		x = x + string(ret[idx])
	}
	return x
}

func main() {
	fmt.Print("Enter IP:PORT ")
	IPAdr := ""
	fmt.Fscan(os.Stdin, &IPAdr)
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	conn, err := net.Dial("tcp", IPAdr)
	if err != nil {
		fmt.Println("Server not found...")
	} else {
		cl_hash_string := get_str()
		key := get_key()
		fmt.Print(cl_hash_string + "\n")
		fmt.Fprintf(conn, cl_hash_string+key+"\n")
		client_protector := Session_protector{cl_hash_string}
		keyy, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Server is not responding.")
		}
		key = client_protector.next_session_key(key)
		for {
			text := ""
			fmt.Fprintf(conn, strings.Replace(text, "\n", "", -1)+key+"\n")
			fmt.Println("Waiting for answer...")
			message, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println("Server is not responding.")
			}
			keyy = ""
			for i := len(message) - 11; i < len(message)-1; i++ {
				keyy = keyy + string(message[i])
			}
			for i := 0; i < len(message)-11; i++ {
				text = text + string(message[i])
			}
			key = client_protector.next_session_key(key)
			fmt.Println("key: ", key, " ", keyy)
			key = client_protector.next_session_key(key)
		}
	}
}
