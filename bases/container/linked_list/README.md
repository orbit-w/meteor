# LinkedList Methods Directory

## Public Methods

- `New[K comparable, V any]() *LinkedList[K, V]`
    - 创建一个新的 LinkedList

- `Init()`
    - 初始化 LinkedList

- `Len() int`
    - 返回 LinkedList 的长度

- `LPush(k K, v V) *Entry[K, V]`
    - 在 LinkedList 的开头插入一个新的元素

- `LPop() *Entry[K, V]`
    - 移除并返回 LinkedList 的第一个元素

- `LPeek() *Entry[K, V]`
    - 返回 LinkedList 的第一个元素，但不移除它

- `Remove(ent *Entry[K, V]) V`
    - 从 LinkedList 中移除一个特定的元素

- `LMove(ent *Entry[K, V])`
    - 将一个特定的元素移动到 LinkedList 的开头

- `RPop() *Entry[K, V]`
    - 返回 LinkedList 的最后一个元素，如果链表为空则返回 nil

- `RPopAt(i int) *Entry[K, V]`
    - RPopAt 返回 LinkedList 的倒数第 i 个元素，如果链表为空或 i 超出范围则返回 nil。
    - i 是从 0 开始的。

- `RPeek() *Entry[K, V]`
    - 返回 LinkedList 的最后一个元素，但不移除它

- `RRange(num int, iter func(k K, v V))`
    - 遍历 LinkedList 的最后 num 个元素