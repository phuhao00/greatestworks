## 游戏玩法相关


# gameplay 协程使用 实践

* slg ,休闲益智，不需要及时响应的一般 交给节点主协程处理<br>
* moba ,fps,一些前卫的玩法设计 每个gameplay 在自己模块开启携程 并发处理

- 总得来说，随着硬件资源的先进，及时反馈，虽然非必需，但是有必要适应未来前卫的玩法设计，避免重构



## for task,condition,title

一般都是 集中写逻辑

但是也可以分散到各个gameplay 写逻辑<br> 
（对应的配表设计可以 增加模块字段区分，<br>
id 还是用数值，不用字符串标识<br>

方便模块任务规划设计，不至于凌乱

）


- condition 一般用作模块，玩法，功能解锁
- task 也是条件的达成，其实某种程度上直接用任务可以达到condition 的功能
- title 可以理解为任务达成的奖励

- 判定依据：
开服天数，时间，等系统数据<br>
玩家自己的数据


## for family-task ,reward-task,daily-circle-task ,active-task,activity-task

- family-task <br>（家族任务）
给予一些家族商店，用于达成称号的一些奖励
- reward-task <br>｛悬赏任务｝
激励玩家参与组建家族玩法
- daily-circle-task<br>（每日任务）
每天的任务-养成升级白票的必要途径<br>
巩固维持玩家对游戏玩法认知度的必要途径<br>
- active-task (活跃任务)
养成升级白票的必要途径<br>
巩固维持玩家对游戏玩法认知度的必要途径<br>
- activity-task （活动任务）<br>
特定的活动，联动，给玩家一些新鲜的玩法体验尝试


##  前端算一遍，后端也算一遍问题

一般来说后端算了，前端就不需要算









