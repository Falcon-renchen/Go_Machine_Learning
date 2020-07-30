package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

//检测输入和输出之间的关系和规则。这就是我们使用神经网络所做的！
//神经网络的节点就像子函数一样，可以调整直到模仿我们提供的输入和输出之间的规则和关系，
//反向传播


//wHidden和bHidden是网络隐层的权重和偏差，
//并且wOut和bOut是网络输出层的权重和偏差

// NeuroNet包含定义受过训练的神经网络的所有信息。
type neuralNet struct {
	config  neuralNetConfig
	wHidden *mat.Dense		//网络隐层的权重
	bHidden *mat.Dense		//网络隐层的偏差
	wOut    *mat.Dense		//网络输出层的权重
	bOut    *mat.Dense		//网络输出层的偏差
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
	//定义我们的输入属性。
	input := mat.NewDense(3, 4, []float64{
		1.0, 0.0, 1.0, 0.0,
		1.0, 0.0, 1.0, 1.0,
		0.0, 1.0, 0.0, 1.0,
	})

	//定义我们的标签。
	labels := mat.NewDense(3, 1, []float64{1.0, 1.0, 0.0})


	//定义我们的网络架构和学习参数。
	config := neuralNetConfig{
		inputNeurons:  4,
		outputNeurons: 1,
		hiddenNeurons: 3,
		numEpochs:     5000,
		learningRate:  0.3,
	}

	// 训练神经网络。
	network := newNetwork(config)
	if err := network.train(input, labels); err != nil {
		log.Fatal(err)
	}

	//输出定义我们网络的权重！
	f := mat.Formatted(network.wHidden, mat.Prefix("          "))
	fmt.Printf("\nwHidden = % v\n\n", f)

	f = mat.Formatted(network.bHidden, mat.Prefix("    "))
	fmt.Printf("\nbHidden = % v\n\n", f)

	f = mat.Formatted(network.wOut, mat.Prefix("       "))
	fmt.Printf("\nwOut = % v\n\n", f)

	f = mat.Formatted(network.bOut, mat.Prefix("    "))
	fmt.Printf("\nbOut = % v\n\n", f)
}

// NewNetwork初始化一个新的神经网络。
func newNetwork(config neuralNetConfig) *neuralNet {
	return &neuralNet{config: config}
}

// Train使用反向传播训练神经网络。
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
		addBHidden := func(_, col int, v float64) float64 {
			return v + bHidden.At(0, col)
		}
		hiddenLayerInput.Apply(addBHidden, &hiddenLayerInput)

		var hiddenLayerActivations mat.Dense
		applySigmoid := func(_, _ int, v float64) float64 {
			return sigmoid(v)
		}
		hiddenLayerActivations.Apply(applySigmoid, &hiddenLayerInput)

		var outputLayerInput mat.Dense
		outputLayerInput.Mul(&hiddenLayerActivations, wOut)
		addBOut := func(_, col int, v float64) float64 {
			return v + bOut.At(0, col)
		}
		outputLayerInput.Apply(addBOut, &outputLayerInput)
		output.Apply(applySigmoid, &outputLayerInput)

		// 完成反向传播。
		var networkError mat.Dense
		networkError.Sub(y, &output)

		var slopeOutputLayer mat.Dense
		applySigmoidPrime := func(_, _ int, v float64) float64 {
			return sigmoidPrime(v)
		}
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

	//定义我们训练有素的神经网络。
	nn.wHidden = wHidden
	nn.bHidden = bHidden
	nn.wOut = wOut
	nn.bOut = bOut

	return nil
}


// sigmoid实现了Sigmoid函数以用于激活函数
func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}


// sigmoidPrime实现sigmoid函数的导数以进行反向传播。
func sigmoidPrime(x float64) float64 {
	return sigmoid(x) * (1.0 - sigmoid(x))
}


// sumAlongAxis对沿特定维度的矩阵求和，并保留 其他维度。
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
