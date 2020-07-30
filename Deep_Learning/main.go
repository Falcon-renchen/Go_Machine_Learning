package main

import (
	"archive/zip"
	"bufio"
	"flag"
	"fmt"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/tensorflow/tensorflow/tensorflow/go/op"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	//使用TensorFlow Go API通过预训练的初始模型（http://arxiv.org/abs/1512.00567）进行图像识别的示例。
	//示例用法：<program> -dir = / tmp / modeldir -image = / path / to / some / jpeg
	//预先训练的模型采用形状为[BATCH_SIZE，IMAGE_HEIGHT，IMAGE_WIDTH，3]的4维张量输入，
	//其中：-BATCH_SIZE允许一次通过图形推断多个图像-IMAGE_HEIGHT是训练模型的图像的高度-IMAGE_WIDTH是训练模型的图像的宽度-3是（R ，G，B）
	//表示为浮点数的像素颜色值。并输出形状为[NUM_LABELS]的向量。
	//output [i]是输入图像被识别为具有第i个标签的概率。
	//一个单独的文件包含与输出的整数索引相对应的字符串标签列表。此示例：
	//-将预训练模型的序列化表示形式加载到图形中-创建会话以在图形上执行操作
	//-将图像文件转换为张量以提供给会话运行以作为输入
	//-执行会话并打印出来可能性最高的标签。
	//要将图像文件转换为适合输入到Inception模型的张量，
	//此示例：-构造另一个TensorFlow图以将图像规范化为适合模型的形式（
	//例如，调整图像大小）-创建并执行一个Session以获得此规范化形式的Tensor。
	modeldir := flag.String("dir", "", "Directory containing the trained model files. The directory will be created and the model downloaded into it if necessary")
	imagefile := flag.String("image", "", "Path of a JPEG-image to extract labels for")
	flag.Parse()
	if *modeldir == "" || *imagefile == "" {
		flag.Usage()
		return
	}
	//1. 要加载模型，我们执行以下操作：

	// 从文件加载序列化的GraphDef。
	modelfile, labelsfile, err := modelFiles(*modeldir)
	if err != nil {
		log.Fatal(err)
	}
	model, err := ioutil.ReadFile(modelfile)
	if err != nil {
		log.Fatal(err)
	}

	//2. 加载深度学习模型的图定义，并使用该图创建一个新的TensorFlow会话，如以下代码所示：

	// 从序列化的形式构造一个内存图。
	graph := tf.NewGraph()
	if err := graph.Import(model, ""); err != nil {
		log.Fatal(err)
	}

	// 创建一个会话以进行图推理。
	session, err := tf.NewSession(graph, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	//3. 使用该模型进行如下推断

	//在* imageFile上运行推理。对于多个图像，可以在循环中（并发地）调用session.Run（）。
	//或者，由于模型接受成批的图像数据作为输入，因此可以成批处理图像。
	tensor, err := makeTensorFromImage(*imagefile)
	if err != nil {
		log.Fatal(err)
	}
	output, err := session.Run(
		map[tf.Output]*tf.Tensor{
			graph.Operation("input").Output(0): tensor,
		},
		[]tf.Output{
			graph.Operation("output").Output(0),
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}


	// output [0] .Value（）是一个包含以下概率的向量
	//“批处理”中每个图像的标签。批次大小为1.查找最可能的标签索引。
	probabilities := output[0].Value().([][]float32)[0]
	printBestLabel(probabilities, labelsfile)
}

func printBestLabel(probabilities []float32, labelsFile string) {
	bestIdx := 0
	for i, p := range probabilities {
		if p > probabilities[bestIdx] {
			bestIdx = i
		}
	}

	//找到最佳匹配。从labelsFile读取字符串，每个标签包含一行。
	file, err := os.Open(labelsFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var labels []string
	for scanner.Scan() {
		labels = append(labels, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Printf("ERROR: failed to read %s: %v", labelsFile, err)
	}
	fmt.Printf("BEST MATCH: (%2.0f%% likely) %s\n", probabilities[bestIdx]*100.0, labels[bestIdx])
}


//将文件名中的图像转换为适合作为Inception模型输入的张量。
func makeTensorFromImage(filename string) (*tf.Tensor, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// DecodeJpeg使用标量字符串值张量作为输入。
	tensor, err := tf.NewTensor(string(bytes))
	if err != nil {
		return nil, err
	}
	//构造图以标准化图像
	graph, input, output, err := constructGraphToNormalizeImage()
	if err != nil {
		return nil, err
	}
	//执行该图以标准化此一张图片
	session, err := tf.NewSession(graph, nil)
	if err != nil {
		return nil, err
	}
	defer session.Close()
	normalized, err := session.Run(
		map[tf.Output]*tf.Tensor{input: tensor},
		[]tf.Output{output},
		nil)
	if err != nil {
		return nil, err
	}
	return normalized[0], nil
}

//初始模型将张量描述的图像作为输入
//特定的归一化格式（特定的图像大小，输入张量的形状，归一化的像素值等）。
//此函数构造一个TensorFlow操作图，该图将JPEG编码的字符串作为输入，并将适合作为输入的张量返回给
//初始模型。
func constructGraphToNormalizeImage() (graph *tf.Graph, input, output tf.Output, err error) {
	//一些特定于预训练模型的常量，位于：https://storage.googleapis.com/download.tensorflow.org/models/inception5h.zip
	// -模型经过训练后，图像缩放为224x224像素。 -使用（值-平均值）/比例将颜色（分别以1字节表示为R，G，B）转换为浮点型。
	const (
		H, W  = 224, 224
		Mean  = float32(117)
		Scale = float32(1)
	)
	//-输入是一个字符串张量，其中字符串是JPEG编码的图像。
	//-初始模型采用4D张量形状
	// [BatchSize，Height，Width，Colors = 3]，其中每个像素为
	//表示为三重浮点数
	//-在每个像素上应用归一化并使用ExpandDims
	//此单个图像是ResizeBilinear的大小为1的“批”。
	s := op.NewScope()
	input = op.Placeholder(s, tf.String)
	output = op.Div(s,
		op.Sub(s,
			op.ResizeBilinear(s,
				op.ExpandDims(s,
					op.Cast(s,
						op.DecodeJpeg(s, input, op.DecodeJpegChannels(3)), tf.Float),
					op.Const(s.SubScope("make_batch"), int32(0))),
				op.Const(s.SubScope("size"), []int32{H, W})),
			op.Const(s.SubScope("mean"), Mean)),
		op.Const(s.SubScope("scale"), Scale))
	graph, err = s.Finalize()
	return graph, input, output, err
}

func modelFiles(dir string) (modelfile, labelsfile string, err error) {
	const URL = "https://storage.googleapis.com/download.tensorflow.org/models/inception5h.zip"
	var (
		model   = filepath.Join(dir, "tensorflow_inception_graph.pb")
		labels  = filepath.Join(dir, "imagenet_comp_graph_label_strings.txt")
		zipfile = filepath.Join(dir, "inception5h.zip")
	)
	if filesExist(model, labels) == nil {
		return model, labels, nil
	}
	log.Println("Did not find model in", dir, "downloading from", URL)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", "", err
	}
	if err := download(URL, zipfile); err != nil {
		return "", "", fmt.Errorf("failed to download %v - %v", URL, err)
	}
	if err := unzip(dir, zipfile); err != nil {
		return "", "", fmt.Errorf("failed to extract contents from model archive: %v", err)
	}
	os.Remove(zipfile)
	return model, labels, filesExist(model, labels)
}

func filesExist(files ...string) error {
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			return fmt.Errorf("unable to stat %s: %v", f, err)
		}
	}
	return nil
}

func download(URL, filename string) error {
	resp, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	return err
}

func unzip(dir, zipfile string) error {
	r, err := zip.OpenReader(zipfile)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		src, err := f.Open()
		if err != nil {
			return err
		}
		log.Println("Extracting", f.Name)
		dst, err := os.OpenFile(filepath.Join(dir, f.Name), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		if _, err := io.Copy(dst, src); err != nil {
			return err
		}
		dst.Close()
	}
	return nil
}
