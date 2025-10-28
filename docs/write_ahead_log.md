# Write Ahead Log

## Log Records

### Common to all

- prevLSN
- transID
- type

### Update records

- pageID
- length
- offset
- beforeImage
- afterImage

## Transaction table

- pageID
- recLSN

## Dirty page table

- transID
- lastLSN
