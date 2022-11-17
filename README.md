# greatestworks
the back-end logic of game 



![](frame.png)





### module 说明


每个模块包含
* Model
  对应的数据存储
* System
  该模块的管理，例如数据的CRUD，有System 的模块，其成员实例不具备 独立处理协程
* Owner
  定义从属模块需要实现的些方法
* Handler
  处理从属模块需要的业务逻辑
* Abstract
  模块成员的抽象，接口定义
* config
  配置,常量的定义等
* Manager
  有 Manager 的模块，其成员实例独自拥有自己的处理协程,例如Player Manager , Scene Manager

### Player

 每个 Player 实例 拥有自己的协程

### Scene

 每个 Scene 实例 拥有自己的协程


  


