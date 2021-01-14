package algorithm

// copyMap creates a copy of the provided map.
func copyMap(m map[int64]*algUser) map[int64]*algUser {
	mCopy := map[int64]*algUser{}
	for k, v := range m {
		mCopy[k] = v
	}
	return mCopy
}
