package storageapplication

import "strconv"

//Video2ImageService convert video to images
type Video2ImageService struct {
	imageProcessingService ImageProcessingService
}

//ProcessVideo ...
func (v2i *Video2ImageService) ProcessVideo(videoData []byte, fileName string) error {

	images, err := video2Images(videoData)
	if err != nil {
		return err
	}
	for i, image := range images {
		v2i.imageProcessingService.ProcessImage(image, fileName+"_"+strconv.Itoa(i))
	}
	return nil
}

func video2Images(videoData []byte) ([][]byte, error) {
	var arrayofarrays [][]byte
	//TODO convert video to Images
	return arrayofarrays, nil
}

//NewVideo2ImageService ...
func NewVideo2ImageService() *Video2ImageService {
	return &Video2ImageService{}
}