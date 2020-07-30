package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gonum/floats"

	"gonum.org/v1/gonum/mat"
)


// // NeuroNet包含定义受过训练的神经网络的所有信息。
type neuralNet struct {
	config  neuralNetConfig
	wHidden *mat.Dense
	bHidden *mat.Dense
	wOut    *mat.Dense
	bOut    *mat.Dense
}


// NeuroNetConfig定义了我们的神经网络架构和学习参数。
type neuralNetConfig struct {
	inputNeurons  int
	outputNeurons int
	hiddenNeurons int
	numEpochs     int
	learningRate  float64
}

func main() {

	f, err := os.Open("Evaluation/train.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = 7

	rawCSVData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// inputsData和labelsData将保存所有最终将用于形成矩阵的浮点值。

	inputsData := make([]float64, 4*len(rawCSVData))
	labelsData := make([]float64, 3*len(rawCSVData))

	// inputsIndex将跟踪输入矩阵值的当前索引。
	var inputsIndex int
	var labelsIndex int

	//将行顺序移动到一片浮点数中。
	for idx, record := range rawCSVData {

		if idx == 0 {
			continue
		}

		//循环遍历float列。
		for i, val := range record {

			//将值转换为浮点数。
			parsedVal, err := strconv.ParseFloat(val, 64)
			if err != nil {
				log.Fatal(err)
			}

			//如果相关，将其添加到labelsData。
			if i == 4 || i == 5 || i == 6 {
				labelsData[labelsIndex] = parsedVal
				labelsIndex++
				continue
			}

			//将float值添加到float切片中。
			inputsData[inputsIndex] = parsedVal
			inputsIndex++
		}
	}

	//形成矩阵。
	inputs := mat.NewDense(len(rawCSVData), 4, inputsData)
	labels := mat.NewDense(len(rawCSVData), 3, labelsData)

	//定义我们的网络架构和学习参数。
	config := neuralNetConfig{
		inputNeurons:  4,
		outputNeurons: 3,
		hiddenNeurons: 3,
		numEpochs:     5000,
		learningRate:  0.3,
	}

	//训练神经网络。
	network := newNetwork(config)
	if err := network.train(inputs, labels); err != nil {
		log.Fatal(err)
	}

	//打开测试数据集文件。
	f, err = os.Open("Evaluation/test.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	//从打开的文件中创建一个新的CSV阅读器。
	reader = csv.NewReader(f)
	reader.FieldsPerRecord = 7

	//读取所有测试CSV记录
	rawCSVData, err = reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// inputsData和labelsData将保存所有最终将用于形成矩阵的浮点值。
	inputsData = make([]float64, 4*len(rawCSVData))
	labelsData = make([]float64, 3*len(rawCSVData))

	// inputsIndex将跟踪输入矩阵值的当前索引。
	inputsIndex = 0
	labelsIndex = 0

	//将行顺序移动到一片浮点数中。
	for idx, record := range rawCSVData {

		if idx == 0 {
			continue
		}

		//循环遍历float列。
		for i, val := range record {

			//将值转换为浮点数。
			parsedVal, err := strconv.ParseFloat(val, 64)
			if err != nil {
				log.Fatal(err)
			}

			//如果相关，将其添加到labelsData。
			if i == 4 || i == 5 || i == 6 {
				labelsData[labelsIndex] = parsedVal
				labelsIndex++
				continue
			}

			//将float值添加到float切片中。
			inputsData[inputsIndex] = parsedVal
			inputsIndex++
		}
	}

	//形成矩阵。
	testInputs := mat.NewDense(len(rawCSVData), 4, inputsData)
	testLabels := mat.NewDense(len(rawCSVData), 3, labelsData)

	//使用经过训练的模型进行预测。
	predictions, err := network.predict(testInputs)
	if err != nil {
		log.Fatal(err)
	}

	//计算模型的准确性。
	var truePosNeg int
	numPreds, _ := predictions.Dims()
	for i := 0; i < numPreds; i++ {

		// 获取标签。
		labelRow := mat.Row(nil, i, testLabels)
		var species int
		for idx, label := range labelRow {
			if label == 1.0 {
				species = idx
				break
			}
		}

		//累积真实的正数/负数。
		if predictions.At(i, species) == floats.Max(mat.Row(nil, i, predictions)) {
			truePosNeg++
		}
	}

	//计算精度（子集精度）。
	accuracy := float64(truePosNeg) / float64(numPreds)

	//将Accuracy值输出到标准输出。
	fmt.Printf("\nAccuracy = %0.2f\n\n", accuracy)
}

// NewNetwork初始化一个新的神经网络。
func newNetwork(config neuralNetConfig) *neuralNet {
	return &neuralNet{config: config}
}

// train使用反向传播训练神经网络。
func (nn *neuralNet) train(x, y *mat.Dense) error {

	//初始化偏差/权重。
	randSource := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSource)

	wHiddenRaw := make([]float64, nn.config.hiddenNeurons*nn.config.inputNeurons)
	bHiddenRaw := make([]float64, nn.config.hiddenNeurons)
	wOutRaw := make([]float64, nn.config.outputNeurons*nn.config.hiddenNeurons)
	bOutRaw := make([]float64, nn.config.outputNeurons)

	for _, param := range [][]float64{wHiddenRaw, bHiddenRaw, wOutRaw, bOutRaw} {
		for i := range param {
			param[i] = randGen.Float64()
		}
	}

	wHidden := mat.NewDense(nn.config.inputNeurons, nn.config.hiddenNeurons, wHiddenRaw)
	bHidden := mat.NewDense(1, nn.config.hiddenNeurons, bHiddenRaw)
	wOut := mat.NewDense(nn.config.hiddenNeurons, nn.config.outputNeurons, wOutRaw)
	bOut := mat.NewDense(1, nn.config.outputNeurons, bOutRaw)

	//定义神经网络的输出。
	var output mat.Dense

	//使用反向传播循环训练我们的模型。
	for i := 0; i < nn.config.numEpochs; i++ {

		//完成前馈过程。
		var hiddenLayerInput mat.Dense
		hiddenLayerInput.Mul(x, wHidden)
		addBHidden := func(_, col int, v float64) float64 { return v + bHidden.At(0, col) }
		hiddenLayerInput.Apply(addBHidden, &hiddenLayerInput)

		var hiddenLayerActivations mat.Dense
		applySigmoid := func(_, _ int, v float64) float64 { return sigmoid(v) }
		hiddenLayerActivations.Apply(applySigmoid, &hiddenLayerInput)

		var outputLayerInput mat.Dense
		outputLayerInput.Mul(&hiddenLayerActivations, wOut)
		addBOut := func(_, col int, v float64) float64 { return v + bOut.At(0, col) }
		outputLayerInput.Apply(addBOut, &outputLayerInput)
		output.Apply(applySigmoid, &outputLayerInput)

		//完成反向传播。
		var networkError mat.Dense
		networkError.Sub(y, &output)

		var slopeOutputLayer mat.Dense
		applySigmoidPrime := func(_, _ int, v float64) float64 { return sigmoidPrime(v) }
		slopeOutputLayer.Apply(applySigmoidPrime, &output)
		var slopeHiddenLayer mat.Dense
		slopeHiddenLayer.Apply(applySigmoidPrime, &hiddenLayerActivations)

		var dOutput mat.Dense
		dOutput.MulElem(&networkError, &slopeOutputLayer)
		var errorAtHiddenLayer mat.Dense
		errorAtHiddenLayer.Mul(&dOutput, wOut.T())

		var dHiddenLayer mat.Dense
		dHiddenLayer.MulElem(&errorAtHiddenLayer, &slopeHiddenLayer)

		// 调整参数。
		var wOutAdj mat.Dense
		wOutAdj.Mul(hiddenLayerActivations.T(), &dOutput)
		wOutAdj.Scale(nn.config.learningRate, &wOutAdj)
		wOut.Add(wOut, &wOutAdj)

		bOutAdj, err := sumAlongAxis(0, &dOutput)
		if err != nil {
			return err
		}
		bOutAdj.Scale(nn.config.learningRate, bOutAdj)
		bOut.Add(bOut, bOutAdj)

		var wHiddenAdj mat.Dense
		wHiddenAdj.Mul(x.T(), &dHiddenLayer)
		wHiddenAdj.Scale(nn.config.learningRate, &wHiddenAdj)
		wHidden.Add(wHidden, &wHiddenAdj)

		bHiddenAdj, err := sumAlongAxis(0, &dHiddenLayer)
		if err != nil {
			return err
		}
		bHiddenAdj.Scale(nn.config.learningRate, bHiddenAdj)
		bHidden.Add(bHidden, bHiddenAdj)
	}

	// 定义我们训练有素的神经网络。
	nn.wHidden = wHidden
	nn.bHidden = bHidden
	nn.wOut = wOut
	nn.bOut = bOut

	return nil
}

//预测基于经过训练的神经网络进行预测。
func (nn *neuralNet) predict(x *mat.Dense) (*mat.Dense, error) {

	//检查以确保我们的NeuroNet值代表训练有素的模型。
	if nn.wHidden == nil || nn.wOut == nil || nn.bHidden == nil || nn.bOut == nil {
		return nil, errors.New("the supplied neurnal net weights and biases are empty")
	}

	// 定义神经网络的输出。
	var output mat.Dense

	// 完成前馈过程。
	var hiddenLayerInput mat.Dense
	hiddenLayerInput.Mul(x, nn.wHidden)
	addBHidden := func(_, col int, v float64) float64 { return v + nn.bHidden.At(0, col) }
	hiddenLayerInput.Apply(addBHidden, &hiddenLayerInput)

	var hiddenLayerActivations mat.Dense
	applySigmoid := func(_, _ int, v float64) float64 { return sigmoid(v) }
	hiddenLayerActivations.Apply(applySigmoid, &hiddenLayerInput)

	var outputLayerInput mat.Dense
	outputLayerInput.Mul(&hiddenLayerActivations, nn.wOut)
	addBOut := func(_, col int, v float64) float64 { return v + nn.bOut.At(0, col) }
	outputLayerInput.Apply(addBOut, &outputLayerInput)
	output.Apply(applySigmoid, &outputLayerInput)

	return &output, nil
}

// sigmoid实现了Sigmoid函数以用于激活函数。
func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

// sigmoidPrime实现sigmoid函数的导数以进行反向传播。
func sigmoidPrime(x float64) float64 {
	return sigmoid(x) * (1.0 - sigmoid(x))
}


// sumAlongAxis对沿特定维度的矩阵求和，并保留另一个维度。
func sumAlongAxis(axis int, m *mat.Dense) (*mat.Dense, error) {

	numRows, numCols := m.Dims()

	var output *mat.Dense

	switch axis {
	case 0:
		data := make([]float64, numCols)
		for i := 0; i < numCols; i++ {
			col := mat.Col(nil, i, m)
			data[i] = floats.Sum(col)
		}
		output = mat.NewDense(1, numCols, data)
	case 1:
		data := make([]float64, numRows)
		for i := 0; i < numRows; i++ {
			row := mat.Row(nil, i, m)
			data[i] = floats.Sum(row)
		}
		output = mat.NewDense(numRows, 1, data)
	default:
		return nil, errors.New("invalid axis, must be 0 or 1")
	}

	return output, nil
}

