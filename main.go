package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const (
	MAXREAD int = 1024 * 512
)

func main() {

	Compress("code.huffman", "code")
	Uncompress("a", "code.huffman")
}

func Compress(dst, src string) {
	tree := StatisticalFrequency(src)
	if tree == nil {
		return
	}
	tree.BuildHuffmanTree()
	preOrderPrint(tree.Root, tree.HuffCode, "")

	fp, err := os.Create(dst)
	defer fp.Close()
	if err != nil {
		log.Fatalln(err)
		return
	}
	var buf bytes.Buffer
	var tmp byte
	var count int
	for key, _ := range tree.HuffCode {
		buf.WriteRune(key)
		buf.WriteString(strconv.Itoa(tree.Frequency[key]))
		buf.WriteRune(' ')
	}
	fp.Write(buf.Bytes())
	buf.Reset()
	tree.FilePoint.Seek(0, 0)
	var context []byte = make([]byte, MAXREAD)
	for {
		if _, err := tree.FilePoint.Read(context); err != nil {
			if err == io.EOF {
				break
			} else {
				return
			}
		} else {
			for _, v := range string(context) {
				curHuffCode := tree.HuffCode[v]
				var num byte = 0
				flag := false
				for _, bit := range curHuffCode {
					num = (num << 1) | byte(bit)
					count++
					if count == 8 {
						tmp = (tmp << (len(curHuffCode))) | byte(num)
						// log.Infof("%b", tmp)
						flag = true
						buf.WriteByte(tmp)
						tmp = 0
						count = 0
						num = 0
					}
				}
				if !flag {
					tmp = tmp<<(len(curHuffCode)) | byte(num)
				} else {
					flag = false
				}
			}
		}
	}
	fp.Write(buf.Bytes())
}

func Uncompress(dst, src string) {
	file, err := os.Open(src)
	defer file.Close()
	if err != nil {
		log.Fatalln(err)
	}
	var reader *bufio.Reader
	var context []byte = make([]byte, MAXREAD)
	reader = bufio.NewReader(file)

	line, err := reader.ReadBytes('\n')
	if err != nil {
		log.Fatalln(err)
	}
	t := reconstructHuffmanTree(string(line))
	if t == nil {
		log.Fatalln("reconstruct error")
		return
	}
	root := t.Root
	var cur byte = 0
	tmp := root

	dstfp, err := os.Create(dst)
	defer dstfp.Close()
	if err != nil {
		log.Fatalln(err)
	}
	var buf bytes.Buffer
	for {
		if _, err := reader.Read(context); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatalln(err)
				return
			}
		} else {
			for _, v := range context {
				cur = v
				for i := 0; i < 8; i++ {
					if cur&1 == 1 {
						if tmp == nil {
							return
						}
						tmp = tmp.Right
						if tmp.Left == nil && tmp.Right == nil {
							buf.WriteRune(tmp.curKey)
							tmp = root
						}
					} else {
						if tmp == nil {
							return
						}
						tmp = tmp.Left
						if tmp.Left == nil && tmp.Right == nil {
							buf.WriteRune(tmp.curKey)
							tmp = root
						}
					}
					cur >>= 1
				}
			}
		}
	}
	dstfp.Write(buf.Bytes())

}

func preOrderPrint(root *HuffmanTreeNode, HuffCode map[rune]string, curPath string) {
	if root == nil {
		return
	}
	if root.Left == nil && root.Right == nil {
		root.Path = curPath
		HuffCode[root.curKey] = curPath
		return
	}
	if root.Left != nil {
		preOrderPrint(root.Left, HuffCode, curPath+"0")
	}
	if root.Right != nil {
		preOrderPrint(root.Right, HuffCode, curPath+"1")
	}
}

func (tree *HuffmanTree) BuildHuffmanTree() {
	for !tree.Heap.IsEmpty() {
		a, err := tree.Heap.Delete()
		if err != nil {
			log.Fatalln(err)
			return
		}
		b, err := tree.Heap.Delete()
		if err != nil {
			log.Infoln("Success build huffman tree")
			tree.Root = a
			return
		}
		root := NewHuffmanTreeNode(a.curVal+b.curVal, ' ')
		if a.curVal > b.curVal {
			a, b = b, a
		}
		root.Left, root.Right = a, b
		tree.Heap.Insert(root)
	}
}

func (tree *HuffmanTree) statisticalFrequency(context string) {
	log.Infoln(context)
	for _, v := range context {
		// log.Infoln(v)
		if v == '0' || v == '\n' {
			continue
		}
		tree.Frequency[v] += 1
	}
}
