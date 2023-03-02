package main


func sampleTree1() *TreeNode {
	var t0, t0_1L, t0_1L_2L, t0_1L_2R, t0_1R TreeNode

	t0_1L_2L.HasToy = false
	t0_1L_2L.Left = nil
	t0_1L_2L.Right = nil

	t0_1L_2R.HasToy = true
	t0_1L_2R.Left = nil
	t0_1L_2R.Right = nil

	t0_1L.HasToy = false
	t0_1L.Left = &t0_1L_2L
	t0_1L.Right = &t0_1L_2R
/* ============================== */
	t0_1R.HasToy = true
	t0_1R.Left = nil
	t0_1R.Right = nil
/* ============================== */
	t0.HasToy = false
	t0.Left = &t0_1L
	t0.Right = &t0_1R

	return &t0
}

func sampleTree2() *TreeNode {
	var t0, t0_1L, t0_1L_2L, t0_1L_2R, t0_1R,
	t0_1R_2R, t0_1R_2L TreeNode

	t0_1L_2L.HasToy = true
	t0_1L_2L.Left = nil
	t0_1L_2L.Right = nil

	t0_1L_2R.HasToy = false
	t0_1L_2R.Left = nil
	t0_1L_2R.Right = nil

	t0_1L.HasToy = true
	t0_1L.Left = &t0_1L_2L
	t0_1L.Right = &t0_1L_2R
/* ============================== */
	t0_1R_2R.HasToy = true
	t0_1R_2R.Left = nil
	t0_1R_2R.Right = nil

	t0_1R_2L.HasToy = true
	t0_1R_2L.Left = nil
	t0_1R_2L.Right = nil

	t0_1R.HasToy = false
	t0_1R.Left = &t0_1R_2L
	t0_1R.Right = &t0_1R_2R
/* ============================== */
	t0.HasToy = false
	t0.Left = &t0_1L
	t0.Right = &t0_1R

	return &t0
}

func sampleTree3() *TreeNode {
	var t0, t0_1L, t0_1R TreeNode

	t0_1L.HasToy = true
	t0_1L.Left = nil
	t0_1L.Right = nil
/* ============================== */
	t0_1R.HasToy = false
	t0_1R.Left = nil
	t0_1R.Right = nil
/* ============================== */
	t0.HasToy = true
	t0.Left = &t0_1L
	t0.Right = &t0_1R

	return &t0
}

func sampleTree4() *TreeNode {
	var t0, t0_1L, t0_1L_2R, t0_1R,
	t0_1R_2R TreeNode

	t0_1L_2R.HasToy = true
	t0_1L_2R.Left = nil
	t0_1L_2R.Right = nil

	t0_1L.HasToy = true
	t0_1L.Left = nil
	t0_1L.Right = &t0_1L_2R
/* ============================== */
	t0_1R_2R.HasToy = true
	t0_1R_2R.Left = nil
	t0_1R_2R.Right = nil

	t0_1R.HasToy = false
	t0_1R.Left = nil
	t0_1R.Right = &t0_1R_2R
/* ============================== */
	t0.HasToy = false
	t0.Left = &t0_1L
	t0.Right = &t0_1R

	return &t0
}