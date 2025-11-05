# B+Trees

- B+Tree Inner Node

  - `Header`
    - `Level`
    - `Slots`
    - `PrevBlockID`
    - `NextBlockID`
    - `HighKey`
    - `LowBlockID`
  - `Keys` - Sorted list
    - `Key`
    - `Position`
  - `Values` - Array
    - `BlockID`

- B+Tree Leaf Node

  - `Header`
    - `Level`
    - `Slots`
    - `PrevBlockID`
    - `NextBlockID`
    - `HighKey`?
  - `Keys` - Sorted list
    - `Key`
    - `Position`
  - `Values` - Array
    - `Value`
