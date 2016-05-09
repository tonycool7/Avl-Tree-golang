package main

import (
    "fmt"
    "bufio"
    "log"
    "os"
    "net"
)

type avlNode struct {
    Key            string
    Height         int
    Lchild, Rchild *avlNode
}

func leftRotate(root *avlNode) *avlNode {
    node := root.Rchild
    // fmt.Println(node.Key)
    root.Rchild = node.Lchild
    node.Lchild = root

    root.Height = max(height(root.Lchild), height(root.Rchild)) + 1
    node.Height = max(height(node.Rchild), height(node.Lchild)) + 1
    return node
}

func leftRigthRotate(root *avlNode) *avlNode {
    root.Lchild = leftRotate(root.Lchild)
    root = rightRotate(root)
    return  root
}

func rightRotate(root *avlNode) *avlNode {
    node := root.Lchild
    root.Lchild = node.Rchild
    node.Rchild = root
    root.Height = max(height(root.Lchild), height(root.Rchild)) + 1
    node.Height = max(height(node.Lchild), height(node.Rchild)) + 1
    return node
}

func rightLeftRotate(root *avlNode) *avlNode {
    root.Rchild = rightRotate(root.Rchild)
    root = leftRotate(root)
    return  root
}

func height(root *avlNode) int {
    if root != nil {
        return root.Height
    }
    return -1
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func insert(root *avlNode, key string) *avlNode {
    if root == nil {
        root = &avlNode{key, 0, nil, nil}
        root.Height = max(height(root.Lchild), height(root.Rchild)) + 1
        return root
    } 

    if key < root.Key {
        root.Lchild = insert(root.Lchild, key)
        if height(root.Lchild)-height(root.Rchild) == 2 {
            if key < root.Lchild.Key {
                root = rightRotate(root) // 左左
            } else {
                root = leftRigthRotate(root) // 左右
            }
        }
    } 

    if key > root.Key {
        root.Rchild = insert(root.Rchild, key)
        if height(root.Rchild)-height(root.Lchild) == 2 {
            if key > root.Rchild.Key {
                root = leftRotate(root) // 右右
            } else {
                root = rightLeftRotate(root) // 右左
            }
        }
    }

    root.Height = max(height(root.Lchild), height(root.Rchild)) + 1
    return root
}


type action func(node *avlNode)

func inOrder(root *avlNode, action action) {
    if root == nil {
        return
    }

    inOrder(root.Lchild, action)
    action(root)
    inOrder(root.Rchild, action)
}


func equal(s1, s2 string) int {
    eq := 0
    if len(s1) > len(s2) {
        s1, s2 = s2, s1
    }
    for key, _ := range s1 {
        if s1[key] == s2[key] {
            eq++
        } else {
            break
        }
    }
    return eq
}

func main() {
    var root *avlNode
    // keys := []string{"tony", "femi", "cool", "jb", "motherfucker", "gabi", "amed"}
    // for _, key := range keys {
    file, err := os.Open("file.txt")
    if err != nil {
        log.Fatal(err)
    }

    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        root = insert(root, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    fmt.Println("Launching server...\n")
    fmt.Println("The dictionary")
    inOrder(root, func(node *avlNode) {
        fmt.Println(node.Key)
         })

    fmt.Println("\n\n")

  // listen on all interfaces
  ln, _ := net.Listen("tcp", ":8081")

  // accept connection on port
  conn, _ := ln.Accept()
  // var num int = 0

  // run loop forever (or until ctrl-c)
  for {
        // will listen for message to process ending in newline (\n)
        message, _ := bufio.NewReader(conn).ReadString('\n')
        // output message received
        fmt.Print("Message Received:", message)
        meslen := len([]rune(message)) - 1
        inOrder(root, func(node *avlNode) {
             temp := node.Key+"\n"
            
                if message == temp[0:meslen]+"\n" {
                    fmt.Println("Client received definition of word '" +message+"\n")
                    conn.Write([]byte(temp))
                }
             })
        
        conn.Write([]byte("not found\n"))

    }
}