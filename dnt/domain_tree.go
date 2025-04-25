// domain to tree

package dnt

import (
	"strings"
)

type domainTreeNode struct {
	name     string
	children map[string]*domainTreeNode
}

type DomainTree struct {
	root *domainTreeNode
}

// NewDomainTree - create domain tree
func NewDomainTree() *DomainTree {
	return &DomainTree{
		root: &domainTreeNode{
			name:     RootDomain,
			children: make(map[string]*domainTreeNode),
		},
	}
}

// Add - add domain into tree
func (d *DomainTree) Add(domain string) {
	if len(domain) < 1 {
		return
	}

	domain = FQD(domain)

	segList := strings.Split(domain, ".")
	node := d.root
	for i := len(segList) - 1; i >= 0; i-- {
		seg := segList[i]
		if seg == "" {
			continue
		}
		next, h := node.children[seg]
		if !h {
			next = &domainTreeNode{
				name:     seg,
				children: make(map[string]*domainTreeNode),
			}
			node.children[seg] = next
		}
		node = next
	}
}

// Match - domain match
func (d *DomainTree) Match(domain string) bool {
	if len(domain) < 1 {
		return false
	}

	node := d.root

	segList := strings.Split(domain, ".")
	for i := len(segList) - 1; i >= 0; i-- {
		seg := segList[i]
		if seg == "" {
			continue
		}
		next, h := node.children[seg]
		if !h {
			// match wildcard domain
			_, h = node.children["*"]
			return h
		}
		node = next
	}
	return true
}
