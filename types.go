package main

import (
	"io"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type HuffmanTreeNode struct {
	Left, Right *HuffmanTreeNode
	Path        string
	curVal      int
	curKey      rune
}

type HuffmanTree struct {
	Root      *HuffmanTreeNode
	FileName  string
	FilePoint *os.File
	Heap      *MinHeap
	Frequency map[rune]int
	HuffCode  map[rune]string
}

func NewHuffmanTreeNode(curVal int, curKey rune) *HuffmanTreeNode {
	return &HuffmanTreeNode{
		curVal: curVal,
		curKey: curKey,
	}
}

func NewHuffmanTree(filePath string) *HuffmanTree {
	tree := &HuffmanTree{
		Root:      NewHuffmanTreeNode(1<<31, rune(' ')),
		FileName:  filePath,
		HuffCode:  make(map[rune]string),
		Frequency: make(map[rune]int),
	}
	return tree
}

func StatisticalFrequency(filePath string) *HuffmanTree {
	tree := NewHuffmanTree(filePath)
	fp, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	tree.FilePoint = fp
	var context []byte = make([]byte, MAXREAD)
	for {
		if _, err := tree.FilePoint.Read(context); err != nil {
			if err == io.EOF {
				break
			} else {
				return nil
			}
		} else {
			tree.statisticalFrequency(string(context))
			context = make([]byte, MAXREAD)
		}
	}
	heap := NewMinHeap(len(tree.Frequency))
	tree.Heap = heap
	for key, value := range tree.Frequency {
		node := NewHuffmanTreeNode(value, key)
		heap.Insert(node)
	}
	heap.BuildHeap()
	return tree
}

func reconstructHuffmanTree(context string) *HuffmanTree {
	next := strings.Split(context, " ")
	heap := NewMinHeap(100)
	tree := NewHuffmanTree("")
	tree.Heap = heap
	for _, v := range next {
		// log.Infoln(v)
		if len(v) <= 1 {
			continue
		}
		char := v[0]
		frequency, err := strconv.Atoi(v[1:])
		if err != nil {
			log.Fatalln(err)
			return nil
		}
		node := NewHuffmanTreeNode(frequency, rune(char))
		heap.Insert(node)
	}
	heap.BuildHeap()
	tree.BuildHuffmanTree()
	return tree
}
