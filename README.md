DNF: A or B or C

given a set of DNFs, split into conjunctions, assign conjunction ID 
each conjunction has a pointer to its parent DNF.
compute the size of each conjunction.

build K-index table: partition all the conjunctions by its size.

0 -> map[key]postingList
1 -> map[key]postingList
2 -> map[key]postingList

posting list is a sorted list of entry,
entry: <conjunction-id, boolean flag(belong or not), score>
entries are sorted by conjunction-id, bool-flag ascendingly
there is a special post-list Z, for zero-size conjunction


