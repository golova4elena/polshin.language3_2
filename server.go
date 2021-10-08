package main

import (
    "net"
    "fmt"
    "bufio"
    "strings" 
    "strconv"
    "math/rand"
    "flag"
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
	if val == 1 {
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
	}
    if val == 2 {
		result := ""
		for idx := 0; idx < len(session_key); idx++ {
			result = result + string(session_key[len(session_key)-idx-1])
		}
		return result
	}	
	if val == 3 {
		return session_key[len(session_key)-5:] + session_key[0:5]
	}
	if val == 4 {
				num := 0
		for i := 1; i < 9; i++ {
			num = num + int(session_key[i]) + 41
		}
		return string(num)
	}
	if val == 5 {
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
	}
		result, _ := strconv.Atoi(session_key)
		return strconv.Itoa(result + val)
}


func (self Session_protector) next_session_key(session_key string) string {	
	if self.__hash == "" {
		fmt.Println("Hash is empty")
		return get_key()
	}
	for idx := 0; idx < len(self.__hash); idx++ {
		i := string(self.__hash[idx])
		if _, err := strconv.Atoi(i); err != nil {
			fmt.Println("Here is letter")
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

func run_connection(conn *net.Conn, id int, point *int) {
	text, serr := bufio.NewReader(*conn).ReadString('\n')
	if serr == nil {
		hash_string := ""
		key := ""
		for i := 0; i < 5; i++ {
			hash_string = hash_string + string(text[i])
		}
		for i := 5; i < 15; i++ {
			key = key + string(text[i])
		}
		fmt.Println(hash_string, key)
		server_protector := Session_protector{strings.Replace(hash_string, "\n", "", -1)}
		keyy := server_protector.next_session_key(key)
		(*conn).Write([]byte(keyy + "\n"))
		for {
			message, err := bufio.NewReader(*conn).ReadString('\n')
			if err == nil {
				key = ""
				text = ""
				for i := len(message) - 11; i < len(message)-1; i++ {
					key = key + string(message[i])
				}
				for i := 0; i < len(message)-11; i++ {
					text = text + string(message[i])
				}
				fmt.Println("Client message id = ", id, string(text), "KEY: ", key)
				new_message := strings.ToUpper(text)
				keyy = server_protector.next_session_key(strings.Replace(key, "\n", "", -1))
				fmt.Print("New key: ", keyy, "\n")
				(*conn).Write([]byte(new_message + keyy + "\n"))
			} else {
				(*conn).Close()
				*point -= 1
				fmt.Println("Client ( id =", id, ") DISCONNECTED!")
				break
			}
		}
	} else {
		(*conn).Close()
		*point -= 1
		fmt.Println("Client ( id =", id, ") DISCONNECTED!")
	}
}

func main() {
	port := flag.String("port", ":8080", "a server listening port")
	n := flag.Int("n", 100, "a number of simultaneous connections")
	flag.Parse()
	fmt.Println("Starting server...")
	var id = 1
	ln, _ := net.Listen("tcp", *port)
	point := 1
	for {
		conn, _ := ln.Accept()
		if point <= *n {
			point = point + 1
			fmt.Println("New client ( id =", id, ") CONNECTED!")
			go run_connection(&conn, id, &point)
			id = id + 1
		} else {
			conn.Close()
		}
	}
}
