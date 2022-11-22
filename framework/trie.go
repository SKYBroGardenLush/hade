package framework

import (
	"errors"
	"strings"
)

type Tree struct {
	root *node
}

func NewTree() *Tree {
	return &Tree{root: newNode(nil)}
}

type node struct {
	isLast   bool                // 该结点是否能成为一个独立的uri，是否自身就是一个终极结点
	segment  string              //uri中的字符串，代表这个节点表示路由中某个段的字符串
	handlers []ControllerHandler //代表这个介蒂安包含的控制器，用于最终加载调用
	children []*node             //代表这个节点下的子节点
	parent   *node
}

func newNode(parent *node) *node {
	return &node{
		isLast:   false,
		segment:  "",
		children: []*node{},
		parent:   parent,
	}
}

//isWildSegment 判断一个segment是否是通用segment，即以:开头
func isWildSegment(segment string) bool {
	return strings.HasPrefix(segment, ":")
}

//filterChildNodes 过滤下一层满足segment规则的子节点
func (n *node) filterChildNodes(segment string) []*node {
	if len(n.children) == 0 {
		return nil
	}

	//如果segment是通配符，则所有下一层子节点都满足需求
	if isWildSegment(segment) {
		return n.children
	}
	nodes := make([]*node, 0, len(n.children))
	//过滤所有下一层的节点
	for _, cNode := range n.children {
		if isWildSegment(cNode.segment) {
			nodes = append(nodes, cNode)
		} else if cNode.segment == segment {
			nodes = append(nodes, cNode)
		}
	}
	return nodes
}

//matchNode 判断路由是否已经在节点的所有子节点树中存在了
func (n *node) matchNode(uri string) *node {

	// 使用分隔符将uri分割成两个部分
	segments := strings.SplitN(uri, "/", 2)

	//第一个部分用于匹配下一层子节点
	segment := segments[0]

	if !isWildSegment(segment) {
		segment = strings.ToUpper(segment)
	}
	//匹配符合的下一层子节点
	cNodes := n.filterChildNodes(segment)

	// 如果当前子节点没有一个符合，那么说明这个uri一定是之前不存在, 直接返回nil
	if cNodes == nil || len(cNodes) == 0 {
		return nil
	}

	//如果只有一个segment，则最后一个标记
	if len(segments) == 1 {

		for _, tn := range cNodes {
			if tn.isLast {
				return tn
			}
		}
		return nil
	}
	for _, tn := range cNodes {
		tnMatch := tn.matchNode(segments[1])
		if tnMatch != nil {
			return tnMatch
		}
	}
	return nil

}
func (t *Tree) AddRouter(uri string, handlers []ControllerHandler) error {
	n := t.root
	if n.matchNode(uri) != nil {
		return errors.New("route exist : " + uri)
	}

	segments := strings.Split(uri, "/")

	for i, segment := range segments {
		if !isWildSegment(segment) {
			segment = strings.ToUpper(segment)
		}
		isLast := i == len(segments)-1
		var objNode *node
		cNodes := n.filterChildNodes(segment)
		if len(cNodes) > 0 {
			for _, cNode := range cNodes {
				if cNode.segment == segment {
					objNode = cNode
					break
				}
			}
		}
		if objNode == nil {
			newChildNode := newNode(n)
			newChildNode.segment = segment
			if isLast {
				newChildNode.isLast = true
				newChildNode.handlers = handlers
			}
			n.children = append(n.children, newChildNode)
			objNode = newChildNode
		}
		n = objNode

	}
	return nil

}

//将uri 解析为params
func (n *node) parseParamsFromEndNode(uri string) map[string]string {

	ret := map[string]string{}
	segments := strings.Split(uri, "/")
	cnt := len(segments)
	cur := n
	for i := cnt - 1; i >= 0; i-- {
		if cur.segment == "" {
			break
		}
		// 如果是通配符节点
		if isWildSegment(cur.segment) {
			//设置 params
			ret[cur.segment[1:]] = segments[i]
		}
		cur = cur.parent
	}
	return ret
}

// FindHandler 匹配uri
func (t *Tree) FindHandler(uri string) []ControllerHandler {
	matchNode := t.root.matchNode(uri)
	if matchNode == nil {

		return nil
	}
	return matchNode.handlers
}
