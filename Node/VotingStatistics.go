package Node

import (
	"log"
	"ppov/MetaData"
	"ppov/utils"
)

type Ticket struct {
	counter map[string]int
}

type MimicTicket struct {
	counter map[string]map[string]int
}

func (node *Node) VotingStatisticsforMimic(item MetaData.BlockGroup, items []*MetaData.BlockGroup) (MetaData.BlockGroup, bool) {
	var THRESHOLD = node.config.VotedNum / 3 * 2
	var THRESHOLD2 = node.config.VotedNum / 2

	// 是一个MimicTicket数组长度为记账节点个数，Ticket里面只有一个属性 counter， counter是一个map，key是区块哈希string，value是int表示票数
	counter_yes := make([]MimicTicket, node.config.WorkerNum)
	counter_no := make([]MimicTicket, node.config.WorkerNum)
	result := Ticket{}
	voteRes := make(map[string]map[string]float64)

	MimicCounterInit(counter_yes)
	MimicCounterInit(counter_no)
	// 遍历所有投票节点的投票
	for _, x := range item.VoteTickets {
		if x.BlockHashes == nil || x.VoteResult == nil {
			continue
		}
		//x.BlockHashes是由每个区块哈希的字节数组组成的切片
		if len(x.BlockHashes) != len(x.VoteResult) {
			log.Println("VotingStatistics : len(x.BlockHashes) != len(x.VoteResult)")
			return item, false
		}
		//i他标注的是第几个区块，同样表示是第几个记账节点产生的，block_hash_str标注具体hash值
		for i := 0; i < len(x.BlockHashes); i++ {
			//检查每个区块hash在counter_yes中存不存在，不存在则赋值为0
			block_hash_str := utils.BytesToHex(x.BlockHashes[i])
			counter_yes[i].counter[x.Voter] = make(map[string]int)
			_, is_exist := counter_yes[i].counter[x.Voter][block_hash_str]
			if !is_exist {
				counter_yes[i].counter[x.Voter][block_hash_str] = 0
			}
			//检查每个区块hash在counter_no中存不存在，不存在则赋值为0
			counter_no[i].counter[x.Voter] = make(map[string]int)
			_, is_exist = counter_no[i].counter[block_hash_str]
			if !is_exist {
				counter_no[i].counter[x.Voter][block_hash_str] = 0
			}
			// 投票中如果对于某个区块表决为1，则相应的投票节点的counter_yes
			if x.VoteResult[i] == 1 {
				counter_yes[i].counter[x.Voter][block_hash_str] += 1
			} else if x.VoteResult[i] == -1 {
				counter_no[i].counter[x.Voter][block_hash_str] += 1
			}
		}
	}

	for _, it := range items {
		for _, x := range it.VoteTickets {
			if x.BlockHashes == nil || x.VoteResult == nil {
				continue
			}
			//x.BlockHashes是由每个区块哈希的字节数组组成的切片
			if len(x.BlockHashes) != len(x.VoteResult) {
				log.Println("VotingStatistics : len(x.BlockHashes) != len(x.VoteResult)")
				return item, false
			}
			for i := 0; i < len(x.BlockHashes); i++ {
				//检查每个区块hash在counter_yes中存不存在，不存在则赋值为0
				block_hash_str := utils.BytesToHex(x.BlockHashes[i])
				counter_yes[i].counter[x.Voter] = make(map[string]int)
				_, is_exist := counter_yes[i].counter[x.Voter][block_hash_str]
				if !is_exist {
					counter_yes[i].counter[x.Voter][block_hash_str] = 0
				}
				//检查每个区块hash在counter_no中存不存在，不存在则赋值为0
				counter_no[i].counter[x.Voter] = make(map[string]int)
				_, is_exist = counter_no[i].counter[x.Voter][block_hash_str]
				if !is_exist {
					counter_no[i].counter[x.Voter][block_hash_str] = 0
				}
				// 投票中如果对于某个区块表决为1，则相应的投票节点的counter_yes
				if x.VoteResult[i] == 1 {
					counter_yes[i].counter[x.Voter][block_hash_str] += 1
				} else if x.VoteResult[i] == -1 {
					counter_no[i].counter[x.Voter][block_hash_str] += 1
				}
			}
		}
	}

	check := 0
	if node.round < 2 {
		for _, x := range counter_yes {
			for voter, block_vote := range x.counter {
				for blockhash, voteNum := range block_vote {
					voteRes[voter] = make(map[string]float64)
					voteRes[voter][blockhash] = float64(voteNum)
				}
			}
		}

		final_blockvote := node.checkMA(voteRes)
		// 如果大于2/3赞成那么就是赞同，result[i].counter[k]表示的是投票i对区块k的表决
		for blockHash, voteNum := range final_blockvote {
			if int(voteNum) > THRESHOLD {
				result.counter[blockHash] = 1
				check += 1
				break
			}
			if node.config.VotedNum-int(voteNum) > THRESHOLD {
				result.counter[blockHash] = -1
				check += 1
				break
			}
		}

		if check != node.config.WorkerNum {
			return item, false
		}
		//fmt.Println("vote statics success")
	} else {
		for _, x := range counter_yes {
			for voter, block_vote := range x.counter {
				for blockhash, voteNum := range block_vote {
					voteRes[voter] = make(map[string]float64)
					voteRes[voter][blockhash] = float64(voteNum)
				}
			}
		}

		final_blockvote := node.checkMA(voteRes)
		// 如果大于2/3赞成那么就是赞同，result[i].counter[k]表示的是投票i对区块k的表决
		for blockHash, voteNum := range final_blockvote {
			if int(voteNum) > THRESHOLD2 {
				result.counter[blockHash] = 1
				check += 1
				break
			}
			if node.config.VotedNum-int(voteNum) > THRESHOLD2 {
				result.counter[blockHash] = -1
				check += 1
				break
			}
		}

		if check != node.config.WorkerNum {
			return item, false
		}
	}

	item.BlockHashes = make([][]byte, node.config.WorkerNum)
	item.VoteResult = make([]int, node.config.WorkerNum)

	for i, bh := range item.BlockHashes {
		for k, v := range result.counter {
			if utils.BytesToHex(bh) == k {
				item.BlockHashes[i], _ = utils.HexToBytes(k)
				item.VoteResult[i] = v
			}
			/*			item.BlockHashes = append(item.BlockHashes, k)
						item.VoteResult = append(item.VoteResult, v)*/
		}
	}

	return item, true
}

func (node *Node) VotingStatistics(item MetaData.BlockGroup) (MetaData.BlockGroup, bool) {
	var THRESHOLD = node.config.VotedNum / 3 * 2
	var THRESHOLD2 = node.config.VotedNum / 2

	counter_yes := make([]Ticket, node.config.WorkerNum)
	counter_no := make([]Ticket, node.config.WorkerNum)
	result := make([]Ticket, node.config.WorkerNum)

	CounterInit(counter_yes)
	CounterInit(counter_no)
	CounterInit(result)

	for _, x := range item.VoteTickets {
		if x.BlockHashes == nil || x.VoteResult == nil {
			continue
		}
		//x.BlockHashes是由每个区块哈希的字节数组组成的切片
		if len(x.BlockHashes) != len(x.VoteResult) {
			log.Println("VotingStatistics : len(x.BlockHashes) != len(x.VoteResult)")
			return item, false
		}
		for i := 0; i < len(x.BlockHashes); i++ {
			//检查每个区块hash在counter_yes中存不存在，不存在则赋值为0
			block_hash_str := utils.BytesToHex(x.BlockHashes[i])
			_, is_exist := counter_yes[i].counter[block_hash_str]
			if !is_exist {
				counter_yes[i].counter[block_hash_str] = 0
			}
			//检查每个区块hash在counter_no中存不存在，不存在则赋值为0
			_, is_exist = counter_no[i].counter[block_hash_str]
			if !is_exist {
				counter_no[i].counter[block_hash_str] = 0
			}
			// 投票中如果对于某个区块表决为1，则相应的投票节点的counter_yes
			if x.VoteResult[i] == 1 {
				counter_yes[i].counter[block_hash_str] += 1
			} else if x.VoteResult[i] == -1 {
				counter_no[i].counter[block_hash_str] += 1
			}
		}
	}

	check := 0
	if node.round < 2 {
		for i, x := range counter_yes {
			//k是区块哈希
			for k, v := range x.counter {
				// 如果大于2/3赞成那么就是赞同，result[i].counter[k]表示的是投票i对区块k的表决
				if v > THRESHOLD {
					result[i].counter[k] = 1
					check += 1
					break
				}
				if counter_no[i].counter[k] > THRESHOLD {
					result[i].counter[k] = -1
					check += 1
					break
				}
			}
		}
		if check != node.config.WorkerNum {
			return item, false
		}
		//fmt.Println("vote statics success")
	} else {
		for i, x := range counter_yes {
			for k, v := range x.counter {
				if v > THRESHOLD2 {
					result[i].counter[k] = 1
					check += 1
					break
				}
				if counter_no[i].counter[k] > THRESHOLD2 {
					result[i].counter[k] = -1
					check += 1
					break
				}
			}
		}
		if check == 0 {
			return item, false
		}
		//fmt.Println("check=", check)
	}

	item.BlockHashes = make([][]byte, node.config.WorkerNum)
	item.VoteResult = make([]int, node.config.WorkerNum)
	for i, x := range result {
		for k, v := range x.counter {
			item.BlockHashes[i], _ = utils.HexToBytes(k)
			item.VoteResult[i] = v
			/*			item.BlockHashes = append(item.BlockHashes, k)
						item.VoteResult = append(item.VoteResult, v)*/
		}
	}
	return item, true
}

func CounterInit(counter []Ticket) {
	for i := 0; i < len(counter); i++ {
		counter[i].counter = make(map[string]int)
	}
}

func MimicCounterInit(counter []MimicTicket) {
	for i := 0; i < len(counter); i++ {
		counter[i].counter = make(map[string]map[string]int)
	}
}
