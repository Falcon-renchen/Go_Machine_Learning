package main

func main()  {
		/*//矩阵运算
		var myvector []float64

		//向 向量添加几个分量。
		myvector = append(myvector,11.0)
		myvector = append(myvector,5.2)
		fmt.Println(myvector)*/

/*
		//矢量运算
		//创建 一个长度为2的矩阵向量
		myvector := mat.NewVecDense(2,[]float64{11.0,5.2})
		fmt.Println(myvector)*/


	/*	//计算A和B的点积
		vectorA := []float64{11.0,5.2,-1.3}
		vectorB := []float64{-7.2,4.2,5.1}
		doProduct := floats.Dot(vectorA,vectorB)	//矩阵乘法
		fmt.Println(doProduct)

		floats.Scale(1.5,vectorA)	// 矩阵A的所有元素乘1.5
		fmt.Println(vectorA)

		normB := floats.Norm(vectorB,2)	//	范数/长度2
		fmt.Println(normB)*/


	/*	myvectorA := mat64.NewVector(3,[]float64{11.0,5.2,-1.3})
		myvectorB := mat64.NewVector(3,[]float64{-7.2,4.2,5.1})
		doProduct := mat64.Dot(myvectorA,myvectorB)	//矩阵乘法
		fmt.Println(doProduct)

		myvectorA.ScaleVec(1.5,myvectorA)							// 将A扩大
		fmt.Println(myvectorA)

		normB := blas64.Nrm2(3,myvectorB.RawVector())				// 将B缩小
		fmt.Println(normB)*/


	/*
		data := []float64{1.2,-5.7,-2.4,7.3}	//切片
		fmt.Println(data)
		a := mat64.NewDense(2,2,data)		// 2✖️2的矩阵
		fa := mat64.Formatted(a,mat64.Prefix("  "))	//	矩阵间隔
		fmt.Printf("A=%v\n",fa)			//输出
	*/

	/*
		//修改矩阵的值
		data := []float64{1.2,-5.7,-2.4,7.3}
		fmt.Println(data)
		a := mat64.NewDense(2,2,data)		//设置矩阵样式
		fa := mat64.Formatted(a,mat64.Prefix(""))
		fmt.Println(fa)			//打印2，2的矩阵
		val := a.At(0,1)	// 从矩阵中抽取出一个值
		fmt.Println(val)

		col := mat64.Col(nil,0,a)		//取出矩阵的第一列的值作为一个新的矩阵
		fmt.Println(col)

		row := mat64.Row(nil,1,a)  	//取出矩阵的第二行的值作为一个新的矩阵
		fmt.Println(row)

		a.Set(0,1,11.2)		// 修改单个值，修改(0,1)的值
		fmt.Println(a)

		a.SetRow(0,[]float64{14.3,-4.2})		//修改第一行的值
		fmt.Println(a)

		a.SetCol(0,[]float64{1.7,-0.3})		//修改第一列的值
		fmt.Println(a)

		fa = mat64.Formatted(a,mat64.Prefix(""))
		fmt.Println(fa)
	*/

/*	a := mat64.NewDense(3,3,[]float64{1,2,3,0,4,5,0,0,6})	//原矩阵
	fa := mat64.Formatted(a,mat64.Prefix(""))
	fmt.Println(fa)
	ft := mat64.Formatted(a.T(), mat64.Prefix(""))	//转置矩阵
	fmt.Println(ft)

	data := mat64.Det(a)
	fmt.Println(data) //矩阵的值

	aInverse := mat64.NewDense(0,0,nil)			// a的逆矩阵
	if err := aInverse.Inverse(a);err!=nil {
		log.Fatal(err)
	}
	fi := mat64.Formatted(aInverse,mat64.Prefix(""))
	fmt.Println(fi)*/

	/*a := mat.NewDense(3,3,[]float64{1,2,3,0,4,5,0,0,6})
	b := mat.NewDense(3,3,[]float64{8,9,10,1,4,2,9,0,2})

	c := mat.NewDense(3,2,[]float64{3,2,1,4,0,8})
	//建立一个空矩阵
	d := mat.NewDense(0,0,nil)
	d.Add(a,b)
	fmt.Println(d)

	f := mat.NewDense(0,0,nil)
	f.Mul(a,c)
	fmt.Println(f)

	g := mat.NewDense(0,0,nil)
	g.Pow(a,5)
	fmt.Println(g)

	h := mat.NewDense(0,0,nil)
	sqrt := func(_,_ int, v float64) float64 {
		return math.Sqrt(v)
	}
	h.Apply(sqrt,a)
	fmt.Println(h)*/

	/*v := mat.NewVecDense(4,[]float64{0,1,2,3})
	matirxPrint1(v)*/

	/*v := mat.NewVecDense(5,[]float64{1,2,3,4,5})
	d := mat.NewVecDense(5,nil)
	d.AddVec(v,v)
	prettyPrintMatrix(d)*/

	/*//加法
	a := mat.NewDense(3,3,[]float64{1,2,3,4,5,6,7,8,9})
	a.Add(a,a)
	matrixPrint2(a)*/

	/*//减法
	a := mat.NewDense(4,2,[]float64{1345, 823, 346, 234, 843, 945, 442, 692})
	b := mat.NewDense(4,2,[]float64{920, 776, 498, 439, 902, 1023, 663, 843})
	var c mat.Dense
	c.Sub(b,a)
	result := mat.Formatted(&c,mat.Prefix(""),mat.Squeeze())
	fmt.Println(result)*/

	//乘法
	/*a := mat.NewDense(3,3,[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9})
	a.Scale(4,a)
	matrixPrint3(a)*/

	/*a := mat.NewDense(2, 2, []float64{1, 2, 3, 4})
	b := mat.NewDense(2, 3, []float64{1, 2, 3, 4, 5, 6})
	var c mat.Dense
	c.Mul(a, b)
	result := mat.Formatted(&c, mat.Prefix(""), mat.Squeeze())
	fmt.Println(result)*/

	/*a := mat.NewDense(1, 3, []float64{0.10, 0.42, 0.37})
	b := mat.NewDense(3, 2, []float64{5, 8, 10, 6, 2, 3})
	var c mat.Dense
	c.Mul(a, b)
	result := mat.Formatted(&c, mat.Prefix("    "), mat.Squeeze())
	fmt.Println(result)*/


/*	a := mat.NewDense(3, 3, []float64{5, 3, 10, 1, 6, 4, 8, 7, 2})
	matrixPrint4(a)
	matrixPrint4(a.T())*/

/*	sparseMatrix := sparse.NewDOK(3, 3)
	sparseMatrix.Set(0, 0, 5)
	sparseMatrix.Set(1, 1, 1)
	sparseMatrix.Set(2, 1, -3)
	fmt.Println(mat.Formatted(sparseMatrix))
	csrMatrix := sparseMatrix.ToCSR()
	fmt.Println(mat.Formatted(csrMatrix))
	cscMatrix := sparseMatrix.ToCSC()
	fmt.Println(mat.Formatted(cscMatrix))*/

	/*sparseMatrix := sparse.NewDOK(4, 4)
	sparseMatrix.Set(0, 2, 1)
	sparseMatrix.Set(1, 0, 2)
	sparseMatrix.Set(2, 3, 3)
	sparseMatrix.Set(3, 1, 4)
	fmt.Print("DOK Matrix:\n", mat.Formatted(sparseMatrix), "\n\n") // Dictionary of Keys
	fmt.Print("CSR Matrix:\n", sparseMatrix.ToCSR(), "\n\n")        // Print CSR version of the matrix*/


	/*sparseMatrix := sparse.NewDOK(4, 4)
	sparseMatrix.Set(0, 2, 1)
	sparseMatrix.Set(1, 0, 2)
	sparseMatrix.Set(2, 3, 3)
	sparseMatrix.Set(3, 1, 4)
	fmt.Print("DOK Matrix:\n", mat.Formatted(sparseMatrix), "\n\n") // Dictionary of Keys
	fmt.Print("CSC Matrix:\n", sparseMatrix.ToCSC(), "\n\n")        // Print CSC version*/
	
}

/*func matirxPrint1(m mat.Matrix) {
	formattedMatrix := mat.Formatted(m,mat.Prefix(""), mat.Squeeze())
	fmt.Printf("%v\n",formattedMatrix)
}*/


/*func prettyPrintMatrix(m mat.Matrix)  {
	formattedM := mat.Formatted(m,mat.Prefix(""),mat.Squeeze())
	fmt.Printf("%v\n",formattedM)
}*/

/*func matrixPrint2(m mat.Matrix)  {
	formattedM := mat.Formatted(m, mat.Prefix(""),mat.Squeeze())
	fmt.Printf("%v\n",formattedM)
}*/

/*func matrixPrint3(m mat.Matrix)  {
	formattedM := mat.Formatted(m,mat.Prefix(""),mat.Squeeze())
	fmt.Println(formattedM)
}*/

/*func matrixPrint4(m mat.Matrix) {
	formattedMatrix := mat.Formatted(m, mat.Prefix(""), mat.Squeeze())
	fmt.Printf("%v\n", formattedMatrix)
}*/

