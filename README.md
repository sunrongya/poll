# 项目
建模练习：调查投票
> 演示DDD + TDD + CQRS + Event Sourcing技术

## 目录
- [简介](#简介)
- [聚合根](#聚合根)
- [值对象](#值对象)
- [调查题目聚合根相关](#调查题目聚合根相关)
- [投票记录聚合根相关](#投票记录聚合根相关)

## 简介
**“调查投票”** 场景（这种场景很多见，比如选择一个你认为最强大的编程语言，可选项：C#、PHP、JAVA、Lisp...）
- 一个调查题目（poll）包含多个选项（choice）
- 一个用户（user）可为某个题目（poll）的多个选项进行投票（比如可以同时选择C#、Lisp）
- 关注题目的投票记录（如投票人不能重复投票、所投票的选项要包含在poll里面）

> 注意：建模时考虑高并发的情况，即多人同时投票。

## 聚合根
- **题目** (ID、标题、选项列表、投票记录ID列表、投票选项统计）
- **投票记录** （ID、题目ID、投票选项ID列表、投票人信息、投票时间）

## 值对象
- **投票人**、**选项** 等

## 调查题目聚合根相关
- **领域服务：** 创建
- **命令：** 创建、投票
- **事件：** 已创建、投票成功、投票失败

## 投票记录聚合根相关
- **领域服务：** 投票
- **命令：** 创建投票、成功投票、失败投票
- **事件：** 已创建投票、投票成功、投票失败

