package portin

//InputPort port in for add image to analize
type InputPort interface {
	ProcessImage(imgData []byte) (string, error)
	//processVideo(videoData []byte) (string, error)

	//processImagePath(imgPath string) (string, error)
	//processVideoPath(videoPath string) (string, error)
}