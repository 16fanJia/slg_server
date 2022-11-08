package router

import (
	"errors"
	"fmt"
	"strings"
)

/*
==========构造前缀树======= 方便 handler 的查询

path 中 * 为通配符 即匹配所有
eg： /'*'/name
	/account/name
*/

var RouteAlreadyExists = errors.New("route already exists ,路由已经存在！")

//树的结构体
type tree struct {
	root *node
}

//树的构造函数
func newTree() *tree {
	return &tree{
		root: newNode(),
	}
}

//findHandler 根据path 寻找 handler
func (t *tree) findHandler(path string) []HandlerFunc {
	//拆分path ' /account/name ' 去除最前面的 /
	if strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
		//路由不区分大小写
		path = strings.ToLower(path)
	} else {
		//打日志 并且退出
		fmt.Println("路由路径格式错误！")
		return nil
	}

	matchNode := t.root.matchNode(path)
	if matchNode == nil {
		return nil
	}
	return matchNode.handlers
}

//add 路由树中添加路由
func (t *tree) addRouter(path string, handlerFunc ...HandlerFunc) {
	n := t.root
	//拆分path ' /account/name ' 去除最前面的 /
	if strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
		//路由不区分大小写
		path = strings.ToLower(path)
	} else {
		//打日志 并且退出
		fmt.Println("路由路径格式错误！")
		return
	}
	//[*,name,fan,jia]
	segments := strings.Split(path, "/")
	for index, segment := range segments {
		isLast := index == len(segments)-1

		//匹配 并 返回目标节点
		var objNode *node
		children := n.matchChildNode(segment)
		if len(children) > 0 {
			for _, cNode := range children {
				//找到segment 相同的子节点
				if cNode.segment == segment {
					objNode = cNode
					break
				}
			}
		}

		//未找到对应的节点
		if objNode == nil {
			cNode := newNode()
			cNode.segment = segment
			if isLast {
				cNode.isLast = true
				cNode.handlers = handlerFunc
			}
			n.children = append(n.children, cNode)
			objNode = cNode
		}
		n = objNode
	}
}

//节点的结构体
type node struct {
	isLast   bool          //是否是路径上的最后一个
	segment  string        //path 中的部分字符串
	handlers []HandlerFunc //处理函数
	children []*node       //子节点
}

func newNode() *node {
	return &node{
		isLast:   false,
		segment:  "/",
		handlers: []HandlerFunc{},
		children: []*node{},
	}
}

//matchChildNode 匹配子节点
func (n *node) matchChildNode(seg string) []*node {
	//当前节点没有孩子节点 则没有目标节点
	if len(n.children) == 0 {
		return nil
	}
	//如果是通配符 * 则所有子节点都满足
	if seg == "*" {
		return n.children
	}

	nodes := make([]*node, 0, len(n.children))
	for _, cNode := range n.children {
		switch cNode.segment {
		case seg:
			nodes = append(nodes, cNode)
		case "*":
			nodes = append(nodes, cNode)
		}
	}
	return nodes
}

//matchNode 根据path 寻找节点
func (n *node) matchNode(path string) *node {
	segments := strings.SplitN(path, "/", 2)
	segment := segments[0]
	//匹配对应的节点
	cNode := n.matchChildNode(segment)
	//未找到对应的节点
	if cNode == nil || len(cNode) == 0 {
		return nil
	}
	//如果只有一个
	if len(segments) == 1 {
		for _, vn := range cNode {
			if vn.isLast {
				return vn
			}
		}
		//都不是最终的节点
		return nil
	}
	for _, vn := range cNode {
		mNode := vn.matchNode("/" + segments[1])
		if mNode != nil {
			return mNode
		}
	}
	return nil
}
