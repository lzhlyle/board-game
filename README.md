# board-game

a board game & ai framework

## 框架博客

- [blog.lzh.today/board-game-framework](http://blog.lzh.today/board-game-framework/)

## 目录 & 文件

- `core/` 游戏基建类，游戏框架（无关 AI）
    - `board_game.go` 游戏框架接口定义
    - `play.go` 游戏主框架，程序运行流程
- `ai/` AI 基建类
    - `alg.go` AI 框架接口定义
- `ai_impl/` AI 默认实现
    - `default_ai_impl.go` AI 默认实现：可实现 10 行生成 AI（如 `gobang_ai.go`）
    - `chess_record_generator.go` 棋谱生成器：枚举所有落子，计算胜平负概率
- `concrete/` 具体游戏实现
    - `tic_tac_toe.go` & `tic_tac_toe_ai.go` 井字棋游戏规则 & 井字棋 AI
    - `gobang.go` & `gobang_ai.go` 五子棋游戏规则 & 五子棋 AI

## 术语定义

- 棋盘：`board`
- 玩家：`player`
- 一步棋：`move`，表示某一玩家在某一回合中的一次决策，不包括决策后带来的影响
    - 决策后的影响：指非玩家主管决策，而由系统自动触发的场面的变化，如围棋/象棋的吃子、飞行棋的跳跃
- 回合：`round`，表示某一玩家的下棋回合，可包括若干步
- 整盘棋：`game`，如：全局开始、全局结束

## 外部依赖

- [jroimartin/gocui](https://github.com/jroimartin/gocui) 实现控制台交互
