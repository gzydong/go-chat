package filesystem

import (
	"testing"
)

func TestName(t *testing.T) {
	// conf := testutil.GetConfig()
	//
	// conf.Filesystem.Local.SSL = false
	// conf.Filesystem.Local.Root = "./data"
	// conf.Filesystem.Local.Endpoint = "127.0.0.1:9000"
	// conf.Filesystem.Local.BucketPublic = "im-static"
	// conf.Filesystem.Local.BucketPrivate = "im-private"
	//
	// conf.Filesystem.Minio.SSL = false
	// conf.Filesystem.Minio.Endpoint = "127.0.0.1:9000"
	// conf.Filesystem.Minio.SecretId = "Q3AM3UQ867SPQQA43P2F"
	// conf.Filesystem.Minio.SecretKey = "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
	// conf.Filesystem.Minio.BucketPublic = "im-static"
	// conf.Filesystem.Minio.BucketPrivate = "im-private"
	//
	// client := NewMinioFilesystem(conf)
	// client := NewLocalFilesystem(conf)

	// err := client.Write("im-private", []byte("hello world"), "filesystem.txt")
	// fmt.Println(err)

	// bt, err := client.GetObject("im-private", "filesystem.txt")
	// fmt.Println(err)
	// fmt.Println(string(bt))
	//
	// err := client.WriteLocal("im-private", "./util.go", "filesystem.txt")
	// fmt.Println(err)

	// err := client.Copy("im-private", "filesystem.txt", "filesystem2.txt")
	// fmt.Println(err)

	// bt, err := client.Stat("im-private", "filesystem2.txt")
	// fmt.Println(err)
	// fmt.Println(bt)

	// value := client.PublicUrl("im-private", "filesystem2.txt")
	// fmt.Println(value)

	// // Make a buffer with 6MB of data
	// buf := bytes.Repeat([]byte("abcdef"), 1024*1024)
	//
	// // Open the file.
	// file, err := os.Open("./node-v18.15.0-linux-x64.tar.xz")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	//
	// defer file.Close()
	//
	// items := make([]ObjectPart, 0)
	//
	// upload, err := client.InitiateMultipartUpload("im-private", "node-v18.15.0-linux-x64.txt")
	// fmt.Println(upload, err)
	//
	// if err != nil {
	// 	return
	// }
	//
	// obj, err := client.PutObjectPart("im-private", "node-v18.15.0-linux-x64.txt", upload, 1, bytes.NewReader(buf[:5*1024*1024]), 5*1024*1024)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	//
	// items = append(items, obj)
	//
	// obj, err = client.PutObjectPart("im-private", "node-v18.15.0-linux-x64.txt", upload, 2, bytes.NewReader(buf[5*1024*1024:]), 1024*1024)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	//
	// items = append(items, obj)
	//
	// // Close the file.
	// err = client.CompleteMultipartUpload("im-private", "node-v18.15.0-linux-x64.txt", upload, items)
	// fmt.Println(err)
}
