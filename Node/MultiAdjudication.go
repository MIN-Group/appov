package Node

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"

	"sync"
)

func (node *Node) checkMA(voteResults map[string]map[string]float64) map[string]float64 {
	m := &sync.Map{}
	n := len(voteResults)
	vis := make(map[string][]string, n)

	tmp := 1 / float64(n)

	blacklist := map[string][]string{}
	//flag := false

	w := make(map[string]map[string]float64, n)
	for nodeID1, _ := range voteResults {
		w[nodeID1] = make(map[string]float64, n)
		for nodeID2, _ := range voteResults {
			w[nodeID1][nodeID2] = tmp
		}
	}

	// node.network.NodeList改成leader集合
	for nodeID1, _ := range voteResults {
		for nodeID2, _ := range voteResults {
			if nodeID2 != nodeID1 {
				vis[nodeID1] = append(vis[nodeID1], nodeID2)
			}
		}
	}
	fmt.Println(vis)
	fmt.Println(w)
	fmt.Println(n)

	pubkey := node.keymanager.GetPubkey()

	txNum := len(voteResults[pubkey])
	ans := make(map[string]float64, txNum)

	for txName, _ := range voteResults[pubkey] {
		x := make(map[string]float64)
		for leader, voteResult := range voteResults {
			x[leader] = voteResult[txName]
			m.Store(leader, voteResult[txName])
		}

		d := make(map[string]float64)
		e := make(map[string]float64)

		pre_sum := make(map[string]float64)

		for k := 0; k < 50; k++ {
			//fmt.Println(x)
			//fmt.Println(blacklist)

			// node.network.NodeList改成leader集合
			for leader1, _ := range voteResults {
				tmp1 := pre_sum[leader1]
				pre_sum[leader1] = 0
				for _, leader2 := range vis[leader1] {
					if k == 0 {
						d[leader1] = math.Max(d[leader1], math.Abs(x[leader1]-x[leader2]))
					}
					//if InSlice(blacklist[i], j) == false {
					//	pre_sum[i-1] = pre_sum[i-1] + math.Abs(x[i-1]-x[j-1])
					//}
				}
				if k > 0 {
					if tmp1 == 0 {
						d[leader1] = 0
					} else {
						d[leader1] = pre_sum[leader1] * d[leader1] / tmp1
					}
				}
			}
			//fmt.Println("门限：", d)
			tmp := make(map[string]float64)
			//tmp := make([]float64, 4, 4)
			for key, v := range x {
				tmp[key] = v
			}

			//tmp := x

			// 改成Leader集合
			for leader1, _ := range voteResults {
				x[leader1] = w[leader1][leader1] * tmp[leader1]
				for _, leader2 := range vis[leader1] {
					if InSlice(blacklist[leader1], leader2) == false {
						x[leader1] = x[leader1] + w[leader1][leader2]*tmp[leader2] + e[leader1]
					}
					if x[leader1] >= 10 {
						x[leader1] = 10
					}
				}
			}

			//
			//for i := 1; i <= 4; i++ {
			//	for _, j := range vis[i] {
			//		if InSlice(blacklist[i], j) == false {
			//			if k > 0 && math.Abs(x[i-1]-x[j-1]) > 1.0*d[i-1] && InSlice(blacklist[i], j) == false {
			//				if flag {
			//					blacklist[i] = append(blacklist[i], j)
			//					p1 := 1.00 / float64(len(vis[i])-len(blacklist[i])+1)
			//					for h := 0; h <= 3; h++ {
			//						w[i-1][h] = p1
			//					}
			//				}
			//				//fmt.Println(len(vis[i])-len(blacklist[i]), "p1:---", p1)
			//			}
			//		}
			//	}
			//}
		}
		var tmplist []float64
		for _, v := range x {
			tmplist = append(tmplist, v)
		}
		ans[txName] = majorityElement(tmplist)
	}
	return ans
}

// InSlice 判断int是否在 slice 中。
func InSlice(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func FileIsExisted(filename string) bool {
	existed := true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		existed = false
	}
	return existed
}

func MakeDir(dir string) error {
	if !FileIsExisted(dir) {
		if err := os.MkdirAll(dir, 0777); err != nil { //os.ModePerm
			fmt.Println("MakeDir failed:", err)
			return err
		}
	}
	return nil
}

func SaveFile(filename string, data []byte) bool {
	var err error
	if len(filename) > 0 && data != nil {
		dir := filepath.Dir(filename)
		if MakeDir(dir) != nil {
			return false
		}

		err = ioutil.WriteFile(filename, data, 0666)
		if err != nil {
			fmt.Println("SaveFile err:", err)
		} else {
			return true
		}
	} else {
		fmt.Println("SaveFile err: wrong params")
	}
	return false
}

func majorityElement(nums []float64) float64 {
	des := make(map[float64]int)
	length := len(nums)
	for _, num := range nums {
		des[num]++               //遍历得到的数字，作为索引进行+1
		if des[num] > length/2 { //每进行一次+1，进行判断，如果出现的次数大于数组长度的一半，则返回此数字
			return num
		}
	}
	return -1
}
