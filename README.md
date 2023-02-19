# WocLang

**WocLang** 是一款以 **GoLang** 为底层语言实现的编程语言，这个语言不求改变什么，
也没有多牛逼，如果非要有个说法，那就用 Linus 大佬的一句话概括：**Just for fun (只是为了好玩)**

## 目前进度
- 词法分析器
  - 已完成初步设计，v2版本基于确定有限自动机设计
  - 关键字数量太少，待后续添加
  - 数值类型暂不支持浮点型
- 语法分析器
  - 解析 var 语句
  - 解析数值字面量
  - 解析前缀表达式（前缀暂时只支持 '-' 和 '!'，后续添加其他支持，目前只做实现）
